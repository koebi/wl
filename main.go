package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nsf/termbox-go"
)

const (
	CursorLeft int = iota
	CursorRight
	Edit int = iota
	Normal
)

var wldir = flag.String("wldir", "/home/koebi/filme/watchlist/", "location of watchlist directory")

type ListEntry struct {
	Name          string
	RecommendedBy map[string]string
	Other         string
	Genre         string
}

type Watchlist []ListEntry

type Prompt struct {
	width  int
	height int
	cursor int
}

func NewPrompt(s *Screen) *Prompt {
	return &Prompt{
		width:  s.Width - 8,
		height: s.Height - 6,
		cursor: 0,
	}
}

type Screen struct {
	Left       []string
	Right      []string
	Width      int
	Height     int
	LeftWidth  int
	RightWidth int
}

func NewScreen() *Screen {
	return &Screen{
		Left:   nil,
		Right:  nil,
		Width:  0,
		Height: 0,
	}
}

type Cursor struct {
	CurrentEntry int
	Side         int
}

func NewCursor() *Cursor {
	return &Cursor{
		CurrentEntry: 0,
		Side:         CursorLeft,
	}
}

func (s *Screen) GetNew(p *Prompt) ListEntry {
	fmt.Fprintf(p, "Foobar")

	return ListEntry{
		Name: "foo",
	}
}

func (wl *Watchlist) Add(entry ListEntry) {
	var file *os.File
	defer file.Close()

	_, err := os.Stat(filepath.Join(*wldir, entry.Name))
	if os.IsNotExist(err) {
		log.Printf("Creating file %s…\n", filepath.Join(*wldir, entry.Name))
		file, err = os.Create(filepath.Join(*wldir, entry.Name))
		if err != nil {
			log.Printf("File creation failed: %v\n", err)
		}
	} else if err != nil {
		log.Printf("File %s seems wrong: %v\n", filepath.Join(*wldir, entry.Name), err)
	} else {
		log.Printf("Movie %s already exists. Please edit it's entry.")
	}

	err = toml.NewEncoder(file).Encode(entry)
	if err != nil {
		log.Printf("Encoding File failed: %v\n", err)
		log.Println(err)
	}

	log.Printf("Watchlist-Entry %s created\n", entry.Name)
	return
}

func (p *Prompt) Write(s []byte) (n int, err error) {
	for x := 4; x < p.width+4; x++ {
		termbox.SetCell(x, p.cursor+4, rune(s[x-4]), termbox.ColorBlack, termbox.ColorBlue)
	}

	p.cursor += 1
	return 0, nil
}

func (p *Prompt) DrawToTerm() {
	for y := 3; y < p.height+4; y++ {
		for x := 4; x < p.width+4; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorBlue, termbox.ColorBlue)
		}
	}
}

func (p *Prompt) EditLine(ev termbox.Event) {

}

func (p *Prompt) HandleKeyPress(k termbox.Key, r rune) error {
	//can I just pass runes to the Terminal, and it will handle them correctly?
	switch {
	case k == termbox.KeyArrowLeft:
	case k == termbox.KeyArrowRight:
	case k == termbox.KeyDelete:
	case k == termbox.KeyHome:
	case k == termbox.KeyEnd:

	}
}

func (p *Prompt) Handle(ev termbox.Event) error {
	switch {
	case event.Type == termbox.EventKey:
		err := p.HandleKeyPress(event.Key, event.Ch)
		return err
	default:
		return fmt.Errorf("%s", "Unknown Event Type, continue")
	}
}

func ParseWatchlist() (wl Watchlist, err error) {
	files, err := ioutil.ReadDir(*wldir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Reading dir %s failed: %v\n", *wldir, err)
		return nil, err
	}

	for _, file := range files {
		entry := ListEntry{Name: file.Name()}

		_, err = toml.DecodeFile(*wldir+file.Name(), &entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Decoding file %s failed: %v\n", (*wldir + file.Name()), err)
			return nil, err
		}

		wl = append(wl, entry)
	}
	return wl, nil
}

