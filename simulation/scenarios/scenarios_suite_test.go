package main_test

/* this is meant to be run on a Diego bosh node */

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
	"sync"

	"github.com/cloudfoundry-incubator/auction/communication/http/auction_http_client"
	"github.com/pivotal-golang/lager"

	"github.com/cloudfoundry-incubator/auction/auctionrunner"
	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"github.com/cloudfoundry-incubator/auction/simulation/auctiondistributor"
	"github.com/cloudfoundry-incubator/auction/simulation/visualization"
	"github.com/cloudfoundry-incubator/auction/util"
	"github.com/cloudfoundry/gunk/workpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"
)

var numCells int
var numAuctioneers int
var concurrentAuctionsPerAuctioneer int
var timeout time.Duration

var auctionDistributor auctiondistributor.AuctionDistributor

var svgReport *visualization.SVGReport
var reports []*visualization.Report

var client auctiontypes.SimulationRepPoolClient
var repAddresses []auctiontypes.RepAddress
var reportName string

func init() {
	flag.IntVar(&numCells, "numCells", 100, "the number of cells")
	flag.IntVar(&numAuctioneers, "numAuctioneers", 0, "the number of auctioneers (0 means use the number of cells)")
	flag.IntVar(&concurrentAuctionsPerAuctioneer, "maxConcurrent", 20, "the maximum number of concurrent auctions to run, per auctioneer")
	flag.Float64Var(&(auctionrunner.DefaultStartAuctionRules.MaxBiddingPoolFraction), "maxBiddingPoolFraction", auctionrunner.DefaultStartAuctionRules.MaxBiddingPoolFraction, "the maximum number of participants in the pool")
	flag.DurationVar(&timeout, "timeout", time.Second, "timeout when waiting for responses from remote calls")
	flag.StringVar(&(auctionrunner.DefaultStartAuctionRules.Algorithm), "algorithm", auctionrunner.DefaultStartAuctionRules.Algorithm, "the auction algorithm to use")
}

func TestAuction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auction Suite")
}

var _ = BeforeSuite(func() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	reportName = fmt.Sprintf("%s-%dcells-%dconc-%.2fpool", auctionrunner.DefaultStartAuctionRules.Algorithm, numCells, concurrentAuctionsPerAuctioneer, auctionrunner.DefaultStartAuctionRules.MaxBiddingPoolFraction)
	if numAuctioneers == 0 {
		numAuctioneers = numCells
	}

	startReport()

	auctioneers := []string{}
	repAddresses = []auctiontypes.RepAddress{}
	for i := 1; i <= numCells; i++ {
		repAddresses = append(repAddresses, auctiontypes.RepAddress{
			RepGuid: fmt.Sprintf("rep-lite-%d", i),
			Address: fmt.Sprintf("http://rep-lite-%d.diego-1.cf-app.com", i),
		})
	}
	for i := 1; i <= numAuctioneers; i++ {
		auctioneers = append(auctioneers, fmt.Sprintf("auctioneer-lite-%d.diego-1.cf-app.com", i))
	}
	client = auction_http_client.New(http.DefaultClient, lager.NewLogger("client"))

	auctionDistributor = auctiondistributor.NewExternalAuctionDistributor(auctioneers, concurrentAuctionsPerAuctioneer, "http")
})

var _ = BeforeEach(func() {
	workers := workpool.NewWorkPool(50)
	wg := &sync.WaitGroup{}
	wg.Add(len(repAddresses))
	for _, repAddress := range repAddresses {
		repAddress := repAddress
		workers.Submit(func() {
			client.Reset(repAddress)
			wg.Done()
		})
	}

	wg.Wait()
	workers.Stop()

	util.ResetGuids()
})

var _ = AfterSuite(func() {
	finishReport()
})

func startReport() {
	svgReport = visualization.StartSVGReport("./"+reportName+".svg", 2, 1)
	svgReport.DrawHeader("Diego Scenario", auctionrunner.DefaultStartAuctionRules, concurrentAuctionsPerAuctioneer)
}

func finishReport() {
	svgReport.Done()
	_, err := exec.LookPath("rsvg-convert")
	if err == nil {
		exec.Command("rsvg-convert", "-h", "2000", "--background-color=#fff", "./"+reportName+".svg", "-o", "./"+reportName+".png").Run()
		// exec.Command("open", "./"+reportName+".png").Run()
	}
	data, err := json.Marshal(reports)
	Ω(err).ShouldNot(HaveOccurred())
	ioutil.WriteFile("./"+reportName+".json", data, 0777)

	summaryBytes, err := ioutil.ReadFile("./summary.csv")
	summary := string(summaryBytes)
	if err != nil {
		summary = "numCells,numAuctioneers,concurrentAuctionsPerAuctioneer,maxBiddingPoolFraction,algorithm,scenario,# auctions,communication,waitTime,biddingTime,distributionScore,nMissing\n"
	}

	summary += fmt.Sprintf("%d,%d,%d,%.2f,%s,%s,%d,%d,%.2f,%.2f,%.4f,%d\n",
		numCells,
		numAuctioneers,
		concurrentAuctionsPerAuctioneer,
		auctionrunner.DefaultStartAuctionRules.MaxBiddingPoolFraction,
		auctionrunner.DefaultStartAuctionRules.Algorithm,
		"10% start",
		reports[0].NAuctions(),
		int64(reports[0].CommStats().Total),
		reports[0].WaitTimeStats().Max,
		reports[0].BiddingTimeStats().Max,
		reports[0].DistributionScore(),
		reports[0].NMissingInstances(),
	)

	summary += fmt.Sprintf("%d,%d,%d,%.2f,%s,%s,%d,%d,%.2f,%.2f,%.4f,%d\n",
		numCells,
		numAuctioneers,
		concurrentAuctionsPerAuctioneer,
		auctionrunner.DefaultStartAuctionRules.MaxBiddingPoolFraction,
		auctionrunner.DefaultStartAuctionRules.Algorithm,
		"cold start",
		reports[1].NAuctions(),
		int64(reports[1].CommStats().Total),
		reports[1].WaitTimeStats().Max,
		reports[1].BiddingTimeStats().Max,
		reports[1].DistributionScore(),
		reports[1].NMissingInstances(),
	)

	ioutil.WriteFile("./summary.csv", []byte(summary), 0666)
}
