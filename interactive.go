package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/McKael/madon"
	"github.com/microcosm-cc/bluemonday"
)

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
			content := " "
			for _, value := range contentArray {
				if len(value) == 0 {
					continue
				}

				if strings.HasPrefix(value, "@") {
					continue
				}

				content += value + " "
			}

			_, err = reply(&noti, content, "봇 테스트 중")
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
