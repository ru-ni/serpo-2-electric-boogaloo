package main

import (
	"bytes"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	//"io/ioutil"
	"strconv"
	"strings"
	"time"
)

/*
how to handle removing numbers from main tree and maintaining per user logs

do user counts in the same way as serpo main tree counts,
then scan all the directories (look into some search library/func?) to find
all occurances of the numbers we're hunting.

optimizations needed to lessen the blow if i can't find a search func;

collect the list of numbers I'm hunting first and see if i can't do them
all at the same time



*/

func hasNum(n string) bool {
	//check that a number exists in the logs
	//first let's check serpo's main tree since it has id's

	serpoCounts := readConfigSubFolder("counts")
	inTree := false
	inUser := false
	id := ""
	for _, v := range serpoCounts {
		if strings.Split(string(v), "-")[0] == n {
			//it exists in serpos tree
			inTree = true
			id = strings.Split(string(v), "-")[1]
		}
	}
	if id != "" {
		userCounts := readFile(id, "counts")
		//check users countlog for number
		for _, v := range userCounts {
			if string(v) == n {
				inUser = true
			}
		}
	}
	if inTree || inUser {
		return true
	} else {
		return false
	}

}
func treeHasNum(n string) bool {
	serpoCounts := readConfigSubFolder("counts")
	for _, v := range serpoCounts {
		if strings.Split(string(v), "-")[0] == n {
			return true
		}
		//return strings.Split(string(v), "-")[0] == n
	}
	return false
}
func getCounts(id string) []string {

	//return strings.Split(readFile(id+"/counts"), "\r\n")FOR USER LOGGING?
	return strings.Split(id, "")
}
func remNum(n string) bool {
	//remove a number from the logs
	return n == n
}
func logCount(m *discordgo.Message) bool {
	/*
		-log count in serpos main tree
		-log count in users folder
	*/
	//serpo log
	return makeFile(configPath[:len(configPath)-1], fmt.Sprintf("counts/%v-%v", m.Content, m.Author.ID), m.ID)
	//make a new file in the config path* with / removed, put it in "counts/number-userID", and store the messsage id inside
	//user folder log

	//return makeFile(m.ID, "/counts/"+m.Content, "") /*setFile(req, before+m.Content+"\r\n") FOR USER LOGGING?
}
func logCounts(archive []*discordgo.Message) bool {
	for _, v := range archive {
		defer logCount(v)
	}
	return true
}

func checkHundy(h []*discordgo.Message) bool {
	/*
		Run through and make sure that no numbers are missed

		first number should be the highest
			then decrement by one each step
	*/
	i := -1
	lastauth := ""
	for x := 0; x != len(h)-1; x++ {

		if i == -1 {
			//This will be the number JUST inputted
			if isNum(h[x].Content) {
				i = getNum(h[x].Content)
				lastauth = h[x].Author.ID
			} else {
				return false
			}
		} else {
			//this will be the rest of the history
			if isNum(h[x].Content) {
				if lastauth == h[x].Author.ID {
					return false
				}
				if i == getNum(h[x].Content)+1 {
					//this step is smaller than the last
					i = getNum(h[x].Content)
					lastauth = h[x].Author.ID
					continue
				} else {
					//fmt.Println(i, getNum(h[x].Content))
					return false
				}
			} else {
				return false
			}
		}

	}
	return true
}

func grabHundy(s *discordgo.Session, m *discordgo.Message) []*discordgo.Message {
	//grab the last 100 to do an integrity check
	hundy, err := s.ChannelMessages(getConfig("countChannel"), 100, "", "", "")
	if err != nil {
		fmt.Println("Err:", err)
	} else {
		return hundy
	}
	return nil
}

