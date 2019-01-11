package main

import (
	"math/rand"
	"net/url"
	"sort"
	"time"
)

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
func random(min int, max int) int {
	return rand.Intn(max-min) + min
}
func randbool() bool {
	return rand.Float32() < 0.5
}

func ranum(max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	time.Sleep(20 * time.Millisecond)
	r1 := rand.New(s1)
	if max == 0 {
		max = 1
	}
	return r1.Intn(max)
}

// isValidUrl tests a string to determine if it is a url or not.
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	} else {
		return true
	}
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func notSpam(id string) bool {
	doit := true // assume something isn't spam
	for _, v := range antispam {
		if id == v { //if we find it in the list
			doit = false // we need to call off our plan
		}
	}
	return doit //doit should be false if ID exists in antispam
}
