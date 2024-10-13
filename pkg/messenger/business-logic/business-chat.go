package business_logic

import (
	"errors"

	"github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/data/chat"
	"github.com/alex988334/messenger/pkg/messenger/data/message"
	"github.com/alex988334/messenger/pkg/messenger/data/user"
	"github.com/alex988334/messenger/pkg/messenger/functions"
)

func (b *BusinessHandler) processExitChat() {

	chUs := b.models[0].(*chat.ChatUser)

	if chUs.User == 0 {
		chUs.User = b.client.GetId()
	}

	ch := chat.NewChat()
	ch.SetActionType(data.ACTION_FIND)
	ch.SetLoadFields("Id", "Author", "Name", "CreateAt", "Status")
	ch.Id = chUs.Chat
	ch.Author = b.client.GetId()

	b.models = []data.IModel{ch}
	if !b.runAdapter() {
		return
	}

	isAuthor := len(b.adapter.GetModelRezult()) > 0

	if chUs.User != b.client.GetId() && !isAuthor {
		b.err = errors.New("ERROR! No access rights for operation")
		return
	}
	if chUs.User == b.client.GetId() && isAuthor {
		b.err = errors.New("ERROR! The assigned author cannot be the current author")
		return
	}
	if chUs.User != b.client.GetId() && isAuthor {
		chUs.SetActionType(data.ACTION_FIND)
		chUs.SetLoadFields("Chat", "User")

		b.models = []data.IModel{chUs}
		if !b.runAdapter() {
			return
		}

		if len(b.adapter.GetModelRezult()) == 0 {
			b.err = errors.New("ERORR! No finded user")
			return
		}

		ch.SetActionType(data.ACTION_UPDATE)
		ch.SetLoadFields()
		ch.Id = chUs.Chat
		ch.Author = chUs.User

		b.models = []data.IModel{ch}
		if !b.runAdapter() {
			return
		}
	}

	chUs.SetActionType(data.ACTION_DELETE)
	chUs.User = b.client.GetId()

	b.models = []data.IModel{chUs}
	if !b.runAdapter() {
		return
	}

	b.AppendModelToResponse(ch, chUs)
	b.generateListConnection([]int{chUs.User}, []int{chUs.Chat}, constants.CLIENT_DISCONNECT)
	b.chatsId = []int{ch.Id}
}

func (b *BusinessHandler) processRemoveChat() {

	ch := b.models[0].(*chat.Chat)

	var actionUserConnections *ActionUserConnections

	findUserChat := chat.NewChatUser()
	findUserChat.SetActionType(data.ACTION_FIND)
	findUserChat.SetLoadFields("User")
	findUserChat.Chat = ch.Id

	b.models = []data.IModel{findUserChat}
	if !b.runAdapter() {
		return
	}

	rez := b.adapter.GetModelRezult()

	if len(rez) > 0 {
		searchChats := make([]data.IModel, 0, len(rez)*7)
		findedId := make([]int, len(rez))

		for i := 0; i < len(rez); i++ {
			user := rez[i][data.MODEL_CHAT_USER].(*chat.ChatUser)
			user.SetArrayFields("User")
			user.AddConditionField("User", data.IN)
			searchChats = append(searchChats, user)
			findedId[i] = user.User
		}

		findUserChat.AddConditionField("Chat", data.NOT_EQUAL)
		findUserChat.SetLoadFields("Chat")
		findUserChat.SetChatUsers(&searchChats)
		if !b.runAdapter() {
			return
		}

		rez = b.adapter.GetModelRezult()
		chatsId := make([]int, 0, len(rez))

		for i := 0; i < len(rez); i++ {
			if id := rez[i][data.MODEL_CHAT_USER].(*chat.ChatUser).Chat; functions.FindInArray(chatsId, id) == -1 {
				chatsId = append(chatsId, id)
			}
		}

		actionUserConnections = &ActionUserConnections{
			UsersId:        findedId,
			ActionChats:    []int{ch.Id},
			ConnectedChats: chatsId,
			Operation:      constants.ALL_CLIENT_DISCONNECT,
		}
	}

	ch.SetActionType(data.ACTION_DELETE)

	b.models = []data.IModel{ch}
	if !b.runAdapter() {
		return
	}

	chUs := chat.NewChatUser()
	chUs.SetActionType(data.ACTION_DELETE)
	chUs.Chat = ch.Id

	b.models = []data.IModel{chUs}
	if !b.runAdapter() {
		return
	}

	ms := message.NewMessageStatus()
	ms.SetActionType(data.ACTION_DELETE)

	m := message.NewMessage()
	m.SetActionType(data.ACTION_FIND)
	m.SetLoadFields("Id")
	m.ChatId = ch.Id
	m.SetSubField("MessageId")
	m.SetSubOperator(data.IN)

	ms.SetMessages(&[]data.IModel{m})

	b.models = []data.IModel{ms}
	if !b.runAdapter() {
		return
	}

	m.SetActionType(data.ACTION_DELETE)
	m.SetLoadFields("")
	m.SetSubField("")
	m.SetSubOperator("")

	b.models = []data.IModel{m}
	if !b.runAdapter() {
		return
	}

	b.AppendModelToResponse(ch)
	b.actionUserConnections = actionUserConnections
	b.chatsId = []int{ch.Id}
}

