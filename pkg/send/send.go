package send

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/slack-go/slack"
)

type Msg struct {
	ChannelId string
	Text      string
	Header    string
	Color     string
}

type Client struct {
	slackClient *slack.Client
}

func NewClient(token string) *Client {
	if token == "" {
		token = os.Getenv("SLACK_TOKEN")
	}
	return &Client{slackClient: slack.New(token)}
}

func (c *Client) SendText(text, channelId string) error {
	if err := c.IsReady(); err != nil {
		return err
	}

	channel, ts, err := c.slackClient.PostMessage(
		channelId,
		slack.MsgOptionText(text, false),
	)
	if err != nil {
		log.Err(err, "post message")
		return err
	}
	log.Info(ts, ":: Posted message to", channel)
	return nil
}

func (c *Client) SendMsg(msg Msg) error {
	return c.Send(msg.Header, msg.Text, msg.Color, msg.ChannelId)
}

func (c *Client) Send(header, text, color, channelId string) error {
	if err := c.IsReady(); err != nil {
		return err
	}
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

	channel, ts, err := c.slackClient.PostMessage(
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

func (c *Client) GetBlock(header, text string) ([]byte, error) {
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

func (c *Client) SendFile(filepath, title, channel string) error {
	if err := c.IsReady(); err != nil {
		return err
	}

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
	file, err := c.slackClient.UploadFileV2(params)
	if err != nil {
		log.Err(err, "upload file error", filepath)
		return err
	}
	log.Info("Uploaded", file.ID, file.Title)
	return nil
}

func (c *Client) SendFiles(folder, title, channel string) []error {
	if err := c.IsReady(); err != nil {
		return []error{err}
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		return []error{err}
	}

	errs := []error{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		err := c.SendFile(path.Join(folder, file.Name()), file.Name(), channel)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (c *Client) IsReady() error {
	if c.slackClient == nil {
		log.Error("Slack Client not initialised.")
		return errors.New("slack client is not ready")
	}
	return nil
}
