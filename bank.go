package main

import (
	// "bytes"
	//"fmt"
	// "github.com/olekukonko/tablewriter"
	"io/ioutil"
	"os"
	// "strconv"
	// "strings"
)

func makeBank(id string) bool {
	err := os.Mkdir("bank/"+id, 0777) //make their bank directory
	os.Chmod("bank/"+id, 0777)        //chmod bc of some gay ass linux shit
	makeFile(id, "money", "100")
	makeFile(id, "swears", "0-Total;\r\n") //register the files
	makeFile(id, "counts", "")
	if err != nil || id == botID {
		return false
		//panic(err)
	} else {
		return true
	}
}

func makeFile(id, name, seed string) bool {
	err := ioutil.WriteFile("bank/"+id+"/"+name, []byte(seed), 0777) //try to deposit money into their account
	os.Chmod("bank/"+id+"/"+name, 0777)
	if err != nil {
		return false
	} else {
		return true
	}
}

func readFile(id, name string) string {
	os.Chmod("bank/"+id+"/"+name, 0777)
	b, err := ioutil.ReadFile("bank/" + id + "/" + name) //read arbitrary files
	if err != nil {
		panic(err)
	}
	return string(b)
}
func setFile(id, name, str string) bool {
	err := ioutil.WriteFile("bank/"+id+"/"+name, []byte(str), 0777) //try to deposit money into their account
	if err != nil {
		return false
	}
	return true
}
func appendFile(id, name, str string, prepend bool) bool {
	//update file by tagging something onto the end
	//prepend tags it onto the front instead
	before := readFile(id, name)
	if prepend {
		return setFile(id, name, str+before)
	} else {
		return setFile(id, name, before+str)
	}
}

func hasFile(id, name string) bool {
	if _, err := os.Stat("bank/" + id + "/" + name); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func hasBank(id string) bool {
	path := "bank/" + id
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		if makeBank(id) { // they didn't have a bank, but dont answer false just yet☺
			return true
		} else {
			return false //couldn't make a bank for whatever reason so they deffo dont have one
		}
	} else {
		// already had one
		//but lets check individually just in case

		if !hasFile(id, "money") {
			return makeFile(id, "money", "100")
		}
		if !hasFile(id, "swears") {

			return makeFile(id, "swears", "0-Total;\r\n")
		}
		if !hasFile(id, "counts") {
			return makeFile(id, "counts", "")
		}
		return true
	}
}
