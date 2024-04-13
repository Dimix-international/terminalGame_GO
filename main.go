package main

import (
	"bytes"
	"fmt"
)

const (
	NOTHING = 0
	WALL    = 1
	PLAYER  = 69
)

type level struct {
	width  int
	height int
	data   [][]byte
}

func newLevel(width, height int) *level {
	data := make([][]byte, height)

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			data[h] = make([]byte, width)
		}
	}

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if w == width-1 || w == 0 {
				data[h][w] = WALL
			}
			if h == height-1 || h == 0 {
				data[h][w] = WALL
			}
		}
	}
	return &level{width: width, height: height, data: data}
}

func (l *level) x() {}

type game struct {
	isRinning bool
	level     *level
}

func newGame(width int, height int) *game {
	return &game{level: newLevel(width, height)}
}

func (g *game) start() {
	g.isRinning = true
	g.loop()
}

func (g *game) loop() {
	for g.isRinning {
		g.update()
		g.render()
	}
}

func (g *game) update() {}

func (g *game) renderLevel() {}

func (g *game) render() {}

func main() {
	width := 80
	height := 18
	level := make([][]byte, height)

	buf := new(bytes.Buffer)
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if level[h][w] == NOTHING {
				buf.WriteString(" ")
			}
			if level[h][w] == WALL {
				buf.WriteString("□")
			}
		}
		buf.WriteString("\n")
	}

	fmt.Println(buf.String())
}
