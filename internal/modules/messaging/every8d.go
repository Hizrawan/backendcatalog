package messaging

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/xinchuantw/hoki-tabloid-backend/internal/config"
)

type Every8dClient struct {
	Username    string
	Password    string
	accessToken string
	BaseURL     string
}

type SendEvent struct {
	Id      *string `json:"id"`
	Message *string `json:"message"`
}

type DirectMessageRequest struct {
	MSG  string `json:"MSG"`
	Dest string `json:"DEST"`
}

type AuthResponse struct {
	Result bool   `json:"Result"`
	Status string `json:"Status"`
	Msg    string `json:"Msg"`
}

func NewEvery8dClient(config config.Every8DConfig) *Every8dClient {
	return &Every8dClient{
		Username: config.Username,
		Password: config.Password,
		BaseURL:  config.BaseURL,
	}
}

func (c *Every8dClient) checkToken() error {
	if c.accessToken == "" {
		accessToken, err := c.getToken()
		if err != nil {
			return err
		}
		c.accessToken = *accessToken
		return nil
	}

	payload := map[string]any{
		"HandlerType": 3,
		"VerifyType":  2,
	}
	var body io.Reader
	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body = bytes.NewBuffer(p)

	req, err := http.NewRequest("POST", joinURL(c.BaseURL, "/API21/HTTP/ConnectionHandler.ashx"), body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	result, err := mapFromResponseBody(response)
	if err != nil {
		panic(err)
	}

	if !result["Result"].(bool) {
		accessToken, err := c.getToken()
		if err != nil {
			return err
		}
		c.accessToken = *accessToken
	}

	return nil
}

func (c *Every8dClient) getToken() (*string, error) {
	payload := map[string]any{
		"HandlerType": 3,
		"VerifyType":  1,
		"UID":         c.Username,
		"PWD":         c.Password,
	}
	var body io.Reader
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body = bytes.NewBuffer(p)

	req, err := http.NewRequest("POST", joinURL(c.BaseURL, "/API21/HTTP/ConnectionHandler.ashx"), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	result, err := mapFromResponseBody(response)
	if err != nil {
		return nil, err
	}

	if !result["Result"].(bool) {
		errorMessage := mapErrorCode(result["Status"].(string))
		if errorMessage == "" {
			errorMessage = result["Msg"].(string)
		}
		return nil, errors.New(errorMessage)
	}

	token := result["Msg"].(string)
	return &token, nil
}

func (c *Every8dClient) SendSMS(to string, message string) (*SendEvent, error) {
	err := c.checkToken()
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("MSG", message)
	params.Add("DEST", to)

	req, err := http.NewRequest("POST", joinURL(c.BaseURL, "/API21/HTTP/SendSMS.ashx"), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	result := string(b)
	split := strings.Split(result, ",")
	if len(split) == 2 {
		errorMessage := mapErrorCode(split[0])
		if errorMessage == "" {
			errorMessage = split[1]
		}
		return nil, errors.New(errorMessage)
	}

	return &SendEvent{
		Id: &split[len(split)-1],
	}, nil
}

func mapErrorCode(errorCode string) string {
	switch errorCode {
	case "-300":
		return "Username and password not provided"
	case "-27":
		return "Destination number required"
	case "-5":
		return "Length of SMS is exceed the maximum length."
	case "-4":
		return "SMS has been retried for over 24-hour."
	case "-3":
		return "Invalid mobile number or mobile number is set as black list."
	case "-2":
		return "API account or password error."
	case "-1":
		return "Invalid parameter error."
	case "101":
	case "107":
		return "Reported from Mobile Carrier: Failed to send SMS due to poor signal or mobile is off-line or mobile error."
	case "102":
		return "Reported from Mobile Carrier: Failed to send SMS due to mobile network error or error of bas station."
	case "103":
		return "Reported from Mobile Carrier: Failed due to invalid mobile number."
	case "104":
		return "Reported from Mobile Carrier: The mobile number is in blacklist."
	case "105":
		return "Reported from Mobile Carrier: Failed due to mobile/handset error."
	case "106":
		return "Reported from Mobile Carrier: Unexpected error."
	case "301":
		return "Out of credit。"
	case "500":
		return "Failed to send international message, please check if internal call is permitted."
	}
	return ""
}

func joinURL(baseUrl string, path string) string {
	parsed, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		panic(err)
	}
	parsed.Path = path

	return parsed.String()
}

func mapFromResponseBody(resp any) (map[string]any, error) {
	var r *http.Response
	switch v := resp.(type) {
	case *http.Response:
		r = v
	case *httptest.ResponseRecorder:
		r = v.Result()
	default:
		return nil, errors.New("MapFromResponseBody expects *http.Response or *httptest.ResponseRecorder")
	}

	body, err := getResponseBody(r)
	if err != nil || body == nil {
		return nil, err
	}

	var m map[string]any
	err = json.Unmarshal(body, &m)
	return m, err
}

func getResponseBody(resp *http.Response) ([]byte, error) {
	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		return io.ReadAll(resp.Body)
	}
	return nil, nil
}
