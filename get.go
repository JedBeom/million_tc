package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JedBeom/soup"
)

func get() {
	// 접속
	resp, err := soup.Get("https://mltd.matsurihi.me/election/")
	if err != nil {
		fmt.Println("Error.")
		os.Exit(1)
	}

	// 파싱
	doc := soup.HTMLParse(resp)
	// 메인
	main := doc.Find("main")
	themesRaw := main.FindAllStrict("div", "class", "row")

	// 주제 이름들
	themeNames := main.FindAll("h3")

	themes := make([]Theme, 0, 3)

	// 아이돌 이름 번역표
	idolReplacer := strings.NewReplacer(idolTable...)
	// 배역 이름 변역표
	roleReplacer := strings.NewReplacer(roleTable...)

	step := 0

	start := time.Now()
	// 주제 개수만큼
	for i, themeRaw := range themesRaw {

		var theme Theme
		theme.Roles = make([]Role, 0, 5)

		rolesRaw := themeRaw.FindAll("div")

		theme.Name = roleReplacer.Replace(themeNames[i].Text())

		// 배역 개수만큼
		for x, roleRaw := range rolesRaw {

			// Find 후에 Replace
			theme.Roles = append(theme.Roles, Role{Name: roleReplacer.Replace(roleRaw.Find("h4").Text())})

			rankList := roleRaw.Find("table").Find("tbody").FindAll("tr")

			// 순위표 안에서 range
			for _, line := range rankList {
				step++
				idolRaw := line.FindAll("td")
				var idol Idol

				rankStr := idolRaw[0].Text()
				idol.Rank, err = strconv.Atoi(rankStr[:len(rankStr)-4])

				idol.Name = idolReplacer.Replace(idolRaw[1].Text())
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
