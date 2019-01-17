package main

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/*
The bank is the central method of persistant storage.



	err := ioutil.WriteFile("bank/"+id+"/money.txt", []byte(strconv.Itoa(100)), 0777) //try to deposit money into their account
	os.Chmod("bank/"+id+"/money.txt", 0777)

*/

func makeBank(id string) bool {
	err := os.Mkdir("bank/"+id, 0777) //make their bank directory
	os.Chmod("bank/"+id, 0777)        //chmod bc of some gay ass linux shit
	makeFile(id, "money", "100")
	makeFile(id, "swears", "Total-0;\r\n") //register the files
	makeFile(id, "counts", "")
	if err != nil || id == "467861483787780107" {
		return false
		//panic(err)
	} else {
		return true
	}
}

func makeFile(id, name, seed string) bool {
	err := ioutil.WriteFile("bank/"+id+"/"+name+".txt", []byte(seed), 0777) //try to deposit money into their account
	os.Chmod("bank/"+id+"/"+name+".txt", 0777)
	if err != nil {
		return false
	} else {
		return true
	}
}

func readFile(path string) string {
	os.Chmod("bank/"+path, 0777)
	b, err := ioutil.ReadFile("bank/" + path) //read arbitrary files
	if err != nil {
		panic(err)
	}
	return string(b)
}
func setFile(path, str string) bool {
	err := ioutil.WriteFile("bank/"+path, []byte(str), 0777) //try to deposit money into their account
	if err != nil {
		panic(err)
	}
	return true
}
func hasFile(path string) bool {
	if _, err := os.Stat("bank/" + path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func hasBank(id string) bool {
	if _, err := os.Stat("bank/" + id); os.IsNotExist(err) {
		// path/to/whatever does not exist
		if makeBank(id) { // they didn't have a bank, but dont answer false just yet
			return true
		} else {
			return false //couldn't make a bank for whatever reason so they deffo dont have one
		}
	} else {
		return true // already had one
	}
}

//count bank funcs
func grabCounts() string {
	out := ""
	files, err := ioutil.ReadDir("./bank/counts/")
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range files {
		out += f.Name() + "\n"
	}
	return out
}
func logCount(id, num string) {

}

//config bank funcs
