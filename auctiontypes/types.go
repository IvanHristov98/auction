package auctiontypes

import (
	"errors"
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/tedsuo/ifrit"
)

var ErrorStackMismatch = errors.New("stack mismatch")
var ErrorInsufficientResources = errors.New("insuccifient resources")
var ErrorNothingToStop = errors.New("nothing to stop")

type AuctionRunner interface {
	ifrit.Runner
	AddLRPStartAuction(models.LRPStartAuction)
	AddLRPStopAuction(models.LRPStopAuction)
}

type AuctionRunnerDelegate interface {
	FetchAuctionRepClients() (map[string]AuctionRep, error)
	DistributedBatch(AuctionResults)
}

type AuctionResults struct {
	SuccessfulStarts []StartAuction
	SuccessfulStops  []StopAuction
	FailedStarts     []StartAuction
	FailedStops      []StopAuction
}

type StartAuction struct {
	LRPStartAuction models.LRPStartAuction
	Winner          string
	Attempts        int

	QueueTime    time.Time
	WaitDuration time.Duration
}

func (s StartAuction) Identifier() string {
	return IdentifierForLRPStartAuction(s.LRPStartAuction)
}

func IdentifierForLRPStartAuction(start models.LRPStartAuction) string {
	return fmt.Sprintf("%s.%d.%s", start.DesiredLRP.ProcessGuid, start.Index, start.InstanceGuid)
}

type StopAuction struct {
	LRPStopAuction models.LRPStopAuction
	Winner         string
	Attempts       int

	QueueTime    time.Time
	WaitDuration time.Duration
}

func (s StopAuction) Identifier() string {
	return fmt.Sprintf("%s.%d", s.LRPStopAuction.ProcessGuid, s.LRPStopAuction.Index)
}

type RepAddress struct {
	RepGuid string
	Address string
}

type AuctionRep interface {
	State() (RepState, error)
	Perform(Work) (Work, error)
}

//simulation-only interface
type SimulationAuctionRep interface {
	AuctionRep

	Reset() error
}

type Work struct {
	Starts []models.LRPStartAuction
	Stops  []models.StopLRPInstance
}

type RepState struct {
	Stack              string
	AvailableResources Resources
	TotalResources     Resources
	LRPs               []LRP
}

type LRP struct {
	ProcessGuid  string
	InstanceGuid string
	Index        int
	MemoryMB     int
	DiskMB       int
}

type Resources struct {
	DiskMB     int
	MemoryMB   int
	Containers int
}
