package coastguard

import (
  docker "github.com/fsouza/go-dockerclient"
  "github.com/coreos/go-etcd/etcd"
  "encoding/json"
  "log"
)

type DockerWatcher struct {
  docker dockerClient
  etcd etcdClient
  events chan *docker.APIEvents
}

type ContainerEvent struct {
  Event string      `json:"event"`
  Name string       `json:"name"`
}

type dockerClient interface {
  AddEventListener(thechan chan<- *docker.APIEvents) error
  InspectContainer(id string) (*docker.Container, error)
}

type etcdClient interface {
  Set(key string, value string, ttl uint64) (*etcd.Response, error)
  Get(key string, sort, recur bool) (*etcd.Response, error)
  AddChild(key string, value string, ttl uint64) (*etcd.Response, error)
}

func NewDockerWatcher(d dockerClient, e etcdClient) *DockerWatcher {
  w := &DockerWatcher{docker: d, etcd: e}
  w.events = make(chan *docker.APIEvents)
  w.docker.AddEventListener(w.events)
  return w
}

func (w *DockerWatcher) Watch() {
  for {
    var name string
    event := <- w.events
    resp, err := w.etcd.Get("/navy/watched/" + event.ID, false, false)
    if err != nil {
      container, err := w.docker.InspectContainer(event.ID)
      if err != nil {
        log.Print(err)
        continue
      }
      name = container.Name[1:] //Strip off starting / from name
      w.etcd.Set("/navy/watched/" + event.ID, name, 0)
    } else {
      name = resp.Node.Value
    }

    etcdEvent := &ContainerEvent{event.Status, name}
    b, err := json.Marshal(etcdEvent)
    if err == nil {
      log.Print("Event: ", string(b))
      w.etcd.AddChild("/navy/events/containers", string(b), 0)
    }
  }
}