func (b *BusinessHandler) processRemoveUserFromChat() {

	chUs := b.models[0].(*chat.ChatUser)
	chUs.SetActionType(data.ACTION_DELETE)

	if chUs.User == b.client.GetId() {
		if !b.runAdapter() {
			return
		}
	} else {
		ch := chat.NewChat()
		ch.SetActionType(data.ACTION_FIND)
		ch.Id = chUs.Chat
		ch.Author = b.client.GetId()
		initConditionOperatorOfSubQuery(ch, data.EXISTS, "Id")

		chUs.SetChats(&[]data.IModel{ch})
		if !b.runAdapter() {
			return
		}
	}
	chUs.SetChats(&[]data.IModel{})
	b.AppendModelToResponse(chUs)

	allConnChats := chat.NewChatUser()
	allConnChats.SetActionType(data.ACTION_FIND)
	allConnChats.User = chUs.User
	allConnChats.SetLoadFields("Chat")

	b.models = []data.IModel{allConnChats}
	if !b.runAdapter() {
		return
	}

	rez := b.adapter.GetModelRezult()
	finded := make([]int, len(rez))
	for i := 0; i < len(rez); i++ {
		finded[i] = rez[i][data.MODEL_CHAT_USER].(*chat.ChatUser).Chat
	}

	b.actionUserConnections = &ActionUserConnections{
		UsersId:        []int{chUs.User},
		ActionChats:    []int{chUs.Chat},
		ConnectedChats: finded,
		Operation:      constants.CLIENT_DISCONNECT,
	}

	b.chatsId = []int{chUs.Chat}
}

func (b *BusinessHandler) processListChats() {

	chUs := b.models[0].(*chat.ChatUser)
	chUs.SetActionType(data.ACTION_FIND)

	ch := chat.NewChat()
	ch.SetActionType(data.ACTION_FIND)
	ch.SetLoadFields("Id", "Author", "Name", "Status", "CreateAt")
	ch.AddLink(data.NewLink(chUs.GetNameModel(), "Chat", ch.GetNameModel(), "Id", data.LINK_WEIGHT_EQUILIBRIUM))
	chUs.SetChats(&[]data.IModel{ch})

	if !b.runAdapter() {
		return
	}
	b.loadModelsFromAdapterRezult()

	rez := b.adapter.GetModelRezult()

	if len(rez) == 0 {
		return
	}

	m := message.NewMessage()

	m.SetActionType(data.ACTION_FIND)
	m.SetLoadFields("Id", "ChatId", "Author", "Message", "FileUrl", "Date", "Time")
	m.SetSortModels([]string{"Id"}, []string{data.DIRECTION_DESC})
	m.SetLimitCountModel(1)

	u := user.NewUser()
	u.SetActionType(data.ACTION_FIND)
	u.SetLoadFields("Id", "Alias")
	u.AddLink(data.NewLink(
		m.GetNameModel(), "Author", u.GetNameModel(), "Id", data.LINK_WEIGHT_EQUILIBRIUM))

	m.SetUsers(&[]data.IModel{u})
	b.models = []data.IModel{m}

	idMessages := make([]int64, 0, len(rez))

	for _, val := range rez {

		m := b.models[0].(*message.Message)
		m.ChatId = val[ch.GetNameModel()].(*chat.Chat).Id

		if !b.runAdapter() {
			return
		}

		models := b.adapter.GetModelRezult()
		if len(models) > 0 {
			idMessages = append(idMessages, models[0][data.MODEL_MESSAGE].(*message.Message).Id)
		}

		b.loadModelsFromAdapterRezult()
	}

	if len(idMessages) > 0 {
		m := message.NewMessage()
		m.SetActionType(data.ACTION_FIND)
		m.SetLoadFields("Id")
		b.models = []data.IModel{m}

		msModels := make([]data.IModel, len(idMessages))

		for ind, val := range idMessages {
			ms := message.NewMessageStatus()
			ms.SetActionType(data.ACTION_FIND)

			if ind == 0 {
				ms.AddLink(data.NewLink(m.GetNameModel(), "Id", ms.GetNameModel(), "MessageId", data.LINK_WEIGHT_EQUILIBRIUM))
			}

			ms.MessageId = val

			ms.UserId = b.client.GetId()

			ms.AddConditionField("MessageId", data.IN)
			ms.SetLoadFields("MessageId", "Status")
			//	ms.AddGroupModels("MessageId", "Status")
			ms.SetArrayFields("MessageId")

			msModels[ind] = ms
		}

		m.SetStatus(&msModels)

		if !b.runAdapter() {
			return
		}
		b.loadModelsFromAdapterRezult(message.NewMessageStatus().GetNameModel())
	}

}

