package main

import (
	"encoding/json"
	"flag"
	"os"
	"path"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/slack-go/slack"
)

func main() {
	header := flag.String("header", "", "Header.")
	text := flag.String("text", "", "Text.")
	color := flag.String("color", "", "Color.")
	channelId := flag.String("channelId", "", "Channel ID.")
	file := flag.String("file", "", "Filepath.")
	fileTitle := flag.String("fileTitle", "", "File title.")
	flag.Parse()

	client := GetClient()
	if *header != "" || *text != "" {
		err := SendMsg(*header, *text, *color, *channelId, client)
		if err != nil {
			os.Exit(2)
		}
	}

	if *file != "" {
		err := SendFile(*file, *fileTitle, *channelId, client)
		if err != nil {
			os.Exit(2)
		}
	}

	log.Info("*** Start")
	log.Info(os.Args)
}

func GetClient() *slack.Client {
	token := os.Getenv("SLACK_TOKEN")
	return slack.New(token)
}

func SendMsg(header, text, color, channelId string, client *slack.Client) error {
	attachment := slack.Attachment{
		Text: text,
		// Fields: []slack.AttachmentField{
		// 	slack.AttachmentField{
		// 		Title: "title",
		// 		Value: "value",
		// 	},
		// },
	}
	if color != "" {
		attachment.Color = color
	}

	channel, ts, err := client.PostMessage(
		channelId,
		slack.MsgOptionText(header, false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		log.Err(err, "post message")
		return err
	}
	log.Info(ts, ":: Posted message to", channel)

	return nil
}

func GetBlock(header, text string, client *slack.Client) ([]byte, error) {
	// prefix = "🍀 *Success!*"
	//     if "{{inputs.parameters.did-fail}}" == "true":
	//       prefix = "💥 *Failed!*"

	//     url = os.getenv('SLACK_WEBHOOK')
	//     headline = prefix + " *{{inputs.parameters.app-name}}*."
	//     info="""{{inputs.parameters.info}}"""
	//     block_text = "*info:* _" + info + "_\n*workflow:* _{{workflow.name}}_\n*status:* _{{workflow.status}}_\n"
	//     attachment = { "color": color,
	//             "blocks": [
	//                 {"type": "section",
	//                 "text": {"text": block_text,
	//                           "type": "mrkdwn"}}
	//             ]}
	//     body = {
	//         "text": headline,
	//         "attachments": [
	//             attachment
	//         ]
	//       }

	var headerSection *slack.SectionBlock
	if header != "" {
		headerObj := slack.NewTextBlockObject("mrkdwn", header, false, false)
		headerSection = slack.NewSectionBlock(headerObj, nil, nil)
	}

	var fieldsSection *slack.SectionBlock
	if text != "" {
		blockText := slack.NewTextBlockObject("mrkdwn", text, false, false)
		fieldsSection = slack.NewSectionBlock(nil, []*slack.TextBlockObject{blockText}, nil)
	}

	msg := slack.NewBlockMessage(
		headerSection,
		fieldsSection,
	)

	body, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		log.Err(err, "marshal message")
		return nil, err
	}

	return body, nil
}

func SendFile(filepath, title, channel string, client *slack.Client) error {
	log.Info("Send file", filepath)
	filename := path.Base(filepath)

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		log.Err(err, "stat file", filepath)
		return err
	}

	params := slack.UploadFileV2Parameters{
		Title:    title,
		File:     filepath,
		FileSize: int(fileInfo.Size()),
		Filename: filename,
		Channel:  channel,
	}
	file, err := client.UploadFileV2(params)
	if err != nil {
		log.Err(err, "upload file error", filepath)
		return err
	}
	log.Info("Uploaded", file.ID, file.Title)
	return nil
}
