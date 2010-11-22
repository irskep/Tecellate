package main

import (
	"fmt"
	"strconv"
	"../ttypes"
)

func serializeGrid(grid [][]int) []byte {
	gridString := make([]byte, len(grid)*len(grid[0]))
	fmt.Printf("%d\n", len(grid)*len(grid[0]))
	for i, row := range(grid) {
		for j, _ := range(row) {
			fmt.Printf("%s", strconv.Itoa(row[j]))
			gridString[i*len(row)+j] = strconv.Itoa(row[j])[0]
		}
	}
	return gridString
}

func simpleGrid(w uint, h uint) *ttypes.Grid {
	grid := new(ttypes.Grid)
	grid.Width = w
	grid.Height = h
	grid.Items = make([]byte, w*h)
	for i := uint(0); i < w; i++ {
		for j := uint(0); j < h; j++ {
			grid.Items[i*w+j] = byte(i*w+j) % 3
		}
	}
	return grid
}
