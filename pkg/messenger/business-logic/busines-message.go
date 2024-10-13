package business_logic

import (
	"errors"
	"strconv"

	"github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/data/chat"
	mess "github.com/alex988334/messenger/pkg/messenger/data/message"
	"github.com/alex988334/messenger/pkg/messenger/data/user"
)

func (b *BusinessHandler) processListMessge() {

	m := b.models[0].(*mess.Message)

	if m.Id == 0 {
		mFind := mess.NewMessage()
		mFind.SetActionType(data.ACTION_FIND)
		mFind.SetLoadFields("Id")
		mFind.AddSortModels(data.DIRECTION_DESC, "Id")
		mFind.SetLimitCountModel(1)

		b.models = []data.IModel{mFind}
		if !b.runAdapter() {
			return
		}

		rez := b.adapter.GetModelRezult()
		if len(rez) > 0 {
			m.Id = rez[0][data.MODEL_MESSAGE].(*mess.Message).Id
		}
		b.models = []data.IModel{m}
	}

	directionSort := data.DIRECTION_DESC
	operator := data.LESS_EQUAL
	if b.GetOperation() == constants.OP_LIST_NEXT_MESSAGES {
		directionSort = data.DIRECTION_ASC
		operator = data.MORE_EQUAL
	}

	m.SetActionType(data.ACTION_FIND)
	m.SetLoadFields("Id", "ChatId", "Author", "ParrentMessage", "Message", "FileUrl", "Date", "Time")
	m.AddSortModels(directionSort, "Id")
	m.SetLimitCountModel(constants.LIST_LIMIT_MESSAGES)
	m.AddConditionField("Id", operator)

	chUs := chat.NewChatUser()
	chUs.SetActionType(data.ACTION_FIND)
	chUs.AddLink(data.NewLink(m.GetNameModel(), "ChatId", chUs.GetNameModel(), "Chat", data.LINK_WEIGHT_EQUILIBRIUM))
	chUs.User = b.client.GetId()

	u := user.NewUser()
	u.SetActionType(data.ACTION_FIND)
	u.SetLoadFields("Id", "Alias")
	u.AddLink(data.NewLink(m.GetNameModel(), "Author", u.GetNameModel(), "Id", data.LINK_WEIGHT_EQUILIBRIUM))

	m.SetChatUsers(&[]data.IModel{chUs})
	m.SetUsers(&[]data.IModel{u})
	if !b.runAdapter() {
		return
	}

	b.loadModelsFromAdapterRezult()

	messageId := make([]int64, 0, len(b.response[data.MODEL_MESSAGE]))
	pIds := make([]int64, 0, len(b.response[data.MODEL_MESSAGE]))
	for _, v := range b.response[data.MODEL_MESSAGE] {

		m := v.(*mess.Message)
		messageId = append(messageId, m.Id)

		if pId := m.ParrentMessage; pId > 0 {
			pIds = append(pIds, pId)
		}
	}

	/*if len(pIds) == 0 {
		return
	}*/

	parentsIdModels := make([]data.IModel, 0, len(pIds))
	for _, pId := range pIds {

		m = mess.NewMessage()
		m.SetActionType(data.ACTION_FIND)
		m.Id = pId
		m.SetLoadFields("Id", "ChatId", "Author", "ParrentMessage", "Message", "FileUrl", "Date", "Time")
		initArrayFieldCondition(m, "Id", data.IN)
		m.AddLink(data.NewLink(u.GetNameModel(), "Id", m.GetNameModel(), "Author", data.LINK_WEIGHT_EQUILIBRIUM))

		parentsIdModels = append(parentsIdModels, m)
	}

	if len(parentsIdModels) > 0 {
		u.ResetLinks()
		u.SetMessages(&parentsIdModels)

		b.models = []data.IModel{u}
		if !b.runAdapter() {
			return
		}

		b.loadModelsFromAdapterRezult()
	}

	if len(messageId) > 0 {
		m := mess.NewMessageStatus()
		m.SetActionType(data.ACTION_FIND)
		m.SetLoadFields("MessageId", "Status")
		m.AddGroupModels("MessageId", "Status")
		m.AddSortModels(data.DIRECTION_ASC, "MessageId")

		msModels := make([]data.IModel, 0, len(messageId))
		for _, id := range messageId {

			ms := mess.NewMessageStatus()
			ms.SetActionType(data.ACTION_FIND)
			ms.SetArrayFields("MessageId")
			ms.AddConditionField("MessageId", data.IN)
			ms.MessageId = id

			msModels = append(msModels, ms)
		}

		m.SetStatus(&msModels)

		b.models = []data.IModel{m}
		if !b.runAdapter() {
			return
		}

		b.loadModelsFromAdapterRezult(m.GetNameModel())
	}
}

