package main

import (
	"strconv"
)

func serializeGrid(grid [][]int) []byte {
	gridString := make([]byte, len(grid)*len(grid[0]))
	for i, row := range(grid) {
		for j, _ := range(row) {
			gridString[i*len(row)+j] = strconv.Itoa(row[j])[0]
		}
	}
	return gridString
}

func simpleGrid() [][]int {
	grid := make([][]int, 100, 100)
	for i, row := range(grid) {
		for j, _ := range(row) {
			row[j] = (i + j) % 3
		}
	}
	return grid
}
