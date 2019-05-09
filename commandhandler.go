package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	/*"strconv"*/
	"bytes"
	"github.com/olekukonko/tablewriter"
	"strings"
	"time"
)

var prefix = "%"
var antispam []string

func sendm(m *discordgo.Message, msg string) {
	ds.ChannelMessageSendEmbed(m.ChannelID, EmbedStr(" ", " ", [][]string{{m.Author.Username, msg, ""}, {"@", time.Now().Format(time.Stamp), "t"}}))
	antispam = append(antispam, m.ID)
}
func sendr(m *discordgo.Message, msg string) {
	ds.ChannelMessageSend(m.ChannelID, msg)
	antispam = append(antispam, m.ID)

}

type Command interface {
	Base(ID string) string                                                  // Returns the base for this command
	Run(session *discordgo.Session, msg *discordgo.Message, split []string) // Runs this command
	Help(specific bool) string                                              // Called when help, specific is true when we called help on this specific command "!help somecommand"
}

// One command
type CommandHelp struct{}

func (c *CommandHelp) Base(ID string) string { return "help" }

func (c *CommandHelp) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	// run this help command
	// You can loop over all the commands in the commands slice and call "Help()" on each of them
	// To list help on all commands here for example, or a specific one if we provided one
	switch len(split) {
	case 2: //-=help CMD?
		for _, cmd := range commands {
			if split[1] == cmd.Base(m.Author.ID) {
				helpmsg := cmd.Help(true)
				sendm(m, helpmsg)

			}
		}
	default:
		var b bytes.Buffer
		data := [][]string{}

		table := tablewriter.NewWriter(&b)
		table.SetHeader([]string{"Name", "Description"})
		table.SetBorder(false) // Set Border to false
		for _, cmd := range commands {
			data = append(data, []string{prefix + cmd.Base(m.Author.ID), cmd.Help(false)})
		}
		table.AppendBulk(data) // Add Bulk Data
		table.SetAlignment(1)
		table.Render()
		sendm(m, "```"+fmt.Sprint(&b)+"```")

	}
}
func (c *CommandHelp) Help(specific bool) string {
	return "You stupid i smort"
}

type CommandEcho struct{}

func (c *CommandEcho) Base(ID string) string {
	if strings.Contains(getConfig("ops"), ID) {

		return "echo"
	}
	return ""
}
func (c *CommandEcho) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	//sendm(m, top10(s, m))
	value := ""
	for i, word := range split {
		if i > 0 {
			value += word
		}
	}
	s.ChannelMessageSend(m.ChannelID, value)

}
func (c *CommandEcho) Help(specific bool) string {
	return "Echoes things back to you"
}

type CommandTest struct{}

func (c *CommandTest) Base(ID string) string {
	if strings.Contains(getConfig("ops"), ID) {

		return "test"
	}
	return ""
}
func (c *CommandTest) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	//sendm(m, top10(s, m))
	// value := ""
	// for i, word := range split {
	// 	if i > 0 {
	// 		value += word
	// 	}
	// }
	// s.ChannelMessageSend(m.ChannelID, value)

}

func (c *CommandTest) Help(specific bool) string {
	switch specific {
	case true:
		return "This is the command I use to test things"
	default:
		return "test command beep boop"
	}
}

type CommandCount struct{}

func (c *CommandCount) Base(ID string) string { return "count" }
func (c *CommandCount) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	//sendm(m, top10(s, m))
	//sendm(m, grabCounts())

	switch len(split) {
	case 1: //%count
		/* return an embed with
		Your Counts
		##
		Your last 5 counts
		xxxx, xxxx, xxxx, xxxx, xxxx
		*/
		sendm(m, commands[3].Help(true))

	case 2: //%count something
		switch split[1] {
		case "stats":
			//
			sendr(m, tallyCounts(grabHundy(s, m)))

		case "help":
			sendm(m, commands[3].Help(true))

		}
	case 3: //%count something x
		switch split[1] {
		case "stats":
			//
			if isNum(split[2]) {

				sendr(m, tallyCounts(grabX(s, m, split[2])))

			}
		case "scan":
			if isNum(split[2]) {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprint(logCounts(grabX(s, m, "1"))))

				//todo
				/*
					check how many counts we have in serpotree
					check what the latest count is in #counting-channel
					check how much of a difference we have
					grabX and log things we were missing
				*/
			}

		}
	}

}

func (c *CommandCount) Help(specific bool) string {
	switch specific {
	case true:
		return "One stop shop for all count stat related things"
	default:
		return "Lists some stats about your counts"
	}
}

type CommandConfig struct{}

