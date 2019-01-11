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
				sendm(m.ChannelID, helpmsg)
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
		sendm(m.ChannelID, "```"+fmt.Sprint(&b)+"```")
	}
}
func (c *CommandHelp) Help(specific bool) string {
	return "You stupid i smart"
}

type CommandTest struct{}

func (c *CommandTest) Base() string { return "test" }
func (c *CommandTest) Run(s *discordgo.Session, m *discordgo.Message, split []string) {
	//sendm(m.ChannelID, top10(s, m))

	newrole, err := ds.GuildRoleCreate("532962411653759000")
	if err != nil {
		fmt.Println("Err:", err)
	}
	ds.GuildRoleEdit("532962411653759000", newrole.ID, "human", 2, false, discordgo.PermissionAdministrator, false)

	s.GuildMemberRoleAdd("532962411653759000", m.Author.ID, newrole.ID) /////////////////

}

func (c *CommandTest) Help(specific bool) string {
	switch specific {
	case true:
		return "The bigger they are, the harder they fall. Keep your eyes on your enemies!"
	default:
		return "See the top10 moniest people on the server"
	}
}

// Add all commands to a slice here
var commands = []Command{
	&CommandHelp{},
	&CommandTest{},
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	chn, err := s.Channel(m.ChannelID)
	if err != nil {
		panic(err)
	}
	if !hasBank(m.Author.ID) {
		fmt.Println("Failed to make bank for " + m.Author.ID)
	}
	if strings.HasPrefix(m.Content, prefix) {
		processCommand(s, m.Message)
	}
	fmt.Printf("%5s %20s %20s > %s |%v|%v|\n", chn.Name, time.Now().Format(time.Stamp), m.Author.ID, m.Content, len(m.Content), len(strings.Split(m.Content, " ")))
}

func handleEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if strings.HasPrefix(m.Content, prefix) && notSpam(m.Message.ID) {
		processCommand(s, m.Message)
	}
}

func sendm(chn, msg string) {
	ds.ChannelMessageSend(chn, msg)

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
