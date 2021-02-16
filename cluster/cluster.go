package cluster

import (
	"github.com/hashicorp/memberlist"
	"io/ioutil"
	"stathat.com/c/consistent"
	"time"
)


//节点接口
type Node interface {
	ShouldProcess(key string) (string, bool)
	Members() []string
	Addr() string
}

//节点结构体
type node struct {
	*consistent.Consistent
	addr string
}

//获取本节点地址
func (n *node) Addr() string {
	return n.addr
}

func New(addr, cluster string) (Node, error) {
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	//丢弃日志
	conf.LogOutput = ioutil.Discard
	//创建一个节点
	l, e := memberlist.Create(conf)
	if e != nil {
		return nil, e
	}
	if cluster == "" {
		cluster = addr
	}
	//节点列表
	clu := []string{cluster}
	//新建节点加入集群
	_, e = l.Join(clu)
	if e != nil {
		return nil, e
	}
	//创建一个环
	circle := consistent.New()
	//每个节点的虚拟节点数量
	circle.NumberOfReplicas = 256
	//每隔一秒将集群列表中的节点,更新到环上
	go func() {
		for {
			//获取集群节点数量
			m := l.Members()
			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			circle.Set(nodes)
			time.Sleep(time.Second)
		}
	}()
	return &node{circle, addr}, nil
}

//chu'li
func (n *node) ShouldProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}
