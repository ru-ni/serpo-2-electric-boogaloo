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
	if err != nil {
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

func bankX(id, target string, dolla int) string { // give it your id, their id, how much you wanna xfer. and it'll give you the reply
	if dolla > 0 {
		if dolla > getMoney(id) {
			return fmt.Sprintf("Stop being so poor, you don't have that kind of money.")
		} else {
			takeMoney(id, dolla)
			putMoney(target, dolla)
			return fmt.Sprintf("%d dolla transfered from <@%s>(new balance: %d) to <@%s>(new balance: %d)", dolla, id, getMoney(id), target, getMoney(target))
		}

	} else {
		return fmt.Sprintf("There's no theivery allowed here, that's for staff only.")
	}
}
func getMoney(id string) int {
	b, err := ioutil.ReadFile("bank/" + id + "/money.txt") //try to see how much money they have
	if err != nil {
		panic(err)
	}
	s := string(b)
	bumhole, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Something went wrong, couldn't find what I was looking for. Try it again maybe? Error: %s", err)
		return 0
	} else {
		return bumhole //returns int bumhole (sounds northern owo)
	}
}

func putMoney(id string, dolla int) {
	startmoney := getMoney(id)
	startmoney += dolla

	err := ioutil.WriteFile("bank/"+id+"/money.txt", []byte(strconv.Itoa(startmoney)), 777) //try to deposit money into their account
	if err != nil {
		panic(err)
	}
}
func takeMoney(id string, dolla int) {
	startmoney := getMoney(id)
	startmoney -= dolla

	err := ioutil.WriteFile("bank/"+id+"/money.txt", []byte(strconv.Itoa(startmoney)), 777) //try to take money from their account
	if err != nil {
		panic(err)
	}
}

func logSwears(id, msg string) {
	if _, err := os.Stat("bank/" + id + "/swears.txt"); os.IsNotExist(err) {
		makeFile(id, "swears", "Total-0;\r\n")
	} else {
		b, err := ioutil.ReadFile("bank/" + id + "/swears.txt") //Get what's there so we don't overwrite
		os.Chmod("bank/"+id+"/swears.txt", 0777)
		if err != nil {
			panic(err)
		}
		s := string(b)
		s += msg
		err1 := ioutil.WriteFile("bank/"+id+"/swears.txt", []byte(s+"\r\n"), 0777) ///try to write to file.
		os.Chmod("bank/"+id+"/swears.txt", 0777)
		if err1 != nil {
			panic(err1)
		}
	}

}
func grabSwears(id string) string {
	b, err := ioutil.ReadFile("bank/" + id + "/swears.txt") //see their swears
	if err != nil {
		panic(err)
	}
	ss := string(b) // swear string
	total := 0
	post := strings.Split(ss, ";")[1] // [1] == the \r\n list of swears

	tally := make(map[string]int) // used to display a tally per swear
	var tb bytes.Buffer
	data := [][]string{}

	table := tablewriter.NewWriter(&tb)
	table.SetHeader([]string{"Word", "Count"})
	table.SetBorder(false) // Set Border to false

	for _, word := range strings.Split(post, "\r\n") {
		tally[strings.Split(word, "-")[0]] += 1
		i, _ := strconv.Atoi(strings.Split(word, "-")[1])
		total += i
	} //run through the swears and add to the tally

	stally := rankByWordCount(tally) //sort this binch
	for _, pair := range stally {
		data = append(data, []string{pair.Key, strconv.Itoa(pair.Value)})
	} //sort the table out

	table.AppendBulk(data) // Add Bulk Data
	table.SetAlignment(1)
	table.Render()

	//kingtime := make(map[string]int)
	err1 := ioutil.WriteFile("bank/"+id+"/swears.txt", []byte("Total-"+strconv.Itoa(total)+";"+post), 0777)
	//write the file back with a re-calculated total number
	os.Chmod("bank/"+id+"/swears.txt", 0777)
	if err1 != nil {
		return fmt.Sprintf("%v", err)
	} else { //return either the error, or the output from the whole thing
		return "```" + fmt.Sprint(&tb) + "```"
	}
}
