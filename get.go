package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/JedBeom/soup"
)

func get() {
	resp, err := soup.Get("https://mltd.matsurihi.me/election/")
	if err != nil {
		fmt.Println("Error.")
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)
	div := doc.Find("main").FindAllStrict("div", "class", "row")

	firstIdols := []string{}

	for _, value := range div {
		roles := value.FindAll("div")

		for _, role := range roles {
			tr := role.Find("table").Find("tbody").FindAll("tr")
			firstIdols = append(firstIdols, tr[0].FindAll("td")[1].Text())
		}
	}

	r := strings.NewReplacer(idolTable...)

	for _, idol := range firstIdols {
		fmt.Println(r.Replace(idol))
	}

}
