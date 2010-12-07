package main

import (
	"log"
	"os"
	"scanner"
	"strconv"
	"ttypes"
)

func simpleGrid(w uint, h uint) *ttypes.Grid {
	grid := new(ttypes.Grid)
	grid.Width = w
	grid.Height = h
	grid.Items = make([]byte, w*h)
	return grid
}

func scanAndReadUint(s *scanner.Scanner) uint {
	s.Scan()
	tok, err := strconv.Atoui(s.TokenText())
	if err != nil { log.Exit(err) }
	return tok
}

func readGridFromFile(path string) (*ttypes.Grid, []ttypes.BotConf) {
	gridFile, err := os.Open(path, os.O_RDONLY, 0)
	if err != nil { log.Exit(err) }
	defer gridFile.Close()
	
	var s scanner.Scanner
	s.Init(gridFile)
	w := scanAndReadUint(&s)
	h := scanAndReadUint(&s)
	grid := simpleGrid(w, h)
	
	botConfs := make([]ttypes.BotConf, 0)
	for i := uint(0); i < w; i++ {
		for j := uint(0); j < h; j++ {
			tok := uint8(scanAndReadUint(&s))
			if tok > 0 {
				botConfs = append(botConfs, ttypes.BotConf{"build/test", i, j})
			}
		}
	}
	
	return grid, botConfs
}
