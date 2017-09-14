package types

import "time"

// ResponseRtmStart represents rtm start message
type ResponseRtmStart struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	URL   string `json:"url"`
	Bot   BotID  `json:"self"`
}

// BotID represents the object storing the user ID
type BotID struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Message represents messages written to and read from web socket.
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`
}

//AccNumRespMsg represents the response from the accountOwnerIDRequest endpoint
type AccNumRespMsg struct {
	Data []struct {
		EnvironmentID     int    `json:"EnvironmentId"`
		OwnerTeamID       int    `json:"OwnerTeamId"`
		GUID              int    `json:"GUID"`
		AccountNumber     string `json:"AccountNumber"`
		AccountName       string `json:"AccountName"`
		RequesterPersonID int    `json:"RequesterPersonId"`
		ClaimRule         string `json:"ClaimRule"`
		Size              string `json:"Size"`
	} `json:"data"`
}

// AdSecurityGroupResp represents response from the adSecurityGroupRequest endpoint
type AdSecurityGroupResp struct {
	Data []struct {
		VP              int         `json:"VP"`
		Name            string      `json:"Name"`
		Tags            interface{} `json:"Tags"`
		OrgName         string      `json:"OrgName"`
		ADSecurityGroup string      `json:"ADSecurityGroup"`
		CreateDate      time.Time   `json:"CreateDate"`
		ID              int         `json:"ID"`
		Director        int         `json:"Director"`
		Disabled        interface{} `json:"Disabled"`
		CostCenter      int         `json:"CostCenter"`
		OnCall          interface{} `json:"OnCall"`
		EmailDistList   string      `json:"EmailDistList"`
		UpdateDate      time.Time   `json:"UpdateDate"`
	} `json:"data"`
}

// ADGroupMemberListResp represents response from looking up members of an AD group
type ADGroupMemberListResp struct {
	Dn          string      `json:"dn"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Info        interface{} `json:"info"`
	Email       string      `json:"email"`
	Type        string      `json:"type"`
	Created     time.Time   `json:"created"`
	Updated     time.Time   `json:"updated"`
	Members     struct {
		Groups []interface{} `json:"groups"`
		Users  []string      `json:"users"`
	} `json:"members"`
	ManagedBy struct {
		Group interface{} `json:"group"`
		User  string      `json:"user"`
	} `json:"managedBy"`
	Groups []string `json:"groups"`
}
