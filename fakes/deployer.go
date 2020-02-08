package fakes

import (
	"fmt"

	"github.com/storyscript/scheduler"
)

type Deployer struct{}

func (c *Deployer) Deploy(story scheduler.Story) error {
	fmt.Println("Deployed: " + story.Name)
	return nil
}
