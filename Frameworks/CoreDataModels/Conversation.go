package CDModels

import "github.com/scylladb/gocqlx/v3/table"

var conversationsByUserMetadata = table.Metadata{
	Name:    "conversations_by_user",
	Columns: []string{"conversation_id", "user_id", "members", "name", "type"},
	PartKey: []string{"user_id"},
	SortKey: []string{"conversation_id"},
}
var conversationsMetadata = table.Metadata{
	Name:    "conversations",
	Columns: []string{"conversation_id", "members", "name", "type"},
	PartKey: []string{"conversation_id"},
	SortKey: []string{},
}

var ConversationsByUserTable = table.New(conversationsByUserMetadata)
var ConversationsTable = table.New(conversationsMetadata)

// Note: A field will not be persisted by adding the `db:"-"` tag or making it unexported.

type ConversationOfUser struct {
	ConversationId int64
	UserId         string
	Members        []string
	Name           string
	Type           string
}
type Conversation struct {
	ConversationId int64
	Members        []int64
	Name           string
	Type           string
}
