package auctionrunner

import (
	"os"
	"time"

	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"

	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry/gunk/workpool"
)

type auctionRunner struct {
	delegate      auctiontypes.AuctionRunnerDelegate
	metricEmitter auctiontypes.AuctionMetricEmitterDelegate
	batch         *Batch
	clock         clock.Clock
	workPool      *workpool.WorkPool
	logger        lager.Logger
}

func New(
	delegate auctiontypes.AuctionRunnerDelegate,
	metricEmitter auctiontypes.AuctionMetricEmitterDelegate,
	clock clock.Clock,
	workPool *workpool.WorkPool,
	logger lager.Logger,
) *auctionRunner {
	return &auctionRunner{
		delegate:      delegate,
		metricEmitter: metricEmitter,
		batch:         NewBatch(clock),
		clock:         clock,
		workPool:      workPool,
		logger:        logger,
	}
}

func (a *auctionRunner) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	close(ready)

	var hasWork chan struct{}
	hasWork = a.batch.HasWork

	for {
		select {
		case <-hasWork:
			logger := a.logger.Session("auction")

			logger.Info("fetching-cell-reps")
			clients, err := a.delegate.FetchCellReps()
			if err != nil {
				logger.Error("failed-to-fetch-reps", err)
				time.Sleep(time.Second)
				hasWork = make(chan struct{}, 1)
				hasWork <- struct{}{}
				break
			}
			logger.Info("fetched-cell-reps", lager.Data{"cell-reps-count": len(clients)})

			hasWork = a.batch.HasWork

			logger.Info("fetching-zone-state")
			fetchStatesStartTime := time.Now()
			zones := FetchStateAndBuildZones(a.workPool, clients)
			a.metricEmitter.FetchStatesCompleted(time.Since(fetchStatesStartTime))
			cellCount := 0
			for zone, cells := range zones {
				logger.Info("zone-state", lager.Data{"zone": zone, "cell-count": len(cells)})
				cellCount += len(cells)
			}
			logger.Info("fetched-zone-state", lager.Data{"cell-state-count": cellCount, "num-failed-requests": len(clients) - cellCount})

			logger.Info("fetching-auctions")
			volAuctions, lrpAuctions, taskAuctions := a.batch.DedupeAndDrain()
			logger.Info("fetched-auctions", lager.Data{
				"lrp-start-auctions":    len(lrpAuctions),
				"task-auctions":         len(taskAuctions),
				"volume-start-auctions": len(volAuctions),
			})
			if len(lrpAuctions) == 0 && len(taskAuctions) == 0 && len(volAuctions) == 0 {
				logger.Info("nothing-to-auction")
				break
			}

			logger.Info("scheduling")
			//Breadcrumb
			auctionRequest := auctiontypes.AuctionRequest{
				Volumes: volAuctions,
				LRPs:    lrpAuctions,
				Tasks:   taskAuctions,
			}

			scheduler := NewScheduler(a.workPool, zones, a.clock)
			auctionResults := scheduler.Schedule(auctionRequest)
			logger.Info("scheduled", lager.Data{
				"successful-volume-start-auctions": len(auctionResults.SuccessfulVolumes),
				"successful-lrp-start-auctions":    len(auctionResults.SuccessfulLRPs),
				"successful-task-auctions":         len(auctionResults.SuccessfulTasks),
				"failed-volume-start-auctions":     len(auctionResults.FailedVolumes),
				"failed-lrp-start-auctions":        len(auctionResults.FailedLRPs),
				"failed-task-auctions":             len(auctionResults.FailedTasks),
			})
			numVolumesFailed := len(auctionResults.FailedVolumes)
			numStartsFailed := len(auctionResults.FailedLRPs)
			numTasksFailed := len(auctionResults.FailedTasks)

			logger.Info("resubmitted-failures", lager.Data{
				"successful-volume-start-auctions":     len(auctionResults.SuccessfulVolumes),
				"successful-lrp-start-auctions":        len(auctionResults.SuccessfulLRPs),
				"successful-task-auctions":             len(auctionResults.SuccessfulTasks),
				"will-not-retry-volume-start-auctions": len(auctionResults.FailedVolumes),
				"will-not-retry-lrp-start-auctions":    len(auctionResults.FailedLRPs),
				"will-not-retry-task-auctions":         len(auctionResults.FailedTasks),
				"will-retry-volume-start-auctions":     numVolumesFailed - len(auctionResults.FailedVolumes),
				"will-retry-lrp-start-auctions":        numStartsFailed - len(auctionResults.FailedLRPs),
				"will-retry-task-auctions":             numTasksFailed - len(auctionResults.FailedTasks),
			})

			a.metricEmitter.AuctionCompleted(auctionResults)
			a.delegate.AuctionCompleted(auctionResults)
		case <-signals:
			return nil
		}
	}
}

func (a *auctionRunner) ScheduleVolumesForAuctions(volumeStarts []models.VolumeStartRequest) {
	a.batch.AddVolumeStarts(volumeStarts)
}

func (a *auctionRunner) ScheduleLRPsForAuctions(lrpStarts []models.LRPStartRequest) {
	a.batch.AddLRPStarts(lrpStarts)
}

func (a *auctionRunner) ScheduleTasksForAuctions(tasks []models.Task) {
	a.batch.AddTasks(tasks)
}
