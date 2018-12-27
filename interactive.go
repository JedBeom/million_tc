package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/McKael/madon"
	"github.com/microcosm-cc/bluemonday"
)

var helpMsg = `사용법: (테마 번호)-(배역 번호)
ex) 외딴 섬 서스펜스 호러의 저택의 여주인 -> 1-4`

var t *template.Template

func init() {
	format := "{{ range . }}{{ .Rank }}등 {{ .Name }} {{ .VoteAmount }}표\n{{ end }}"

	t = template.Must(template.New("format").Parse(format))
}

/*
func RunOld() {
	themes := get()

	format := "{{ range . }}{{ .Rank }}등 {{ .Name }} {{ .VoteAmount }}표\n{{ end }}"

	t := template.Must(template.New("foramt").Parse(format))

	var tpl bytes.Buffer

	t.Execute(&tpl, themes[0].Roles[3].Idols)

	st, err := toot(tpl.String(), "밀리시타")
	if err != nil {
		panic(err)
	}

	fmt.Println(st)

}
*/

func Run() {

	var (
		events   = make(chan madon.StreamEvent)
		stopChan <-chan bool
		doneChan chan bool
	)

	err := mc.StreamListener("user", "", events, stopChan, doneChan)
	if err != nil {
		fmt.Println(err)
		return
	}

Streamer:
	for {
		event := <-events
		if event.Event == "notification" {
			noti := event.Data.(madon.Notification)

			if noti.Account.Bot == true {
				continue
			}

			if noti.Type != "mention" {
				continue
			}

			contentRaw := bluemonday.StrictPolicy().Sanitize(noti.Status.Content)

			contentArray := strings.Split(contentRaw, " ")
			content := ""
			for _, value := range contentArray {
				if len(value) == 0 {
					continue
				}

				if strings.HasPrefix(value, "@") {
					continue
				}

				content += value + " "
			}

			if len(content) != 0 {
				content = content[:len(content)-1]
			}

			requestsStr := strings.Split(content, "-")
			if len(requestsStr) != 2 {
				//fmt.Println("len:", len(requestsStr))
				_, err = reply(&noti, helpMsg, "사용법")
				continue
			}

			requestsInt := []int{}
			for _, value := range requestsStr {
				n, err := strconv.Atoi(value)

				if err != nil {
					_, err = reply(&noti, helpMsg, "사용법")
					//fmt.Println("value:", value)
					continue Streamer
				}

				requestsInt = append(requestsInt, n)

			}

			if requestsInt[0] < 1 || requestsInt[0] > 3 || requestsInt[1] < 1 || requestsInt[1] > 5 {
				_, err = reply(&noti, "첫번째 숫자는 1부터 3, 두번째 숫자는 1부터 5만 있어요!", "")
			}

			themes, err := get()
			if err != nil {
				_, err = reply(&noti, "죄송해요. 크롤러가 잠시 잠자는 중이에요. 나중에 다시 시도 해주세요!", "")
				continue
			}

			//fmt.Println(requestsInt)
			//fmt.Println(len(themes))

			theme := themes[requestsInt[0]-1]
			//fmt.Println(len(theme.Roles))
			role := theme.Roles[requestsInt[1]-1]

			var tpl bytes.Buffer

			t.Execute(&tpl, role.Idols)

			roleFullName := theme.Name + "-" + role.Name

			ranks := roleFullName + " 역\n" + tpl.String()

			_, err = reply(&noti, ranks, roleFullName)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
