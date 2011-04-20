package datagram

import "agents/wifi/lib/queue"

type DataGramQueue queue.Queue

func NewDataGramQueue() *DataGramQueue { return (*DataGramQueue)(queue.NewQueue()) }
func (self *DataGramQueue) Len() int { return ((*queue.Queue)(self)).Len() }
func (self *DataGramQueue) Empty() bool { return ((*queue.Queue)(self)).Empty() }
func (self *DataGramQueue) Queue(m *DataGram) { ((*queue.Queue)(self)).Queue(m) }
func (self *DataGramQueue) QueueFront(m *DataGram) { ((*queue.Queue)(self)).QueueFront(m) }
func (self *DataGramQueue) Peek() *DataGram {
    return ((*queue.Queue)(self)).Peek().(*DataGram)
}
func (self *DataGramQueue) Dequeue() (*DataGram, bool) {
    q := (*queue.Queue)(self)
    dg, ok := q.Dequeue()
    return dg.(*DataGram), ok
}

func (self *DataGramQueue) Clean() {
    q := (*queue.Queue)(self)
    for e := q.List.Front(); e != nil; {
        m := e.Value.(*DataGram)
        m.DecTTL()
        if m.SendTTL == 0 || m.TTL == 0 {
            next_e := e.Next()
            q.List.Remove(e)
            e = next_e
            if e == nil { break }
        } else {
            e = e.Next()
        }
    }
}
