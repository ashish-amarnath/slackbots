package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAccNumFromRoleArn(t *testing.T) {
	Convey("getAccNumFromRoleArn should parse out account number from a valid AWS role ARN", t, func() {
		testRoleArn := "arn:aws:iam::123456789012:role/foo/bar/foo-bar-mysuperawesomerole"
		expected := "123456789012"
		actual := getAccNumFromRoleArn(testRoleArn)
		So(actual, ShouldResemble, expected)
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
		Convey("Shoudl parse a valid JSON string into AdSecurityGroupResp", func() {
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
