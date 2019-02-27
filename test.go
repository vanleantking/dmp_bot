package main

import (
	// "github.com/andrewdruzhinin/go-skype/skype"
	"fmt"

	"./skype"
	"./utils"
)

func main() {
	SSkype := &skype.SkypeService{}
	SSkype.InitClient(utils.MIRCOSOFT_APP_ID, utils.MIRCOSOFT_APP_PASSWORD)
	er := SSkype.Authenticate(utils.AUTHENTICATE_URL)
	if er != nil {
		fmt.Println("Error: Can not login to bot service", er.Error())
		panic(er.Error())
	}
	fmt.Println(SSkype.Authorization)
	conver := skype.Conversation{
		Bot:     skype.ChannelAccount{ID: "ea228964-bb55-42c7-9bcc-ea90bbe53b92", Name: "ureka_dmp_bot"},
		IsGroup: false,
		Members: []skype.ChannelAccount{
			skype.ChannelAccount{ID: "ea228964-bb55-42c7-9bcc-ea90bbe53b92", Name: "ureka_dmp_bot"}},
		TopicName: "Start conversation"}
	converresp, err := SSkype.StartConversation(utils.MESSAGE_TRANSFER_URL+utils.START_CONVERSATION, conver)
	fmt.Println(converresp, err)

	message := skype.Messages{Type: "non-reply", Text: "hello world"}
	SSkype.SendMessage(message)
}