func grabX(s *discordgo.Session, m *discordgo.Message, x string) []*discordgo.Message {
	//grab all for an archive
	lastid := ""      //used for pagification of results
	goal := getNum(x) //where we wanna stop
	current := 0      //so we know when to stop
	archive, err := s.ChannelMessages(getConfig("countChannel"), 100, "", "", "")
	if err != nil {
		fmt.Println("Err:", err)
	} else {
		//first 100 loaded fine so we know what our ID and number is for pagification
		lastid = archive[len(archive)-1].ID
		current = getNum(archive[len(archive)-1].Content)
		for current > goal {
			//grab it all
			next, err := s.ChannelMessages(getConfig("countChannel"), 100, lastid, "", "")
			if err != nil {
				fmt.Println("Err:", err)
			} else {
				lastid = next[len(next)-1].ID
				current = getNum(next[len(next)-1].Content)
				for _, thing := range next {
					//add things to the initial archive
					archive = append(archive, thing)
				}
			}
		}
	}
	return archive
}
func handleCount(s *discordgo.Session, chn *discordgo.Channel, m *discordgo.Message) {

	if len(m.Attachments) != 0 {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		return
	} //stop people sending images with their numbers
	if strings.HasPrefix(m.Content, "0") {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		return
	} //stops 0000000000000000000251 from being valid

	if !isNum(m.Content) {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		return
	}
	hundy := grabHundy(s, m)

	switch checkHundy(hundy) {
	case true:
		//alls good in the chain
		var oldkings []string
		for i, msg := range hundy { //figure out who to remove counter from and who to apply it to. (can't just do -1 bc of discord lag?)
			if i < 5 {
				//check the last 5 messages to find out who we need to un-king.
				if msg.Author.ID != m.Author.ID { //if we find outselves then just skip it to save api timeouts
					oldkings = append(oldkings, msg.Author.ID)
				}
			} else {
				break
			}
		}

		for _, ID := range oldkings {
			s.GuildMemberRoleRemove(chn.GuildID, ID, getConfig("countRole"))

		}
		//add the role
		s.GuildMemberRoleAdd(chn.GuildID, hundy[0].Author.ID, getConfig("countRole"))
		//log the count
		logCount(m)
	case false:
		for !checkHundy(hundy) {
			//Something is awry

			//first we'll try deleting the input message
			s.ChannelMessageDelete(m.ChannelID, hundy[0].ID)
			hundy = grabHundy(s, m)
		}
	}

}

func tallyCounts(archive []*discordgo.Message) string {
	//var tallies []PairList //holds all the stats to be sorted later

	tCount := make(map[string]int) //Tally of # of counts per username
	tKoth := make(map[string]int)  //Length of time someone was last counter
	tResp := make(map[string]int)  //How fast someone snatches the crown (avg)
	var lastTimeObject time.Time
	for i, msg := range archive {
		tCount[msg.Author.Username] += 1
		//tCount done
		if i+1 < len(archive)-1 && i-1 > 0 {
			/*
				We need to make sure that we're "in the middle" so that
				our +1 and -1 checks don't go out of bounds.

				step through the entire archive and note the following:
					¬timestamps for i, i-1, i+1
					avdiff	= i-1.sub(i)
					diff	= i.sub(i+1)
			*/
			t, err := dateparse.ParseLocal(fmt.Sprint(msg.Timestamp))
			if err != nil {
				panic(err.Error())
			}
			tplus1, err := dateparse.ParseLocal(fmt.Sprint(archive[i+1].Timestamp))
			if err != nil {
				panic(err.Error())
			}
			tminus1, err := dateparse.ParseLocal(fmt.Sprint(archive[i-1].Timestamp))
			if err != nil {
				panic(err.Error())
			}
			avdiff := tminus1.Sub(t)
			diff := t.Sub(tplus1)
			tKoth[msg.Author.Username] += int(diff.Seconds())   //get the diff and add the seconds
			tResp[msg.Author.Username] += int(avdiff.Seconds()) //Fill this map and then divide it by pair.Value later

		} //The 1st and last message should be negligable for stats
		lasttobj, err := dateparse.ParseLocal(fmt.Sprint(msg.Timestamp))
		if err != nil {
			panic(err.Error())
		} else {
			lastTimeObject = lasttobj
		}
	}
	sCount := rankByWordCount(tCount)
	sKoth := rankByWordCount(tKoth)
	sResp := rankByWordCount(tResp)
	var b bytes.Buffer
	data := [][]string{}

	table := tablewriter.NewWriter(&b)
	table.SetHeader([]string{"Counts", "Nerds", "Total Time", "Avg Resp Time", "Counts per day"})

	for i, pair := range sCount {
		total, _ := time.ParseDuration(strconv.Itoa(tKoth[pair.Key]) + "s")
		average, _ := time.ParseDuration(strconv.Itoa(tResp[pair.Key]/pair.Value) + "s")
		//table.AddRow(pair.Key, pair.Value, total, average)
		//output += fmt.Sprintf("%-10v %6v|\tTotal time: %-10v %10v :Avg resp\n", pair.Key, strconv.Itoa(pair.Value), total, average)
		countsPerDay := float64(pair.Value) / (time.Since(lastTimeObject).Hours() / 24)
		data = append(data, []string{strconv.Itoa(pair.Value), pair.Key, fmt.Sprint(total), fmt.Sprint(average), fmt.Sprintf("%.2f", countsPerDay)})
		if i > 5 {
			break
		}

	}
	table.SetFooter([]string{sCount[0].Key, "-Kings-", sKoth[0].Key, sResp[len(sResp)-1].Key, sCount[0].Key}) // Add Footer
	table.SetBorder(false)

	table.AppendBulk(data) // Add Bulk Data
	table.SetAlignment(1)
	table.Render()

	return strings.Replace(fmt.Sprintf("```%v```", &b), " ", "​ ", -1)
}