func (wl Watchlist) GetWidth() int {
	max := 0
	for _, entry := range wl {
		if len(entry.Name) >= max {
			max = len(entry.Name)
		}
	}
	return max
}

func PadLeft(s string, i int) string {
	var padded string

	if len(s) > i-3 {
		padded = "| " + s[:i-3] + " "
	} else {
		padded = "| " + s
		padded += strings.Repeat(" ", (i - len(padded)))
	}
	return padded
}

func PadRight(s string, i int) string {
	var padded string

	if len(s) > i-3 {
		padded = " " + s[:i-3] + " |"
	} else {
		padded = " " + s
		padded += strings.Repeat(" ", (i - 2 - len(padded)))
		padded += " |"

	}
	return padded
}

func (s *Screen) SetLeft(wl Watchlist) {
	for i := 0; i < len(wl); i++ {
		s.Left = append(s.Left, PadLeft(wl[i].Name, s.LeftWidth))
	}
	for i := len(s.Left); i < s.Height; i++ {
		s.Left = append(s.Left, PadLeft(" ", s.LeftWidth))
	}
}

func (s *Screen) SetRight(entry ListEntry) {
	s.Right = []string{}
	s.Right = append(s.Right, PadRight("[Genre]", s.RightWidth))
	s.Right = append(s.Right, PadRight(" "+entry.Genre, s.RightWidth))

	s.Right = append(s.Right, PadRight(" ", s.RightWidth))
	s.Right = append(s.Right, PadRight("[Recommended By]", s.RightWidth))

	for person, date := range entry.RecommendedBy {
		s.Right = append(s.Right, PadRight("  "+person+" on "+date, s.RightWidth))
	}

	s.Right = append(s.Right, PadRight("", s.RightWidth))
	s.Right = append(s.Right, PadRight("[Other Info]", s.RightWidth))
	s.Right = append(s.Right, PadRight(" "+entry.Other, s.RightWidth))

	for i := len(s.Right); i < s.Height; i++ {
		s.Right = append(s.Right, PadRight(" ", s.RightWidth))
	}
}

func (s *Screen) Init(wl Watchlist) {
	s.SetWidth(wl)

	s.SetLeft(wl)
	s.SetRight(wl[0])
}

