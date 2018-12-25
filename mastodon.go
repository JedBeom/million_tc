package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/McKael/madon"
)

type Config struct {
	ClientKey    string `json:"client_key"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
}

var (
	config Config
	mc     *madon.Client
)

func init() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	userToken := madon.UserToken{
		AccessToken: config.AccessToken,
		CreatedAt:   time.Now().UnixNano(),
		Scope:       "read write",
		TokenType:   "urn:ietf:wg:oauth:2.0:oob",
	}

	mc, err = madon.RestoreApp("MillionTC", "uri.life", config.ClientKey, config.ClientSecret, &userToken)

	if err != nil {
		panic(err)
	}
}

func toot(content, cw string) (st *madon.Status, err error) {
	status := madon.PostStatusParams{
		Text:       content,
		Visibility: "unlisted",
	}

	status.SpoilerText = cw

	st, err = mc.PostStatus(status)
	return
}

func reply(replyTo int64, content, cw string) (st *madon.Status, err error) {
	status := madon.PostStatusParams{
		Text:        content,
		Visibility:  "unlisted",
		SpoilerText: cw,
		InReplyTo:   replyTo,
	}

	st, err = mc.PostStatus(status)
	return
}
