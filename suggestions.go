package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

/*
	this has to handle new messages (suggetions)
	as well as reactions getting added?

	first check the reactions
		if it has reactions then just monitor the vote
	then check the message content
		if its formatted correctly, add the reactions
		if not delete it
*/

func handleSuggestions(s *discordgo.Session, chn *discordgo.Channel, m *discordgo.Message) {
	//valid new suggestion has been posted
	/*start a vote
	Votes span 24h and are split into 5 checks
		every 4.8h the vote is checked and the highest voted option is logged
		after 5 logs:
			best 3 out of 5 wins.
	If no best option won, it will fail.
	*/

	if len(m.Reactions) == 0 {
		//this should be a fresh message
		if strings.HasPrefix(strings.ToLower(m.Content), "suggestion:") {
			//goodly formatted
			/*
				start the vote
			*/
			s.MessageReactionAdd(chn.ID, m.ID, "👍")
			s.MessageReactionAdd(chn.ID, m.ID, "🤔")
			s.MessageReactionAdd(chn.ID, m.ID, "👎")
			go startVote(s, m)
		} else {
			//badly formatted
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}

	}
}
func startVote(s *discordgo.Session, m *discordgo.Message) {
	/*
		Create a file in configpath/suggestions/
		with the format of {MessageID-NextTimestamp}

		Then every 4.8 hours (24h/5) a function will append
		that file with the results of getVotes(m)

		after the final loop it will average the results
		the final winer will have at least 3/5. If nothing
		wins, then the vote fails by default.

		Title: m.ID
		Content;
			STAGE [1-5]
			int(time+4.8h)-0,0,0
			int(time+9.6h)-0,0,0
			int(time+14.4h)-0,0,0
			int(time+19.2h)-0,0,0
			int(time+24h)-0,0,0
		Pathos:
			strings.Split(Content,"\r\n")[STAGE] == current line

			Upon starting the vote, write the file and then
			populate the timestamp sections so have our 'checkpoints'

			Then we'll want to sleep for the duration between time.Now()
			and our current stages' timestamp before checking the vote,
			logging it, and moving to the next stage.


	*/
	//timestamps
	gen := getTime(string(m.Timestamp))
	//seed the initial file

	//now we write the checkpoint and vote boilerplate

	var stage, until int

	for i := 0; i <= getNum(getConfig("suggestionCheckpoints")); i++ {
		var err error
		m, err = s.ChannelMessage(m.ChannelID, m.ID)
		if err != nil {
			doPanic(err, fmt.Sprintf("SuggestionCheckpoint message refresh error: ", stage, until))
		}
		if i == 0 {
			makeFile(botID, "suggestions/"+m.ID, popVote(0, gen.Unix()))
			setStage(m.ID, stage)
			stage, until = checkVote(m.ID)
			fmt.Println("Made file and set the stage", stage, until)
			time.Sleep(time.Second * 1)
		} else {
			//this should be fine to do, we write the nth file
			appendFile(botID, "suggestions/"+m.ID, popVote(i, gen.Unix()), false)
			setStage(m.ID, stage)
			//now we just update our info so we know when to wait for
			//we might also need to do some safety checking here based on stage
			stage, until = checkVote(m.ID)
			fmt.Println("Appended to the file and set the stage", stage, until)

			untilTime := getTime(strconv.Itoa(until))
			durationUntil := getDuration(untilTime, time.Now())
			fmt.Println("Sleeping until ", untilTime, durationUntil)

			time.Sleep(time.Second * time.Duration(durationUntil))
			//we're woke.

			//time to getVotes(m) and replace 0,0,0 with the results.
			setFile(botID, "suggestions/"+m.ID, insertVotes(m))
			fmt.Println("Woke up and set the votes")
		}

	}
	//finish up the vote

	/*	we now have a filled out suggestion file.
		5 tallies over 24hrs, majority wins.

		lets loop over the rawSplit
	*/
	result := finshVote(m)
	fmt.Println("result:", result, (result == ""))
	if result == "" { //no strong winner

	}
	response := ""
	switch result {
	case "":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	case "UP":
		//someone got that bread
		response = fmt.Sprintf("%v, A vote just won by majority.\r\n```%v```", getConfig("staffRole"), m.Content)
	case "DOWN":
		//no clear winner here
		response = fmt.Sprintf("A vote just lost by majority.\r\n```%v```", m.Content)
		s.ChannelMessageDelete(m.ChannelID, m.ID)

	case "HMM":
		//fgsfds, tf is this guy on about
		//lets still make a message in metaChannel
		//but ping the author for details
		response = fmt.Sprintf("<@%v>, hey! That suggestion confused some people, maybe explain it here?", m.Author.ID)
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
	s.ChannelMessageSend(getConfig("metaChannel"), response)
}

func finshVote(m *discordgo.Message) string {
	raw := readFile(botID, "suggestions/"+m.ID)
	rawSplit := strings.Split(raw, "\r\n")
	outMap := make(map[string]int)
	for _, line := range rawSplit {
		if len(line) < 10 {
			continue
		}
		//fmt.Println(strings.Split(line, "-"))
		outMap[countVotes(strings.Split(line, "-")[1])] += 1
	}
	//now we just need to find the largest thing in outMap
	record := 0
	var output string
	for opt := range outMap {
		if outMap[opt] > record {
			record = outMap[opt]
			output = opt
		}
	}
	return output
}
func countVotes(n string) string {
	//n == "5,2,2\r\n"
	var nice []int //number slice NEEDS TO BE A MAP?
	for _, num := range strings.Split(strings.Replace(n, "\r\n", "", -1), ",") {
		if isNum(num) {
			nice = append(nice, getNum(num))
		}
	}
	var output string
	if nice[0] > nice[1] && nice[0] > nice[2] {
		//up won
		output += "UP"
	}
	if nice[1] > nice[0] && nice[1] > nice[2] {
		//down won
		output += "DOWN"
	}
	if nice[2] > nice[0] && nice[2] > nice[1] {
		//hmm won
		output += "HMM"
	}
	return output
}
func insertVotes(m *discordgo.Message) string {
	/* STAGE int(n)
	int(timestamp)-0,0,0
	.. */
	raw := readFile(botID, "suggestions/"+m.ID)
	rawSplit := strings.Split(raw, "\r\n")
	stage, _ := checkVote(m.ID)
	//should have all we need to log the stuff
	var output string
	for i, line := range rawSplit {
		if len(line) < 5 {
			continue
		}
		if i != stage {
			output += line + "\r\n"
		} else {
			//line we want!
			pre := strings.Split(line, "-")[0]
			up, down, hmm := getVotes(m.Reactions)
			output += fmt.Sprintf("%v-%v,%v,%v\r\n", pre, up, down, hmm) //maybe add newline at end
		}
	}
	//fmt.Println("insert raw", raw, "\r\noutput", output)
	return output
}
func setStage(id string, n int) {
	//sets the "STAGE X" text based on number of lines
	raw := readFile(botID, "suggestions/"+id)
	rawSplit := strings.Split(raw, "\r\n")
	stage := len(rawSplit) - 1

	var output string
	for i, line := range rawSplit {
		if i == 0 {
			output += "STAGE " + strconv.Itoa(stage) + "\r\n"
		} else {
			output += line + "\r\n"
		}
	}
	//fmt.Println("set raw", raw, "\r\noutput", output)
	setFile(botID, "suggestions/"+id, output)
}
func checkVote(id string) (int, int) {
	//in this function we'll check how a vote is doing
	//we'll return the stage it's at, and the next checkpoint.
	raw := readFile(botID, "suggestions/"+id)
	rawSplit := strings.Split(raw, "\r\n")
	stage := getNum(strings.Split(rawSplit[0], " ")[1])
	checkpoint := getNum(strings.Split(rawSplit[stage], "-")[0])
	return stage, checkpoint
}
func voteSleep(t int) {
	time.Sleep(time.Second * time.Duration(t))
}
func popVote(stage int, gen int64 /*time.Time*/) string {
	switch stage {
	case 0:
		return fmt.Sprintf("STAGE %v", stage)
	default:
		duration, err := strconv.Atoi(getConfig("suggestionDelay"))
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%v-%v", int(gen)+(stage*(duration*60)), "0,0,0")
	}
}

func getVotes(votes []*discordgo.MessageReactions) (int, int, int) {
	up, down, hmm := 0, 0, 0
	for i := 0; i < len(votes); i++ {
		//fmt.Println(votes[i].Emoji)
		if strings.ToLower(votes[i].Emoji.Name) == "🤔" {
			hmm = votes[i].Count
		}
		if strings.ToLower(votes[i].Emoji.Name) == "👍" {
			up = votes[i].Count
		}
		if strings.ToLower(votes[i].Emoji.Name) == "👎" {
			down = votes[i].Count
		}
	}
	// fmt.Println("REACTIONS", up, down, hmm)
	return up, down, hmm
}
