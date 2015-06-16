package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type listEntry struct {
	RecommendedBy string
	Comment       string
	Date          time.Time
}

var watchlist = flag.String("watchlistlocation", "/home/koebi/filme/watchlist/", "location of watchlist directory")

type Scanner struct {
	*bufio.Scanner
}

func NewScanner() *Scanner {
	return &Scanner{
		bufio.NewScanner(os.Stdin),
	}
}

func (s *Scanner) getString(prompt string) (string, error) {
	fmt.Println(prompt)

	if !s.Scan() {
		return "", s.Err()
	}
	return strings.TrimSpace(s.Text()), nil
}

func main() {
	flag.Parse()
	movie := flag.Arg(0)
	s := NewScanner()

	// ask for movie, if not provided as an arg
	for movie == "" {

		fmt.Println("Want to add/alter a movie on your watchlist? Give a name.")
		if !s.Scan() {
			movie = s.Text()
		}
		if strings.Contains(movie, "/") {
			fmt.Println("Movie names cannot contain /.")
		}
	}

	// decode given file, create, if not
	var wl listEntry
	var file os.File

	_, err := os.Stat(*watchlist + movie)
	if os.IsNotExist(err) {
		file, err := os.OpenFile(*watchlist+movie, os.O_CREATE, 0)
		defer file.Close()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}
	} else if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	} else {
		file, err := os.Open(*watchlist + movie)
		defer file.Close()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}

		_, err = toml.DecodeFile(*watchlist+movie, wl)
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}
	}

	// get alternative info
	choice, err := s.getString("Add alternative Information?\n[r]ecommended  [o]ther\n [Enter] to save.")
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}

	switch {
	case choice == "r":
		wl.RecommendedBy, err = s.getString("Who recommended the movie?")
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}

	case choice == "o":
		wl.Comment, err = s.getString("Other Infos?")
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}
	}

	// encode to file
	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(wl)
	if err != nil {
		fmt.Println("Something went wrong", err)
		return
	}

	file.Write([]byte(buf.String()))
	if err != nil {
		fmt.Println("Something went wrong", err)
		return
	}

	fmt.Printf("Watchlist-Entry %s created", movie)
	return
}
