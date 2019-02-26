package skype

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SkypeService struct {
	Client        *http.Client
	AppID         string
	AppPassword   string
	Authorization string
	Headers       map[string]string
	Messages
}

type Messages struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AuthenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	ExExpires   int    `json:"ext_expires_in"`
	AccessToken string `json:"access_token"`
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

func (SSkype *SkypeService) InitClient(id string, password string) {
	client := &http.Client{}
	SSkype.Client = client
	SSkype.AppID = id
	SSkype.AppPassword = password
}

func (SSkype *SkypeService) Authenticate(url string) error {
	var headers = map[string]string{
		"host":         "login.microsoftonline.com",
		"content-type": "application/x-www-form-urlencoded"}

	SSkype.setHeader(headers)
	requeststr := fmt.Sprintf("client_id=%s&scope=https://api.botframework.com/.default&grant_type=client_credentials&client_secret=%s")
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

func (SSkype *SkypeService) SendMessage(message Messages) {

}