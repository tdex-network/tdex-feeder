package application

import (
	"errors"
	"sync"

	"github.com/tdex-network/tdex-feeder/internal/domain"
)

type TdexFeeder interface {
	Start() error
	Stop()
	IsRunning() bool
}

type tdexFeeder struct {
	feeds    []domain.Feed
	targets  []domain.Target
	stopChan chan bool
	running  bool
	locker   sync.Locker
}

func NewTdexFeeder(feeds []domain.Feed, targets []domain.Target) TdexFeeder {
	return &tdexFeeder{
		feeds:    feeds,
		targets:  targets,
		stopChan: make(chan bool),
		running:  false,
		locker:   &sync.Mutex{},
	}
}

// Start observe all the feeds chan (using merge function)
// and push the results to all targets
func (t *tdexFeeder) Start() error {
	if t.IsRunning() {
		return errors.New("the feeder is already started")
	}

	t.running = true
	marketPriceChannel := merge(t.feeds...)

	for t.IsRunning() {
		select {
		case <-t.stopChan:
			t.running = false
			break
		case marketPrice := <-marketPriceChannel:
			for _, target := range t.targets {
				target.Push(marketPrice)
			}
		}
	}

	return nil
}

func (t *tdexFeeder) Stop() {
	t.stopChan <- true
}

func (t *tdexFeeder) IsRunning() bool {
	t.locker.Lock()
	defer t.locker.Unlock()
	return t.running
}

// merge gathers several feeds into a unique channel
func merge(feeds ...domain.Feed) <-chan domain.MarketPrice {
	mergedChan := make(chan domain.MarketPrice)
	var wg sync.WaitGroup

	wg.Add(len(feeds))
	for _, feed := range feeds {
		c := feed.GetMarketPriceChan()
		go func(c <-chan domain.MarketPrice) {
			for marketPrice := range c {
				mergedChan <- marketPrice
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(mergedChan)
	}()

	return mergedChan
}
