package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

func get() {
	// 접속
	resp, err := soup.Get("https://mltd.matsurihi.me/election/")
	if err != nil {
		fmt.Println("Error.")
		os.Exit(1)
	}

	total := append(idolTable, roleTable...)
	totalReplacer := strings.NewReplacer(total...)

	resp = totalReplacer.Replace(resp)

	// 파싱
	doc := soup.HTMLParse(resp)
	start := time.Now()

	// 메인
	main := doc.Find("main")
	themesRaw := main.FindAllStrict("div", "class", "row")

	// 주제 이름들
	themeNames := main.FindAll("h3")

	themes := make([]Theme, 0, 3)

	step := 0

	// 주제 개수만큼
	for i, themeRaw := range themesRaw {

		var theme Theme
		theme.Roles = make([]Role, 0, 5)

		rolesRaw := themeRaw.FindAll("div")

		theme.Name = themeNames[i].Text()

		// 배역 개수만큼
		for x, roleRaw := range rolesRaw {

			// Find 후에 Replace
			theme.Roles = append(theme.Roles, Role{Name: roleRaw.Find("h4").Text()})

			rankList := roleRaw.Find("table").Find("tbody").FindAll("tr")

			// 순위표 안에서 range
			for _, line := range rankList {
				step++
				idolRaw := line.FindAll("td")
				var idol Idol

				rankStr := idolRaw[0].Text()
				idol.Rank, err = strconv.Atoi(rankStr[:len(rankStr)-4])

				idol.Name = idolRaw[1].Text()
				idol.VoteAmount, err = strconv.Atoi(idolRaw[2].Text())

				if err != nil {
					fmt.Println(err)
				}

				theme.Roles[x].Idols = append(theme.Roles[x].Idols, idol)

			}

		}

		themes = append(themes, theme)
	}

	fmt.Println(time.Now().Sub(start))

	fmt.Println("Step:", step)

	js, _ := json.MarshalIndent(&themes[0], "", "    ")
	fmt.Println(string(js))

}
