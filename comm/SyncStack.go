package comm

import (
	"container/list"
	"fmt"
	"sync"
)

type SyncStack struct {

	list *list.List  //栈链表
	cond *sync.Cond  //条件互斥量
}

func NewSyncStack() (ss *SyncStack) {

	ss = &SyncStack{
		list: list.New(),
		cond: sync.NewCond(new(sync.Mutex)),
	}
	return
}


func (this *SyncStack) Init() {
	this.cond.L.Lock()
	this.list.Init()
	this.cond.L.Unlock()
}

//将元素压入栈中并释放信号
func (this *SyncStack) Push(v interface{}) {
	//fmt.Println("Push try lock")
	this.cond.L.Lock()
	this.list.PushBack(v)
	//fmt.Println("Push element ok, sending signal")
	this.cond.Signal()
	//fmt.Println("Push element send signal ok")
	this.cond.L.Unlock()
	//fmt.Println("Push unlock ok")
}

//将元素弹出栈，如果元素数为0则等待Push信号
func (this *SyncStack) Pop() (v interface{}){
	//fmt.Println("Pop try lock")
	this.cond.L.Lock()
	//fmt.Println("Pop try get element")
	if this.list.Len() == 0 {
		fmt.Println("Pop element nil, wait signal")
		this.cond.Wait() //解锁后等待信号，获得信号后再加锁
	}
	e := this.list.Front()
	v = e.Value
	//fmt.Println("Pop element ok")
	this.cond.L.Unlock()
	//fmt.Println("Pop unlock ok")
	return
}

func (this *SyncStack) Size() int {

	return this.Size()
}
