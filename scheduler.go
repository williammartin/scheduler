package scheduler

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type Publisher interface {
	Watch() (<-chan PublishEvent, error)
}

type Deployer interface {
	Deploy(story Story) error
}

type Stopper struct {
	stopCh chan<- struct{}
}

func (s Stopper) Stop() {
	close(s.stopCh)
}

type Scheduler struct {
	Publisher Publisher
	Deployer  Deployer
}

func (s *Scheduler) Start() (Stopper, error) {
	publishes, err := s.Publisher.Watch()
	if err != nil {
		return Stopper{}, errors.Wrap(err, "failed to watch publishes")
	}

	stopCh := make(chan struct{})
	go func() {
		for {
			select {
			case event := <-publishes:
				if event.Error != nil {
					fmt.Fprintln(os.Stderr, errors.Wrap(event.Error, "failed to get publish event"))
					continue
				}

				if err := s.Deployer.Deploy(event.Story); err != nil {
					fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed deploy story"))
					continue
				}
			case <-stopCh:
				return
			}
		}
	}()

	return Stopper{stopCh: stopCh}, nil
}

type PublishEvent struct {
	Story Story
	Error error
}

type Story struct {
	Name    string
	Payload string
}
