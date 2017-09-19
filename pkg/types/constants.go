package types

// Constants
const (
	SlackRtmURLFmt              = "https://slack.com/api/rtm.start?token=%s"
	SlackAPIServerURL           = "https://api.slack.com/"
	MessageType                 = "message"
	RequestKube2IamBotReq       = "!requestKube2iam"
	RequestKube2IamBotReqFormat = "```!requestKube2iam <namespace> <roleArn> <cluster>```"
	ApproveKube2IamBotReq       = "!approveKube2iam"
	ApproveKube2IamBotReqFormat = "```!approveKube2iam <namespace> <roleArn> <cluster>```"
	Kube2IamBotReqLength        = 4
	AWSMetaDataServerAccRsrcEp  = "dev_read/accounts?AccountNumber"
	ADSecurityGroupEndPoint     = "dev_read/teams?ID"
	AccountNumberIndexInRoleArn = 4
	WhichKubectl                = "which kubectl"
)
