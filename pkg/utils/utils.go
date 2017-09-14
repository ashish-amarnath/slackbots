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
	return fmt.Sprintf("[ID=%d, Type=%s, Text=%s, Channel=%s, User=%s]",
		msg.ID, msg.Type, msg.Text, msg.Channel, msg.User)
}

// GetBotType parses message text to extract bot type
func getBotType(msgText string) string {
	return strings.Split(msgText, " ")[0]
}

// ProcessBotRquest processes the request based on the request type
func ProcessBotRquest(msgText, adGroupMemberlistURL, metadataServerURL, metadataServerAPIKey string) string {
	glog.V(9).Infof("msgTxt: %s\n", msgText)

	botType := getBotType(msgText)
	var botResp string
	if botType == types.ValidateKube2IamBotReq {
		botResp = cmd.ProcessValidateKube2IamReq(adGroupMemberlistURL, metadataServerURL, metadataServerAPIKey, msgText)
	} else if botType == types.ApplysKube2IamBotReq {
		botResp = cmd.ApplyKube2IamReq(msgText)
	}

	return botResp
}
