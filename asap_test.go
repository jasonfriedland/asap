package asap

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"os"
	"testing"
)

// Key fixtures, generate RSA key + Base64-ecoded PKCS8
var (
	validKey, _      = rsa.GenerateKey(rand.Reader, 2048) // minimum size
	shortKey, _      = rsa.GenerateKey(rand.Reader, 32)   //nolint:gosec
	pkcs8ValidKey, _ = x509.MarshalPKCS8PrivateKey(validKey)
	b64ValidKey      = base64.StdEncoding.EncodeToString(pkcs8ValidKey)
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		envPairs map[string]string
		wantErr  bool
	}{
		{
			"Missing env case, error",
			nil,
			true,
		},
		{
			"Valid key but missing kid, error",
			map[string]string{
				"ASAP_PRIVATE_KEY": "data:application/pkcs8;base64," + b64ValidKey,
			},
			true,
		},
		{
			"Borked key, error",
			map[string]string{
				"ASAP_PRIVATE_KEY": "data:application/pkcs8;kid=service%2Ftest;base64,lolnotakey",
			},
			true,
		},
		{
			"Valid key, no error",
			map[string]string{
				"ASAP_PRIVATE_KEY": "data:application/pkcs8;kid=service%2Ftest;base64," + b64ValidKey,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.envPairs {
				os.Setenv(k, v)
			}
			if _, err := NewClient(); (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_AuthToken(t *testing.T) {
	type fields struct {
		Kid        string
		Issuer     string
		Audience   []string
		PrivateKey *rsa.PrivateKey
		Expiry     uint64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Empty case, error",
			fields{},
			true,
		},
		{
			"Empty fields, with valid key, no error",
			fields{
				PrivateKey: validKey,
			},
			false,
		},
		{
			"Simple case, no error",
			fields{
				Kid:        "hello/123",
				PrivateKey: validKey,
				Issuer:     "my/service",
				Audience:   []string{"your/service"},
				Expiry:     60,
			},
			false,
		},
		{
			"Short key case, error",
			fields{
				Kid:        "hello/123",
				PrivateKey: shortKey,
				Issuer:     "my/service",
				Audience:   []string{"your/service"},
				Expiry:     60,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Kid:        tt.fields.Kid,
				Issuer:     tt.fields.Issuer,
				Audience:   tt.fields.Audience,
				PrivateKey: tt.fields.PrivateKey,
				Expiry:     tt.fields.Expiry,
			}
			if _, err := c.AuthToken(); (err != nil) != tt.wantErr {
				t.Errorf("Client.AuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
