package utils

import (
	"testing"

	"github.com/ashish-amarnath/slackbots/pkg/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStringifyMessage(t *testing.T) {
	Convey("StringifyMessage", t, func() {
		Convey("Should return expected string version of a message object", func() {
			var testMsg types.Message
			testMsg.Channel = "utchannel"
			testMsg.ID = 12
			testMsg.Text = "unit test message"
			testMsg.Type = "message"
			testMsg.User = "unit-tester"
			expectedString := "[ID=12, Type=message, Text=unit test message, Channel=utchannel, User=unit-tester]"
			actualString := StringifyMessage(testMsg)
			So(actualString, ShouldEqual, expectedString)
		})
		Convey("Should use default values", func() {
			var testMsg types.Message
			expectedString := "[ID=0, Type=, Text=, Channel=, User=]"
			actualString := StringifyMessage(testMsg)
			So(actualString, ShouldEqual, expectedString)
		})
	})
}

func TestRunBashCmd(t *testing.T) {
	Convey("RunBashCmd", t, func() {
		Convey("should successfully run invoked with a valid bash command string", func() {
			validBashCmd := "echo this should pass"
			expected := "this should pass"
			actual, err := RunBashCmd(validBashCmd)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})
		Convey("should fail when invoked with an invalid bash command string", func() {
			invalidBashCmd := "whatwasithinking"
			expected := ""
			actual, err := RunBashCmd(invalidBashCmd)
			So(err, ShouldNotBeNil)
			So(actual, ShouldResemble, expected)
		})
	})
}

func TestStringifySlackUser(t *testing.T) {
	Convey("StringifySlackUser", t, func() {
		Convey("Should return expected string version of a SlackUser object", func() {
			var to types.SlackUser
			to.ID = "UCRAY7US3R"
			to.Profile.FirstName = "John"
			to.Profile.LastName = "Doe"
			to.Profile.Email = "john.doe@johndoe.com"
			expected := "[ID=UCRAY7US3R, FirstName=John, LastName=Doe, Email=john.doe@johndoe.com]"
			actual := StringifySlackUser(to)
			So(actual, ShouldEqual, expected)
		})
		Convey("Should use default values", func() {
			var to types.SlackUser
			expected := "[ID=, FirstName=, LastName=, Email=]"
			actual := StringifySlackUser(to)
			So(actual, ShouldEqual, expected)
		})
	})
}

func TestStringifyADUser(t *testing.T) {
	Convey("StringifyADUser", t, func() {
		Convey("Should return expected string version of a SlackUser object", func() {
			var to types.ADUser
			to.LanID = "XYZ7"
			to.FirstName = "John"
			to.LastName = "Doe"
			to.Email = "john.doe@johndoe.com"
			expected := "[LanID=XYZ7, FirstName=John, LastName=Doe, Email=john.doe@johndoe.com]"
			actual := StringifyADUser(to)
			So(actual, ShouldEqual, expected)
		})
		Convey("Should use default values", func() {
			var to types.ADUser
			expected := "[LanID=, FirstName=, LastName=, Email=]"
			actual := StringifyADUser(to)
			So(actual, ShouldEqual, expected)
		})
	})
}

func TestGetBotReqType(t *testing.T) {
	Convey("GetBotReqType should parse and return the type of request", t, func() {
		msg := "@superbot !doSomethingAwesome foo arn:aws:iam::123456789012:role/superawesome-powerful-Role3 hydrogen"
		expected := "!doSomethingAwesome"
		actual := GetBotReqType(msg)
		So(actual, ShouldResemble, expected)
	})
}

func TestGetBotReqParams(t *testing.T) {
	Convey("GetBotReqParams should construct a BotReqParam struct for the supplied params", t, func() {
		var expected types.BotReqParams
		expected.ADGroupLookupURL = "https://adGrpLkp/api/v1/usr/get"
		expected.ADUserLookupURL = "https://adUsrLkp/api/v1/usr/get"
		expected.AWSAPIKey = "masterkey"
		expected.AWSMetadataServerURL = "https://jibberish.execute-api.us-west-81.amazonaws.com"
		expected.KubeConfig = "/User/craycrayuser/.kube/config"
		expected.Message = "@superbot !doSomethingAwesome foo arn:aws:iam::123456789012:role/superawesome-powerful-Role3 hydrogen"
		expected.SlackUser = "UCRAY7Q"

		actual := GetBotReqParams(expected.ADGroupLookupURL, expected.ADUserLookupURL, expected.AWSMetadataServerURL, expected.AWSAPIKey, expected.KubeConfig, expected.Message, expected.SlackUser)
		So(actual, ShouldResemble, expected)
	})
}

func TestStringifyBotReqParams(t *testing.T) {
	Convey("StringifyBotReqParams should return a right representation of a BotReqParams struct", t, func() {
		var to types.BotReqParams
		to.ADGroupLookupURL = "https://adGrpLkp/api/v1/usr/get"
		to.ADUserLookupURL = "https://adUsrLkp/api/v1/usr/get"
		to.AWSAPIKey = "blahziblahziblah"
		to.AWSMetadataServerURL = "https://jibberish.execute-api.us-west-81.amazonaws.com"
		to.KubeConfig = "/User/craycrayuser/.kube/config"
		to.Message = "@superbot !doSomethingAwesome foo arn:aws:iam::123456789012:role/superawesome-powerful-Role3 hydrogen"
		to.SlackUser = "UCRAY7Q"

		expected := `[ADGroupLkpURL=[https://adGrpLkp/api/v1/usr/get], ADUserLkpUrl=[https://adUsrLkp/api/v1/usr/get], MetadataServerURL=[https://jibberish.execute-api.us-west-81.amazonaws.com], MetadataServerAPIKey=[blahziblahziblah], kubeConfig=[/User/craycrayuser/.kube/config], Message=[@superbot !doSomethingAwesome foo arn:aws:iam::123456789012:role/superawesome-powerful-Role3 hydrogen], SlackUser=[UCRAY7Q]`
		actual := StringifyBotReqParams(to)
		So(actual, ShouldResemble, expected)
	})
}
