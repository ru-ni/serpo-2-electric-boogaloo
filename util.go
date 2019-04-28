package main

import (
	"fmt"
	"image"
	"math/rand"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

//contains but for a slice instead
func sliceContains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
func mapContains(that map[string]int, this string) bool {
	if _, ok := that[this]; ok {
		return true
	} else {
		return false
	}
}
func isNum(n string) bool {
	_, err := strconv.Atoi(n)
	if err != nil {
		return false
	} else {
		return true
	}
}

func isWorthy(num int) bool {
	switch true {
	case num%500 == 0:
		return true
	case num%1000 == 0:
		return true
	}

	return false
}
func getNum(n string) int {
	i, err := strconv.Atoi(n)
	if err != nil {
		return -1
	} else {
		return i
	}
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

//used to try and only operate on messages once
func notSpam(id string) bool {
	doit := true // assume something isn't spam
	for _, v := range antispam {
		if id == v { //if we find it in the list
			doit = false // we need to call off our plan
		}
	}
	return doit //doit should be false if ID exists in antispam
}

//Pairs used to sort the counting leaderboards
type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
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

// get the local file as a golang image
func getImage(filePath string) image.Image {
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println("Cannot read file:", err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Cannot decode file:", err)
	}
	return img
}
