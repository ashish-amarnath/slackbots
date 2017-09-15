package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ashish-amarnath/slackbots/pkg/types"
	"github.com/golang/glog"
)

func whichKubectl() (loc string, err error) {
	loc, err = RunBashCmd(types.WhichKubectl)
	return
}

// RunBashCmd runs a supplied bash command
func RunBashCmd(cmd string) (res string, err error) {
	glog.V(8).Infof("Running [%s]\n", cmd)
	toRun := exec.Command("bash", "-c", cmd)
	var stderr bytes.Buffer
	toRun.Stderr = &stderr
	out, err := toRun.Output()
	if err != nil {
		res = ""
		glog.Infof("stderr: %s", stderr.String())
	}
	res = strings.TrimSpace(string(out))
	return
}

func getKubeCtlBaseCmd(kubeconfig, cluster string) (baseCmd string, err error) {
	var kcLoc string
	kcLoc, err = whichKubectl()
	baseCmd = fmt.Sprintf("%s --user %s_sudo --context=%s --kubeconfig=%s", kcLoc, cluster, cluster, kubeconfig)
	return
}

// ApplyUpdatedNamespaceMetadata applies the supplied namespace metadata to the supplied namespace in the supplied cluster
func ApplyUpdatedNamespaceMetadata(kubeConfig, cluster, metadataJSON string) (err error) {
	kcBaseCmd, err := getKubeCtlBaseCmd(kubeConfig, cluster)
	applyCmd := fmt.Sprintf("%s | %s apply -f", metadataJSON, kcBaseCmd)
	_, err = RunBashCmd(applyCmd)
	return
}

// GetNamespaceDefnJSON fetches the current namespace definition in JSON format
func GetNamespaceDefnJSON(kubeConfig, cluster, namespace string) (json string, err error) {
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

// GetBotReqType parses message text to extract bot type
func GetBotReqType(msgText string) string {
	return strings.Split(msgText, " ")[0]
}
