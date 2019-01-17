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

type Command interface {
	Base() string                                                           // Returns the base for this command
	Run(session *discordgo.Session, msg *discordgo.Message, split []string) // Runs this command
	Help(specific bool) string                                              // Called when help, specific is true when we called help on this specific command "!help somecommand"
}

// One command
type CommandHelp struct{}

func (c *CommandHelp) Base() string { return "help" }
func (c *CommandHelp) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	// run this help command
	// You can loop over all the commands in the commands slice and call "Help()" on each of them
	// To list help on all commands here for example, or a specific one if we provided one
	switch len(split) {
	case 2: //-=help CMD?
		for _, cmd := range commands {
			if split[1] == cmd.Base() {
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
			data = append(data, []string{cmd.Base(), cmd.Help(false)})
		}
		table.AppendBulk(data) // Add Bulk Data
		table.SetAlignment(1)
		table.Render()
		sendm(m, "```"+fmt.Sprint(&b)+"```")
	}
}
func (c *CommandHelp) Help(specific bool) string {
	return "You stupid i smart"
}

type CommandTest struct{}

func (c *CommandTest) Base() string { return "test" }
func (c *CommandTest) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	//sendm(m, top10(s, m))

}

func (c *CommandTest) Help(specific bool) string {
	switch specific {
	case true:
		return "This is the command I use to test things"
	default:
		return "beep boop"
	}
}

type CommandConfig struct{}

func (c *CommandConfig) Base() string { return "cfg" }
func (c *CommandConfig) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	//sendm(m, top10(s, m))

	switch len(split) {
	case 1: //%cfg
		sendm(m, fmt.Sprint(getConfig("all")))
	case 2: //%cfg something
		sendm(m, fmt.Sprint(getConfig(split[1])))
	default: //%cfg [set/add/...] something value
		if len(split) >= 4 {
			value := ""
			for i, word := range split {
				if i > 2 {
					value += word
				}
			}
			switch split[1] {
			case "set":
				//overwrite it
				fmt.Println(setConfig(split[2], value))
			case "add":
				//append something to it
				fmt.Println(appendConfig(split[2], value))
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

type CommandCount struct{}

func (c *CommandCount) Base() string { return "count" }
func (c *CommandCount) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	//sendm(m, top10(s, m))
	sendm(m, grabCounts())

	switch len(split) {
	case 1: //%count
		/* return an embed with
		Your Counts
		##
		Your last 5 counts
		xxxx, xxxx, xxxx, xxxx, xxxx
		*/
	case 2: //%count something
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

//////////////////////////////////////////////////////////////////////////
// Add all commands to a slice here
var commands = []Command{
	&CommandHelp{},
	&CommandTest{},
	&CommandCount{},
	&CommandConfig{},
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	chn, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
	}
	if !hasBank(m.Author.ID) {
		fmt.Println("Failed to make bank for " + m.Author.ID)
	}
	if strings.HasPrefix(m.Content, prefix) {
		processCommand(s, m.Message)
	}
	if isNum(m.Content) || m.ChannelID == getConfig("countChannel") { //potential way to handle the counts dynamically?
		handleCount(s, m.Message)
	}
	fmt.Printf("%5s %20s %20s > %s |%v|%v|\n", chn.Name, time.Now().Format(time.Stamp), m.Author.ID, m.Content, len(m.Content), len(strings.Split(m.Content, " ")))
}

func handleEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if strings.HasPrefix(m.Content, prefix) && notSpam(m.Message.ID) {
		processCommand(s, m.Message)
	}
}

func sendm(m *discordgo.Message, msg string) {
	ds.ChannelMessageSendEmbed(m.ChannelID, EmbedStr(" ", " ", [][]string{{m.Author.Username, msg, ""}, {"@", time.Now().Format(time.Stamp), "t"}}))
}

func processCommand(s *discordgo.Session, m *discordgo.Message) {
	presplit := strings.TrimPrefix(m.Content, prefix) //clear the prefix before we split the string
	split := strings.Split(presplit, " ")             //split the sting by the spaces.

	for _, cmd := range commands { //Loop through the commands and run whichever command matches.
		if split[0] == cmd.Base() {
			cmd.Run(s, m, split)
			antispam = append(antispam, m.ID)
		}
	}
}