func (b *BusinessHandler) processStatusMessage() {

	ms := b.models[0].(*mess.MessageStatus)

	if ms.Status != constants.MESSAGE_READED && ms.Status != constants.MESSAGE_DELIVERED {
		b.err = errors.New("Not support message status! Status " + ms.Status)
		return
	}

	//	search all statuses before change rows in data base
	msFinded := mess.NewMessageStatus()
	msFinded.SetActionType(data.ACTION_FIND)
	msFinded.MessageId = ms.MessageId
	msFinded.SetLoadFields("MessageId", "UserId", "Status")

	b.models = []data.IModel{msFinded}
	if !b.runAdapter() {
		return
	}

	beforeСhanges := map[string]bool{}
	rez := b.adapter.GetModelRezult()

	for i := 0; i < len(rez); i++ {
		beforeСhanges[rez[i][ms.GetNameModel()].(*mess.MessageStatus).Status] = true
	}

	if len(rez) == 0 {
		b.err = errors.New("Message not found, messsge id: " + strconv.Itoa(int(ms.MessageId)))
		return
	}

	//	update row in data base with security finded row
	ms.SetActionType(data.ACTION_UPDATE)
	ms.Date = b.DateRegistration()
	ms.Time = b.TimeRegistration()
	//ms.AddConditionField("MessageId", data.LESS_EQUAL)
	//	ms.AddConditionField("UserId", data.EQUAL)

	m := mess.NewMessage()
	m.SetActionType(data.ACTION_UPDATE)
	m.AddLink(data.NewLink(ms.GetNameModel(), "MessageId", m.GetNameModel(), "Id", data.LINK_WEIGHT_EQUILIBRIUM))

	cu := chat.NewChatUser()
	cu.SetActionType(data.ACTION_UPDATE)
	cu.User = ms.UserId
	cu.AddLink(data.NewLink(m.GetNameModel(), "ChatId", cu.GetNameModel(), "Chat", data.LINK_WEIGHT_EQUILIBRIUM))

	m.SetChatUsers(&[]data.IModel{cu})
	ms.SetMessages(&[]data.IModel{m})

	b.models = []data.IModel{ms}
	if !b.runAdapter() {
		return
	}

	//	update all rows where message id less id in input data, upgrade all statuses for user to DELIVERED
	ms1 := mess.NewMessageStatus()
	ms1.SetActionType(data.ACTION_UPDATE)
	ms1.Status = constants.MESSAGE_DELIVERED
	ms1.UserId = ms.UserId
	ms1.AddConditionField("UserId", data.EQUAL)

	ms2 := mess.NewMessageStatus()
	ms2.SetActionType(data.ACTION_FIND)
	ms2.MessageId = ms.MessageId
	ms2.UserId = ms.UserId
	ms2.Status = constants.MESSAGE_CREATED
	ms2.SetSubField("MessageId")
	initConditionOperatorOfSubQuery(ms2, data.IN, "MessageId")
	ms2.AddConditionField("MessageId", data.LESS)
	ms2.AddConditionField("UserId", data.EQUAL)
	ms2.AddConditionField("Status", data.EQUAL)

	ms1.SetStatus(&[]data.IModel{ms2})
	b.models = []data.IModel{ms1}

	if !b.runAdapter() {
		return
	}

	//	update all rows where message id less id in input data, upgrade all statuses for user to READED
	if ms.Status == constants.MESSAGE_READED {
		ms1.Status = constants.MESSAGE_READED
		ms2.Status = constants.MESSAGE_DELIVERED

		if !b.runAdapter() {
			return
		}
	}

	//	search all statuses after change rows in data base
	b.models = []data.IModel{msFinded}
	if !b.runAdapter() {
		return
	}

	rez = b.adapter.GetModelRezult()
	if len(rez) == 0 {
		b.err = errors.New("Message not found, messsge id: " + strconv.Itoa(int(ms.MessageId)))
		return
	}
	afterChange := map[string]bool{}
	for i := 0; i < len(rez); i++ {
		afterChange[rez[i][ms.GetNameModel()].(*mess.MessageStatus).Status] = true
	}

	//	comparison statuses
	_, createB := beforeСhanges[constants.MESSAGE_CREATED]
	_, createA := afterChange[constants.MESSAGE_CREATED]
	_, deliverB := beforeСhanges[constants.MESSAGE_DELIVERED]
	_, deliverA := afterChange[constants.MESSAGE_DELIVERED]
	_, readedB := beforeСhanges[constants.MESSAGE_READED]
	_, readedA := afterChange[constants.MESSAGE_READED]

	if (createB && createA) || (deliverA && deliverB && !createB) ||
		(readedB && readedA && !createB && !deliverB) {
		//	exit if status stage not change
		return
	}

	change := ""
	if createB && !createA {
		change = constants.MESSAGE_DELIVERED
	}

	if deliverB && !deliverB {
		change = constants.MESSAGE_READED
	}

	if createB && !createA && !deliverA {
		change = constants.MESSAGE_READED
	}

	if change == "" {
		//	exit if status stage not change
		return
	} else {
		ms.Status = change
	}

	m = mess.NewMessage()
	m.SetActionType(data.ACTION_FIND)
	m.Id = ms.MessageId
	m.SetLoadFields("Id", "ChatId")
	b.models = []data.IModel{m}

	if !b.runAdapter() {
		return
	}
	ms.UserId = 0

	b.loadModelsFromAdapterRezult()
	b.AppendModelToResponse(ms)

	b.chatsId = []int{b.adapter.GetModelRezult()[0][m.GetNameModel()].(*mess.Message).ChatId}
}

