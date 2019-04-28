package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
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

func handleSuggestion(s *discordgo.Session, chn *discordgo.Channel, m *discordgo.Message) {
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
			go startVote(m)
		} else {
			//badly formatted
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}

	}
}
func startVote(m *discordgo.Message) {
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
	gen := getTime(m.ID).Unix()
	s1 := gen + 288
	s2 := s1 + 288
	s3 := s2 + 288
	s4 := s3 + 288
	s5 := s4 + 288
	if s5 < 0 {

	}
	//seed the initial file
	makeFile(botID, "suggestions/"+m.ID,
		fmt.Sprintf("STAGE %v\r\n", 1))

	//now we should finish up then wait
	appendFile(botID, "suggestions/"+m.ID, popVote(1, gen), false)
}
func voteSleep(t int) {
	time.Sleep(time.Second * time.Duration(t))
}
func popVote(stage int, gen int64 /*time.Time*/) string {
	switch stage {
	case 0:
		return fmt.Sprintf("STAGE %v\r\n", stage)
	default:
		return fmt.Sprintf("%v-%v\r\n", int(gen)+(stage*288), "0,0,0")
	}
}

func getVotes(votes []*discordgo.MessageReactions) (int, int, int) {
	up, down, hmm := 0, 0, 0
	for i := 0; i < len(votes); i++ {
		fmt.Println(votes[i].Emoji)
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
	return up, down, hmm
}
