package main

import (
	"math/rand/v2"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

type LifeGame struct {
	board    [][]bool
	row, col int
	time     int
}

func (p *LifeGame) isEmpty() bool {
	for r := 0; r < p.row; r++ {
		for c := 0; c < p.col; c++ {
			if p.board[r][c] {
				return false
			}
		}
	}
	return true
}

func (p *LifeGame) equals(other *LifeGame) bool {
	if p.row != other.row || p.col != other.col {
		return false
	}
	for r := 0; r < p.row; r++ {
		for c := 0; c < p.col; c++ {
			if p.board[r][c] != other.board[r][c] {
				return false
			}
		}
	}
	return true
}

func (p *LifeGame) Init(row, col int) *LifeGame {
	p.row = row
	p.col = col
	p.board = make([][]bool, row)
	for r := range row {
		p.board[r] = make([]bool, col)
		for c := range col {
			if rand.Float64() < 0.3 {
				p.board[r][c] = true
			}
		}
	}
	p.time = 0
	return p
}

func (p *LifeGame) Print() {
	var sb strings.Builder

	sb.WriteString("\033[H\033[2J")

	for r := 0; r < p.row; r++ {
		for c := 0; c < p.col; c++ {
			if p.board[r][c] {
				sb.WriteString("██")
				continue
			}
			sb.WriteString("  ")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("generation: ")
	sb.WriteString(strconv.Itoa(p.time))
	_, err := os.Stdout.WriteString(sb.String())
	if err != nil {
		panic(err)
	}
}

func (p *LifeGame) count_now_alive(r, c int) (i int) {
	i = 0
	if r < 0 {
		r += p.row
	}
	if p.row <= r {
		r -= p.row
	}
	if c < 0 {
		c += p.col
	}
	if p.col <= c {
		c -= p.col
	}
	if p.board[r][c] {
		i = 1
	}
	return i
}
func (p *LifeGame) is_dead_or_alive(r, c int) (b bool) {
	count := p.count_now_alive(r-1, c-1) +
		p.count_now_alive(r-1, c) +
		p.count_now_alive(r-1, c+1) +
		p.count_now_alive(r, c-1) +
		p.count_now_alive(r, c) +
		p.count_now_alive(r, c+1) +
		p.count_now_alive(r+1, c-1) +
		p.count_now_alive(r+1, c) +
		p.count_now_alive(r+1, c+1)
	switch count {
	case 3:
		b = true
	case 4:
		b = p.board[r][c]
	default:
		b = false
	}
	return b
}

func generate_gen(game *LifeGame, ch chan<- *LifeGame) {
	var previous *LifeGame
	for {
		next := new(LifeGame).Init(game.row, game.col)
		for r := 0; r < game.row; r++ {
			for c := 0; c < game.col; c++ {
				next.board[r][c] = game.is_dead_or_alive(r, c)
			}
		}
		next.time = game.time + 1

		// Check if all cells are dead
		if next.isEmpty() {
			ch <- next
			close(ch)
			return
		}

		// Check if state is the same as previous (oscillation period 1)
		if previous != nil && next.equals(previous) {
			ch <- next
			close(ch)
			return
		}

		// Check if state is the same as current (static state)
		if next.equals(game) {
			ch <- next
			close(ch)
			return
		}

		ch <- next
		previous = game
		game = next
	}
}

func terminalSize() (int, int) {
	col, row, err := term.GetSize(0)
	if err != nil {
		println("Error getting terminal size:", err.Error())
		os.Exit(1)
	}
	return row - 1, col / 2
}

func main() {
	game := new(LifeGame).Init(terminalSize())

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		print("\033[?25h") // show cursor
		println("\nGame interrupted. Exiting...")
		os.Exit(0)
	}()

	ch := make(chan *LifeGame)
	go generate_gen(game, ch)
	print("\033[?25l") // hide cursor
	for next := range ch {
		// print clear screen
		next.Print()
		time.Sleep(60 * time.Millisecond)
	}
	print("\033[?25h") // show cursor
	println("\nGame stabilized. Exiting...")
}
