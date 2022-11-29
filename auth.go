package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/KarpelesLab/rest"
)

const clientID = "oaap-wyslxk-7aqv-gbva-57y7-einuc2a4"

type authInfo struct {
	token    *rest.Token
	name     string
	filepath string
}

func checkLogin() (*authInfo, error) {
	// let's check if we have a token stored in config
	auth := &authInfo{
		name: "default",
	}
	if v := os.Getenv("SHELLS_PROFILE"); v != "" {
		auth.name = v
	}

	if err := auth.init(); err != nil {
		return nil, err
	}
	if err := auth.load(); err != nil {
		// attempt to do auth
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("no login information found, logging in...")
		} else {
			log.Printf("failed to load login (%s), logging in...", err)
		}
		err = auth.login()
		if err != nil {
			return nil, err
		}
		err = auth.save()
		if err != nil {
			return nil, err
		}
	}

	return auth, nil
}

func (auth *authInfo) init() error {
	cnf, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to locate conf dir: %w", err)
	}

	cnf = filepath.Join(cnf, "shells-cli")
	os.MkdirAll(cnf, 0700) // make sure dir exists
	cnf = filepath.Join(cnf, "auth-"+auth.name+".json")
	auth.filepath = cnf
	return nil
}

func (auth *authInfo) load() error {
	// no error, file exists. Load it
	data, err := os.ReadFile(auth.filepath)
	if err != nil {
		return fmt.Errorf("failed to read auth: %w", err)
	}
	if err := json.Unmarshal(data, &auth.token); err != nil {
		return err
	}
	auth.token.ClientID = clientID
	return nil
}

func (auth *authInfo) save() error {
	if auth.token == nil {
		return os.ErrNotExist
	}

	// save token
	data, err := json.Marshal(auth.token)
	if err != nil {
		return err
	}

	return os.WriteFile(auth.filepath, data, 0600)
}

func (auth *authInfo) login() error {
	// prepare to login
	// we need a realtime token
	var res map[string]interface{}
	err := rest.Apply(context.Background(), "OAuth2/App/"+clientID+":token_create", "POST", map[string]interface{}{}, &res)
	if err != nil {
		return err
	}
	tok, ok := res["polltoken"].(string)
	if !ok {
		return fmt.Errorf("failed to fetch polltoken")
	}

	// see: https://www.shells.com/.well-known/openid-configuration?pretty
	tokuri := url.QueryEscape("polltoken:" + tok)
	fulluri := fmt.Sprintf("https://www.shells.com/_rest/OAuth2:auth?response_type=code&client_id=%s&redirect_uri=%s&scope=profile", clientID, tokuri)

	if u, ok := res["xox"].(string); ok {
		fulluri = u
	}

	log.Printf("Please open this URL in order to access shells:\n%s", fulluri)

	// wait for login to complete
	for {
		var res map[string]interface{}
		err := rest.Apply(context.Background(), "OAuth2/App/"+clientID+":token_poll", "POST", map[string]interface{}{"polltoken": tok}, &res)
		if err != nil {
			return err
		}

		v, ok := res["response"]
		if !ok {
			time.Sleep(time.Second) // just in case
			continue
		}

		resp, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid response from api, response of invalid type")
		}

		code, ok := resp["code"].(string)
		if !ok {
			return fmt.Errorf("invalid response from api, response not containing code")
		}

		log.Printf("fetching auth token...")

		// https://www.shells.com/_special/rest/OAuth2:token
		httpresp, err := http.PostForm("https://www.shells.com/_special/rest/OAuth2:token", url.Values{"client_id": {clientID}, "grant_type": {"authorization_code"}, "code": {code}})
		if err != nil {
			return fmt.Errorf("while fetching token: %w", err)
		}
		defer httpresp.Body.Close()

		if httpresp.StatusCode != 200 {
			return fmt.Errorf("invalid status code from server: %s", httpresp.Status)
		}

		body, err := io.ReadAll(httpresp.Body)
		if err != nil {
			return fmt.Errorf("while reading token: %w", err)
		}

		// decode token
		err = json.Unmarshal(body, &auth.token)
		if err != nil {
			return fmt.Errorf("while decoding token: %w", err)
		}
		auth.token.ClientID = clientID

		return nil
	}
}

func (auth *authInfo) Apply(ctx context.Context, p, m string, arg map[string]interface{}, target interface{}) error {
	err := rest.Apply(auth.token.Use(ctx), p, m, arg, target)
	if err != nil {
		return err
	}
	auth.save() // perform save just in case token was updated
	return nil
}
