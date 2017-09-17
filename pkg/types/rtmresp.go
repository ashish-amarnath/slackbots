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
		Name            string `json:"Name"`
		OrgName         string `json:"OrgName"`
		ADSecurityGroup string `json:"ADSecurityGroup"`
		ID              int    `json:"ID"`
		Director        int    `json:"Director"`
		CostCenter      int    `json:"CostCenter"`
		EmailDistList   string `json:"EmailDistList"`
	} `json:"data"`
}

// ADGroupMemberListResp represents response from looking up members of an AD group
type ADGroupMemberListResp struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Type        string `json:"type"`
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

// KubernetesNamespace represents k8s namespace metadata
type KubernetesNamespace struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations struct {
			CloudTeamID                                 string `json:"cloud-team-id"`
			ContactEmail                                string `json:"contact-email"`
			CostCenter                                  string `json:"cost-center"`
			Kube2IamAllowedRoles                        string `json:"kube2iam.beta.nordstrom.net/allowed-roles"`
			KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
			KubernetesIoChangeCause                     string `json:"kubernetes.io/change-cause"`
		} `json:"annotations"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Name              string    `json:"name"`
		ResourceVersion   string    `json:"resourceVersion"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		Finalizers []string `json:"finalizers"`
	} `json:"spec"`
	Status struct {
		Phase string `json:"phase"`
	} `json:"status"`
}
