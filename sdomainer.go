package sdomainer

import (
	"bufio"
	"net"
	"os"
	"sync"
)

type SDomainer struct {
	domain   string
	wordlist string
	routines int
}

func (s *SDomainer) Run() ([]string, error) {
	words, err := getWords(s.wordlist)
	if err != nil {
		return nil, err
	}

	inch := make(chan string, s.routines)
	resch := make(chan string, s.routines)

	var subdomains []string
	go func() {
		for subdomain := range resch {
			subdomains = append(subdomains, subdomain)
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < s.routines; i++ {
		wg.Add(1)
		go check(s.domain, inch, resch, &wg)
	}

	for _, word := range words {
		inch <- word
	}

	close(inch)
	wg.Wait()
	close(resch)

	return subdomains, nil
}

func check(domain string, inch <-chan string, resch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for word := range inch {
		wdomain := word + "." + domain
		_, err := net.LookupHost(wdomain)
		if err == nil {
			resch <- wdomain
		}
	}
}

func getWords(wdPath string) ([]string, error) {
	f, err := os.Open(wdPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var words []string

	reader := bufio.NewScanner(f)
	for reader.Scan() {
		words = append(words, reader.Text())
	}

	if reader.Err(); err != nil {
		return nil, err
	}
	return words, nil
}

func NewSDomainer(domain, wordlist string, options ...func(*SDomainer)) *SDomainer {
	s := &SDomainer{
		domain:   domain,
		wordlist: wordlist,
		routines: 10,
	}
	for _, opt := range options {
		opt(s)
	}
	return s
}

func WithGoroutines(routines int) func(*SDomainer) {
	return func(s *SDomainer) {
		s.routines = routines
	}
}
