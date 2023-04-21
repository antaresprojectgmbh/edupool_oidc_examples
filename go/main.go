package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	oauth2Config *oauth2.Config

	// should be a random string per authentication request to prevent CSRF attacks
	// https://auth0.com/docs/secure/attack-protection/state-parameters
	oauthStateString = "p2IHueZueYNDr1gsYrh9yEXtL2u5X8Ii"

	verifier *oidc.IDTokenVerifier
)

func init() {
	// create a new oidc provider.
	// OIDC gets all needed metadata automatically from here: https://oidc.edupool.cloud/.well-known/openid-configuration
	provider, err := oidc.NewProvider(context.Background(), "https://oidc.edupool.cloud/")
	if err != nil {
		panic(err)
	}

	oauth2Config = &oauth2.Config{
		RedirectURL:  "http://127.0.0.1:9010/callback",                                    // the URL where we redirect the user after successfull login
		ClientID:     "yourClientID",                                                      // the client id
		ClientSecret: "yourClientSecret",                                                  // the client secret
		Scopes:       []string{oidc.ScopeOpenID, "offline", "profile", "antares.context"}, // define "scopes" and request only the data you need.
		// Email -> The users Email
		// Profile -> Name, Surname, Fullname and user role (Techer or Student)
		// antares.context -> infos about the customer like schoolid and state
		// we can add more custom scopes if needed
		Endpoint: provider.Endpoint(),
	}

	verifier = provider.Verifier(&oidc.Config{ClientID: oauth2Config.ClientID})
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleOAuth2Callback)
	fmt.Println(http.ListenAndServe(":9010", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	// just a super simple page to initiate the login.
	var htmlIndex = `<html>
<body>
	<a href="/login">Log In</a>
</body>
</html>`

	fmt.Fprintln(w, htmlIndex)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// initiate the login and redirect to the ID-Provieders AuthCode URL
	http.Redirect(w, r, oauth2Config.AuthCodeURL(oauthStateString), http.StatusFound)
}

// this handler (http://127.0.0.1:9010/callback) is called after a successfull login.
// you can check the data and proceed as you like
func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	// Verify state and errors.

	// Exchange converts an authorization code into a token.
	// It is used after a resource provider redirects the user back to the Redirect URI (the URL obtained from AuthCodeURL).
	oauth2Token, err := oauth2Config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		panic(err)
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		panic("missing token")
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		panic(err)
	}

	// Extract custom claims
	var claims any
	if err := idToken.Claims(&claims); err != nil {
		panic(err)
	}

	// print the claims in a readable way
	jc, _ := json.MarshalIndent(claims, "", "  ")
	fmt.Println(string(jc))
}
