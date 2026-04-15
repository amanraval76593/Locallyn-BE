package constants

type IdentityType string

const (
	IdentityPublic       IdentityType = "PUBLIC"
	IdentityPseudonymous IdentityType = "PSEUDONYMOUS"
	IdentityAnonymous    IdentityType = "ANONYMOUS"
)

type PostType string

const (
	PostTypeIncident  PostType = "INCIDENT"
	PostTypeBroadcast PostType = "BROADCAST"
)

type FeedbackType string

const (
	FeedbackHelpful    FeedbackType = "HELPFUL"
	FeedbackMisleading FeedbackType = "MISLEADING"
)