func (c *CommandConfig) Base(ID string) string {
	if strings.Contains(getConfig("ops"), ID) {

		return "cfg"
	}
	return ""
}
func (c *CommandConfig) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	//sendm(m, top10(s, m))

	switch len(split) {
	case 1: //%cfg
		sendm(m, commands[4].Help(true))

	case 2: //%cfg something
		switch split[1] {
		case "all":
			sendm(m, fmt.Sprint(getConfig("all")))
			return

		case "help":
			sendm(m, commands[4].Help(true))
			return

		}

		sendm(m, fmt.Sprint(getConfig(split[1])))
	case 3:
		switch split[1] {
		case "create":
			sendm(m, fmt.Sprint(createConfig(split[2])))

		}
	default: //%cfg [set/add/...] something values+
		if len(split) >= 4 {
			value := ""
			for i, word := range split {
				if i > 2 {
					value += word
				}
			}
			switch split[1] {
			case "set", "overwrite":
				//overwrite it
				sendm(m, fmt.Sprint(setConfig(split[2], value)))

			case "add", "+", "append":
				//append something to it
				sendm(m, fmt.Sprint(appendConfig(split[2], value)))

			}
		}

	}
}

func (c *CommandConfig) Help(specific bool) string {
	switch specific {
	case true:
		return "[cfg] displays all registered configs. [cfg [append/set] key value] writes a thing to it"
	default:
		return "View and modify serpo's config"
	}
}

type CommandSwearjar struct{}

func (c *CommandSwearjar) Base(ID string) string { return "jar" }
func (c *CommandSwearjar) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	//list all current swear words
	//view personal swearjars
	//view eachothers swearjars
	//

	switch len(split) {
	case 1: //%jar
		sendm(m, commands[5].Help(true))

		return //help screen thing
	case 2: //%jar [view]
		switch split[1] {
		case "help": //might as well just cover this
			sendm(m, commands[5].Help(true))

		case "view": //view their swearlog file.
			sendm(m, readFile(m.Author.ID, "swears"))

		}
		return
	case 3: //%jar [view] @mention
		switch split[1] {
		case "view":
			if len(m.Mentions) != 0 {
				sendm(m, readFile(m.Author.ID, "swears"))

			}
		}
		return
	}
}

func (c *CommandSwearjar) Help(specific bool) string {
	switch specific {
	case true:
		return "List all your swears, all of friends swears, or a bad boy leaderboard"
	default:
		return "Command for swearjar info"
	}
}

//////////////////////////////////////////////////////////////////////////
// Add all commands to a slice here
var commands = []Command{
	&CommandHelp{},
	&CommandEcho{},
	&CommandTest{},
	&CommandCount{},
	&CommandConfig{},
	&CommandSwearjar{},
}

func handler(s *discordgo.Session, m *discordgo.Message, chn *discordgo.Channel) {
	/*
			this function acts as a funnel for handleMessage handleEdit and
			all other handlers that i wind up adding. I should look to my
			command handler code for something to base this off, perhaps.

			the control flow should be like
			serpo.go
				dg.addhandler(handleX)
		ur here>	handler(session, messageObject, channelObject)



	*/

	if chn.GuildID != "532962411653759000" {
		return
	} /////////EOF TEST SEGREGATOR
	presplit := strings.TrimPrefix(m.Content, prefix) //clear the prefix before we split the string
	split := strings.Split(presplit, " ")             //split the sting by the spaces.

	if !m.Author.Bot {
		if !hasBank(m.Author.ID) {
			fmt.Println("Failed to make bank for " + m.Author.ID)
		}
		//handle any swears that might be in it
		go handleSwears(s, m, split)
	}
	//detect a %command and process it as long as it hasn't been processed already
	if strings.HasPrefix(m.Content, prefix) && notSpam(m.ID) {
		go processCommand(s, chn, m, split)
	}

	//If something changed in the count channel then we'll just check its still unbroken
	if m.ChannelID == getConfig("countChannel") {
		go handleCount(s, chn, m)
	}

	if m.ChannelID == getConfig("suggestionChannel") {
		go handleSuggestions(s, chn, m)
	}
	fmt.Printf("%5s %20s %20s > %s |%v|%v|\n", chn.Name, time.Now().Format(time.Stamp), m.Author.ID, split, m.Author.Bot, "<Bot")
}
func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	chn, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
		return
	}
	go handler(s, m.Message, chn)

}

func handleEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	chn, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if m.Author != nil { //anti doubletap crash
		go handler(s, m.Message, chn)
	}
}
func handleDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	chn, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
		return
	}
	//bugged out when i sent m.Message
	//so we'll send the message one before that
	last1, err := s.ChannelMessages(m.ChannelID, 1, "", "", "")
	if err != nil {
		fmt.Println("Err:", err)
	} else {
		go handler(s, last1[0], chn)
	}
}

func processCommand(s *discordgo.Session, chn *discordgo.Channel, m *discordgo.Message, split []string) {

	for _, cmd := range commands { //Loop through the commands and run whichever command matches.
		if split[0] == cmd.Base(m.Author.ID) { //send ID for elevation checks on "mod" commands

			cmd.Run(s, m, split)
			antispam = append(antispam, m.ID)
		}
	}
}
