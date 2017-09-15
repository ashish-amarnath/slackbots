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
