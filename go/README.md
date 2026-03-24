# EduPool OAuth2 Client in Golang

This repository contains a Golang program that demonstrates how to implement a basic OpenID Connect client client using the `coreos/go-oidc/v3` library, specifically for the EduPool Identity Provider. The client authenticates users via the EduPool OIDC server, fetches their profile information, and displays it in the console.

## Usage

**Configure the script:**

Open the script in a text editor and replace `yourClientID` and `yourClientSecret` with your actual client ID and secret provided by EduPool.

**Start the server:**

`go run main.go`

The server will be accessible at [http://127.0.0.1:9010](http://127.0.0.1:9010/).

1.  Open the URL [http://127.0.0.1:9010](http://127.0.0.1:9010/) in your browser.
2.  Click the "Log In" link.
3.  You will be redirected to the EduPool OIDC server. Log in with your credentials.
4.  After successful authentication, your user information will be displayed in the console.

## Optional Hub Shortcuts

You can optionally pass `land` or `context` during login initiation. If omitted, the default login behavior stays unchanged.

- `http://127.0.0.1:9010/login` → default behavior
- `http://127.0.0.1:9010/login?land=nrw` → preselect land in Hub
- `http://127.0.0.1:9010/login?context=hh/HH` → jump directly to a specific location login in Hub

## Components

This script uses the following libraries:

-   [coreos/go-oidc/v3](https://github.com/coreos/go-oidc): A Go implementation of OpenID Connect (OIDC) for clients.
-   [golang.org/x/oauth2](https://github.com/golang/oauth2): A Go library for OAuth2.
