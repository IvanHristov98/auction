package auction_http_client_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/runtime-schema/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run", func() {
	var lrpStartAuction models.LRPStartAuction

	BeforeEach(func() {
		lrpStartAuction = models.LRPStartAuction{
			InstanceGuid: "instance-guid",
			Index:        1,
		}

		auctionRepA.RunReturns(nil)
		auctionRepB.RunReturns(errors.New("oops"))
	})

	It("should tell the rep to run ", func() {
		client.Run("A", lrpStartAuction)

		Ω(auctionRepA.RunCallCount()).Should(Equal(1))
		Ω(auctionRepA.RunArgsForCall(0)).Should(Equal(lrpStartAuction))

		client.Run("B", lrpStartAuction)

		Ω(auctionRepB.RunCallCount()).Should(Equal(1))
		Ω(auctionRepB.RunArgsForCall(0)).Should(Equal(lrpStartAuction))

		//these should not panic/explode
		client.Run("RepThat500s", lrpStartAuction)
		client.Run("RepThatErrors", lrpStartAuction)

	})

	PIt("what about errors?", func() {

	})
})
