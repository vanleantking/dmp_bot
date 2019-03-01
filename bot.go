package main

import (
	"fmt"

	"./skype"
	"./utils"
)

var bot *skype.Bot

func main() {
	bot = skype.InitBot(":9443", utils.MIRCOSOFT_APP_ID, utils.MIRCOSOFT_APP_PASSWORD)

	er := bot.ConnectorAPI.Authenticate(utils.AUTHENTICATE_URL)
	if er != nil {
		fmt.Println("Error at authenticate")
		panic(er.Error())
	}

	srv := bot.SetupEndpointHandler(handleActivity)
	panic(srv.ListenAndServe())
}

func handleActivity(activity *skype.Activity) {
	bot.ConnectorAPI.SActivity = activity
	if activity.Type == "message" {
		// // hard coding an auth token is no good practice! I am just doing this to make this example more simple.
		// if err := skype.SendReplyMessage(activity, "Good evening. Nice to meet you!", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ii1zeE1KTUxDSURXTVRQdlp5SjZ0eC1DRHh3MCIsImtpZCI6Ii1zeE1KTUxDSURXTVRQdlp5SjZ0eC1DRHh3MCJ9.eyJhdWQiOiJodHRwczovL2FwaS5ib3RmcmFtZXdvcmsuY29tIiwiaXNzIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvZDZkNDk0MjAtZjM5Yi00ZGY3LWExZGMtZDU5YTkzNTg3MWRiLyIsImlhdCI6MTU1MTMzOTA2NywibmJmIjoxNTUxMzM5MDY3LCJleHAiOjE1NTEzNDI5NjcsImFpbyI6IjQySmdZT2prZXYxWmpIMUtMQ2UzRURlMzFBTWZBQT09IiwiYXBwaWQiOiI2MzJjZmRkOC1lMmRiLTQzZWYtOTY1My04MjA1NzI5YTEwZjkiLCJhcHBpZGFjciI6IjEiLCJpZHAiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9kNmQ0OTQyMC1mMzliLTRkZjctYTFkYy1kNTlhOTM1ODcxZGIvIiwidGlkIjoiZDZkNDk0MjAtZjM5Yi00ZGY3LWExZGMtZDU5YTkzNTg3MWRiIiwidXRpIjoiaEx0X3owek1WMENkSWE0Mm5KY2tBQSIsInZlciI6IjEuMCJ9.dgcMSnJ87cVhy7F_0BIjhnEWPeiyE_OfTChNcBHcV_wmerebFfZR4AJL18gBoDLMcTI-FtyxI6yChMMxQzQXZuF3iF7RtWtweA_FyDe3RvM8UfgKRjik0w3g0Hp9EF5KhXjyiChT-SzwzoZZORnq0NmrAgEurU1i6e757rXqQDzZeMwJsdpII43TPx06th5NaeviBL63Y1UHD21qZvzgqPzR26oGQnplkdwPZWxwaCfR8cT4ceOzPg69bW8mKvalXGRP_etTZyZqb8J3CXR3K1_RFj_gkqjNuxH0LfmIuoV4MhZLV7a7vNE6n2Kl9VDFbVNIEdV5BuNZC5uBRuTnAQ"); err != nil {
		// 	panic(err.Error())
		// } else {
		// 	fmt.Println("Successfully sent response message to skype user: " + activity.From.Name)
		// }
		conversation := skype.Conversation{
			Bot:     skype.ChannelAccount{ID: "28:632cfdd8-e2db-43ef-9653-8205729a10f9", Name: "ureka_dmp_bot"},
			IsGroup: false,
			Members: []skype.ChannelAccount{
				skype.ChannelAccount{ID: "29:18Rd5oZglmOLXNTG_l90wjTcDqrj79RZbhbSNG3XleU4", Name: "van le"}},
			TopicName: "Start conversation"}
		bot.ConnectorAPI.BeginConversation(conversation, "tin reply")
	}
}
