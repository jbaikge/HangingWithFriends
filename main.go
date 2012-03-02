package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	MaxLength = 8
	MinLength = 6
)

var (
	debug    = make(chan string)
	done     = make(chan bool)
	letters  string
	verbose  bool
	wordChan = make(chan string)
	words    = make(map[int][]string)
)

func debugHandler() {
	for msg := range debug {
		if verbose {
			log.Println(msg)
		}
	}
}

func filterWords() {
	for word := range wordChan {
		if validWord(word) {
			words[len(word)] = append(words[len(word)], word)
		}
	}
}

func filterWordsWithin(filename string) error {
	go filterWords()
	go sendWordList(filename)
	return nil
}

func printWords(wordList map[int][]string) {
	for i := MinLength; i <= MaxLength; i++ {
		fmt.Printf("%d-letter words\n", i)
		for n, w := range wordList[i] {
			fmt.Printf("%4d. %s\n", n+1, w)
		}
	}
}

func sendWordList(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	r := bufio.NewReader(file)
	for {
		line, err := r.ReadString(byte('\n'))
		if err != nil {
			break
		}
		wordChan <- strings.TrimSpace(line)
	}
	done <- true
	return nil
}

func validWord(word string) bool {
	if len(word) < MinLength || len(word) > MaxLength {
		//debug <- fmt.Sprintf("Invalid length for %s: %d", word, len(word))
		return false
	}
	if !strings.ContainsAny(word, letters) {
		//debug <- fmt.Sprintf("Does not contain required letters: %s [%s]", word, letters)
		return false
	}
	for _, l := range strings.Split(letters, "") {
		i := strings.Index(word, l)
		if i == -1 {
			//debug <- fmt.Sprintf("Does not contain letter: %s [%s]", word, l)
			continue
		}
		word = strings.Replace(word, l, "~", 1)
		//debug <- fmt.Sprint(word)
	}
	return strings.Repeat("~", len(word)) == word
}

func main() {
	var filename string
	flag.StringVar(&filename, "f", "/usr/share/dict/words", "Dictionary list to use")
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.Parse()

	//go debugHandler()

	letters = flag.Arg(0)

	if letters == "" {
		log.Panic("Please provide letters to search for")
	}

	err := filterWordsWithin(filename)
	if err != nil {
		log.Panic("Unable to get words.", err)
	}

	<-done

	printWords(words)
}
