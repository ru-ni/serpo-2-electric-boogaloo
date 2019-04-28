package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"

	"os"
	"os/signal"
	"syscall"
)

var ds *discordgo.Session

var botID string

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	//TODO
	//cycle through our custom database to make sure
	//that all configs are present for getConfig()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	ds = dg
	botID = dg.State.User.ID
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(handleMessage)
	dg.AddHandler(handleEdit)
	dg.AddHandler(handleDelete)
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
