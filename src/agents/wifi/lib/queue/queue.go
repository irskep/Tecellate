package queue

import "container/list"

type Queue struct {
    List *list.List
}

func NewQueue() *Queue {
    return &Queue{
        List:list.New(),
    }
}

func (self *Queue) Len() int {
    return self.List.Len()
}

func (self *Queue) Empty() bool {
    return self.List.Len() == 0
}

func (self *Queue) Queue(m interface{}) {
    self.List.PushBack(m)
}

func (self *Queue) QueueFront(m interface{}) {
    self.List.PushFront(m)
}

func (self *Queue) Dequeue() (interface{}, bool) {
    front := self.List.Front()
    if front == nil { return nil, false }
    m := front.Value
    self.List.Remove(front)
    return m, true
}
