package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ashish-amarnath/slackbots/pkg/types"
	"github.com/golang/glog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAccNumFromRoleArn(t *testing.T) {
	Convey("getAccNumFromRoleArn", t, func() {
		Convey("should parse out account number from a valid AWS role ARN", func() {
			testRoleArn := "arn:aws:iam::123456789012:role/foo/bar/foo-bar-mysuperawesomerole"
			expected := "123456789012"
			actual, err := getAccNumFromRoleArn(testRoleArn)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})
		Convey("should report error when called with invalid role ARN", func() {
			invalidRoleArn := "arn-aws-iam--123456789012-role/foo/bar/foo-bar-mysuperawesomerole"
			actual, err := getAccNumFromRoleArn(invalidRoleArn)
			So(actual, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetAccountOwnerIDEndpoint(t *testing.T) {
	Convey("getAccountOwnerIDEndpoint should return the correct endpoing to get account owner id", t, func() {
		testSrvrURL := "https://myawesome-metadata-server.com"
		testAccNum := "123456789012"
		actual := getAccountOwnerIDEndpoint(testSrvrURL, testAccNum)
		expecterd := `https://myawesome-metadata-server.com/dev_read/accounts?AccountNumber=123456789012`

		So(actual, ShouldResemble, expecterd)
	})
}

func TestParseAccOwnerResponse(t *testing.T) {
	Convey("parseAccOwnerResponse", t, func() {
		Convey("Should fail when called with invalid JSON bytes", func() {
			invalidJSON := `{"data": [{"EnvironmentId": 7, "OwnerTeamId": 11, "GUID": 11534399, "AccountNumber": "123456789012", "AccountName": "unittest", "RequesterPersonId": 10, "ClaimRule":`
			_, err := parseAccOwnerResponse([]byte(invalidJSON))
			So(err, ShouldNotBeNil)
		})
		Convey("Should parse a valid JSON string into AccNumRespMsg", func() {
			validJSON := `{"data": [{"EnvironmentId": 7, "OwnerTeamId": 11, "GUID": 11534399, "AccountNumber": "123456789012", "AccountName": "unittest", "RequesterPersonId": 10, "ClaimRule": "@RuleName = \"Role mapping for unittest Account\"\nc:[Type == \"http://temp/variable\", Value =~ \"(?i)^AWS-Prod_k8s-\"]\n=> issue(Type = \"https://aws.amazon.com/SAML/Attributes/Role\", Value = RegExReplace(c.Value, \"AWS-Prod_k8s-\", \"arn:aws:iam::184475218016:saml-provider/NORD,arn:aws:iam::184475218016:role/NORD-Prod_k8s-\"));", "Size": "TEAMXL"}]}`
			actual, err := parseAccOwnerResponse([]byte(validJSON))
			So(err, ShouldBeNil)
			So(len(actual.Data), ShouldEqual, 1)
			So(actual.Data[0].EnvironmentID, ShouldEqual, 7)
			So(actual.Data[0].OwnerTeamID, ShouldEqual, 11)
			So(actual.Data[0].GUID, ShouldEqual, 11534399)
			So(actual.Data[0].AccountNumber, ShouldResemble, "123456789012")
			So(actual.Data[0].AccountName, ShouldResemble, "unittest")
			So(actual.Data[0].RequesterPersonID, ShouldEqual, 10)
			So(actual.Data[0].Size, ShouldResemble, "TEAMXL")
		})
	})
}

func TestParseKubernetesNamespaceMetadata(t *testing.T) {
	Convey("parseKubernetesNamespaceMetadata", t, func() {
		Convey("Should fail when invoeked with invalid JSIN bytes", func() {
			invalidJSON := `{
				"apiVersion": "v1",
				"kind": "Namespace",
				"metadata": {
				  "annotations": {
					"contact-email": "unit@test.com",
					"cost-center": "",
					"kube2iam.beta.nordstrom.net/allowed-roles": "[\"arn:aws:iam::123456789012:role/superpowerfulrole1\",\"arn:aws:iam::123456789012:role/superpowerfulrole2",\"arn:aws:iam::123456789012:role/superpowerfulrole3",\"arn:aws:iam::123456789012:role/superpowerfulrole4\"]\n",
					"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Namespace\",\"metadata\":{\"annotations\":{\"contact-email\":\"unit@test.com\",\"cost-center\":\"\",\"kube2iam.beta.nordstrom.net/allowed-roles\":\"[\\\"arn:aws:iam::123456789012:role/foo/k8s/fooS3AndKmsStack-fooFlinkForS3AndKMS-9336B1OLABTR\\\"]\\n\",\"kubernetes.io/change-cause\":\"kubectl edit ns foo --user=admin\",\"slack-channel-events\":\"\",\"slack-channel-urgent\":\"\",\"slack-channel-users\":\"#foo\"},\"creationTimestamp\":null,\"name\":\"foo\",\"namespace\":\"\",\"selfLink\":\"/api/v1/namespacesfoo\"},\"spec\":{\"finalizers\":[\"kubernetes\"]},\"status\":{\"phase\":\"Active\"}}\n",
					"kubernetes.io/change-cause": "kubectl edit ns foo --context=cluster --user=admin"
			  }`
			_, err := parseKubernetesNamespace([]byte(invalidJSON))
			So(err, ShouldNotBeNil)
		})
		Convey("Should parse a valid JSON string into KubernetesNamespaceMetadata", func() {
			validJSON := `{
				"apiVersion": "v1",
				"kind": "Namespace",
				"metadata": {
					"annotations": {
						"contact-email": "unit@test.com",
						"cost-center": "",
						"kube2iam.beta.nordstrom.net/allowed-roles": "[\"arn:aws:iam::123456789012:role/foo/k8s/superawesomerole1\",\"arn:aws:iam::123456789012:role/foo/k8s/foo/k8s/superawesomerole2\",\"arn:aws:iam::123456789012:role/NonProd_DSBIA/k8s/superawesomerole3\",\"arn:aws:iam::123456789012:role/NonProd_DSBIA/k8s/superawesomerole2\"]\n",
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Namespace\",\"metadata\":{\"annotations\":{\"contact-email\":\"unit.test@unittest.com\",\"cost-center\":\"\",\"kube2iam.beta.nordstrom.net/allowed-roles\":\"[\\\"arn:aws:iam::123456789012:role/foo/k8s/fooS3AndKmsStack-fooFlinkForS3AndKMS-9336B1OLABTR\\\"]\\n\",\"kubernetes.io/change-cause\":\"kubectl edit ns foo --user=admin\",\"slack-channel-events\":\"\",\"slack-channel-urgent\":\"\",\"slack-channel-users\":\"#foo\"},\"creationTimestamp\":null,\"name\":\"foo\",\"namespace\":\"\",\"selfLink\":\"/api/v1/namespacesfoo\"},\"spec\":{\"finalizers\":[\"kubernetes\"]},\"status\":{\"phase\":\"Active\"}}\n",
						"kubernetes.io/change-cause": "kubectl edit ns foo --context=cluster --user=admin",
						"slack-channel-events": "",
						"slack-channel-urgent": "",
						"slack-channel-users": "#foo"
					},
					"creationTimestamp": "2017-06-30T00:15:38Z",
					"name": "foo",
					"resourceVersion": "193695467",
					"selfLink": "/api/v1/namespacesfoo",
					"uid": "3e95e64d-5d29-11e7-8024-0607ec9cbe90"
				},
				"spec": {
					"finalizers": [
						"kubernetes"
					]
				},
				"status": {
					"phase": "Active"
				}
			}`
			var expected types.KubernetesNamespace
			expected.APIVersion = "v1"
			expected.Kind = "Namespace"
			expected.Metadata.Annotations.ContactEmail = "unit@test.com"
			actual, err := parseKubernetesNamespace([]byte(validJSON))
			So(err, ShouldBeNil)
			So(actual.APIVersion, ShouldResemble, expected.APIVersion)

		})
	})
}

func TestDoHTTPRequest(t *testing.T) {
	Convey("doHTTPRequest should returh with error when unable to successfully process the HTTP request", t, func() {
		url := "foobar.baz"
		apiKey := "supersecret"
		actual, err := doHTTPRequest(url, apiKey)
		So(actual, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestRunRawCurlCommands(t *testing.T) {
	Convey("runRawCurlCommands should return with error when unable to successfully run the curl cmd", t, func() {
		url := "foobar.baz"
		atcual, err := runRawCurlCommands(url)
		So(atcual, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})
}

func TestGetAWSAccountOwnerID(t *testing.T) {
	Convey("getAWSAccountOwnerID should return error when unable to process request sucessfully", t, func() {
		glog.Errorf("Expected Error:\n")
		actual, err := getAWSAccountOwnerID("https://example.com", "open-key", "123456789012")
		So(actual, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})
}

func TestParseAdSecGrpResponse(t *testing.T) {
	Convey("parseAdSecGrpResponse", t, func() {
		Convey("Should fail when called with invalid JSON bytes", func() {
			invalidJSON := `{"data": [{"VP": 201, "Name": "k8s", "Tags": null, "OrgName": "Engineering Platform", "ADSecurityGroup": `
			_, err := parseAdSecGrpResponse([]byte(invalidJSON))
			So(err, ShouldNotBeNil)
		})
		Convey("Should parse a valid JSON string into AdSecurityGroupResp", func() {
			validJSON := `{"data": [{"Name": "ut", "OrgName": "unit test", "ADSecurityGroup": "unittestAdmins", "ID": 11, "Director": 123, "CostCenter": 12345, "EmailDistList": "unit@test.com"}]}`
			actual, err := parseAdSecGrpResponse([]byte(validJSON))
			So(err, ShouldBeNil)
			So(actual.Data[0].ADSecurityGroup, ShouldResemble, "unittestAdmins")
			So(actual.Data[0].CostCenter, ShouldEqual, 12345)
			So(actual.Data[0].Director, ShouldResemble, 123)
			So(actual.Data[0].EmailDistList, ShouldResemble, "unit@test.com")
			So(actual.Data[0].ID, ShouldEqual, 11)
			So(actual.Data[0].Name, ShouldResemble, "ut")
			So(actual.Data[0].OrgName, ShouldResemble, "unit test")
		})
	})
}

func TestGetOwnerADSecurityGroup(t *testing.T) {
	Convey("getOwnerADSecurityGroup should return with error when unable to find AD security group for the AWS account number", t, func() {
		url := "foobar.baz"
		apiKey := "supersecret"
		owner := "ADMINS"
		actual, err := getOwnerADSecurityGroup(url, apiKey, owner)
		So(actual, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})
}

func TestParseADGroupMemberListResp(t *testing.T) {
	Convey("parseADGroupMemberListResp", t, func() {
		Convey("Should fail when called with invalid JSON bytes", func() {
			invalidJSON := `{"name":"codeNinjas","description":"super awesome group","email":"codeninjas`
			_, err := parseADGroupMemberListResp([]byte(invalidJSON))
			So(err, ShouldNotBeNil)
		})
		Convey("Should parse a valid JSON string into ADGroupMemberListResp", func() {
			validJSON := `{"name":"codeNinjas","description":"super awesome group","email":"codeninjas@ninjaing.com","type":"unittest","updated":"2017-08-28T17:24:46.000Z","members":{"groups":[],"users":["ninja1","ninja2","ninja3","ninja4","ninja5"]},"managedBy":{"group":null,"user":"ninjaLeader"},"groups":["Ninja-Team1","Ninja-Team2","Ninja-Team3"]}`
			actual, err := parseADGroupMemberListResp([]byte(validJSON))
			So(err, ShouldBeNil)
			So(actual.Description, ShouldResemble, "super awesome group")
			So(actual.Name, ShouldResemble, "codeNinjas")
			So(actual.Email, ShouldResemble, "codeninjas@ninjaing.com")
			So(actual.Type, ShouldResemble, "unittest")
			So(strings.Join(actual.Members.Users, ","), ShouldResemble, "ninja1,ninja2,ninja3,ninja4,ninja5")
		})
	})
}

func TestParseADUserResp(t *testing.T) {}

func TestGetAdGrpMembers(t *testing.T) {
	Convey("getAdGrpMembers should return with error when unable to get members of an AD group", t, func() {
		url := "myadserver.foo"
		adGrp := "ADMINS"
		actual, err := getAdGrpMembers(url, adGrp)
		So(actual, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestGetRoleOwners(t *testing.T) {
	Convey("getRoleOwners", t, func() {
		Convey("should return with error when unable to find owners of an AWS role", func() {
			adGrpURL := "my-awesome-adsrvr.foo"
			msdURL := "super-awesome-mdsSrv.foo"
			mdsAPIKey := "topsecret"
			testRole := "arn:aws:iam::123456789012:role/superawesome-powerful-Role3"
			actual, err := getRoleOwners(adGrpURL, msdURL, mdsAPIKey, testRole)
			So(actual, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
		Convey("should return with error when called with an invalid AWS role", func() {
			adGrpURL := "my-awesome-adsrvr.foo"
			msdURL := "super-awesome-mdsSrv.foo"
			mdsAPIKey := "topsecret"
			testRole := "arn-aws-iam--123456789012-role/superawesome-powerful-Role3"
			actual, err := getRoleOwners(adGrpURL, msdURL, mdsAPIKey, testRole)
			So(actual, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetADUsrLookupEp(t *testing.T) {
	Convey("getADUsrLookupEp should return the correct AD user lookup endpoint", t, func() {
		expected := `https://adUsrLkp/api/v1/usr/get/doe%2c%20john`
		actual := getADUsrLookupEp("john", "doe", "https://adUsrLkp/api/v1/usr/get")
		So(actual, ShouldResemble, expected)
	})
}

func TestGetADUserByCN(t *testing.T) {
	Convey("getADUserByCN should return with error when unable to get the requested AD user", t, func() {
		fname := "John"
		lname := "Doe"
		email := "john.doe@johndoe.com"
		adUsrURL := "https://adUsrLkp/api/v1/usr/get"
		_, err := getADUserByCN(fname, lname, email, adUsrURL)
		So(err, ShouldNotBeNil)
	})
}

func TestGetADUserForSlackUser(t *testing.T) {
	Convey("getADUserForSlackUser return error when unable to get AD user corresponding to the supplied slack user", t, func() {
		testSlackUsr := "U725Q5UAY"
		adUsrURL := "https://adUsrLkp/api/v1/usr/get"
		_, err := getADUserForSlackUser(testSlackUsr, adUsrURL)
		So(err, ShouldNotBeNil)
	})
}

func TestIsRequestorOwner(t *testing.T) {
	Convey("isRequestorOwner", t, func() {
		Convey("should return true when requestor is owner", func() {
			var testADUsr types.ADUser
			testADUsr.FirstName = "jOhN"
			testADUsr.LastName = "dOe"
			testADUsr.Email = "JOHN.DOE@jOhnDoE.cOm"
			var owners []string

			owners = append(owners, "Foe, John")
			owners = append(owners, "bar, foo")
			owners = append(owners, "Doe, Jane")
			owners = append(owners, "DoE, JoHn")
			actual := isRequestorOwner(testADUsr, owners)
			So(actual, ShouldBeTrue)
		})
		Convey("should return false when requestor is not owner", func() {
			var testADUsr types.ADUser
			testADUsr.FirstName = "jOhN"
			testADUsr.LastName = "dOe"
			testADUsr.Email = "JOHN.DOE@jOhnDoE.cOm"
			var owners []string

			owners = append(owners, "Foe, John")
			owners = append(owners, "bar, foo")
			owners = append(owners, "Doe, Jane")
			owners = append(owners, "DoE, James")
			actual := isRequestorOwner(testADUsr, owners)
			So(actual, ShouldBeFalse)
		})
	})
}

func TestIsRequestValid(t *testing.T) {
	Convey("isRequestValid", t, func() {
		Convey("should return true for a valid request", func() {
			var validReq types.BotReqParams
			validReq.ADGroupLookupURL = "https://adGrpLkp/api/v1/usr/get"
			validReq.ADUserLookupURL = "https://adUsrLkp/api/v1/usr/get"
			validReq.AWSAPIKey = "blahziblahziblah"
			validReq.AWSMetadataServerURL = "https://jibberish.execute-api.us-west-81.amazonaws.com"
			validReq.KubeConfig = "/User/craycrayuser/.kube/config"
			validReq.Message = "@superbot !doSomethingAwesome foo arn:aws:iam::123456789012:role/superawesome-powerful-Role3 hydrogen"
			validReq.SlackUser = "UCRAY7Q"

			actual := isRequestValid(validReq)
			So(actual, ShouldBeTrue)
		})
		Convey("should return false for an invalid request", func() {
			var invalidReq types.BotReqParams
			actual := isRequestValid(invalidReq)
			So(actual, ShouldBeFalse)
		})
	})
}

func TestRequestKube2IamReq(t *testing.T) {
	Convey("RequestKube2IamReq", t, func() {
		Convey("should return error when unable to get owners of role ARN", func() {
			var validReq types.BotReqParams
			validReq.ADGroupLookupURL = "https://adGrpLkp/api/v1/usr/get"
			validReq.ADUserLookupURL = "https://adUsrLkp/api/v1/usr/get"
			validReq.AWSAPIKey = "blahziblahziblah"
			validReq.AWSMetadataServerURL = "https://jibberish.execute-api.us-west-81.amazonaws.com"
			validReq.KubeConfig = "/User/craycrayuser/.kube/config"
			validReq.Message = "@superbot !doSomethingAwesome foo arn:aws:iam::123456789012:role/superawesome-powerful-Role3 hydrogen"
			validReq.SlackUser = "UCRAY7Q"

			expected := `Failed to get owners of awsRoleArn=arn:aws:iam::123456789012:role/superawesome-powerful-Role3. err=doHttpRequest to getAWSAccountOwnerID url=https://jibberish.execute-api.us-west-81.amazonaws.com/dev_read/accounts?AccountNumber=123456789012 failed, err=unexpected end of JSON input`
			actual := RequestKube2IamReq(validReq)
			So(actual, ShouldResemble, expected)
		})
		Convey("should return error when called with invalid request", func() {
			var invalidReq types.BotReqParams
			expected := fmt.Sprintf("ERROR:\n Request should be of the form \n %s Order is important. Received ```%s```", types.RequestKube2IamBotReqFormat, invalidReq.Message)
			actual := RequestKube2IamReq(invalidReq)
			So(actual, ShouldResemble, expected)
		})
	})
}

func TestAddNewKube2IamRole(t *testing.T) {
	Convey("addNewKube2IamRole", t, func() {
		Convey("should add a new role to existing empty roles", func() {
			testRole := "arn:aws:iam::123456789012:role/superawesome-powerful-Role3"
			currentRole := "[]"

			expected := `["arn:aws:iam::123456789012:role/superawesome-powerful-Role3"]`
			actual := addNewKube2IamRole(currentRole, testRole)
			So(actual, ShouldResemble, expected)
		})
		Convey("should add new role to existing roles", func() {
			current := `["arn:aws:iam::123456789012:role/superawesome-powerful-Role3","arn:aws:iam::123456789012:role/superawesome-powerful-Role1"]`
			newRole := "arn:aws:iam::123456789012:role/superawesome-powerful-Role2"

			expected := `["arn:aws:iam::123456789012:role/superawesome-powerful-Role3","arn:aws:iam::123456789012:role/superawesome-powerful-Role1","arn:aws:iam::123456789012:role/superawesome-powerful-Role2"]`
			actual := addNewKube2IamRole(current, newRole)
			So(actual, ShouldResemble, expected)
		})
		Convey("should not add duplicate roles", func() {
			current := `["arn:aws:iam::123456789012:role/superawesome-powerful-Role3","arn:aws:iam::123456789012:role/superawesome-powerful-Role1","arn:aws:iam::123456789012:role/superawesome-powerful-Role2"]`
			newRole := `arn:aws:iam::123456789012:role/superawesome-powerful-Role3`

			actual := addNewKube2IamRole(current, newRole)
			So(actual, ShouldResemble, current)
		})
	})
}

func TestApproveKube2IamReq(t *testing.T) {
	Convey("ApproveKube2IamReq", t, func() {
		Convey("should return error when unable to get owners of role ARN", func() {
			var validReq types.BotReqParams
			validReq.ADGroupLookupURL = "https://adGrpLkp/api/v1/usr/get"
			validReq.ADUserLookupURL = "https://adUsrLkp/api/v1/usr/get"
			validReq.AWSAPIKey = "blahziblahziblah"
			validReq.AWSMetadataServerURL = "https://jibberish.execute-api.us-west-81.amazonaws.com"
			validReq.KubeConfig = "/User/craycrayuser/.kube/config"
			validReq.Message = "@superbot !doSomethingAwesome foo arn:aws:iam::123456789012:role/superawesome-powerful-Role3 hydrogen"
			validReq.SlackUser = "UCRAY7Q"

			expected := `Failed to get owners of awsRoleArn=arn:aws:iam::123456789012:role/superawesome-powerful-Role3. err=doHttpRequest to getAWSAccountOwnerID url=https://jibberish.execute-api.us-west-81.amazonaws.com/dev_read/accounts?AccountNumber=123456789012 failed, err=unexpected end of JSON input`
			actual := ApproveKube2IamReq(validReq)
			So(actual, ShouldResemble, expected)
		})
		Convey("should return error when called with invalid request", func() {
			var invalidReq types.BotReqParams
			expected := fmt.Sprintf("ERROR:\n Request should be of the form \n %s Order is important. Received ```%s```", types.ApproveKube2IamBotReqFormat, invalidReq.Message)
			actual := ApproveKube2IamReq(invalidReq)
			So(actual, ShouldResemble, expected)
		})
	})
}

func TestGetRespMsg(t *testing.T) {
	Convey("getRespMsg should return a copy of supplied message as a response", t, func() {
		var req types.Message
		req.ID = 1010
		req.Channel = "hit-channel"
		req.Text = "this is the most liked message"
		req.Type = "unit-test"
		req.User = "SUPER-HIP"

		actual := getRespMsg(req)

		So(actual.Text, ShouldNotResemble, req.Text)
		So(actual.ID, ShouldEqual, req.ID)
		So(actual.Channel, ShouldResemble, req.Channel)
		So(actual.Type, ShouldResemble, req.Type)
		So(actual.User, ShouldResemble, req.User)
	})
}
