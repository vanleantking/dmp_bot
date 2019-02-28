package skype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type MessageListener interface {
	handle()
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

const (
	unexpectedHttpStatusCodeTemplate = "The microsoft servers returned an unexpected http status code: %v"

	requestTokenUrl      = "https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token"
	replyMessageTemplate = "%vv3/conversations/%v/activities/%v"
)

func RequestAccessToken(microsoftAppId string, microsoftAppPassword string) (TokenResponse, error) {
	var tokenResponse TokenResponse
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", microsoftAppId)
	values.Set("client_secret", microsoftAppPassword)
	values.Set("scope", "https://api.botframework.com/.default")
	if response, err := http.PostForm(requestTokenUrl, values); err != nil {
		return tokenResponse, err
	} else if response.StatusCode == http.StatusOK {
		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&tokenResponse)
		return tokenResponse, err
	} else {
		return tokenResponse, fmt.Errorf(unexpectedHttpStatusCodeTemplate, response.StatusCode)
	}
}

func SendReplyMessage(activity *Activity, message, authorizationToken string) error {
	responseActivity := &Activity{
		Type:         activity.Type,
		From:         activity.Recipient,
		Conversation: activity.Conversation,
		Recipient:    activity.From,
		Text:         message,
		ReplyToID:    activity.ID,
	}
	replyUrl := fmt.Sprintf(replyMessageTemplate, activity.ServiceURL, activity.Conversation.ID, activity.ID)
	fmt.Println("Send reply message, ", replyUrl)
	return SendActivityRequest(responseActivity, replyUrl, authorizationToken)
}

func SendActivityRequest(activity *Activity, replyUrl, authorizationToken string) error {
	client := &http.Client{}
	fmt.Println("Enter here,: ", authorizationHeaderValuePrefix+authorizationToken)
	if jsonEncoded, err := json.Marshal(*activity); err != nil {
		return err
	} else {
		fmt.Println("prepare new request", authorizationToken)
		req, err := http.NewRequest(
			http.MethodPost,
			replyUrl,
			bytes.NewBuffer(*&jsonEncoded),
		)
		if err == nil {
			req.Header.Set("Authorization", authorizationToken)
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(*&req)
			if err == nil {
				defer resp.Body.Close()
				var statusCode int = resp.StatusCode
				if statusCode == http.StatusOK || statusCode == http.StatusCreated ||
					statusCode == http.StatusAccepted || statusCode == http.StatusNoContent {
					return nil
				} else {
					fmt.Println(statusCode)
					panic("Stop here, ")
					return fmt.Errorf(unexpectedHttpStatusCodeTemplate, statusCode)
				}
			} else {
				panic("Stop here, " + err.Error())
				return err
			}
		} else {
			return err
		}
	}
}
