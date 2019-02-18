package slackhelper

import (
	"github.com/nlopes/slack"
	"github.com/sazo/slackbot/pkg/flow"
)

func SendSummary(rtm *slack.RTM, f *flow.Flow){
	msg := ":wave: <@"+f.UserID+"> just joined the team! Welcome! :tada::tada::tada: \n\n"
	for _, question := range f.Questions {
		msg += "*"+question.Question+"* \n"
		msg += ""+question.Answer+" \n\n"
	}
	rtm.SendMessage(rtm.NewOutgoingMessage(msg, "CG68D5RU5"))
}