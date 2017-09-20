package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ashish-amarnath/slackbots/cmd"
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
	adLookupServerURL       *string
	slackbotToken           *string
	kubeconfig              *string
)

func printUsage() {
	fmt.Println("Usage:")
}

func main() {
	helpFlag = flag.Bool("help", false, "")
	awsMetadataServerAPIKey = flag.String("apikey", "", "API key to use to engage AWS meta-data service")
	awsMetadataServerURL = flag.String("metadataServerURL", "", "URL for AWS metadata server to get account info")
	adGroupMemberLookupURL = flag.String("adgrouplookupurl", "", "URL for the AD group member list service.")
	adLookupServerURL = flag.String("adLookupServerURL", "", "URL to lookup AD user")
	slackbotToken = flag.String("slackbotToken", "", "Slack generated token for the bot")
	kubeconfig = flag.String("kubeconfig", "", "Path to the kubeconfig for kubectl to use")
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
			glog.Errorf("Failed to read message sent to slackbot. err=%s\n", err.Error())
			continue
		}
		if msg.Type != types.MessageType || !strings.HasPrefix(msg.Text, "<@"+slackConn.UserID+">") {
			glog.V(9).Infof("Ignoring message %s\n", utils.StringifyMessage(msg))
			continue
		}

		go cmd.ProcessBotRquest(slackConn, msg, *adGroupMemberLookupURL, *awsMetadataServerURL, *awsMetadataServerAPIKey, *kubeconfig, *adLookupServerURL)
	}
}
