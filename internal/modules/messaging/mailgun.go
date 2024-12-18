package messaging

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	ustrings "github.com/xinchuantw/hoki-tabloid-backend/utils/strings"
)

type MailgunSendEvent struct {
	ID      *string `json:"id"`
	Message *string `json:"message"`
}

type MailgunClient struct {
	BaseURL string
	APIKey  string
}

func NewMailgunClient(baseURL string, apiKey string) *MailgunClient {
	return &MailgunClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

func (c *MailgunClient) SendVerificationEmail(to string, name string, verificationLink string, lang string) (*MailgunSendEvent, error) {
	f := url.Values{}
	f.Set("from", "HOKI tabloid Mobile App <app@hokishoptaiwan.com>")
	f.Set("to", fmt.Sprintf("%s <%s>", name, to))
	f.Set("template", fmt.Sprintf("verification_link_%s", ustrings.NormalizeLangForTemplate(lang)))

	v := struct {
		Name string `json:"name"`
		Link string `json:"link"`
	}{
		Name: name,
		Link: verificationLink,
	}
	vars, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	f.Set("h:X-Mailgun-Variables", string(vars))

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "messages")
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(f.Encode()))
	if err != nil {
		return nil, err
	}
	req.Close = true

	cred := fmt.Sprintf("api:%s", c.APIKey)
	cred64 := base64.StdEncoding.EncodeToString([]byte(cred))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", cred64))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var evt MailgunSendEvent
	err = json.NewDecoder(res.Body).Decode(&evt)
	if err != nil {
		return nil, err
	}
	return &evt, nil
}

func (c *MailgunClient) SendLoginCode(to string, name string, code string, lang string) (*MailgunSendEvent, error) {
	f := url.Values{}
	f.Set("from", "HOKI tabloid Mobile App <app@hokishoptaiwan.com>")
	f.Set("to", fmt.Sprintf("%s <%s>", name, to))
	f.Set("template", fmt.Sprintf("verification_code_%s", ustrings.NormalizeLangForTemplate(lang)))

	v := struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}{
		Name: name,
		Code: code,
	}
	vars, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	f.Set("h:X-Mailgun-Variables", string(vars))

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "messages")
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(f.Encode()))
	if err != nil {
		return nil, err
	}
	req.Close = true

	cred := fmt.Sprintf("api:%s", c.APIKey)
	cred64 := base64.StdEncoding.EncodeToString([]byte(cred))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", cred64))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var evt MailgunSendEvent
	err = json.NewDecoder(res.Body).Decode(&evt)
	if err != nil {
		return nil, err
	}
	return &evt, nil
}
