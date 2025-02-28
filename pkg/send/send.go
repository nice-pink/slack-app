package send

import (
	"encoding/json"
	"os"
	"path"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/slack-go/slack"
)

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

func SendFiles(folder, title, channel string, client *slack.Client) []error {
	files, err := os.ReadDir(folder)
	if err != nil {
		return []error{err}
	}

	errs := []error{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		err := SendFile(path.Join(folder, file.Name()), file.Name(), channel, client)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
