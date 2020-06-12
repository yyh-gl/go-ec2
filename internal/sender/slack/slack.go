package slack

import (
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/yyh-gl/go-ec2/internal/sender"
)

type (
	Sender struct {
		webHook   string
		channel   string
		username  string
		iconEmoji string
	}
)

func NewClient(sender interface{}) sender.Sender {
	// TODO: yaml読み取りデータのバリデーションチェック
	s := sender.(map[interface{}]interface{})

	return &Sender{
		webHook:   s["web_hook"].(string),
		channel:   s["channel"].(string),
		username:  s["username"].(string),
		iconEmoji: s["icon_emoji"].(string),
	}
}

func (s Sender) Send(materials sender.Materials) error {
	if len(materials) == 0 {
		// TODO: message for when no instance
		return nil
	}

	p := slack.Payload{
		Username:    s.username,
		IconUrl:     "",
		IconEmoji:   ":" + s.iconEmoji + ":",
		Channel:     s.channel,
		Text:        "",
		LinkNames:   "true",
		Attachments: createAttachments(materials),
		UnfurlLinks: false,
		UnfurlMedia: false,
		Markdown:    false,
	}

	slack.Send(s.webHook, "", p)
	return nil
}

func createAttachments(materials sender.Materials) []slack.Attachment {
	attachments := make([]slack.Attachment, len(materials))
	for i, m := range materials {
		tmp := make([]string, 3)
		tmp[0] = "ID: " + m.ID
		tmp[1] = "Name: no name"
		if m.Name != "" {
			tmp[1] = "Name: " + m.Name
		}
		tmp[2] = "State: " + m.State
		text := strings.Join(tmp, "\n") + "\n"

		f := slack.Field{
			Title: "",
			Value: text,
		}

		c := "stopped"
		switch m.State {
		case "running":
			c = "good"
		case "pending":
			c = "warning"
		case "stopping", "stopped", "shutting-down", "terminated":
			c = "danger"
		}
		a := slack.Attachment{Color: &c}
		a.AddField(f)
		attachments[i] = a
	}
	return attachments
}
