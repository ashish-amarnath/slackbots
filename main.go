package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ashish-amarnath/slackbots/pkg/slack"
	"github.com/ashish-amarnath/slackbots/pkg/types"
	"github.com/ashish-amarnath/slackbots/pkg/utils"
	"github.com/golang/glog"
)

var (
	helpFlag                *bool
	awsMetadataServerAPIKey *string
	awsMetadataServerURL    *string
	adGroupMemberLookupURL  *string
	slackbotToken           *string
)

func printUsage() {
	fmt.Println("Usage:")
}

func main() {
	helpFlag = flag.Bool("help", false, "")
	awsMetadataServerAPIKey = flag.String("apikey", "", "API key to use to engage AWS meta-data service")
	awsMetadataServerURL = flag.String("metadataServerURL", "", "URL for AWS metadata server to get account info")
	adGroupMemberLookupURL = flag.String("adgrouplookupurl", "", "URL for the AD group member list service.")
	slackbotToken = flag.String("slackbotToken", "", "Slack generated token for the bot")
	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(1)
	}

	slackConn := slack.NewSlackServerConn(*slackbotToken)

	glog.V(1).Infoln("Slackbot listening for messages to process...")
	for {
		msg, err := slackConn.ReadMessage()
		if err != nil {
			glog.Fatalf("Failed to read message sent to slackbot. err=%s\n", err)
		}
		if msg.Type != types.MessageType {
			glog.V(9).Infof("Ignoring messages of type %s\n", msg.Type)
			continue
		}

		glog.V(10).Infof("Processing message [%s]\n", utils.StringifyMessage(msg))
		botResp := utils.ProcessBotRquest(msg.Text, *adGroupMemberLookupURL, *awsMetadataServerURL, *awsMetadataServerAPIKey)
		glog.V(4).Infof("Bot response=%s", botResp)
		msg.Text = botResp
		glog.V(8).Infoln(utils.StringifyMessage(msg))
		slackConn.SendMessage(msg)
	}
}
