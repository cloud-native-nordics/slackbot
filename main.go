package main

import (
	"fmt"
	"log"
	"os"
	"github.com/sazo/slackbot/pkg/flow"
	"github.com/sazo/slackbot/pkg/slackhelper"
	"github.com/nlopes/slack"
)

func main() {
	api := slack.New(
		os.Getenv("SLACK_TOKEN"),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	fm := &flow.FlowManager{
		Questions: []*flow.Question{
			&flow.Question{
				Order: 1,
				Question: "Where are you located (e.g. your country and city)?",
				IsFirst: true,				
			},
			&flow.Question{
				Order: 2,
				Question: "Tell us a little about your self, and your experience with Cloud Native (e.g. your background, where you work, which technologies are you using/looking at)",
				IsFirst: false,
			},
			&flow.Question{
				Order: 3,
				Question: "Thank you, and the last question: What are you hoping to get from this community?",
				IsFirst: false,
			},
		},
	}

	rtm := api.NewRTM()
	go rtm.ManageConnection()
	
	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

			if(ev.BotID != ""){
				continue
			}
			
			userID := ev.User
			isInFlow := fm.IsInFlow(userID) 
			if(!isInFlow){
				fm.AddNew(ev.Channel, userID)
			}
			flowObj, err := fm.Get(userID)

			if(err != nil){
				panic(err)
			}
			
			q, _ := flowObj.GetPreviousQuestion()
			if q != nil {
				q.Answer = ev.Text
			}
			currentQuestion, _ := flowObj.GetCurrentQuestion()
			flowObj.NextQuestion()
			

			if(currentQuestion == nil){ // true if no more question
				fm.Remove(userID)
				rtm.PostMessage(
					flowObj.ChannelID,
					slack.MsgOptionText("Thats all thank you! :wave:", false),
					slack.MsgOptionUser(flowObj.UserID),
					slack.MsgOptionAsUser(true),
				);
				slackhelper.SendSummary(rtm, flowObj)
			}else{
				rtm.PostMessage(
					flowObj.ChannelID,
					slack.MsgOptionText(currentQuestion.Question, false),
					slack.MsgOptionUser(flowObj.UserID),
					slack.MsgOptionAsUser(true),
				);
			}
		
		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		case *slack.RTMEvent:
			fmt.Printf("RTMEvent: %v\n", ev)

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}

}