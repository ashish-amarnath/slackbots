package cmd

import (
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

func TestGetAWSAccountOwnerID(t *testing.T) {
	Convey("getAWSAccountOwnerID should return error when unable to process request sucessfully", t, func() {
		glog.Errorf("Expected Error:\n")
		actual, err := getAWSAccountOwnerID("https://example.com", "open-key", "123456789012")
		So(actual, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
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
					"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Namespace\",\"metadata\":{\"annotations\":{\"contact-email\":\"unit@test.com\",\"cost-center\":\"\",\"kube2iam.beta.nordstrom.net/allowed-roles\":\"[\\\"arn:aws:iam::123456789012:role/foo/k8s/fooS3AndKmsStack-fooFlinkForS3AndKMS-9336B1OLABTR\\\"]\\n\",\"kubernetes.io/change-cause\":\"kubectl edit ns foo --user=athens_sudo\",\"slack-channel-events\":\"\",\"slack-channel-urgent\":\"\",\"slack-channel-users\":\"#foo\"},\"creationTimestamp\":null,\"name\":\"foo\",\"namespace\":\"\",\"selfLink\":\"/api/v1/namespacesfoo\"},\"spec\":{\"finalizers\":[\"kubernetes\"]},\"status\":{\"phase\":\"Active\"}}\n",
					"kubernetes.io/change-cause": "kubectl edit ns foo --context=cluster --user=cluster_sudo"
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
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Namespace\",\"metadata\":{\"annotations\":{\"contact-email\":\"techdsbiadtlkdvops@nordstrom.com\",\"cost-center\":\"\",\"kube2iam.beta.nordstrom.net/allowed-roles\":\"[\\\"arn:aws:iam::123456789012:role/foo/k8s/fooS3AndKmsStack-fooFlinkForS3AndKMS-9336B1OLABTR\\\"]\\n\",\"kubernetes.io/change-cause\":\"kubectl edit ns foo --user=athens_sudo\",\"slack-channel-events\":\"\",\"slack-channel-urgent\":\"\",\"slack-channel-users\":\"#foo\"},\"creationTimestamp\":null,\"name\":\"foo\",\"namespace\":\"\",\"selfLink\":\"/api/v1/namespacesfoo\"},\"spec\":{\"finalizers\":[\"kubernetes\"]},\"status\":{\"phase\":\"Active\"}}\n",
						"kubernetes.io/change-cause": "kubectl edit ns foo --context=cluster --user=cluster_sudo",
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
