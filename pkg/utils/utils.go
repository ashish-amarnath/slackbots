package utils

import (
	"fmt"
	"strings"

	"github.com/ashish-amarnath/slackbots/cmd"
	"github.com/ashish-amarnath/slackbots/pkg/types"
	"github.com/golang/glog"
)

// StringifyMessage returns a string representation of a message
func StringifyMessage(msg types.Message) string {
	return fmt.Sprintf("[ID=%d, Type=%s, Text=%s, Channel=%s]",
		msg.ID, msg.Type, msg.Text, msg.Channel)
}

// GetBotType parses message text to extract bot type
func getBotType(msgText string) string {
	return strings.Split(msgText, " ")[0]
}

func ProcessBotRquest(msgText, metadataServerURL, metadataServerAPIKey string) string {
	glog.V(2).Infof("msgTxt: %s\n", msgText)
	botType := getBotType(msgText)
	var botResp string
	if botType == types.ValidateKube2IamBotReq {
		botResp = cmd.ProcessValidateKube2IamReq(metadataServerURL, metadataServerAPIKey, msgText)
	} else if botType == types.ApplysKube2IamBotReq {
		botResp = cmd.ApplyKube2IamReq(msgText)
	}

	return botResp
}
