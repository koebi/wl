package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type listEntry struct {
	RecommendedBy string
	Comment       string
	Date          string
}

var wldir = flag.String("wldir", "/home/koebi/filme/watchlist/", "location of watchlist directory")

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
		movie, err := s.getString("Want to add/alter a movie on your watchlist? Give a name.")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Getting your response went wrong: %v\n", err)
			return
		}

		if strings.Contains(movie, "/") {
			fmt.Println("Movie names cannot contain /.")
		}
	}

	// decode given file, create, if not
	wl := map[string]interface{}{"Date": time.Now().Format("02. Jan 2006")}
	var file *os.File
	defer file.Close()

	_, err := os.Stat(*wldir + movie)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stdout, "Creating file %s…\n", filepath.Join(*wldir, movie))
		file, err = os.Create(filepath.Join(*wldir, movie))
		if err != nil {
			fmt.Fprintf(os.Stderr, "File creation failed: %v\n", err)
			return
		}
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "File %s seems wrong: %v\n", filepath.Join(*wldir, movie), err)
		return
	} else {
		fmt.Fprintf(os.Stdout, "Opening file %s…\n", filepath.Join(*wldir, movie))
		file, err = os.OpenFile(filepath.Join(*wldir, movie), os.O_RDWR, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Opening file failed: %v\n", err)
			return
		}

		_, err = toml.DecodeFile(filepath.Join(*wldir, movie), wl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Decoding file %s failed: %v\n", filepath.Join(*wldir, movie), err)
			return
		}
	}

	// get alternative info
	choice, err := s.getString("Add alternative Information?\n[r]ecommended  [o]ther [Enter] to save.")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Getting your choice failed: %v\n", err)
		return
	}

	switch {
	case choice == "r":
		wl["RecommendedBy"], err = s.getString("Who recommended the movie?")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Getting recommendation failed: %v\n", err)
			return
		}

	case choice == "o":
		wl["Comment"], err = s.getString("Other Infos?")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Getting information failed: %v\n", err)
			return
		}
	}

	// encode to file
	err = toml.NewEncoder(file).Encode(wl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encoding File failed: %v\n", err)
		fmt.Println(err)
		return
	}

	fmt.Printf("Watchlist-Entry %s created\n", movie)
	return
}
