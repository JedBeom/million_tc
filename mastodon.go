package main

import (
	"fmt"
	"time"

	"github.com/McKael/madon"
)

func mast() {
	userToken := madon.UserToken{
		AccessToken: "",
		CreatedAt:   time.Now().UnixNano(),
		Scope:       "read write",
		TokenType:   "urn:ietf:wg:oauth:2.0:oob",
	}

	mc, err := madon.RestoreApp("MillionTC", "uri.life", "", "", &userToken)

	if err != nil {
		fmt.Println(err)
	}

	status := madon.PostStatusParams{
		Text:       "Test",
		Visibility: "unlisted",
	}
	st, err := mc.PostStatus(status)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(st)
}
