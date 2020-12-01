package domain

import "errors"


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
}

func NewTdexFeeder(feeds []Feed, targets []Target) TdexFeeder {
	return &tdexFeeder{
		feeds: feeds,
		targets: targets,
		stopChan: make(chan bool),
		running: false,
	}
}

func (t *tdexFeeder) Start() error {
	if t.IsRunning() {
		return errors.New("the feeder is already started")
	}

	t.running = true
	marketPriceChannel := merge(t.feeds...)

	for t.running {
		select {
		case <-t.stopChan:
			t.running = false
			break;
		case marketPrice := <-marketPriceChannel:
			for _, target := range t.targets {
				err := target.Push(marketPrice)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	return nil
}

func (t *tdexFeeder) Stop() {
	t.stopChan <- true
}

func (t *tdexFeeder) IsRunning() bool {
	return t.running
}



