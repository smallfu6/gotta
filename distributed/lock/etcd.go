package main

/*
	基于 etcd 实现分布式阻塞锁
	TODO: etcd

	etcd中没有像ZooKeeper那样的Sequence节点, 所以其锁实现和基于ZooKeeper
	实现的有所不同;
	下列代码中使用的etcdsync的加锁流程如下:
	- 先检查/lock路径下是否有值, 如果有值, 说明锁已经被别人抢了;
	- 如果没有值, 那么写入自己的值; 如果写入成功返回, 说明加锁成功;
		如果写入时节点被其他节点写入过了, 那么会导致加锁失败;
	- 监视/lock下的事件, 此时陷入阻塞;
	- 当/lock路径下发生事件时, 当前进程被唤醒;
	检查发生的事件是否是删除事件(说明锁持有者主动解锁)或者过期事件(说明锁过
		期失效), 如果是的话, 那么回到1继续走抢锁流程;

*/

import (
	"log"

	"github.com/zieckey/etcdsync"
)

func main() {
	m, err := etcdsync.New("/lock", 10, []string{"http://127.0.0.1:2379"})
	if m == nil || err != nil {
		log.Printf("etcdsync.New failed")
		return
	}

	err = m.Lock()
	if err != nil {
		log.Printf("etcdsync.Lock failed")
		return
	}

	log.Printf("etcdsync.Lock OK")
	log.Printf("Get the lock. Do something here.")

	err = m.Unlock()
	if err != nil {
		log.Printf("etcdsync.Unlock failed")
	} else {
		log.Printf("etcdsync.Unlock OK")
	}
}
