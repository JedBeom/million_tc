package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/McKael/madon"
	"github.com/microcosm-cc/bluemonday"
)

var (
	t *template.Template

	events   = make(chan madon.StreamEvent)
	stopChan = make(chan bool)
	doneChan = make(chan bool)

	helpMsg = `사용법: (테마 번호)-(배역 번호)
ex) 외딴 섬 서스펜스 호러의 저택의 여주인 -> 1-4`
)

func init() {
	format := "\n{{ range . }}{{ .Rank }}등 {{ .Name }} {{ .VoteAmount }}표\n{{ end }}"

	t = template.Must(template.New("format").Parse(format))
}

func Run() {

	startedTime := time.Now()

	err := mc.StreamListener("user", "", events, stopChan, doneChan)
	if err != nil {
		fmt.Println(err)
		return
	}

Streamer:
	for {
		event := <-events

		// Catch only notifications
		if event.Event == "notification" {
			noti := event.Data.(madon.Notification)

			// Avoid bot loop
			if noti.Account.Bot == true {
				continue
			}

			// React only mention
			if noti.Type != "mention" {
				continue
			}

			// Remove HTML tags
			contentRaw := bluemonday.StrictPolicy().Sanitize(noti.Status.Content)

			contentArray := strings.Split(contentRaw, " ")
			content := ""

			// Remove @...
			for _, value := range contentArray {
				if len(value) == 0 {
					continue
				}

				if string(value[0]) == "@" {
					continue
				}

				content += value + " "
			}

			// No words, no reaction
			if len(content) != 0 {
				content = content[:len(content)-1]
			}

			// print uptime
			if content == "업타임" {
				uptime := time.Now().Sub(startedTime)

				reply(&noti, "uptime: "+uptime.String(), "")
				continue
			}

			// Split with -
			requestsStr := strings.Split(content, "-")
			if len(requestsStr) != 2 {
				//fmt.Println("len:", len(requestsStr))
				_, err = reply(&noti, helpMsg, "사용법")
				continue
			}

			// Convert to integer
			requestsInt := []int{}
			for _, value := range requestsStr {
				if value == "NaN" {
					reply(&noti, "저 괴롭히지 마세요... 흑흑...", "")
					continue Streamer
				}

				n, err := strconv.Atoi(value)

				if err != nil {
					reply(&noti, helpMsg, "사용법")
					//fmt.Println("value:", value)
					continue Streamer
				}

				requestsInt = append(requestsInt, n)

			}

			// Prevent index out of range
			if requestsInt[0] < 1 || requestsInt[0] > 3 || requestsInt[1] < 1 || requestsInt[1] > 5 {
				reply(&noti, "첫번째 숫자는 1부터 3, 두번째 숫자는 1부터 5만 있어요!", "")
				continue
			}

			// Get themes
			themes, err := get()
			if err != nil {
				reply(&noti, "죄송해요. 크롤러가 잠시 잠자는 중이에요. 나중에 다시 시도 해주세요!", "")
				continue
			}

			//fmt.Println(requestsInt)
			//fmt.Println(len(themes))

			theme := themes[requestsInt[0]-1]
			//fmt.Println(len(theme.Roles))
			role := theme.Roles[requestsInt[1]-1]

			// Template
			var tpl bytes.Buffer

			// Execute
			t.Execute(&tpl, role.Idols)

			cwText := theme.Name + ": " + role.Name

			_, err = reply(&noti, tpl.String(), cwText)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

// Restart streamer if session closed
func restarter() {
	for {
		fmt.Println("Restarter Loop started")
		_, ok := <-doneChan
		if !ok {
			fmt.Println("Restarting...")

			stopChan = make(chan bool)
			doneChan = make(chan bool)

			err := mc.StreamListener("user", "", events, stopChan, doneChan)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		time.Sleep(time.Second)

	}
}

/*
func Tester() {
	time.Sleep(time.Second)
	fmt.Println("Closing stopChan")
	close(stopChan)
}
*/
