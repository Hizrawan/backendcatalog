package testlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/app"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/models"
	"github.com/xinchuantw/hoki-tabloid-backend/utils/database"
	"gopkg.in/guregu/null.v4"
)

// TestRequestContext contains information required to create a test request,
// including base URL, signing keys for access tokens, and handle to the
// database client.
type TestRequestContext struct {
	BaseURL    string
	SigningKey jwk.RSAPrivateKey
	VerifyKey  jwk.RSAPublicKey
	DB         database.Queryer
}

// NewTestRequestContext will return a TestRequestContext instance to be used
// during testing to speed up tasks related to test request creation
func NewTestRequestContext(app *app.Registry) *TestRequestContext {
	return &TestRequestContext{
		BaseURL:    app.Config.AppURL,
		SigningKey: app.SigningKey,
		VerifyKey:  app.VerifyKey,
		DB:         app.DB,
	}
}

// AuthenticateRequest will add the Authorization header containing bearer token
// for the provided models.Customer instance. It will then return the request or
// return any errors that occurred.
func (t TestRequestContext) AuthenticateRequest(r *http.Request, user models.JWTAuthenticatable) (*http.Request, error) {
	err := AuthenticateRequest(t.DB, r, user, t.SigningKey)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// URL will return the absolute URL given a path according to this request
// context.
func (t TestRequestContext) URL(path string) string {
	return URL(t.BaseURL, path)
}

// NewEmptyRequest will return a *http.Request instance for the given method
// and path with empty payload. The path will be converted to absolute URL
// using the context's base URL. Content-Type header will be set to
// application/json
func (t TestRequestContext) NewEmptyRequest(method string, path string) (*http.Request, error) {
	r, err := http.NewRequest(method, t.URL(path), nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// NewJSONRequest will return a *http.Request instance for the given method,
// path, and payload. The path will be converted to absolute URL using the
// context's base URL, while the payload will be JSON encoded with proper
// Content-Type header.
func (t TestRequestContext) NewJSONRequest(method string, path string, payload any) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		p, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(p)
	}

	r, err := http.NewRequest(method, t.URL(path), body)
	if err != nil {
		return nil, err
	}

	if payload != nil {
		r.Header.Set("Content-Type", "application/json")
	}
	return r, nil
}

// NewMultipartRequest will return a *http.Request instance for the given method,
// path, and payload function. The path will be converted to absolute URL using
// the context's base URL. The provided payload function will be executed and
// written to the request body. The correct Content-Type header will be added to
// the request.
func (t TestRequestContext) NewMultipartRequest(method string, path string, fn func(w *multipart.Writer)) (*http.Request, error) {
	body := bytes.Buffer{}
	w := multipart.NewWriter(&body)
	fn(w)
	if err := w.Close(); err != nil {
		return nil, err
	}

	r, err := http.NewRequest(method, t.URL(path), &body)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", w.FormDataContentType())
	return r, nil
}

func AuthenticateRequest(db database.Queryer, req *http.Request, user models.JWTAuthenticatable, sk jwk.RSAPrivateKey) error {
	lifetime := 1 * time.Hour
	token, err := user.IssueAccessToken(time.Now().Add(lifetime))
	if err != nil {
		return err
	}
	tokenStr, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, sk))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+string(tokenStr))

	switch u := user.(type) {
	case *models.Admin:
		t := models.AdminAccessToken{
			Model: models.Model{
				ID: token.JwtID(),
			},
			AdminID:   u.ID,
			ExpiredAt: null.TimeFrom(time.Now().Add(lifetime)),
		}
		err = t.Insert(db)
	case *models.System:
		t := models.SystemAccessToken{
			Model: models.Model{
				ID: token.JwtID(),
			},
			SystemID:  u.ID,
			ExpiredAt: null.TimeFrom(time.Now().Add(lifetime)),
		}
		err = t.Insert(db)

	default:
		return fmt.Errorf("AuthenticateRequest expects user to be *models.Admin, *model.System, ")
	}

	return err
}

// URL will return an absolute URL based on the provided base URL and path
func URL(baseUrl string, path string) string {
	parsed, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		panic(err)
	}
	parsed.Path = path

	return parsed.String()
}

// GetResponseBody will return the body of a *http.Response as a []byte
func GetResponseBody(resp *http.Response) ([]byte, error) {
	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		return io.ReadAll(resp.Body)
	}
	return nil, nil
}

// MapFromResponseBody will return the body of a *http.Response or a
// *httptest.ResponseRecorder as a map, or an error if any occurs
func MapFromResponseBody(resp any) (map[string]any, error) {
	var r *http.Response
	switch v := resp.(type) {
	case *http.Response:
		r = v
	case *httptest.ResponseRecorder:
		r = v.Result()
	default:
		return nil, errors.New("MapFromResponseBody expects *http.Response or *httptest.ResponseRecorder")
	}

	body, err := GetResponseBody(r)
	if err != nil || body == nil {
		return nil, err
	}

	var m map[string]any
	err = json.Unmarshal(body, &m)
	return m, err
}

func (t TestRequestContext) NewQueryParam(req *http.Request, q map[string]any) (*http.Request, error) {
	return NewQueryParam(req, q)
}

// NewQueryParam will add query param to request for the given query param.
// It will then return the request or return any errors that occurs
func NewQueryParam(req *http.Request, query map[string]any) (*http.Request, error) {
	q := url.Values{}

	for key, el := range query {
		var val string
		switch el := el.(type) {
		case string:
			val = el
		case int, int8, int16, int32, int64:
			val = strconv.FormatInt(el.(int64), 10)
		case float32, float64:
			val = strconv.FormatFloat(el.(float64), 'f', -1, 64)
		case bool:
			if el {
				val = "true"
			} else {
				val = "false"
			}
		default:
			b, err := json.Marshal(el)
			if err != nil {
				return nil, err
			}
			val = string(b)
		}

		q.Add(key, val)
	}

	req.URL.RawQuery = q.Encode()
	return req, nil
}
