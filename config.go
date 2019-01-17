package main

import (
	"fmt"
	"io/ioutil"
)

var configPath = "467861483787780107/"

func getConfig(n string) string {
	files, err := ioutil.ReadDir("./bank/" + configPath)
	var slice []string
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range files {
		if f.Name() == n {
			slice = append(slice, readFile(configPath+n))
		} else if n == "all" {
			slice = append(slice, f.Name())
		}

	}
	return fmt.Sprint(slice)[1 : len(fmt.Sprint(slice))-1]
}
func setConfig(n, seed string) bool {
	if hasFile(configPath + n) {
		return setFile(configPath+n, seed)
	} else {
		return false
	}
}
func appendConfig(n, str string) bool {
	if hasFile(configPath + n) {
		return setFile(configPath+n, readFile(configPath+n)+"\r\n"+str)
	} else {
		return false
	}
}
func makeConfig(n string) {
	if !hasFile("467861483787780107" + n) { //Check serpo's userid folder to see if the config entry already exists
		//if not lets make it
		makeFile("467861483787780107", n, "")

	}
}
