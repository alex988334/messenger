package business_logic

import (
	"encoding/json"
	"errors"

	cons "github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/data/chat"
	"github.com/alex988334/messenger/pkg/messenger/data/message"
	"github.com/alex988334/messenger/pkg/messenger/data/user"
)

type ResponseServer struct {
	Status        StatusRequest
	Chat          []chat.Chat
	ChatUser      []chat.ChatUser
	Message       []message.Message
	MessageStatus []message.MessageStatus
	User          []user.User
	BlackList     []user.BlackList
	UserPhone     []user.UserPhone
}

type StatusRequest struct {
	Operation int
	Status    int
	Message   string
}

/*
func encodeModel[T StatusRequest | *data.IModel | map[string][]byte](model T) []byte {

		d, _ := json.Marshal(model)
		return d
	}
*/
func (b *BusinessHandler) GenerateResponse() ([]byte, []int, error) {

	/*m := map[string]any{}*/

	r := ResponseServer{}

	if b.err != nil {
		r.Status = *generateStatusModel(b.err, b.opearation, "")
		enc, _ := json.Marshal(r)
		return enc, nil, errors.New("")
	} else {
		r.Status = *generateStatusModel(nil, b.opearation, "")
	}

	for k, v := range b.response {

		switch k {
		case data.MODEL_CHAT:
			models := make([]chat.Chat, len(v))
			for key, val := range v {
				models[key] = *val.(*chat.Chat)
			}
			r.Chat = models

		case data.MODEL_CHAT_USER:
			models := make([]chat.ChatUser, len(v))
			for key, val := range v {
				models[key] = *val.(*chat.ChatUser)
			}
			r.ChatUser = models

		case data.MODEL_MESSAGE:
			models := make([]message.Message, len(v))
			for key, val := range v {
				models[key] = *val.(*message.Message)
			}
			r.Message = models

		case data.MODEL_STATUS_MESSAGE:
			models := make([]message.MessageStatus, len(v))
			for key, val := range v {
				models[key] = *val.(*message.MessageStatus)
			}
			r.MessageStatus = models

		case data.MODEL_USER:
			models := make([]user.User, len(v))
			for key, val := range v {
				models[key] = *val.(*user.User)
			}
			r.User = models

		case data.MODEL_USER_PHONE:
			models := make([]user.UserPhone, len(v))
			for key, val := range v {
				models[key] = *val.(*user.UserPhone)
			}
			r.UserPhone = models

		case data.MODEL_BLACK_LIST:
			models := make([]user.BlackList, len(v))
			for key, val := range v {
				models[key] = *val.(*user.BlackList)
			}
			r.BlackList = models
		}
	}

	enc, _ := json.Marshal(r)
	return enc, b.GetSendChatsId(), nil
}

func generateStatusModel(err error, operation int, message string) *StatusRequest {

	if err != nil {
		return &StatusRequest{
			Operation: operation,
			Status:    cons.STATUS_ERROR,
			Message:   err.Error(),
		}
	}

	return &StatusRequest{
		Operation: operation,
		Status:    cons.STATUS_ACCEPT,
		Message:   message,
	}
}
