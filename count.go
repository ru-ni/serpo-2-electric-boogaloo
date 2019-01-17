package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func isNum(n string) bool {
	_, err := strconv.Atoi(n)
	if err != nil {
		return false
	} else {
		return true
	}
}

func isWorthy(num int) bool {
	switch true {
	case num%500 == 0:
		return true
	case num%1000 == 0:
		return true
	}

	return false
}
func getNum(n string) int {
	i, err := strconv.Atoi(n)
	if err != nil {
		return -1
	} else {
		return i
	}
}

func checkHundy(h []*discordgo.Message) bool {
	/*
		Run through and make sure that no numbers are missed

		first number should be the highest
			then decrement by one each step
	*/
	i := -1
	lastauth := ""
	for x := 0; x != len(h); x++ {
		fmt.Println(i, getNum(h[x].Content))
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
	hundy, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
	if err != nil {
		fmt.Println("Err:", err)
	} else {
		return hundy
	}
	return nil
}
func handleCount(s *discordgo.Session, m *discordgo.Message) {

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

	if !checkHundy(grabHundy(s, m)) {
		//Something is awry

		//first we'll try deleting the input message
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		fmt.Println(checkHundy(grabHundy(s, m)))
	}
}
