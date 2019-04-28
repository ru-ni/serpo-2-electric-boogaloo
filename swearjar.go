package main

import (
	//"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func handleSwears(s *discordgo.Session, m *discordgo.Message, split []string) bool {
	//called once on every message seen by serpo
	//and on message updates,
	//so we have to be careful about doupledipping
	//
	//Loop over the split and check if any words
	// are in the swearMap, if they are then it
	// tallies them up in inputMap for later logging
	inputMap := make(map[string]int)
	swearMap := grabSwears()
	for _, word := range split {

		returnTest := strings.Split(word, "\r\n")
		if len(returnTest) != 1 { //There's some stacked swears in there
			//loop over each thing in the stack
			for i := 0; i < len(returnTest); i++ {
				if _, ok := swearMap[returnTest[i]]; ok {
					inputMap[returnTest[i]] += 1
				}
			}
		}
		//we've taken care of all our edgecases
		if _, ok := swearMap[word]; ok {
			inputMap[word] += 1
		}
	}
	//sendm(m, fmt.Sprint(inputMap))
	//Now we need to log the swears we've found
	antispam = append(antispam, m.ID)
	return logSwears(m.Author.ID, inputMap)
}

func grabSwears() map[string]int {
	return mapSwears(getConfig("swearList"))
}
func mapSwears(rawlist string) map[string]int {
	//take a \r\n separated list
	//and fill a map with swear:value
	swearMap := make(map[string]int)
	for _, line := range strings.Split(rawlist, "\r\n") {
		kvSplit := strings.Split(line, "-")
		key := kvSplit[0]
		value := getNum(kvSplit[1])
		swearMap[key] += value

	}
	return swearMap
}

//TODO: change to logSwears and do it in batch to reduce io
func logSwears(id string, inputMap map[string]int) bool {
	if len(inputMap) < 1 {
		return false //called on a junk thing
	}
	//grab the users swearjar
	//append new swear to the end
	//update the total heading
	//save new jar
	//path := id + "/swears"
	//time to 'render' the swears for logging
	swearbatch := ""
	for i, v := range inputMap {
		swearbatch += i + "-" + strconv.Itoa(v) + "\r\n"
	}
	appendFile(botID, "swears", swearbatch, false)
	//got that over with, now for some housekeeping fluff
	//swearmap := grabSwears()
	file := readFile(botID, "swears")

	statMap := parseJar(file)

	return setFile(botID, "swears", statMap["total"]+"-Total;"+statMap["payload"])

}

func parseJar(log string) map[string]string {
	//take a swearjar formatted string
	//and spit out some statistic

	/*
		0-Total;
		blah-5
		shid-1
	*/
	//all we really need so far
	statMap := make(map[string]string)

	totOld := strings.Split(log, "-")[0]

	statMap["oldTotal"] = totOld

	totNew := 0
	for i, line := range strings.Split(log, "\r\n") {
		if i == 0 {
			continue
		}
		if strings.Contains(line, "-") {
			candidate := strings.Split(line, "-")[1]

			if isNum(candidate) {
				totNew += getNum(candidate)
			}
		}

	}
	statMap["total"] = strconv.Itoa(totNew)
	statMap["payload"] = strings.Split(log, ";")[1]
	return statMap
}
