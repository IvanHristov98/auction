package auctionrunner

import (
	"sync"

	"github.com/cloudfoundry-incubator/auction/auctiontypes"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/pivotal-golang/clock"
)

type Batch struct {
	volAuctions  []auctiontypes.VolumeAuction
	lrpAuctions  []auctiontypes.LRPAuction
	taskAuctions []auctiontypes.TaskAuction
	lock         *sync.Mutex
	HasWork      chan struct{}
	clock        clock.Clock
}

func NewBatch(clock clock.Clock) *Batch {
	return &Batch{
		volAuctions: []auctiontypes.VolumeAuction{},
		lrpAuctions: []auctiontypes.LRPAuction{},
		lock:        &sync.Mutex{},
		clock:       clock,
		HasWork:     make(chan struct{}, 1),
	}
}

func (b *Batch) AddVolumeStarts(starts []models.VolumeStartRequest) {
	auctions := []auctiontypes.VolumeAuction{}
	now := b.clock.Now()
	for _, start := range starts {
		for _, i := range start.Indices {
			vs := start.VolumeSet
			auctions = append(auctions, auctiontypes.VolumeAuction{
				VolumeSetGuid:    vs.VolumeSetGuid,
				Stack:            vs.Stack,
				SizeMB:           vs.SizeMB,
				ReservedMemoryMB: vs.ReservedMemoryMB,
				Index:            int(i),
				AuctionRecord: auctiontypes.AuctionRecord{
					QueueTime: now,
				}})
		}
	}

	b.lock.Lock()
	b.volAuctions = append(b.volAuctions, auctions...)
	b.claimToHaveWork()
	b.lock.Unlock()
}

func (b *Batch) AddLRPStarts(starts []models.LRPStartRequest) {
	auctions := make([]auctiontypes.LRPAuction, 0, len(starts))
	now := b.clock.Now()
	for _, start := range starts {
		for _, i := range start.Indices {
			auctions = append(auctions, auctiontypes.LRPAuction{
				DesiredLRP: start.DesiredLRP,
				Index:      int(i),
				AuctionRecord: auctiontypes.AuctionRecord{
					QueueTime: now,
				}})
		}
	}

	b.lock.Lock()
	b.lrpAuctions = append(b.lrpAuctions, auctions...)
	b.claimToHaveWork()
	b.lock.Unlock()
}

func (b *Batch) AddTasks(tasks []models.Task) {
	auctions := make([]auctiontypes.TaskAuction, 0, len(tasks))
	now := b.clock.Now()
	for _, t := range tasks {
		auctions = append(auctions, auctiontypes.TaskAuction{
			Task: t,
			AuctionRecord: auctiontypes.AuctionRecord{
				QueueTime: now,
			},
		})
	}

	b.lock.Lock()
	b.taskAuctions = append(b.taskAuctions, auctions...)
	b.claimToHaveWork()
	b.lock.Unlock()
}

func (b *Batch) DedupeAndDrain() ([]auctiontypes.VolumeAuction, []auctiontypes.LRPAuction, []auctiontypes.TaskAuction) {
	b.lock.Lock()
	volAuctions := b.volAuctions
	lrpAuctions := b.lrpAuctions
	taskAuctions := b.taskAuctions
	b.volAuctions = []auctiontypes.VolumeAuction{}
	b.lrpAuctions = []auctiontypes.LRPAuction{}
	b.taskAuctions = []auctiontypes.TaskAuction{}
	select {
	case <-b.HasWork:
	default:
	}
	b.lock.Unlock()

	dedupedVolAuctions := []auctiontypes.VolumeAuction{}
	presentVolAuctions := map[string]bool{}
	for _, startAuction := range volAuctions {
		id := startAuction.Identifier()
		if presentVolAuctions[id] {
			continue
		}
		presentVolAuctions[id] = true
		dedupedVolAuctions = append(dedupedVolAuctions, startAuction)
	}

	dedupedLRPAuctions := []auctiontypes.LRPAuction{}
	presentLRPAuctions := map[string]bool{}
	for _, startAuction := range lrpAuctions {
		id := startAuction.Identifier()
		if presentLRPAuctions[id] {
			continue
		}
		presentLRPAuctions[id] = true
		dedupedLRPAuctions = append(dedupedLRPAuctions, startAuction)
	}

	dedupedTaskAuctions := []auctiontypes.TaskAuction{}
	presentTaskAuctions := map[string]bool{}
	for _, taskAuction := range taskAuctions {
		id := taskAuction.Identifier()
		if presentTaskAuctions[id] {
			continue
		}
		presentTaskAuctions[id] = true
		dedupedTaskAuctions = append(dedupedTaskAuctions, taskAuction)
	}

	return dedupedVolAuctions, dedupedLRPAuctions, dedupedTaskAuctions
}

func (b *Batch) claimToHaveWork() {
	select {
	case b.HasWork <- struct{}{}:
	default:
	}
}
