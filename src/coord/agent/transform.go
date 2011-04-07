package agent

import geo "coord/geometry"

type Transform interface {
    Turn() uint64
    Position() *geo.Point
    Energy() Energy
    Alive() bool
    Wait() uint16
}