func (s *Screen) DrawVerticalBounds() {
	var padding []string
	padding = append(padding, "+"+strings.Repeat("-", s.LeftWidth-1)+"+"+strings.Repeat("-", s.RightWidth-1)+"+")
	padding = append(padding, "|"+strings.Repeat(" ", s.LeftWidth-1)+"|"+strings.Repeat(" ", s.RightWidth-1)+"|")

	// draw first two lines
	for y := 0; y < 2; y++ {
		for x := 0; x < s.Width; x++ {
			termbox.SetCell(x, y, rune(padding[y][x]), termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	padding = append(padding, "|"+strings.Repeat(" ", s.LeftWidth-1)+"|"+strings.Repeat(" ", s.RightWidth-1)+"|")
	padding = append(padding, "+"+strings.Repeat("-", s.LeftWidth-1)+"+"+strings.Repeat("-", s.RightWidth-1)+"+")
	padding = append(padding, strings.Repeat(" ", s.Width))

	// draw last three lines
	for x := 0; x < s.Width; x++ {
		termbox.SetCell(x, s.Height-3, rune(padding[2][x]), termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, s.Height-2, rune(padding[3][x]), termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, s.Height-1, rune(padding[4][x]), termbox.ColorWhite, termbox.ColorDefault)
	}
}

func (s *Screen) DrawToTerm(c *Cursor) error {
	s.DrawVerticalBounds()

	for y := 0; y < s.Height-5; y++ {
		draw := ' '
		for x := 0; x < s.LeftWidth; x++ {
			bg := termbox.ColorDefault
			fg := termbox.ColorWhite
			if y == c.CurrentEntry && c.Side == CursorLeft {
				if x > 1 {
					bg = termbox.ColorYellow
					fg = termbox.ColorBlack
				}
			}

			if y < len(s.Left) {
				draw = rune(s.Left[y][x])
			}

			termbox.SetCell(x, y+2, draw, fg, bg)
		}
		termbox.SetCell(s.LeftWidth, y+2, '|', termbox.ColorWhite, termbox.ColorDefault)
		draw = ' '
		for x := s.LeftWidth + 1; x < s.Width; x++ {
			bg := termbox.ColorDefault
			fg := termbox.ColorWhite
			if y == c.CurrentEntry && c.Side == CursorRight {
				if x < s.Width-2 {
					bg = termbox.ColorYellow
					fg = termbox.ColorBlack
				}
			}

			if y < len(s.Right) {
				draw = rune(s.Right[y][x-(s.LeftWidth+1)])
			}

			termbox.SetCell(x, y+2, draw, fg, bg)
		}
	}

	if err := termbox.Flush(); err != nil {
		return err
	}
	return nil
}

func (s *Screen) SetWidth(wl Watchlist) {
	s.Width, s.Height = termbox.Size()

	if s.Width%2 == 0 {
		s.LeftWidth = (s.Width / 2) - 1
	} else {
		s.LeftWidth = (s.Width - 1) / 2
	}

	s.RightWidth = s.Width - s.LeftWidth - 1
}

func (c *Cursor) MoveUp() {
	if c.Side == CursorLeft {
		c.CurrentEntry -= 1
	}

	c.Sanitize()
}

func (c *Cursor) MoveDown() {
	if c.Side == CursorLeft {
		c.CurrentEntry += 1
	}

	c.Sanitize()
}

func (c *Cursor) Switch() {
	if c.Side == CursorLeft {
		c.Side = CursorRight
	} else {
		c.Side = CursorLeft
	}

	c.Sanitize()
}

func (c *Cursor) Sanitize() {
	if c.Side != CursorLeft {
		if c.Side != CursorRight {
			c.Side = CursorLeft
		}
	}

	if c.CurrentEntry < 0 {
		c.CurrentEntry = 0
	}
}

func main() {
	wl, err := ParseWatchlist()
	if err != nil {
		log.Fatal(err)
	}

	err = termbox.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Initializing termbox failed: %v\n", err)
		return
	}
	defer termbox.Close()

	termbox.HideCursor()

	timestampQueue := time.NewTicker(500 * time.Millisecond)
	keyQueue := make(chan termbox.Event)
	cursor := NewCursor()
	screen := NewScreen()
	prompt := NewPrompt(screen)
	mode := Normal

	screen.Init(wl)

	go func() {
		for {
			keyQueue <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case event := <-keyQueue:
			if mode == Edit {
				prompt.Handle(event)
			} else {
				if event.Type == termbox.EventKey {
					switch {
					case event.Key == termbox.KeyArrowUp:
						cursor.MoveUp()
						screen.SetRight(wl[cursor.CurrentEntry])
					case event.Key == termbox.KeyArrowDown:
						cursor.MoveDown()
						screen.SetRight(wl[cursor.CurrentEntry])
					case event.Key == termbox.KeyArrowLeft || event.Key == termbox.KeyArrowRight:
						cursor.Switch()
					case event.Ch == 'q' || event.Ch == 'Q':
						return
					case event.Ch == '+':
						entry := screen.GetNew(prompt)
						wl.Add(entry)
						screen.SetLeft(wl)
					default:
						fmt.Fprintf(os.Stderr, "Cannot parse event %v\n", event)
					}
				}
			}
		case <-timestampQueue.C:
			screen.DrawToTerm(cursor)

		}
	}
}
