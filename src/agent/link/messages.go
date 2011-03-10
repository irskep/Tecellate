package link

import geo "coord/geometry"

// stub until we have this in coord/game
// heirarchy
//         Object
//            |
//      +-----+-----+
//      |           |
//    Agent        Item
type GameObject interface{}
type GameItem   interface{}
type GameAgent  interface{}

type InventoryItem struct {
    Item GameItem
    Amt  uint8
}

type Vision interface {
    Objects() <-chan GameObject // all objects (agents and items) in vision
    Agents() <-chan GameAgent   // all agents in vision
    Items() <-chan GameItem     // all items in vision
    Look(geo.Point) <-chan GameObject  // what is on sqaure (x,y)
                                       // with the current position of the agent
                                       // taken as the origin.
                                       // a closed channel if out of range.
}

type Move interface {
    Move() geo.Point // relative to current position
}

type Audio interface {
    Hear() []byte
}

type Broadcast interface {
    Message() (uint8, []byte)
}

type Inventory interface {
    Items() <-chan InventoryItem
}
