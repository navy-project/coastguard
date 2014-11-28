package main

import (
  "bitbucket.org/navy-project/coastguard/coastguard"
  docker "github.com/fsouza/go-dockerclient"
  "github.com/coreos/go-etcd/etcd"
  "os"
  "log"
)

func main() {
  dockerClient := setupDocker()
  etcdClient := setupEtcd()
  watcher := coastguard.NewDockerWatcher(dockerClient, etcdClient)
  watcher.Watch()
}

func setupDocker() *docker.Client {
  dockerAddr := "unix:///var/run/docker.sock"
  client, err := docker.NewClient(dockerAddr)
  if err != nil {
    panic(err)
  }
  log.Println("Listening To Docker: ", dockerAddr)
  return client
}

func setupEtcd() *etcd.Client {
  etcdserver := "http://" + os.Getenv("ETCD_PORT_4001_TCP_ADDR") + ":" + os.Getenv("ETCD_PORT_4001_TCP_PORT")
  log.Println("Connected To Etcd: ", etcdserver)
  return etcd.NewClient([]string{etcdserver})
}
