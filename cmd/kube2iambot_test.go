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
