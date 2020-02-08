package fakes

import (
	"fmt"
	"os"
	"time"

	"github.com/storyscript/scheduler"
)

type Publisher struct{}

func (c *Publisher) Watch() (<-chan scheduler.PublishEvent, error) {
	eventCh := make(chan scheduler.PublishEvent)
	ticker := time.NewTicker(time.Second * 10)

	go func() {
		count := 0
		for range ticker.C {
			story := scheduler.Story{
				Name:    "foo",
				Payload: fmt.Sprintf(`{"count": %d }`, count),
			}
			eventCh <- scheduler.PublishEvent{Story: story, Error: nil}
			count++

			fmt.Fprintln(os.Stdout, "Publishing Story: ", story)
		}
	}()

	return eventCh, nil
}
