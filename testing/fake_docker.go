package testing

import (
  docker "github.com/fsouza/go-dockerclient"
  "errors"
)

type FakeDocker struct {
  Listener chan<- *docker.APIEvents
  Containers map[string]*docker.Container
}

/* Methods provided by Docker Client */
func (c *FakeDocker) AddEventListener(thechan chan<- *docker.APIEvents) error {
  c.Listener = thechan
  return nil
}

func (c *FakeDocker) InspectContainer(id string) (*docker.Container, error) {
  stored, ok := c.Containers[id]
  if !ok { 
    return nil, errors.New("Missing Container")
  }
  return stored, nil
}

/* Methods Purely For Test Purposes */

func NewFakeDocker() *FakeDocker {
  dockerClient := &FakeDocker{}
  dockerClient.Containers = make(map[string]*docker.Container)
  return dockerClient
}

func (c *FakeDocker) SetContainer(sha string, container *docker.Container) {
  c.Containers[sha] = container
}

func (c *FakeDocker) SendEvent(event,sha string) {
  c.Listener <- &docker.APIEvents{Status:event, ID:sha}
}

