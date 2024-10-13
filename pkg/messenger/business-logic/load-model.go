package business_logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	cons "github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	chat "github.com/alex988334/messenger/pkg/messenger/data/chat"
	mess "github.com/alex988334/messenger/pkg/messenger/data/message"
	user "github.com/alex988334/messenger/pkg/messenger/data/user"
)

func getErrorLoad(nameModel string) error {
	return errors.New("ERROR! Request data was not loaded into the " + nameModel + "!")
}

func parseModel[T *StatusRequest | *data.IModel](data string, model T) bool {

	err := json.Unmarshal([]byte(data), model)
	if err != nil {
		fmt.Println("ERROR parseModel(): model:", model, ", error =>", err)
		return false
	} else {
		return true
	}
}

func LoadModels(requestData *map[string]string, clientData IClient) ([]data.IModel, error) {

	request := *requestData
	models := []data.IModel{}

	var st StatusRequest = StatusRequest{}

	if !parseModel(request[data.MODEL_STATUS], &st) {
		return nil, getErrorLoad("StatusRequest")
	}

	switch st.Operation {
	case cons.OP_STATUS_MESSAGE:
		ms := mess.NewMessageStatus()
		ms.SetOperation(cons.OP_STATUS_MESSAGE)

		var m data.IModel = ms
		if !parseModel(request[data.MODEL_STATUS_MESSAGE], &m) {
			return nil, getErrorLoad(ms.GetNameModel())
		}
		ms.SetUserId(clientData.GetId())
		models = append(models, ms)

	case cons.OP_NEW_MESSAGE:
		mes := mess.NewMessage()
		mes.SetOperation(cons.OP_NEW_MESSAGE)

		var m data.IModel = mes
		if !parseModel(request[data.MODEL_MESSAGE], &m) {
			return nil, getErrorLoad(mes.GetNameModel())
		}
		mes.SetAuthor(clientData.GetId())
		models = append(models, mes)

	case cons.OP_LIST_USERS:
		cu := chat.NewChatUser()
		cu.SetOperation(cons.OP_LIST_USERS)

		var m data.IModel = cu
		if !parseModel(request[data.MODEL_CHAT_USER], &m) {
			return nil, getErrorLoad(cu.GetNameModel())
		}
		cu.SetUser(clientData.GetId())
		models = append(models, cu)

	case cons.OP_CREATE_NEW_CHAT:
		ch := chat.NewChat()
		ch.SetOperation(cons.OP_CREATE_NEW_CHAT)

		var m data.IModel = ch
		if !parseModel(request[data.MODEL_CHAT], &m) {
			return nil, getErrorLoad(ch.GetNameModel())
		}
		ch.SetAuthor(clientData.GetId())
		models = append(models, ch)

	case cons.OP_WRITEN:
		chUs := chat.NewChatUser()
		chUs.SetOperation(cons.OP_WRITEN)

		var m data.IModel = chUs
		if !parseModel(request[data.MODEL_CHAT_USER], &m) {
			return nil, getErrorLoad(chUs.GetNameModel())
		}
		chUs.SetUser(clientData.GetId())
		models = append(models, chUs)

	case cons.OP_SYSTEM:
	case cons.OP_SEARCH_USER:
		u := user.NewUser()
		u.SetOperation(cons.OP_SEARCH_USER)

		var m data.IModel = u
		if !parseModel(request[data.MODEL_USER], &m) {
			ph := user.NewUserPhone()
			ph.SetOperation(cons.OP_SEARCH_USER)

			m = ph
			if parseModel(request[data.MODEL_USER_PHONE], &m) {
				models = append(models, ph)
			}
		} else {
			if u.Alias == "" {
				return nil, errors.New("ERROR! Alias is empty")
			}
			models = append(models, u)
		}

		if len(models) == 0 {
			return nil, errors.New("ERROR! Request data was not loaded")
		}
	case cons.OP_GET_CHATS:
		cu := chat.NewChatUser()
		cu.SetOperation(cons.OP_GET_CHATS)
		cu.User = clientData.GetId()
		models = append(models, cu)

	case cons.OP_LIST_PREVIOUS_MESSAGES:
		mes := mess.NewMessage()
		mes.SetOperation(cons.OP_LIST_PREVIOUS_MESSAGES)

		var m data.IModel = mes
		if !parseModel(request[data.MODEL_MESSAGE], &m) {
			return nil, getErrorLoad(mes.GetNameModel())
		}
		models = append(models, mes)

	case cons.OP_LIST_NEXT_MESSAGES:
		mes := mess.NewMessage()
		mes.SetOperation(cons.OP_LIST_NEXT_MESSAGES)

		var m data.IModel = mes
		if !parseModel(request[data.MODEL_MESSAGE], &m) {
			return nil, getErrorLoad(mes.GetNameModel())
		}
		models = append(models, mes)

	case cons.OP_EXIT_CHAT:
		chUs := chat.NewChatUser()
		chUs.SetOperation(cons.OP_EXIT_CHAT)

		var m data.IModel = chUs
		if !parseModel(request[data.MODEL_CHAT_USER], &m) {
			return nil, getErrorLoad(chUs.GetNameModel())
		}

		if chUs.Chat == 0 {
			return nil, getErrorLoad(chUs.GetNameModel())
		}

		models = append(models, chUs)

	case cons.OP_REMOVE_USER:
		chUs := chat.NewChatUser()
		chUs.SetOperation(cons.OP_REMOVE_USER)

		var m data.IModel = chUs
		if !parseModel(request[data.MODEL_CHAT_USER], &m) {
			return nil, getErrorLoad(chUs.GetNameModel())
		}

		if chUs.User == 0 || chUs.Chat == 0 {
			return nil, getErrorLoad(chUs.GetNameModel())
		}
		models = append(models, chUs)

	case cons.OP_ADD_USER:
		var chUs data.IModel = chat.NewChatUser()
		if !parseModel(request[data.MODEL_CHAT_USER], &chUs) {
			return nil, getErrorLoad(chUs.GetNameModel())
		}

		chUs.SetOperation(cons.OP_ADD_USER)
		models = append(models, chUs)

	case cons.OP_REMOVE_CHAT:
		ch := chat.NewChat()
		ch.SetOperation(cons.OP_REMOVE_CHAT)

		var m data.IModel = ch
		if !parseModel(request[data.MODEL_CHAT], &m) {
			return nil, getErrorLoad(ch.GetNameModel())
		}
		ch.SetAuthor(clientData.GetId())
		models = append(models, ch)

	case cons.OP_BLOCK_USERS:
		bl := user.NewBlackList()
		bl.SetOperation(cons.OP_BLOCK_USERS)

		var m data.IModel = bl
		if !parseModel(request[data.MODEL_BLACK_LIST], &m) {
			return nil, getErrorLoad(bl.GetNameModel())
		}
		bl.SetUser(clientData.GetId())
		models = append(models, bl)

	case cons.OP_UNLOOCK_USERS:
		bl := user.NewBlackList()
		bl.SetOperation(cons.OP_UNLOOCK_USERS)

		var m data.IModel = bl
		if !parseModel(request[data.MODEL_BLACK_LIST], &m) {
			return nil, getErrorLoad(bl.GetNameModel())
		}
		bl.SetUser(clientData.GetId())
		models = append(models, bl)

	case cons.OP_BLACK_LIST_USERS:
		bl := user.NewBlackList()
		bl.SetOperation(cons.OP_BLACK_LIST_USERS)

		var m data.IModel = bl
		if !parseModel(request[data.MODEL_BLACK_LIST], &m) {
			return nil, getErrorLoad(bl.GetNameModel())
		}
		bl.SetUser(clientData.GetId())
		models = append(models, bl)

	case cons.OP_GET_FILE:
	case cons.OP_MY_DATA:
		u := user.NewUser()
		u.SetOperation(cons.OP_MY_DATA)
		u.Id = clientData.GetId()
		models = append(models, u)

	case cons.OP_AUTORIZATE:
		u := user.NewUser()
		u.SetOperation(cons.OP_AUTORIZATE)

		var m data.IModel = u
		if !parseModel(request[data.MODEL_USER], &m) {
			return nil, getErrorLoad(u.GetNameModel())
		}
		models = append(models, u)
	case cons.OP_REGISTRATION:
		u := user.NewUser()
		u.SetOperation(cons.OP_REGISTRATION)

		var m data.IModel = u
		if !parseModel(request[data.MODEL_USER], &m) {
			return nil, getErrorLoad(u.GetNameModel())
		}

		if u.Login == "" || u.Alias == "" || u.Email == "" || u.PassHash == "" {
			return nil, getErrorLoad(u.GetNameModel())
		}
		models = append(models, u)
	default:
		return nil, errors.New("ERROR! Unsupported type operation! operation => " + strconv.Itoa(st.Operation))
	}

	return models, nil
}

