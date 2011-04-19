package message

import "agents/wifi/lib/queue"

type MessageQueue queue.Queue

func NewMessageQueue() *MessageQueue { return (*MessageQueue)(queue.NewQueue()) }
func (self *MessageQueue) Len() int { return ((*queue.Queue)(self)).Len() }
func (self *MessageQueue) Empty() bool { return ((*queue.Queue)(self)).Empty() }
func (self *MessageQueue) Queue(m *Message) { ((*queue.Queue)(self)).Queue(m) }
func (self *MessageQueue) QueueFront(m *Message) { ((*queue.Queue)(self)).QueueFront(m) }
func (self *MessageQueue) Peek() *Message {
    return ((*queue.Queue)(self)).Peek().(*Message)
}
func (self *MessageQueue) Dequeue() (*Message, bool) {
    q := (*queue.Queue)(self)
    dg, ok := q.Dequeue()
    return dg.(*Message), ok
}
