package Conversation

import (
	"ConnectServer/Frameworks/CoreData"
	CDModels "ConnectServer/Frameworks/CoreDataModels"
	"ConnectServer/Types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/scylladb/gocqlx/qb"
)

func FetchManyConversationsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	encoder := *json.NewEncoder(writer)

	accountId := request.Context().Value("tokenSubject").(string)
	log.Println("Fetching conversations for UID:", accountId)

	var conversations []CDModels.ConversationOfUser
	query := CoreData.StorageServicesDatabaseSession.Query(CDModels.ConversationsByUserTable.Select()).BindMap(qb.M{"user_id": accountId})
	if err := query.SelectRelease(&conversations); err != nil {
		log.Fatal(err)
	}
	log.Println("Fetched conversations for UID:", accountId, conversations)

	conversationsMappedToProvidable := []Types.ConversationProvidableData{}

	for _, conversation := range conversations {
		membersOfConversationMappedToProvidable := []Types.ConversationMemberProvidableData{}
		for _, conversationMember := range conversation.Members {
			mappedMember := Types.ConversationMemberProvidableData{
				AccountId:  conversationMember,
				FirstName:  "Unknown TODO",
				LastName:   "Unknown TODO",
				MiddleName: "Unknown Optional TODO",
			}
			membersOfConversationMappedToProvidable = append(membersOfConversationMappedToProvidable, mappedMember)
		}
		mappedConversation := Types.ConversationProvidableData{
			ConversationId: strconv.FormatInt(conversation.ConversationId, 10),
			//TODO: Members
			Members: membersOfConversationMappedToProvidable,
			Name:    conversation.Name,
			Type:    conversation.Type,
		}
		conversationsMappedToProvidable = append(conversationsMappedToProvidable, mappedConversation)
	}

	response := Types.ConversationFetchManyResponse{
		Conversations: conversationsMappedToProvidable,
		HttpResponse: Types.HttpResponse{
			Success: true,
		},
	}
	encoder.Encode(response)
}
