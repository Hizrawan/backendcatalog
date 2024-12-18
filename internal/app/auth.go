package app

import (
	"github.com/xinchuantw/hoki-tabloid-backend/internal/config"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/modules/authentication"
	mobilebe "github.com/xinchuantw/hoki-tabloid-backend/internal/modules/authentication/mobile_be"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/modules/authentication/xinchuanauth"
)

func NewAuthModule(config config.AuthConfig) authentication.Auth {
	return authentication.Auth{
		XinchuanAuth: NewXinchuanAuthClient(config.XinchuanAuth),
		MobileBEAuth: NewMobileBEClient(config.MobileBEAuth),
	}
}

func NewXinchuanAuthClient(config config.XinchuanAuthConfig) *xinchuanauth.Client {
	baseUrl := config.BaseURL
	if baseUrl == "" {
		baseUrl = "https://auth.xinchuan.tw/"
	}

	return xinchuanauth.NewClient(xinchuanauth.Options{
		BaseURL:      baseUrl,
		ClientID:     config.ID,
		ClientSecret: config.Secret,
		RedirectURI:  config.Callback,
	})
}

func NewMobileBEClient(config config.MobileBEAuthConfig) *mobilebe.Client {
	baseUrl := config.BaseURL
	if baseUrl == "" {
		baseUrl = "base url mobile be"
	}

	return mobilebe.NewClient(mobilebe.Options{
		BaseURL: config.BaseURL,
		Secret:  config.Secret,
	})
}
