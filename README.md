asap
====

![test](https://github.com/jasonfriedland/asap/actions/workflows/actions.yml/badge.svg)

A simple client package for [ASAP](https://s2sauth.bitbucket.io/) authentication.

Installation
------------

```shell
go get github.com/jasonfriedland/asap
```

Environment Variables
---------------------

```shell
ASAP_PRIVATE_KEY=data:application/pkcs8;kid=webapp%2Fabc123;base64,...
ASAP_ISSUER=services/webapp
ASAP_AUDIENCE=webapp,webapp-service
```

Usage
-----

Ensure the relevant environment variables are set. Then:

```golang
import "github.com/jasonfriedland/asap"

client, _ := asap.NewClient()
token, _ := client.AuthToken()

fmt.Printf("Bearer %s", token)
```
