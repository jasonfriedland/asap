package asap

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/atlassian/go-asap"
	"github.com/kelseyhightower/envconfig"
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

// asapConfig holds ASAP config from the environment.
type asapConfig struct {
	PrivateKey string `split_words:"true"`
	Issuer     string
	Audience   string
}

// NewClient populates a new Client from the values set in the environment.
func NewClient() (*Client, error) {
	// ASAP config from env
	var conf asapConfig
	err := envconfig.Process("asap", &conf)
	fmt.Println(conf)
	if err != nil {
		return nil, err
	}
	// Key from environment is a PKCS8 encoded data URL
	dataURL, err := dataurl.DecodeString(conf.PrivateKey)
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
		return nil, fmt.Errorf("kid not found in dataURL")
	}
	// Return a new asapConfig
	return &Client{
		Kid:        kid,
		Issuer:     conf.Issuer,
		Audience:   strings.Split(conf.Audience, ","),
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