func (b *BusinessHandler) processNewChat() {

	chModel := b.models[0].(*chat.Chat)

	userMod := user.NewUser()
	userMod.Id = chModel.Author
	userMod.Status = constants.USER_UNLOCK
	initExistConditionOfSubQuery(userMod, data.EXISTS, "Id")

	chModel.Id = 0
	chModel.SetActionType(data.ACTION_SAVE)
	chModel.CreateAt = b.DateRegistration()
	chModel.Status = data.CHAT_STATUS_ACTIVE

	chModel.SetUsers(&[]data.IModel{userMod})
	if !b.runAdapter() {
		return
	}

	chModel.Id = int(b.adapter.GetIdOfLastInsertRow())

	chUser := chat.NewChatUser()
	chUser.SetActionType(data.ACTION_SAVE)
	chUser.Chat = chModel.Id
	chUser.User = chModel.Author

	b.models = []data.IModel{chUser}
	if !b.runAdapter() {
		return
	}

	b.AppendModelToResponse(chModel, chUser)

	chUser = chat.NewChatUser()
	chUser.SetActionType(data.ACTION_FIND)
	chUser.User = chModel.Author
	chUser.SetLoadFields("Chat")

	b.models = []data.IModel{chUser}
	if !b.runAdapter() {
		return
	}

	rez := b.adapter.GetModelRezult()
	connected := make([]int, 0, len(rez))

	for i := 0; i < len(rez); i++ {
		connected = append(connected, rez[i][chUser.GetNameModel()].(*chat.ChatUser).Chat)
	}

	b.actionUserConnections = &ActionUserConnections{
		UsersId:        []int{chUser.User},
		ActionChats:    []int{chModel.Id},
		ConnectedChats: connected,
		Operation:      constants.CLIENT_CONNECT,
	}

	b.chatsId = []int{chModel.Id}
}

func (b *BusinessHandler) processAddUserInChat() {

	chUs := b.models[0].(*chat.ChatUser)
	chUs.SetActionType(data.ACTION_SAVE)

	ch := chat.NewChat()
	ch.Id = chUs.Chat
	ch.Author = b.client.GetId()
	initExistConditionOfSubQuery(ch, data.EXISTS, "Id")

	us := user.NewUser()
	us.Id = chUs.User
	initExistConditionOfSubQuery(us, data.EXISTS, "Id")

	bl := user.NewBlackList()
	bl.User = chUs.User
	bl.BlockedUser = b.client.GetId()
	initExistConditionOfSubQuery(bl, data.NOT_EXISTS, "User")

	chUs.SetUsers(&[]data.IModel{us})
	chUs.SetChats(&[]data.IModel{ch})
	chUs.SetBlackLists(&[]data.IModel{bl})

	if !b.runAdapter() {
		return
	}

	chUs.SetBlackLists(&[]data.IModel{})

	findUserConn := chat.NewChatUser()
	findUserConn.SetActionType(data.ACTION_FIND)
	findUserConn.User = chUs.User
	findUserConn.SetLoadFields("User", "Chat")

	u := user.NewUser()
	u.SetActionType(data.ACTION_FIND)
	u.SetLoadFields("Id", "Alias")
	u.AddLink(data.NewLink(findUserConn.GetNameModel(), "User",
		u.GetNameModel(), "Id", data.LINK_WEIGHT_EQUILIBRIUM))
	findUserConn.SetUsers(&[]data.IModel{u})

	b.models = []data.IModel{findUserConn}

	if !b.runAdapter() {
		return
	}

	rez := b.adapter.GetModelRezult()
	connected := make([]int, 0, len(rez))
	var usr *user.User

	for i := 0; i < len(rez); i++ {
		if i == 0 {
			usr = rez[i][u.GetNameModel()].(*user.User)
		}
		connected = append(connected, rez[i][chUs.GetNameModel()].(*chat.ChatUser).Chat)
	}

	b.actionUserConnections = &ActionUserConnections{
		UsersId:        []int{findUserConn.User},
		ActionChats:    []int{chUs.Chat},
		ConnectedChats: connected,
		Operation:      constants.CLIENT_CONNECT,
	}

	b.chatsId = []int{chUs.Chat}
	b.response[chUs.GetNameModel()] = []data.IModel{chUs}
	b.response[usr.GetNameModel()] = []data.IModel{usr}
}