func (b *BusinessHandler) processNewMessage() {

	m := b.models[0].(*mess.Message)

	m.SetActionType(data.ACTION_SAVE)
	m.Date = b.DateRegistration()
	m.Time = b.TimeRegistration()

	cu := chat.NewChatUser()
	cu.SetActionType(data.ACTION_FIND)
	cu.Chat = m.ChatId
	cu.User = m.Author
	initExistConditionOfSubQuery(cu, data.EXISTS, "Chat")

	m.SetChatUsers(&[]data.IModel{cu})
	if !b.runAdapter() {
		return
	}

	m.Id = b.adapter.GetIdOfLastInsertRow()
	m.SetChatUsers(&[]data.IModel{})
	b.AppendModelToResponse(m)
	b.chatsId = []int{m.ChatId}

	cu = chat.NewChatUser()
	cu.Chat = m.ChatId
	cu.SetLoadFields("Chat", "User")
	cu.SetActionType(data.ACTION_FIND)

	b.models = []data.IModel{cu}
	if !b.runAdapter() {
		return
	}

	rez := b.adapter.GetModelRezult()
	idUsers := make([]int, 0, len(rez))
	for _, v := range rez {
		idUsers = append(idUsers, v[cu.GetNameModel()].(*chat.ChatUser).User)
	}

	for _, v := range idUsers {

		ms := mess.NewMessageStatus()
		ms.MessageId = m.Id
		ms.Time = b.TimeRegistration()
		ms.Date = b.DateRegistration()
		ms.Status = constants.MESSAGE_CREATED
		ms.UserId = v
		ms.SetActionType(data.ACTION_SAVE)

		b.models = []data.IModel{ms}
		if !b.runAdapter() {
			return
		}
	}
}
