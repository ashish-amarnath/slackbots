package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"github.com/ashish-amarnath/slackbots/pkg/types"
	"github.com/golang/glog"
	"golang.org/x/net/websocket"
)

// ServerConn represents an RTM connection to slack
type ServerConn struct {
	URL    string
	conn   *websocket.Conn
	userID string
	msgID  uint64
}

func parseRtmStartResponse(respBytes []byte) (respJSON types.ResponseRtmStart, err error) {
	err = json.Unmarshal(respBytes, &respJSON)
	return
}

func getSlackRTMURL(token string) string {
	return fmt.Sprintf(types.SlackRtmUrlFmt, token)
}

func startSlackRTM(token string) (wsURL, userID string, err error) {
	if token == "" {
		err = fmt.Errorf("expected non-empty slackbot integration token, got [%s]", token)
		return
	}
	rtmURL := getSlackRTMURL(token)
	glog.V(3).Infof("Contacting slack rtm server at %s\n", rtmURL)

	resp, err := http.Get(rtmURL)
	if err != nil || resp.StatusCode != 200 {
		glog.Fatalf("Request to RTM server at %s failed with %d", rtmURL, resp.StatusCode)
		return
	}

	rBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		glog.Fatalf("Failed to read response body from RTM server at %s. err=%s\n", rtmURL, err)
		return
	}

	glog.V(8).Infof("RTM response body[\n %s\n]\n", rBody)
	var respJSON types.ResponseRtmStart
	respJSON, err = parseRtmStartResponse(rBody)
	if err != nil {
		glog.Fatalf("Slack RTM error:%s [Server=%s]\n", err, rtmURL)
		return
	}

	glog.V(3).Infoln("Successfully unmarshalled RTMStart Response.")
	glog.V(5).Infof("rtmStartResp.OK=%t\n", respJSON.Ok)
	glog.V(5).Infof("rmtStartResp.Error=%s\n", respJSON.Error)
	glog.V(5).Infof("rtmStartResp.URL=%s\n", respJSON.URL)
	glog.V(5).Infof("rtmStartResp.Self.ID=%s\n", respJSON.Bot.ID)
	glog.V(5).Infof("rtmStartResp.Self.Name=%s\n", respJSON.Bot.Name)

	if !respJSON.Ok {
		err = fmt.Errorf("Slack RTM error=%s", respJSON.Error)
		glog.Fatalln(err)
		return
	}

	wsURL = respJSON.URL
	userID = respJSON.Bot.ID
	glog.V(1).Infof("Initiated RTM session to slackbot at %s as user %s", wsURL, userID)
	return
}

func getSlackConn(webSockURL string) *websocket.Conn {
	conn, err := websocket.Dial(webSockURL, "", types.SlackAPIServerURL)
	if err != nil {
		glog.Fatalf("Failed to dial to URL=%s err=%s\n", webSockURL, err)
		return nil
	}
	glog.V(1).Infof("Successfully connected to slackbot at %s\n", webSockURL)
	return conn
}

// ReadMessage reads a message sent to the slackbot
func (s *ServerConn) ReadMessage() (m types.Message, err error) {
	err = websocket.JSON.Receive(s.conn, &m)
	return
}

func (s *ServerConn) getNextMessageID() uint64 {
	return atomic.AddUint64(&s.msgID, 1)
}

// SendMessage sends a message from the slack bot
func (s *ServerConn) SendMessage(m types.Message) error {
	m.ID = s.getNextMessageID()
	glog.V(4).Infof("Responding with message id:%d\n", m.ID)
	return websocket.JSON.Send(s.conn, m)

}

//NewSlackServerConn creates and returns a new connection to the slackbot identfied by the token
func NewSlackServerConn(token string) *ServerConn {
	rtmURL, user, err := startSlackRTM(token)
	if err != nil {
		glog.Fatalf("Failed to start slack RTM, err=%s\n", err)
	}
	wsConn := getSlackConn(rtmURL)
	return &ServerConn{
		URL:    rtmURL,
		userID: user,
		conn:   wsConn,
		msgID:  0,
	}
}
