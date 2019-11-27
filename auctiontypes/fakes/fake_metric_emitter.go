// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"
	"time"

	"code.cloudfoundry.org/auction/auctiontypes"
)

type FakeAuctionMetricEmitterDelegate struct {
	AuctionCompletedStub        func(auctiontypes.AuctionResults)
	auctionCompletedMutex       sync.RWMutex
	auctionCompletedArgsForCall []struct {
		arg1 auctiontypes.AuctionResults
	}
	FailedCellStateRequestStub        func()
	failedCellStateRequestMutex       sync.RWMutex
	failedCellStateRequestArgsForCall []struct {
	}
	FetchStatesCompletedStub        func(time.Duration) error
	fetchStatesCompletedMutex       sync.RWMutex
	fetchStatesCompletedArgsForCall []struct {
		arg1 time.Duration
	}
	fetchStatesCompletedReturns struct {
		result1 error
	}
	fetchStatesCompletedReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeAuctionMetricEmitterDelegate) AuctionCompleted(arg1 auctiontypes.AuctionResults) {
	fake.auctionCompletedMutex.Lock()
	fake.auctionCompletedArgsForCall = append(fake.auctionCompletedArgsForCall, struct {
		arg1 auctiontypes.AuctionResults
	}{arg1})
	fake.recordInvocation("AuctionCompleted", []interface{}{arg1})
	auctionCompletedStubCopy := fake.AuctionCompletedStub
	fake.auctionCompletedMutex.Unlock()
	if auctionCompletedStubCopy != nil {
		auctionCompletedStubCopy(arg1)
	}
}

func (fake *FakeAuctionMetricEmitterDelegate) AuctionCompletedCallCount() int {
	fake.auctionCompletedMutex.RLock()
	defer fake.auctionCompletedMutex.RUnlock()
	return len(fake.auctionCompletedArgsForCall)
}

func (fake *FakeAuctionMetricEmitterDelegate) AuctionCompletedCalls(stub func(auctiontypes.AuctionResults)) {
	fake.auctionCompletedMutex.Lock()
	defer fake.auctionCompletedMutex.Unlock()
	fake.AuctionCompletedStub = stub
}

func (fake *FakeAuctionMetricEmitterDelegate) AuctionCompletedArgsForCall(i int) auctiontypes.AuctionResults {
	fake.auctionCompletedMutex.RLock()
	defer fake.auctionCompletedMutex.RUnlock()
	argsForCall := fake.auctionCompletedArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeAuctionMetricEmitterDelegate) FailedCellStateRequest() {
	fake.failedCellStateRequestMutex.Lock()
	fake.failedCellStateRequestArgsForCall = append(fake.failedCellStateRequestArgsForCall, struct {
	}{})
	fake.recordInvocation("FailedCellStateRequest", []interface{}{})
	failedCellStateRequestStubCopy := fake.FailedCellStateRequestStub
	fake.failedCellStateRequestMutex.Unlock()
	if failedCellStateRequestStubCopy != nil {
		failedCellStateRequestStubCopy()
	}
}

func (fake *FakeAuctionMetricEmitterDelegate) FailedCellStateRequestCallCount() int {
	fake.failedCellStateRequestMutex.RLock()
	defer fake.failedCellStateRequestMutex.RUnlock()
	return len(fake.failedCellStateRequestArgsForCall)
}

func (fake *FakeAuctionMetricEmitterDelegate) FailedCellStateRequestCalls(stub func()) {
	fake.failedCellStateRequestMutex.Lock()
	defer fake.failedCellStateRequestMutex.Unlock()
	fake.FailedCellStateRequestStub = stub
}

func (fake *FakeAuctionMetricEmitterDelegate) FetchStatesCompleted(arg1 time.Duration) error {
	fake.fetchStatesCompletedMutex.Lock()
	ret, specificReturn := fake.fetchStatesCompletedReturnsOnCall[len(fake.fetchStatesCompletedArgsForCall)]
	fake.fetchStatesCompletedArgsForCall = append(fake.fetchStatesCompletedArgsForCall, struct {
		arg1 time.Duration
	}{arg1})
	fake.recordInvocation("FetchStatesCompleted", []interface{}{arg1})
	fetchStatesCompletedStubCopy := fake.FetchStatesCompletedStub
	fake.fetchStatesCompletedMutex.Unlock()
	if fetchStatesCompletedStubCopy != nil {
		return fetchStatesCompletedStubCopy(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.fetchStatesCompletedReturns
	return fakeReturns.result1
}

func (fake *FakeAuctionMetricEmitterDelegate) FetchStatesCompletedCallCount() int {
	fake.fetchStatesCompletedMutex.RLock()
	defer fake.fetchStatesCompletedMutex.RUnlock()
	return len(fake.fetchStatesCompletedArgsForCall)
}

func (fake *FakeAuctionMetricEmitterDelegate) FetchStatesCompletedCalls(stub func(time.Duration) error) {
	fake.fetchStatesCompletedMutex.Lock()
	defer fake.fetchStatesCompletedMutex.Unlock()
	fake.FetchStatesCompletedStub = stub
}

func (fake *FakeAuctionMetricEmitterDelegate) FetchStatesCompletedArgsForCall(i int) time.Duration {
	fake.fetchStatesCompletedMutex.RLock()
	defer fake.fetchStatesCompletedMutex.RUnlock()
	argsForCall := fake.fetchStatesCompletedArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeAuctionMetricEmitterDelegate) FetchStatesCompletedReturns(result1 error) {
	fake.fetchStatesCompletedMutex.Lock()
	defer fake.fetchStatesCompletedMutex.Unlock()
	fake.FetchStatesCompletedStub = nil
	fake.fetchStatesCompletedReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeAuctionMetricEmitterDelegate) FetchStatesCompletedReturnsOnCall(i int, result1 error) {
	fake.fetchStatesCompletedMutex.Lock()
	defer fake.fetchStatesCompletedMutex.Unlock()
	fake.FetchStatesCompletedStub = nil
	if fake.fetchStatesCompletedReturnsOnCall == nil {
		fake.fetchStatesCompletedReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.fetchStatesCompletedReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeAuctionMetricEmitterDelegate) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.auctionCompletedMutex.RLock()
	defer fake.auctionCompletedMutex.RUnlock()
	fake.failedCellStateRequestMutex.RLock()
	defer fake.failedCellStateRequestMutex.RUnlock()
	fake.fetchStatesCompletedMutex.RLock()
	defer fake.fetchStatesCompletedMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeAuctionMetricEmitterDelegate) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ auctiontypes.AuctionMetricEmitterDelegate = new(FakeAuctionMetricEmitterDelegate)
