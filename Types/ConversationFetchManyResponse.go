package Types

type ConversationFetchManyResponse struct {
	Conversations []ConversationProvidableData `json:"conversations"`
	HttpResponse
}
