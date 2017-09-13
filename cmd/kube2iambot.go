package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/ashish-amarnath/slackbots/pkg/types"
	"github.com/golang/glog"
)

func getRoleAccNum(arnName string) string {
	roleArnParts := strings.Split(arnName, ":")

	return roleArnParts[4]
}

func getAccountAWSAccountOwnerID(metadataServerURL, metadataServerAPIKey, awsAccountNumber string) string {
	reqEp := fmt.Sprintf("%s/%s=%s", metadataServerURL, types.AWSMetaDataServerAccRsrcEp, awsAccountNumber)
	glog.V(1).Infof("Meta-data server url=%s\n", reqEp)
	req, err := http.NewRequest("GET", reqEp, nil)
	if err != nil {
		err := fmt.Errorf("failed to create new HTTP request to meta-data server at %s to get account ownerID for account number=%s err=%s", metadataServerURL, awsAccountNumber, err)
		glog.Error(err)
		return err.Error()
	}
	req.Header.Set("X-Api-Key", metadataServerAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		err := fmt.Errorf("failed to get account ownerID for account number=%s err=%s httpsStatusCode=%d[%s]", awsAccountNumber, err, resp.StatusCode, resp.Status)
		glog.Error(err)
		return err.Error()
	}

	rBody, err := ioutil.ReadAll(resp.Body)
	glog.V(1).Infof("AccountOwnerId: \n[%s]\n", rBody)

	resp.Body.Close()

	return ""
}

func getAWSAccountOwners(accountNumber string) string {
	kube2iamValidateScript := fmt.Sprintf("%s %s %s", types.Kube2IamValidateScriptLoc, accountNumber, types.AWSMetadataServerAPIKey)
	glog.V(4).Infof("running command %s\n", kube2iamValidateScript)
	cmd := exec.Command("bash", "-c", kube2iamValidateScript)
	out, _ := cmd.Output()
	return string(out)
}

func getAWSAccOwners(roleAccOwnerID string) string {
	glog.V(2).Infof("roleAccOwnerID=%s\n", roleAccOwnerID)
	return ""
}

func ProcessValidateKube2IamReq(metaDataServerURL, metadataServerAPIKey, msgText string) string {
	msgTxtArr := strings.Split(msgText, " ")
	awsRoleArn := msgTxtArr[3]
	awsAccountNumber := getRoleAccNum(awsRoleArn)
	roleAccOwnerID := getAccountAWSAccountOwnerID(metaDataServerURL, metadataServerAPIKey, awsAccountNumber)

	return getAWSAccOwners(roleAccOwnerID)
}

// ProcessValidateKube2IamReq validates kube2iam request
func ProcessValidateKube2IamReq_v1(msgText string) string {
	msgTxtArr := strings.Split(msgText, " ")
	namespace := msgTxtArr[1]
	requestorName := msgTxtArr[2]
	awsRoleArn := msgTxtArr[3]
	awsAccountNumber := getRoleAccNum(awsRoleArn)

	glog.V(4).Infof("Processing kube2iam request for namespace=%s from requestor=%s, roleArn=%s, awsAccountNumber=%s", namespace, requestorName, awsRoleArn, awsAccountNumber)
	resp := strings.Join(strings.Split(strings.TrimSpace(getAWSAccountOwners(awsAccountNumber)), "\n"), ";")
	glog.V(3).Infof("account owners: %s", resp)
	return getAWSAccOwners(awsAccountNumber)
}

// ApplyKube2IamReq applies kube2iam annotations to namespaces
func ApplyKube2IamReq(msgText string) string {
	msgTxtArr := strings.Split(msgText, " ")
	namespace := msgTxtArr[1]
	awsRoleArn := msgTxtArr[2]

	resp := fmt.Sprintf("Allowing pods in namespace=%s to assume role=%s", namespace, awsRoleArn)
	glog.V(1).Infof(resp)

	return resp
}
