package main

import (
	// "bytes"
	// "encoding/json"
	"fmt"

	// "io/ioutil"
	"net/http"

	"./skype"
)

func main() {
	activity := &skype.Activity{}

	activity.Type = "event"
	activity.Text = "another test message"
	activity.From = skype.ChannelAccount{ID: skype.DMP_BOT_SKYPE, Name: skype.DMP_BOT}
	activity.ServiceURL = skype.WebchatURL
	client := &http.Client{}

	sskype := skype.SkypeService{Client: client,
		AppID:       skype.MIRCOSOFT_APP_ID,
		AppPassword: skype.MIRCOSOFT_APP_PASSWORD}

	er := sskype.Authenticate(skype.AUTHENTICATE_URL)
	if er != nil {
		fmt.Println("Error at authenticate")
		panic(er.Error())
	}
	conversation := skype.Conversation{
		Bot:     skype.ChannelAccount{ID: activity.From.ID, Name: activity.From.Name},
		IsGroup: false,
		Members: []skype.ChannelAccount{
			skype.ChannelAccount{ID: skype.DMP_GROUP_ID, Name: skype.DMP_GROUP_NAME}},
		TopicName: "Start conversation"}
	fmt.Println(conversation)
	sskype.BeginConversation(conversation, activity.Text)

}
