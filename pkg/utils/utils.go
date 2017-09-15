package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ashish-amarnath/slackbots/cmd"
	"github.com/ashish-amarnath/slackbots/pkg/types"
	"github.com/golang/glog"
)

func whichKubectl() (loc string, err error) {
	loc, err = RunBashCmd(types.WhichKubectl)
	return
}

// RunBashCmd runs a supplied bash command
func RunBashCmd(cmd string) (res string, err error) {
	toRun := exec.Command("bash", "-c", cmd)
	var stderr bytes.Buffer
	toRun.Stderr = &stderr
	out, err := toRun.Output()
	if err != nil {
		glog.Infof("stderr: %s", stderr.String())
	}
	res = string(out)
	return
}

func getKubeCtlBaseCmd(kubeconfig, cluster string) (baseCmd string, err error) {
	var kcLoc string
	kcLoc, err = whichKubectl()
	baseCmd = fmt.Sprintf("%s --user %s_sudo --context=%s --kubeconfig=%s", kcLoc, cluster, cluster, kubeconfig)
	return
}

func getNamespaceDefnJSON(kubeConfig, cluster, namespace string) (json string, err error) {
	var kcBaseCmd string
	kcBaseCmd, err = getKubeCtlBaseCmd(kubeConfig, cluster)
	bashCmd := fmt.Sprintf(" get namespace %s --export=true -ojson", namespace)
	json, err = RunBashCmd(kcBaseCmd + bashCmd)
	return
}

// StringifyMessage returns a string representation of a message
func StringifyMessage(msg types.Message) string {
	return fmt.Sprintf("[ID=%d, Type=%s, Text=%s, Channel=%s, User=%s]",
		msg.ID, msg.Type, msg.Text, msg.Channel, msg.User)
}

// GetBotType parses message text to extract bot type
func getBotReqType(msgText string) string {
	return strings.Split(msgText, " ")[0]
}

// ProcessBotRquest processes the request based on the request type
func ProcessBotRquest(botReq, adLookupServerURL, metadataServerURL, metadataServerAPIKey, kubeconfig, cluster string) string {
	glog.V(9).Infof("msgTxt: %s\n", botReq)

	botReqType := getBotReqType(botReq)
	var botResp string
	if botReqType == types.ValidateKube2IamBotReq {
		botResp = cmd.ProcessValidateKube2IamReq(adLookupServerURL, metadataServerURL, metadataServerAPIKey, botReq)
	} else if botReqType == types.ApplysKube2IamBotReq {
		botResp = cmd.ApplyKube2IamReq(botReq, kubeconfig, cluster)
	} else if botReqType == types.RejectKube2IamBotReq {
		botResp = cmd.RejectKube2IamReq(botReq)
	} else {
		glog.V(6).Infof("Unknown botReq %s", botReqType)
	}

	return botResp
}
