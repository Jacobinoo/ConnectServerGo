package CDModels

import "github.com/scylladb/gocqlx/v3/table"

var messagesByUserMetadata = table.Metadata{
	Name:    "messages_by_user",
	Columns: []string{"user_id", "conversation_id", "message_id", "author_id", "content"},
	PartKey: []string{"user_id", "conversation_id"},
	SortKey: []string{"message_id"},
}

var MessagesByUserTable = table.New(messagesByUserMetadata)

// Note: A field will not be persisted by adding the `db:"-"` tag or making it unexported.

type Message struct {
	UserId         int64
	ConversationId int64
	MessageId      int64
	AuthorId       int64
	Content        string
}

// Tutorial

//Loading a single row into struct

// message := CoreDataModels.Message{
// }
// q := session.Query(CoreDataModels.MessagesByUserTable.Get()).BindStruct(message)
// if err := q.GetRelease(&message); err != nil {
// 	log.Fatal(err)
// }
// log.Print(message)

//Insert a row from struct

// message := CoreDataModels.Message{
// 	UserId:         6,
// 	ConversationId: 1,
// 	MessageId:      1,
// 	AuthorId:       1,
// 	Content:        "Elo",
// }
// q := session.Query(CoreDataModels.MessagesByUserTable.Insert()).BindStruct(message)
// if err := q.ExecRelease(); err != nil {
// 	log.Fatal(err)
// }

//Load all rows to slice

// var messages []Message
// q := session.Query(personTable.Select()).BindMap(qb.M{"user_id": 1, "conversation_id": 1})
// if err := q.SelectRelease(&messages); err != nil {
// 	log.Fatal(err)
// }
// log.Print(messages)
