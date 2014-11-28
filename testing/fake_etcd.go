package testing

import (
  "github.com/coreos/go-etcd/etcd"
  "errors"
)

type FakeEtcd struct {
  keys map[string]*etcd.Response
}

func (c *FakeEtcd) Get(key string, sort, recur bool) (*etcd.Response, error) {
  stored, ok := c.keys[key]
  if !ok { 
    return nil, errors.New("Missing Key")
  }
  return stored, nil
}

func (c *FakeEtcd) Set(key string, value string, ttl uint64) (*etcd.Response, error) {
  response := &etcd.Response{Node: &etcd.Node{Value: value}}
  c.keys[key] = response
  return response, nil
}

func (c *FakeEtcd) AddChild(key string, value string, ttl uint64) (*etcd.Response, error) {
  var parent *etcd.Response
  parent, ok := c.keys[key]
  if !ok {
    parent, _ = c.Set(key, "A Dir", 0)
  }
  response := &etcd.Response{Node: &etcd.Node{Value: value}}
  parent.Node.Nodes = append(parent.Node.Nodes, response.Node)
  return response, nil
}

/* Method purely for test purposes */

func NewFakeEtcd() *FakeEtcd {
  etcdClient := &FakeEtcd{}
  etcdClient.keys = make(map[string]*etcd.Response)
  return etcdClient
}
