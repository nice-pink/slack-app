package main

import (
	"flag"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/slack-app/pkg/send"
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

	client := send.NewClient()
	if *header != "" || *text != "" {
		err := client.SendMsg(*header, *text, *color, *channelId)
		if err != nil {
			os.Exit(2)
		}
	}

	if *file != "" {
		err := client.SendFile(*file, *fileTitle, *channelId)
		if err != nil {
			os.Exit(2)
		}
	}

	if *folder != "" {
		errs := client.SendFiles(*folder, *fileTitle, *channelId)
		if len(errs) > 0 {
			os.Exit(2)
		}
	}

	log.Info("*** Start")
	log.Info(os.Args)
}
