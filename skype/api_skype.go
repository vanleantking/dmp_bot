package skype

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	unexpectedHttpStatusCodeTemplate = "The microsoft servers returned an unexpected http status code: %v"
)

type SkypeService struct {
	Client        *http.Client
	AppID         string
	AppPassword   string
	Authorization string
	Headers       map[string]string
	SActivity     *Activity
}

type AuthenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	ExExpires   int    `json:"ext_expires_in"`
	AccessToken string `json:"access_token"`
}

type Conversation struct {
	Bot       ChannelAccount   `json:"bot"`
	IsGroup   bool             `json:"isGroup"`
	Members   []ChannelAccount `json:"members"`
	TopicName string           `json:"topicName"`
}

type CoversationResp struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Error            string         `json:"error"`
	ErrorDescription string         `json:"error_description"`
	ErrorCodes       map[int]string `json:"error_codes"`
	TimeStamp        string         `json:"timestamp"`
	TraceID          string         `json:"trace_id"`
	CorrelationID    string         `json:"correlation_id"`
}

func (SSkype *SkypeService) setHeader(headers map[string]string) {
	if SSkype.Authorization != "" {
		headers["Authorization"] = SSkype.Authorization
	}
	SSkype.Headers = map[string]string{}
	for key, value := range headers {
		SSkype.Headers[key] = value
	}
}

func (SSkype *SkypeService) InitClient() {
	client := &http.Client{}
	SSkype.Client = client
}

func (SSkype *SkypeService) Authenticate(url string) error {
	var headers = map[string]string{
		"host":         "login.microsoftonline.com",
		"content-type": "application/x-www-form-urlencoded"}

	SSkype.setHeader(headers)
	requeststr := fmt.Sprintf("client_id=%s&scope=https://api.botframework.com/.default&grant_type=client_credentials&client_secret=%s", SSkype.AppID, SSkype.AppPassword)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requeststr)))

	// Set header request
	if len(SSkype.Headers) > 0 {
		for key, value := range SSkype.Headers {
			req.Header.Set(key, value)
		}
	}
	authen := AuthenResponse{}

	resp, err := SSkype.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("string body request: ", string(body))
	if err := json.Unmarshal(body, &authen); err != nil {
		var ErrorResp = ErrorResponse{}
		err = json.Unmarshal(body, &ErrorResp)
		if err == nil {
			return errors.New(ErrorResp.ErrorDescription)
		}
		return err
	}

	SSkype.Authorization = authen.TokenType + " " + authen.AccessToken

	return nil
}

func (SSkype *SkypeService) SendActivity(message string) {
	responseActivity := &Activity{
		Type:         "message",
		From:         SSkype.SActivity.Recipient,
		Conversation: SSkype.SActivity.Conversation,
		Recipient:    SSkype.SActivity.From,
		Text:         message,
		ReplyToID:    SSkype.SActivity.ID,
	}

	replyUrl := fmt.Sprintf(SEND_MESSAGE,
		SSkype.SActivity.ServiceURL,
		SSkype.SActivity.Conversation.ID,
		SSkype.SActivity.ID)

	var headers = map[string]string{
		"content-type": "application/json"}

	SSkype.setHeader(headers)
	activitystr, err := json.Marshal(responseActivity)
	if err != nil {
		panic(err.Error())
	}
	req, err := http.NewRequest("POST", replyUrl, bytes.NewBuffer(activitystr))

	// Set header request
	if len(SSkype.Headers) > 0 {
		for key, value := range SSkype.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := SSkype.Client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("string body request: ", string(body))
}

func (SSkype *SkypeService) requestConversation(conversation Conversation, message string) error {

	var headers = map[string]string{
		"content-type": "application/json"}

	SSkype.setHeader(headers)
	con, err := json.Marshal(conversation)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", MESSAGE_TRANSFER_URL+START_CONVERSATION, bytes.NewBuffer(con))
	if err != nil {
		panic(err.Error())
	}

	// Set header request
	if len(SSkype.Headers) > 0 {
		for key, value := range SSkype.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := SSkype.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var statusCode int = resp.StatusCode
	if statusCode == http.StatusOK || statusCode == http.StatusCreated ||
		statusCode == http.StatusAccepted || statusCode == http.StatusNoContent {
		var conResp CoversationResp
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}
		fmt.Println("string body request: ", string(body))
		if err := json.Unmarshal(body, &conResp); err != nil {
			var ErrorResp = ErrorResponse{}
			err = json.Unmarshal(body, &ErrorResp)
			if err == nil {
				return errors.New(ErrorResp.ErrorDescription)
			}
			return err
		}
		activity := &Activity{
			Type:         "message",
			From:         conversation.Bot,
			Conversation: ConversationAccount{ID: conResp.ID, IsGroup: conversation.IsGroup, Name: conversation.TopicName},
			Recipient:    conversation.Members[0],
			Text:         message,
			ServiceURL:   MESSAGE_TRANSFER_URL}
		SSkype.SActivity = activity
		return nil

	}
	return fmt.Errorf(unexpectedHttpStatusCodeTemplate, statusCode)
}

func (SSkype *SkypeService) BeginConversation(conversation Conversation, message string) error {
	err := SSkype.requestConversation(conversation, message)
	if err != nil {
		return err
	}
	SSkype.SendActivity(message)
	return nil
}

func (SSkype *SkypeService) MakeRequest(activity *Activity) {

	activitystr, err := json.Marshal(activity)
	if err != nil {
		panic(err.Error())
	}

	req, err := http.NewRequest("POST", "http://localhost:9443"+ACTIONHOOK, bytes.NewBuffer(activitystr))
	req.Header.Set("Authorization", SSkype.Authorization)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err.Error())
	}

	resp, err := SSkype.Client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("string body request: ", string(body), resp.StatusCode)
}
