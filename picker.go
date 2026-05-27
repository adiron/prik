package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const numLines = 5
const brailleStart = 0x2800
const hlStart = "\x1b[1;33;4m"
const hlEnd = "\x1b[0m"
const markNo = "[ ]"
const markYes = "[x]"

var moveSide = map[uint8]uint8 {
	0: 3,
	1: 4,
	2: 5,
	3: 0,
	4: 1,
	5: 2,
	6: 7,
	7: 6,
}

var moveUp = map[uint8]uint8 {
 0: 6,
 1: 0,
 2: 1,
 6: 2,
 3: 7,
 4: 3,
 5: 4,
 7: 5,
}

var moveDown = map[uint8]uint8 {
 6: 0,
 0: 1,
 1: 2,
 2: 6,
 7: 3,
 3: 4,
 4: 5,
 5: 7,
}

var flipV = map[uint8]uint8 {
	0: 6,
	1: 2,
	2: 1,
	6: 0,
	3: 7,
	4: 5,
	5: 4,
	7: 3,
}

type Key int

const (
	KeyEnter Key = iota
	KeyBreak
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeySpace
	KeyA
	KeyFlipH
	KeyFlipV
	KeyInvert
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	KeyUnknown
)

type Picker struct {
	tty      *os.File
	fd       int
	oldState *term.State
	cursor   uint8
	dots     uint8
}
func NewPicker() (*Picker, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("open /dev/tty: %w", err)
	}
	fd := int(tty.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		tty.Close()
		return nil, fmt.Errorf("raw mode: %w", err)
	}
	return &Picker{tty: tty, fd: fd, oldState: oldState, cursor: 255, dots: 0b0}, nil
}

func (p *Picker) Close() {
	p.clearInline()
	fmt.Fprint(p.tty, "\x1b[?25h") // Show cursor
	term.Restore(p.fd, p.oldState)
	p.tty.Close()
}

type Result struct {
	Output string
	Code   int
}

func (p *Picker) Run() Result {
	p.initInline()
	defer p.Close()

	p.draw()
	for {
		key := p.readKey()
		switch key {
		case KeyEnter:
			return Result{Output: string(rune(brailleStart + rune(p.dots))), Code: 0}
		case KeyBreak:
			return Result{Code: 1}
		case KeyA:
			if p.dots == 0b11111111 {
				p.dots = 0b0
			} else {
				p.dots = 0b11111111
			}
		case Key1, Key2, Key3, Key4, Key5, Key6, Key7, Key8:
			p.dots = p.dots ^ (1 << uint8(key - Key1))
		case KeyLeft:
			if p.cursor == 255 {
					p.cursor = 3
				break
			}
			p.cursor = moveSide[p.cursor]

		case KeyRight:
			if p.cursor == 255 {
				p.cursor = 0
				break
			}
			p.cursor = moveSide[p.cursor]

		case KeyDown:
			if p.cursor == 255 {
				p.cursor = 0
				break
			}
			p.cursor = moveDown[p.cursor]
		case KeyUp:
			if p.cursor == 255 {
				p.cursor = 6
				break
			}
			p.cursor = moveUp[p.cursor]

		case KeySpace:
			if p.cursor == 255 {
				break
			}
			p.dots = p.dots ^ (1 << p.cursor)

		case KeyFlipH:
			var newDots uint8 = 0
			for before, after := range moveSide {
				if (p.dots >> before) & 0b01 == 1 {
					newDots |= 0b01 << after
				}
			}
			p.dots = newDots

		case KeyFlipV:
			var newDots uint8 = 0
			for before, after := range flipV {
				if (p.dots >> before) & 0b01 == 1 {
					newDots |= 0b01 << after
				}
			}
			p.dots = newDots

		case KeyInvert:
			p.dots = ^p.dots

		}
		p.draw()
	}
}

func dotMark(dots uint8, nth uint8, highlight bool) string {
	mark := markNo
	if (dots >> nth) & 0b01 == 1 {
		mark = markYes
	}

	if highlight {
		return hlStart + mark + hlEnd
	}

	return mark
}

