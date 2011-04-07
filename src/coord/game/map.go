package game

type Map struct {
    Values [][]int
    Width uint
    Height uint
}

func NewMap(w uint, h uint) *Map {
    return &Map{make([][]int, w, h), w, h}
}

func (self *Map) Copy() *Map {
    newMap := NewMap(self.Width, self.Height)
    for i := uint(0); i < self.Width; i++ {
        for j := uint(0); j < self.Height; j++ {
            newMap.Values[i][j] = self.Values[i][j]
        }
    }
    return newMap
}
