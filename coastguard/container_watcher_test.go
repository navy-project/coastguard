package coastguard_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "bitbucket.org/navy-project/coastguard/coastguard"
  ct "bitbucket.org/navy-project/coastguard/testing"
  docker "github.com/fsouza/go-dockerclient"
  //"github.com/coreos/go-etcd/etcd"
)

var _ = Describe("Container Watcher", func() {
  var dockerClient *ct.FakeDocker
  var etcdClient *ct.FakeEtcd
  var subject *coastguard.DockerWatcher

  BeforeEach(func() {
    dockerClient = ct.NewFakeDocker()
    etcdClient = ct.NewFakeEtcd()

    subject = coastguard.NewDockerWatcher(dockerClient, etcdClient)
    go subject.Watch()
  })

  Describe("When an event is recieved from docker", func() {
    Describe("When the container hasn't been seen before", func() {
      It("Remembers the container name in Etcd", func() {
        dockerClient.SetContainer("unknown_sha", &docker.Container{Name: "/the_name"})

        dockerClient.SendEvent("someevent", "unknown_sha")

        resp, err := etcdClient.Get("/navy/watched/unknown_sha", false, false)
        Expect(err).Should(BeNil())
        Expect(resp.Node.Value).To(Equal("the_name"))
      })
    })

    Describe("When the container has been seen before", func() {
      It("Uses the name from Etcd", func() {
        etcdClient.Set("/navy/watched/known_sha", "the_name", 0)
        dockerClient.SendEvent("someevent", "known_sha")

        dir, err := etcdClient.Get("/navy/events/containers", false, false)
        Expect(err).Should(BeNil())
        Expect(len(dir.Node.Nodes)).To(Equal(1))
        event := dir.Node.Nodes[0].Value
        Expect(event).To(Equal("{\"event\":\"someevent\",\"name\":\"the_name\"}"))
      })
    })

    It("Publishes the event to the etcd events stream", func() {
      dockerClient.SetContainer("unknown_sha", &docker.Container{Name: "/the_name"})
      dockerClient.SendEvent("someevent", "unknown_sha")

      dir, err := etcdClient.Get("/navy/events/containers", false, false)
      Expect(err).Should(BeNil())
      Expect(len(dir.Node.Nodes)).To(Equal(1))
      event := dir.Node.Nodes[0].Value
      Expect(event).To(Equal("{\"event\":\"someevent\",\"name\":\"the_name\"}"))
    })
  })
})
