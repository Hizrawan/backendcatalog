package app

import (
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/config"
)

func NewSigningKey(config *config.PrivateConfig) (jwk.RSAPrivateKey, jwk.RSAPublicKey, error) {
	rawKey := config.SigningKey
	private, err := jwk.ParseKey([]byte(rawKey))
	if err != nil {
		return nil, nil, err
	}
	if _, ok := private.(jwk.RSAPrivateKey); !ok {
		return nil, nil, fmt.Errorf("signing key must be an RSA private key")
	}

	public, err := private.PublicKey()
	if err != nil {
		return nil, nil, err
	}

	return private.(jwk.RSAPrivateKey), public.(jwk.RSAPublicKey), nil
}
