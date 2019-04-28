package main

import (
	//"fmt"
	"github.com/araddon/dateparse"
	"time"
)

func getTime(timestamp string) time.Time {
	t, err := dateparse.ParseAny(timestamp)
	if err != nil {
		doPanic(err, "Couldn't parse timestamp into time object.")
	}
	return t

}
func getDuration(t1, t2 time.Time) float64 {
	//t2 sub t1 (bigger - larger = duration)
	return float64(t2.Sub(t1).Seconds())
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
