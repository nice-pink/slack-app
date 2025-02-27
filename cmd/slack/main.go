package main

import (
	"flag"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/slack-app/pkg/msg"
)

func main() {
	header := flag.String("header", "", "Header.")
	text := flag.String("text", "", "Text.")
	color := flag.String("color", "", "Color.")
	channelId := flag.String("channelId", "", "Channel ID.")
	file := flag.String("file", "", "Filepath.")
	folder := flag.String("folder", "", "Filepath.")
	fileTitle := flag.String("fileTitle", "", "File title.")
	flag.Parse()

	client := msg.GetClient()
	if *header != "" || *text != "" {
		err := msg.SendMsg(*header, *text, *color, *channelId, client)
		if err != nil {
			os.Exit(2)
		}
	}

	if *file != "" {
		err := msg.SendFile(*file, *fileTitle, *channelId, client)
		if err != nil {
			os.Exit(2)
		}
	}

	if *folder != "" {
		errs := msg.SendFiles(*folder, *fileTitle, *channelId, client)
		if len(errs) > 0 {
			os.Exit(2)
		}
	}

	log.Info("*** Start")
	log.Info(os.Args)
}
