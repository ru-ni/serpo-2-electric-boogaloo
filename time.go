package main

import (
	//"fmt"
	"github.com/araddon/dateparse"
	"strconv"
	"time"
)

func getTime(timestamp string) time.Time {
	t, err := dateparse.ParseLocal(timestamp)
	if err != nil {
		doPanic(err, "Couldn't parse timestamp into time object.")
	}
	return t

}

func unixToTime(stamp int) time.Time {
	i, err := strconv.ParseInt("1405544146", 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm
}

func getDuration(t1, t2 time.Time) float64 {
	return float64(t1.Sub(t2).Seconds())
}

/*
	if count > 2 {

	t, errt := dateparse.ParseLocal(fmt.Sprint(star.Timestamp))
	if errt != nil {
		panic(errt.Error())
	} else {
		if int(t.Sub(time.Now()).Seconds()) < -(60 * 60 * 24) {
		} else {
			err := s.ChannelMessagePin(p.ChannelID, p.MessageID)
			if err != nil {
				s.ChannelMessageSend("468602601290727454", fmt.Sprint(err))
			} else {

				s.MessageReactionAdd(p.ChannelID, p.MessageID, "📍")
			}
		}
	}*/
