asap
====

[![Build Status](https://travis-ci.com/jasonfriedland/asap.svg?branch=main)](https://travis-ci.com/jasonfriedland/asap)

A simple client package for [ASAP](https://s2sauth.bitbucket.io/) authentication.

Installation
------------

    go get github.com/jasonfriedland/asap

Environment Variables
---------------------

    ASAP_PRIVATE_KEY=data:application/pkcs8;kid=webapp%2Fabc123;base64,...
    ASAP_ISSUER=services/webapp
    ASAP_AUDIENCE=webapp,webapp-service

Usage
-----

Ensure the relevant environment variables are set. Then:

```go
import "github.com/jasonfriedland/asap"

client, _ := asap.NewClient()
token, _ := client.AuthToken()

fmt.Printf("Bearer %s", token)
```
