package main

import (
	"bufio"
	"fmt"
	"slices"
	"strings"
	"sync"
)

type Input struct {
    mu      sync.Mutex
    buf     []rune
    history []string
    idx     int
}

const prompt = "> "

func (in *Input) redraw() {
	fmt.Printf("\r\033[K%s%s", prompt, string(in.buf))
}

func (in *Input) Println(s string) {
	in.mu.Lock()
	defer in.mu.Unlock()
	fmt.Printf("\r\033[K%s\r\n", s)
	in.redraw()
}

func (in *Input) ReadLine(r *bufio.Reader) (string, error) {
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}

		in.mu.Lock()
		switch {
		case ch == 27:
			next, _, _ := r.ReadRune()
			code, _, _ := r.ReadRune()
			if next == '[' && code == 'A' && in.idx > 0 {
				in.idx--
				in.buf = []rune(in.history[in.idx])
				in.redraw()
			} else if next == '[' && code == 'B' && in.idx < len(in.history) {
				in.idx++
				if in.idx == len(in.history) {
					in.buf = nil
				} else {
					in.buf = []rune(in.history[in.idx])
				}
				in.redraw()
			}
		case ch == '\r' || ch == '\n':
			line := string(in.buf)
			in.buf = nil
			fmt.Print("\r\n")
			if strings.TrimSpace(line) != "" {
				if slices.Contains(in.history, line) {
					var new []string
					for _, l := range in.history {
						if l != line {
							new = append(new, l)
						}
					}
					new = append(new, line)
					in.history = new
				} else {
					in.history = append(in.history, line)
				}
			}
			in.idx = len(in.history)
			in.mu.Unlock()
			return line, nil
		case ch == 3 || ch == 4:
			in.mu.Unlock()
			return "", fmt.Errorf("quit")
		case ch == 127 || ch == 8:
			if len(in.buf) > 0 {
				in.buf = in.buf[:len(in.buf)-1]
				in.redraw()
			}
		case ch >= 32:
			in.buf = append(in.buf, ch)
			in.redraw()
		}
		in.mu.Unlock()
	}
}