/*
func operationBlockUnblockUser(operation int, authorId int, usersId *[]int) []data.IModel {

	usId := *usersId
	models := make([]data.IModel, len(usId))

	for i := 0; i < len(usId); i++ {
		bl := user.NewBlackList()
		bl.SetOperation(operation)
		bl.SetUser(authorId)
		bl.SetBlockedUser(usId[i])
		models[i] = bl
	}

	return models
}

func operationChatUser(operation int, chatId int, authorId int, usersId *[]int) []data.IModel {

	var users []int = *usersId
	var models []data.IModel = make([]data.IModel, len(users)+1)

	us := user.NewUser()
	us.SetOperation(operation)
	us.SetId(authorId)
	models[0] = us

	for i := 1; i < len(models); i++ {
		cu := chat.NewChatUser()
		cu.SetOperation(operation)
		cu.SetChat(chatId)
		cu.SetUser(users[i-1])
		models[i] = *cu
	}

	return models
}

/***
func LoadModels(requestData *map[string]interface{}) []data.IModel {

	request := *requestData
	var models []data.IModel = make([]data.IModel, 1)
	var action int = int(request[cons.KEY_ACTION].(float64))

	switch action {
	case cons.OP_STATUS_MESSAGE:
		ms := mess.NewMessageStatus()
		ms.SetOperation(cons.OP_STATUS_MESSAGE)
		ms.SetUserId(request[cons.KEY_AUTHOR].(int))
		ms.SetMessageId(request[cons.KEY_ID].(int64))
		ms.SetStatusMessage(request[cons.KEY_STATUS].(string))
		models = append(models, ms)

	case cons.OP_INPUT_MESSAGE:
	case cons.OP_NEW_MESSAGE:
		mes := mess.NewMessage()
		mes.SetOperation(cons.OP_NEW_MESSAGE)
		mes.SetMessage(request[cons.KEY_MESSAGE].(string))
		mes.SetAuthor(request[cons.KEY_AUTHOR].(int))
		mes.SetChatId(request[cons.KEY_CHAT].(int))
		models = append(models, mes)

	case cons.OP_SET_USER_NAME:
	case cons.OP_LIST_USERS:
		cu := chat.NewChatUser()
		cu.SetOperation(cons.OP_LIST_USERS)
		cu.SetChat(request[cons.KEY_CHAT].(int))
		models = append(models, cu)

	case cons.OP_CREATE_NEW_CHAT:
		ch := chat.NewChat()
		ch.SetOperation(cons.OP_CREATE_NEW_CHAT)
		ch.SetName(request[cons.KEY_CHAT_NAME].(string))
		ch.SetAuthor(request[cons.KEY_AUTHOR].(int))
		models = append(models, ch)

		var usersId []int = request[cons.KEY_USERS].([]int)
		for _, v := range usersId {
			us := user.NewUser()
			us.SetOperation(cons.OP_CREATE_NEW_CHAT)
			us.SetId(v)
			models = append(models, us)
		}

	case cons.OP_WRITEN:
		us := user.NewUser()
		us.SetOperation(cons.OP_WRITEN)
		us.SetId(request[cons.KEY_AUTHOR].(int))
		models = append(models, us)

		chUs := chat.NewChatUser()
		chUs.SetChat(request[cons.KEY_CHAT].(int))
		chUs.SetOperation(cons.OP_WRITEN)
		models = append(models, chUs)

	case cons.OP_SYSTEM:
	case cons.OP_ERROR_NAME:
	case cons.OP_SEARCH_USER:
	case cons.OP_GET_CHATS:
		cu := chat.NewChatUser()
		cu.SetOperation(cons.OP_GET_CHATS)
		cu.SetUser(request[cons.KEY_AUTHOR].(int))
		models = append(models, cu)

	case cons.OP_GET_HISTORY_MESSAGE:
		chUs := chat.NewChatUser()
		chUs.SetOperation(cons.OP_GET_HISTORY_MESSAGE)
		chUs.SetUser(request[cons.KEY_AUTHOR].(int))
		chUs.SetChat(request[cons.KEY_CHAT].(int))
		models = append(models, chUs)

		mess := mess.NewMessage()
		mess.SetOperation(cons.OP_GET_HISTORY_MESSAGE)
		mess.SetChatId(request[cons.KEY_CHAT].(int))
		mess.SetId(request[cons.KEY_ID].(int64))
		models = append(models, mess)

	case cons.OP_EXIT_CHAT:
		ch := chat.NewChat()
		ch.SetOperation(cons.OP_EXIT_CHAT)
		ch.SetId(request[cons.KEY_CHAT].(int))
		models = append(models, ch)

		chUs := chat.NewChatUser()
		chUs.SetOperation(cons.OP_EXIT_CHAT)
		chUs.SetUser(request[cons.KEY_USERS].(int))
		chUs.SetChat(request[cons.KEY_CHAT].(int))
		models = append(models, chUs)

	case cons.OP_REMOVE_USER:
		usersId := (request[cons.KEY_USERS].([]int))
		models = operationChatUser(cons.OP_REMOVE_USER, request[cons.KEY_CHAT].(int),
			request[cons.KEY_AUTHOR].(int), &usersId)

	case cons.OP_ADD_USER:
		usersId := (request[cons.KEY_USERS].([]int))
		models = operationChatUser(cons.OP_ADD_USER, request[cons.KEY_CHAT].(int),
			request[cons.KEY_AUTHOR].(int), &usersId)

	case cons.OP_REMOVE_CHAT:
		us := user.NewUser()
		us.SetOperation(cons.OP_REMOVE_CHAT)
		us.SetId(request[cons.KEY_AUTHOR].(int))
		models = append(models, us)

		ch := chat.NewChat()
		ch.SetOperation(cons.OP_REMOVE_CHAT)
		ch.SetId(request[cons.KEY_ID].(int))
		models = append(models, ch)

	case cons.OP_BLOCK_USERS:
		usersId := request[cons.KEY_USERS].([]int)
		models = operationBlockUnblockUser(cons.OP_BLOCK_USERS, request[cons.KEY_AUTHOR].(int),
			&usersId)

	case cons.OP_UNLOOCK_USERS:
		usersId := request[cons.KEY_USERS].([]int)
		models = operationBlockUnblockUser(cons.OP_UNLOOCK_USERS, request[cons.KEY_AUTHOR].(int),
			&usersId)

	case cons.OP_BLACK_LIST_USERS:
		bl := user.NewBlackList()
		bl.SetOperation(cons.OP_BLACK_LIST_USERS)
		bl.SetUser(request[cons.KEY_AUTHOR].(int))
		models = append(models, bl)

	case cons.OP_GET_FILE:
	case cons.OP_HAVE_MESSAGE:
	case cons.OP_MY_DATA:
	default:
	}

	return models
}
****/
