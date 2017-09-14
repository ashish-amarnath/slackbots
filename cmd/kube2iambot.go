package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/ashish-amarnath/slackbots/pkg/types"
	"github.com/golang/glog"
)

func getAccNumFromRoleArn(arnName string) string {
	roleArnParts := strings.Split(arnName, ":")
	return roleArnParts[types.AccountNumberIndexInRoleArn]
}

func getAccountOwnerIDEndpoint(metadataServerURL, accNum string) string {
	return fmt.Sprintf("%s/%s=%s", metadataServerURL, types.AWSMetaDataServerAccRsrcEp, accNum)
}

func parseAccOwnerResponse(raw []byte) (respObj types.AccNumRespMsg, err error) {
	err = json.Unmarshal(raw, &respObj)
	return
}

func doHTTPRequest(url, apiKey string) (raw []byte, err error) {
	raw = nil
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err := fmt.Errorf("failed to create request to url=%s err=%s", url, err)
		glog.Error(err)
		return nil, err
	}
	req.Header.Set("X-Api-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		err := fmt.Errorf("request to url=%s failed err=%s, httpStatusCode=%d(%s)", url, err, resp.StatusCode, resp.Status)
		glog.Error(err)
		return nil, err
	}

	raw, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	glog.V(1).Infof("%s\n%s", url, raw)
	return raw, nil
}

func getAWSAccountOwnerID(baseURL, apiKey, awsAccNum string) string {
	url := getAccountOwnerIDEndpoint(baseURL, awsAccNum)

	rBody, err := doHTTPRequest(url, apiKey)
	respJSON, err := parseAccOwnerResponse(rBody)
	if err != nil {
		err := fmt.Errorf("doHttpRequest to getAWSAccountOwnerID url=%s failed, err=%s", url, err)
		glog.Error(err)
		return err.Error()
	}

	return fmt.Sprintf("%d", respJSON.Data[0].OwnerTeamID)
}

func parseAdSecGrpResponse(raw []byte) (respObj types.AdSecurityGroupResp, err error) {
	err = json.Unmarshal(raw, &respObj)
	return
}

func getOwnerADSecurityGroup(baseURL, apiKey, ownerTeadID string) string {
	url := fmt.Sprintf("%s/%s=%s", baseURL, types.ADSecurityGroupEndPoint, ownerTeadID)

	rBody, err := doHTTPRequest(url, apiKey)
	if err != nil {
		err := fmt.Errorf("doHttpRequest to url=%s failed, err=%s", url, err)
		glog.Error(err)
		return err.Error()
	}
	respJSON, err := parseAdSecGrpResponse(rBody)
	if err != nil {
		err := fmt.Errorf("failed to parse response from end point %s, err=%s", url, err)
		glog.Error(err)
		return err.Error()
	}

	return respJSON.Data[0].ADSecurityGroup
}

func parseADGroupMemberListResp(raw []byte) (respJSON types.ADGroupMemberListResp, err error) {
	err = json.Unmarshal(raw, &respJSON)
	return
}

func getAdGrpMembers(adGroupMemberlistURL, adSecGrp string) string {
	curlCmd := fmt.Sprintf("curl -s %s/%s", adGroupMemberlistURL, adSecGrp)
	glog.V(5).Infof("curlCmd: %s", curlCmd)
	cmd := exec.Command("bash", "-c", curlCmd)
	out, err := cmd.Output()
	if err != nil {
		glog.Errorf("failed to successfully run [%s] err=%s", curlCmd, err)
		return err.Error()
	}
	adGrpMemberListResp, err := parseADGroupMemberListResp(out)
	return strings.Join(adGrpMemberListResp.Members.Users, ", ")
}

// ProcessValidateKube2IamReq validates kube2iam request
func ProcessValidateKube2IamReq(adGrpListURL, mdsURL, mdsAPIKey, msg string) string {
	msgParts := strings.Split(msg, " ")
	awsRoleArn := msgParts[1]

	awsAccountNumber := getAccNumFromRoleArn(awsRoleArn)
	roleAccOwnerID := getAWSAccountOwnerID(mdsURL, mdsAPIKey, awsAccountNumber)
	adSecGrp := getOwnerADSecurityGroup(mdsURL, mdsAPIKey, roleAccOwnerID)
	return getAdGrpMembers(adGrpListURL, adSecGrp)
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
