package business_logic

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/data/chat"
	"github.com/alex988334/messenger/pkg/messenger/data/user"
	"github.com/alex988334/messenger/pkg/messenger/functions"
)

func (b *BusinessHandler) processRegistration() {

	us := b.models[0].(*user.User)

	/*if len(us.Login) < 6 || len(us.Alias) < 6 {
		b.err = errors.New("ERROR! the length of the login, nickname cannot be less than 6 characters")
		return
	}
	if len(us.PassHash) < 12 || len(us.Email) < 12 {
		b.err = errors.New("ERROR! the length of the password or email cannot be less than 12 characters")
		return
	}
	*/
	passHash, err := bcrypt.GenerateFromPassword([]byte(us.PassHash), constants.PASSHASH_COST)
	if err != nil {
		fmt.Println("ERROR! Failed to create password hash:", err)
		b.err = errors.New("ERROR! Failed to create password hash, please change password")
		return
	}

	us.PassHash = string(passHash)
	us.AuthKey = functions.RandomString(32)
	us.CreateAt = b.registrationTime.GetUnix()
	us.UpdateAt = b.registrationTime.GetUnix()
	us.SetActionType(data.ACTION_SAVE)

	if !b.runAdapter() {
		return
	}

	us.PassHash = ""
	b.AppendModelToResponse(us)
}

func (b *BusinessHandler) processAutorizate() {

	us := b.models[0].(*user.User)
	us.SetActionType(data.ACTION_FIND)
	us.SetLoadFields("Id", "Login", "PassHash")

	password := us.PassHash
	us.PassHash = ""

	if !b.runAdapter() {
		return
	}

	us.PassHash = password

	rez := b.adapter.GetModelRezult()
	if len(rez) == 0 {
		b.err = errors.New("ERROR autorizate! User not found")
		return
	}

	usUpdate := rez[0][us.GetNameModel()].(*user.User)

	if usUpdate.Login != us.Login || bcrypt.CompareHashAndPassword([]byte(usUpdate.PassHash), []byte(us.PassHash)) != nil {
		b.err = errors.New("ERROR autorization! Login or password is invalid")
		return
	}

	usUpdate.SetActionType(data.ACTION_UPDATE)
	usUpdate.AuthKey = functions.RandomString(32)
	usUpdate.AddConditionField("Id", data.EQUAL)

	b.models = []data.IModel{usUpdate}
	if !b.runAdapter() {
		return
	}
	usUpdate.Login = ""
	usUpdate.PassHash = ""

	b.client.SetId(usUpdate.Id)
	b.AppendModelToResponse(usUpdate)

	b.generateListConnection([]int{usUpdate.Id}, []int{}, constants.CLIENT_CONNECT_ALL_CHATS)
}

func (b *BusinessHandler) processMyData() {

	us := b.models[0].(*user.User)
	us.SetActionType(data.ACTION_FIND)
	us.SetLoadFields("Id", "Login", "Alias", "Email", "Status", "CreateAt", "UpdateAt", "Avatar")

	ph := user.NewUserPhone()
	ph.SetActionType(data.ACTION_FIND)
	ph.AddLink(data.NewLink(us.GetNameModel(), "Id", ph.GetNameModel(), "UserId", data.LINK_WEIGHT_EQUILIBRIUM))
	ph.SetLoadFields("UserId", "Phone")

	us.SetUsers(&[]data.IModel{ph})
	if !b.runAdapter() {
		return
	}

	b.loadModelsFromAdapterRezult()
}

func (b *BusinessHandler) processSearchUser() {

	prepare := false
	if u, ok := b.models[0].(*user.User); ok {
		u.SetActionType(data.ACTION_FIND)
		u.AddLoadFields("Id", "Alias")
		prepare = true
	}

	if up, ok := b.models[0].(*user.UserPhone); ok {
		up.SetActionType(data.ACTION_FIND)
		up.AddLoadFields("UserId", "Phone")
		prepare = true
	}

	if !prepare {
		b.err = errors.New("ERROR! Data is not loading")
		return
	}

	if !b.runAdapter() {
		return
	}

	b.loadModelsFromAdapterRezult()
}

func (b *BusinessHandler) processListUsers() {

	chUs := b.models[0].(*chat.ChatUser)
	chUs.SetActionType(data.ACTION_FIND)
	initExistConditionOfSubQuery(chUs, data.EXISTS, "Chat")

	chusers := chat.NewChatUser()
	chusers.SetActionType(data.ACTION_FIND)
	chusers.SetLoadFields("Chat", "User")
	chusers.Chat = chUs.Chat

	us := user.NewUser()
	us.SetActionType(data.ACTION_FIND)
	us.SetLoadFields("Id", "Alias")
	us.AddLink(data.NewLink(chusers.GetNameModel(), "User", us.GetNameModel(), "Id", data.LINK_WEIGHT_EQUILIBRIUM))

	chusers.SetUsers(&[]data.IModel{us})
	chusers.SetChatUsers(&[]data.IModel{chUs})

	b.models = []data.IModel{chusers}
	if !b.runAdapter() {
		return
	}

	b.loadModelsFromAdapterRezult()
}
