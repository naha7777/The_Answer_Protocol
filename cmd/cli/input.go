package main

import (
	"bufio"
	"fmt"
	"slices"
	"strings"
	"sync"
	"unicode/utf8"
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

// print a server message
func (in *Input) Println(s string) {
	in.mu.Lock()
	defer in.mu.Unlock()
	// delete the line, write the message and display the line "> "
	fmt.Printf("\r\033[K%s\r\n", s)
	in.redraw()
}

func (in *Input) ReadLine(r *bufio.Reader) (string, error) {
	for {
		// read a keyboard key
		ch, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}

		in.mu.Lock()
		switch {
		// 27 = ECHAP = beginning of a special keyboard key sequency
		// arrow return 3 bytes : ESC '[' 'A' (up) or  ESC '[' 'B' (down)
		case ch == 27:
			next, _, _ := r.ReadRune() // '['
			code, _, _ := r.ReadRune() // 'A' or 'B'
			if next == '[' && code == 'A' && in.idx > 0 {
				// if up arrow : go up on the history
				in.idx--
				in.buf = []rune(in.history[in.idx])
				in.redraw()
			} else if next == '[' && code == 'B' && in.idx < len(in.history) {
				// else if down arrow, go down on history
				in.idx++
				if in.idx == len(in.history) {
					in.buf = nil // back on the last entry : empty line
				} else {
					in.buf = []rune(in.history[in.idx])
				}
				in.redraw()
			}
		// ENTRY : line validation
		case ch == '\r' || ch == '\n':
			line := string(in.buf)
			in.buf = nil
			fmt.Print("\r\n") // go to the next line on the screen
			if strings.TrimSpace(line) != "" {
				if slices.Contains(in.history, line) {
					// if the line already exist in history, remove it and put
					// it on the end of history
					var filtered []string
					for _, l := range in.history {
						if l != line {
							filtered = append(filtered, l)
						}
					}
					filtered = append(filtered, line)
					in.history = filtered
				} else {
					in.history = append(in.history, line)
				}
			}
			in.idx = len(in.history) // new line position
			in.mu.Unlock()
			return line, nil
		// CTRL+C (3) or CTRL+D (4) : quit
		case ch == 3 || ch == 4:
			in.mu.Unlock()
			return "", fmt.Errorf("quit")
		// BACKSPACE : delete the last character
		case ch == 127 || ch == 8:
			if len(in.buf) > 0 {
				in.buf = in.buf[:len(in.buf)-1]
				in.redraw()
			}
		// printable character : add it to buffer
		case ch >= 32:
			// limit : 1000 bytes
			if len(string(in.buf))+utf8.RuneLen(ch) < 1000 {
				in.buf = append(in.buf, ch)
				in.redraw()
			}
		}
		in.mu.Unlock()
	}
}
