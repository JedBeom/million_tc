package main

import (
	"context"
	"log"

	mastodon "github.com/JedBeom/go-mastodon"
)

func main() {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       "https://planet.moe",
		ClientID:     "e009ba56d0d4eff4454cb91826825f12b17c7aa7e0e340615689cb8792099ecc",
		ClientSecret: "48f035c2f345f7fa0e693f1565e63525d7c1226b95ae44b0cc52f37ec7e6c2ee",
	})
	err := c.AuthenticateToken(context.Background(), "4ba503602045522f9f4311e4698fa3fe0b12000816eee950d8ddd0ffbcd3171a", "urn:ietf:wg:oauth:2.0:oob")
	if err != nil {
		log.Fatal(err)
	}
}
