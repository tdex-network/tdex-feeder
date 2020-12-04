package domain

import (
	"errors"
	"sync"
)

type TdexFeeder interface {
	Start() error
	Stop()
	IsRunning() bool
}

type tdexFeeder struct {
	feeds []Feed
	targets []Target
	stopChan chan bool
	running bool
	locker sync.Locker
}

func NewTdexFeeder(feeds []Feed, targets []Target) TdexFeeder {
	return &tdexFeeder{
		feeds: feeds,
		targets: targets,
		stopChan: make(chan bool),
		running: false,
		locker: &sync.Mutex{},
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
			break;
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