func (p *Picker) draw() {
	char := rune(rune(p.dots) + brailleStart)
	// --- Line 1
	// Begin. Do the return thing to bring the cursor back and erase the line
	fmt.Fprint(p.tty, "\r\x1b[2K")
	fmt.Fprintf(
		p.tty, 
		"1 %s %s 4   move  ←/→ ↓/↑",
		dotMark(p.dots, 0, p.cursor == 0),
		dotMark(p.dots, 3, p.cursor == 3),
	)
	// End. Move the cursor one line down
	fmt.Fprint(p.tty, "\x1b[1B")

	// --- Line 2
	fmt.Fprint(p.tty, "\r\x1b[2K")
	fmt.Fprintf(
		p.tty, 
		"2 %s %s 5   toggle  space",
		dotMark(p.dots, 1, p.cursor == 1),
		dotMark(p.dots, 4, p.cursor == 4),
	)
	fmt.Fprint(p.tty, "\x1b[1B")

	// --- Line 3
	fmt.Fprint(p.tty, "\r\x1b[2K")
	fmt.Fprintf(
		p.tty, 
		"3 %s %s 6   togg. all  a",
		dotMark(p.dots, 2, p.cursor == 2),
		dotMark(p.dots, 5, p.cursor == 5),
	)
	fmt.Fprint(p.tty, "\x1b[1B")

	// --- Line 4
	fmt.Fprint(p.tty, "\r\x1b[2K")
	fmt.Fprintf(
		p.tty, 
		"7 %s %s 8   flip H/V  invert i",
		dotMark(p.dots, 6, p.cursor == 6),
		dotMark(p.dots, 7, p.cursor == 7),
	)
	fmt.Fprint(p.tty, "\x1b[1B")

	// --- Line 5
	fmt.Fprint(p.tty, "\r\x1b[2K")
	fmt.Fprintf(
		p.tty,
		"out: %c (0x%x)  ok  enter",
		char,
		char,
	)

	// Move back to top of inline area for next draw
	fmt.Fprintf(p.tty, "\x1b[%dA", numLines-1)
}

func (p *Picker) initInline() {
	fmt.Fprint(p.tty, "\x1b[?25l") // Hide cursor
	for range numLines {
		fmt.Fprint(p.tty, "\r\n")
	}
	fmt.Fprintf(p.tty, "\x1b[%dA", numLines)
}

func (p *Picker) clearInline() {
	for i := range numLines {
		fmt.Fprint(p.tty, "\r\x1b[2K")
		if i < numLines-1 {
			fmt.Fprint(p.tty, "\x1b[1B")
		}
	}
	if numLines > 1 {
		fmt.Fprintf(p.tty, "\x1b[%dA", numLines-1)
	}
}

func (p *Picker) readKey() Key {
	var buf [8]byte
	n, err := p.tty.Read(buf[:])
	if err != nil {
		return KeyUnknown
	}
	b := buf[:n]

	if n == 1 {
		if b[0] >= '1' && b[0] <= '8' {
			return Key1 + Key(b[0] - '1')
		}

		switch b[0] {
		case '\r', '\n':
			return KeyEnter
		case 0x03, 'q':
			return KeyBreak
		case ' ':
			return KeySpace
		case 'a':
			return KeyA
		case 'j':
			return KeyDown
		case 'k':
			return KeyUp
		case 'h':
			return KeyLeft
		case 'l':
			return KeyRight
		case 'H':
			return KeyFlipH
		case 'V':
			return KeyFlipV
		case 'i':
			return KeyInvert
		}
	}

	// Arrow key handling
	if n >= 3 && b[0] == 0x1b && b[1] == '[' {
		switch b[2] {
		case 'A':
			return KeyUp
		case 'B':
			return KeyDown
		case 'C':
			return KeyRight
		case 'D':
			return KeyLeft
		}
	}

	// Escape (if not part of arrows etc.)
	if n == 1 && b[0] == '\x1b' {
		return KeyBreak
	}

	return KeyUnknown
}
