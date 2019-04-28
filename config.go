package main

//todo?
//go from using single file key:values and
//start involving some directories as categories
import (
	"fmt"
	"io/ioutil"
)

var configPath = botID + "/"

func getConfig(n string) string {
	//for n = `all`
	/* read the serpo config path and loop over all entries.
	and add them to our output slice.
	*/
	//for n = x
	/* read the serpo config path and loop over all entries.
	if n == x then all file contents to output slice.
	*/
	files, err := ioutil.ReadDir("./bank/" + configPath)
	var slice []string //holds our output
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range files {
		if f.Name() == n { //if we find `n` then we'll post its value
			slice = append(slice, readFile(botID, n))
		} else if n == "all" { //if `n` is "all", we'll just list all the names
			slice = append(slice, f.Name())
		}

	}
	output := fmt.Sprint(slice)[1 : len(fmt.Sprint(slice))-1] //returns all but 1st and last chars
	if len(output) < 1 {
		//this means someone wanted something that doesn't exist
		//should we make it, log it, or something
	}
	return output
}
func setConfig(n, seed string) bool {
	if hasFile(botID, n) {
		return setFile(botID, n, seed)
	} else {
		return false
	}
}
func readConfigSubFolder(req string) string {
	out := ""
	files, err := ioutil.ReadDir("./bank/" + botID + req + "/")
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range files {
		out += f.Name() + "\r\n"
	}
	return out
}
func createConfig(n string) bool {
	return makeFile(botID, n, "")
}
func appendConfig(n, str string) bool {
	if hasFile(botID, n) {
		return setFile(botID, n, readFile(botID, n)+"\r\n"+str)
	} else {
		return false
	}
}
func makeConfig(n string) {
	if !hasFile(botID, n) { //Check serpo's userid folder to see if the config entry already exists
		//if not lets make it
		makeFile(botID, n, "")

	}
}
