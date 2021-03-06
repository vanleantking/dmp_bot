package main

import (
	"encoding/json"
	"fmt"

	"./skype"
)

var (
	bot      *skype.Bot
	contacts []skype.Contact
)

func main() {
	bot = skype.InitBot(":9443", skype.MIRCOSOFT_APP_ID, skype.MIRCOSOFT_APP_PASSWORD)

	er := bot.ConnectorAPI.Authenticate(skype.AUTHENTICATE_URL)
	if er != nil {
		fmt.Println("Error at authenticate")
		panic(er.Error())
	}
	contacts, _ = skype.LoadContacts("./skype/contacts.json")

	srv := bot.SetupEndpointHandler(handleActivity)
	panic(srv.ListenAndServe())
}

func handleActivity(activity *skype.Activity) {
	// Currently only support on skype chat, not web chat
	if activity.ServiceURL != skype.WebchatURL {
		switch activity.Type {
		case skype.Message:
			bot.ConnectorAPI.SActivity = activity

			// check for personal user
			if isExist, _ := skype.CheckContactExist(contacts, activity.From.ID); !isExist {
				contact := skype.Contact{
					Id:         activity.From.ID,
					Name:       activity.From.Name,
					IsGroup:    false,
					ServiceURL: activity.ServiceURL}
				contacts = append(contacts, contact)
				skype.SaveContacts("./skype/contacts.json", contacts)
			}
			// check for group (if message send from user in group)
			if activity.From.ID != activity.Conversation.ID {
				if isExist, _ := skype.CheckContactExist(contacts, activity.Conversation.ID); !isExist {
					contact := skype.Contact{
						Id:         activity.Conversation.ID,
						Name:       "Group",
						IsGroup:    true,
						ServiceURL: activity.ServiceURL}
					contacts = append(contacts, contact)
					skype.SaveContacts("./skype/contacts.json", contacts)
				}
			}
			bot.ConnectorAPI.SendActivity(skype.MessageReply)
			fmt.Println("Successfully sent response message to skype user: " + activity.From.Name)

		case skype.Event:
			fmt.Println("Enter activity event type, ", activity)
			conversation := skype.Conversation{
				Bot:     skype.ChannelAccount{ID: activity.From.ID, Name: activity.From.Name},
				IsGroup: false,
				Members: []skype.ChannelAccount{
					skype.ChannelAccount{ID: skype.DMP_GROUP_ID, Name: skype.DMP_GROUP_NAME}},
				TopicName: "Start conversation"}
			fmt.Println(conversation)
			bot.ConnectorAPI.BeginConversation(conversation, activity.Text)

		case skype.ConversationUpdate:
			bot.ConnectorAPI.SActivity = activity
			// in case remove
			if len(activity.MembersRemoved) > 0 {
				// bot was removed from contact list
				if isExist, _ := skype.CheckInChannelList(activity.MembersRemoved, "28:632cfdd8-e2db-43ef-9653-8205729a10f9"); isExist {
					if exist, index := skype.CheckContactExist(contacts, "28:632cfdd8-e2db-43ef-9653-8205729a10f9"); exist {
						contacts = append(contacts[:index], contacts[index+1:]...)
						skype.SaveContacts("./skype/contacts.json", contacts)
					}
				}
			}
			// in case member was added
			if len(activity.MembersAdded) > 0 {
				// bot was add to contact list
				if isExist, _ := skype.CheckInChannelList(activity.MembersAdded, "28:632cfdd8-e2db-43ef-9653-8205729a10f9"); isExist {
					if exist, _ := skype.CheckContactExist(contacts, "28:632cfdd8-e2db-43ef-9653-8205729a10f9"); !exist {
						contact := skype.Contact{
							Id:         activity.Conversation.ID,
							Name:       "Group",
							IsGroup:    true,
							ServiceURL: activity.ServiceURL}
						contacts = append(contacts, contact)
						skype.SaveContacts("./skype/contacts.json", contacts)
						bot.ConnectorAPI.SendActivity(skype.MessageThankyou)
					}
				} else {
					bot.ConnectorAPI.SendActivity(skype.MessageWelcome)
				}

			}
			bytes, _ := json.MarshalIndent(activity, "", "  ")
			fmt.Println(string(bytes))
			fmt.Println("conversation update")
		}
	}
}
