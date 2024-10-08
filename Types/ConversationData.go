package Types

// ConversationProvidableData struct contains providable metadata of a conversation.
//
// Info: Providable data is data that the API is able to serve to the user, often transformed, so that it does not expose database schema or unnecessary fields.
type ConversationProvidableData struct {
	ConversationId string                             `json:"conversationId"`
	Members        []ConversationMemberProvidableData `json:"members"`
	Name           string                             `json:"name"`
	Type           string                             `json:"type"`
}

// ConversationMemberProvidableData struct contains providable metadata of a specific conversation member.
//
// Info: Providable data is data that the API is able to serve to the user, often transformed, so that it does not expose database schema or unnecessary fields.
type ConversationMemberProvidableData struct {
	AccountId  string `json:"accountId"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}
