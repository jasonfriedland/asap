package asap

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"bitbucket.org/atlassian/go-asap"
	"github.com/vincent-petithory/dataurl"
)

// Client holds the config that represents an ASAP JWT.
type Client struct {
	Kid        string
	Issuer     string
	Audience   []string
	PrivateKey *rsa.PrivateKey
	Expiry     uint64
}

// NewASAPConfig populates our config from the environment.
func NewClient() (*Client, error) {
	// Key from environment is a PKCS8 encoded data URL
	dataURL, err := dataurl.DecodeString(os.Getenv("ASAP_PRIVATE_KEY"))
	if err != nil {
		return nil, err
	}
	// RSA key portion is the Data
	rsaPrivateKey, err := x509.ParsePKCS8PrivateKey(dataURL.Data)
	if err != nil {
		return nil, err
	}
	// Key ID
	kid, ok := dataURL.Params["kid"]
	if !ok {
		return nil, errors.New("kid not found in dataURL query string")
	}
	// Return a new asapConfig
	return &Client{
		Kid:        kid,
		Issuer:     os.Getenv("ASAP_ISSUER"),
		Audience:   strings.Split(os.Getenv("ASAP_AUDIENCE"), ","),
		PrivateKey: rsaPrivateKey.(*rsa.PrivateKey),
	}, nil
}

// AuthToken generates a unique Bearer token, which must be generated
// per-request.
func (c *Client) AuthToken() (string, error) {
	if c.PrivateKey == nil {
		return "", fmt.Errorf("nil PrivateKey")
	}
	provider := asap.NewMicrosProvisioner(c.Audience, time.Minute)
	token, err := provider.Provision()
	if err != nil {
		return "", err
	}
	headerValue, err := token.Serialize(c.PrivateKey)
	if err != nil {
		return "", err
	}
	return string(headerValue), nil
}
