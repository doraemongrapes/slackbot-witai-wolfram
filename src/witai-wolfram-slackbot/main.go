package main

import (
	"log"
	"os"

	wolfram "github.com/Krognol/go-wolfram"
	wit "github.com/christianrondeau/go-wit"
	"github.com/nlopes/slack"
)

var (
	slackClient   *slack.Client
	witClient     *wit.Client
	wolframClient *wolfram.Client
)

func main() {
	slackClient = slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))
	witClient = wit.NewClient(os.Getenv("WITAI_ACCESS_TOKEN"))
	wolframClient = &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}

	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			go handleMessage(ev)
		}
	}

}

func handleMessage(ev *slack.MessageEvent) {
	//fmt.Printf("%v\n", ev)
	resp, err := witClient.Message(ev.Msg.Text)
	if err != nil {
		log.Printf("unable to get wit.ai response: %v", err)
		return
	}

	var (
		confidenceThreshold = 0.5
		topEntity           wit.MessageEntity
		topEntityKey        string
	)

	for entityKey, entityList := range resp.Entities {
		for _, entity := range entityList {
			if entity.Confidence > confidenceThreshold && entity.Confidence > topEntity.Confidence {
				topEntity = entity
				topEntityKey = entityKey
			}
		}
	}

	replyToUser(ev, topEntityKey, topEntity)
}

func replyToUser(ev *slack.MessageEvent, entityKey string, entity wit.MessageEntity) {
	switch entityKey {
	case "greetings":
		//		fmt.Println(ev.User, ev.Username)

		slackClient.PostMessage(ev.User, slack.MsgOptionText("Hi", false), slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			AsUser: true,
		}))
		return
	case "wolfram_search_query":
		res, err := wolframClient.GetSpokentAnswerQuery(entity.Value.(string), wolfram.Metric, 1000)
		if err != nil {
			log.Printf("unable to get wolfram result: %v", err)
			return
		}

		slackClient.PostMessage(ev.User, slack.MsgOptionText(res, false), slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			AsUser: true,
		}))
		return
	}

	slackClient.PostMessage(ev.User, slack.MsgOptionText("-\\___(0_0)___//-",false), slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
		AsUser: true,
	}))
}
