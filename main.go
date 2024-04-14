package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	NOTHING     = 0
	WALL        = 1
	PLAYER      = 69
	MAX_SAMPLES = 100
)

type input struct {
	pressedKey byte
}

func (i *input) update() {

	// ch := make(chan byte)
	// tick := time.NewTicker(time.Millisecond * 2)

	// go func() {
	// 	b := make([]byte, 1)
	// 	os.Stdin.Read(b) //blocking until stdin has stuff in buffer
	// 	ch <- b[0]
	// 	close(ch)
	// }()

	// for {
	// 	select {
	// 	case <-tick.C:
	// 		return
	// 	case key, ok := <-ch:
	// 		if !ok {
	// 			break
	// 		}
	// 		i.pressedKey = key
	// 	}
	// }

}

type position struct {
	x, y int
}

type player struct {
	pos     position
	level   *level
	reverse bool
	input   *input
}

func (p *player) update() {
	if p.reverse {
		p.pos.x -= 1
		if p.pos.x == 1 {
			p.pos.x += 1
			p.reverse = false
		}
		return
	}

	p.pos.x += 1
	if p.pos.x == p.level.width-2 {
		p.pos.x -= 1
		p.reverse = true
	}
}

type stats struct {
	start  time.Time
	frames int
	fps    float64
}

func newStats() *stats {
	return &stats{
		fps:   69,
		start: time.Now(),
	}
}

func (s *stats) update() {
	s.frames++
	if s.frames == MAX_SAMPLES {
		s.fps = float64(s.frames) / time.Since(s.start).Seconds()
		s.frames = 0
		s.start = time.Now()
	}
}

type level struct {
	width  int
	height int
	data   [][]int
}

func newLevel(width, height int) *level {
	data := make([][]int, height)

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			data[h] = make([]int, width)
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

func (l *level) set(pos position, v int) {
	l.data[pos.y][pos.x] = v
}

type game struct {
	isRinning bool
	level     *level
	drawBuf   *bytes.Buffer
	stats     *stats
	player    *player
	input     *input
}

func newGame(width int, height int) *game {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	var (
		lvl = newLevel(width, height)
		inp = &input{}
	)

	return &game{
		level:   lvl,
		drawBuf: new(bytes.Buffer),
		stats:   newStats(),
		player:  &player{level: lvl, pos: position{x: 2, y: 5}, input: inp},
		input:   inp,
	}
}

func (g *game) start() {
	g.isRinning = true
	g.loop()
}

func (g *game) loop() {
	for g.isRinning {
		g.input.update()
		g.update()
		g.render()
		g.stats.update()
		time.Sleep(time.Millisecond * 13) //limit fps
	}
}

func (g *game) update() {
	g.level.set(g.player.pos, NOTHING)
	g.player.update()
	g.level.set(g.player.pos, PLAYER)
}

func (g *game) renderLevel() {
	for h := 0; h < g.level.height; h++ {
		for w := 0; w < g.level.width; w++ {
			if g.level.data[h][w] == NOTHING {
				g.drawBuf.WriteString(" ")
			}
			if g.level.data[h][w] == WALL {
				g.drawBuf.WriteString("□")
			}
			if g.level.data[h][w] == PLAYER {
				g.drawBuf.WriteString("♚")
			}
		}
		g.drawBuf.WriteString("\n")
	}
}

func (g *game) render() {
	g.drawBuf.Reset()
	fmt.Fprint(os.Stdout, "\033[2J\033[1;1H")

	g.renderLevel()
	g.renderStats()
	fmt.Fprint(os.Stdout, g.drawBuf.String())
}

func (g *game) renderStats() {
	g.drawBuf.WriteString("--STATS\n")
	g.drawBuf.WriteString(fmt.Sprintf("FPS: %.2f", g.stats.fps))
}

func main() {
	width := 80
	height := 15
	newGame(width, height).start()
}
