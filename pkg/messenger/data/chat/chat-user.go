package chat

import (
	"errors"

	"github.com/alex988334/messenger/pkg/messenger/data"
)

type ChatUser struct {
	data.IModel
	Chat        int
	User        int
	SessionHash string
}

func NewChatUser() *ChatUser {

	var cu ChatUser = ChatUser{}

	m := data.NewModel(data.MODEL_CHAT_USER)
	m.SetValidator(cu.createValidator(m))
	m.AddConditionField("Chat", data.EQUAL)
	m.AddConditionField("User", data.EQUAL)
	m.SetComparisonHandler(cu.getComparisonHandler())

	cu.IModel = m
	return &cu
}

func (cu *ChatUser) getComparisonHandler() func(data.IModel) bool {

	return func(model data.IModel) bool {

		if m, ok := model.(*ChatUser); !ok || cu.Chat != m.Chat || cu.User != m.User ||
			cu.SessionHash != m.SessionHash {
			return false
		}

		return true
	}
}

func (cu *ChatUser) createValidator(m *data.Model) func() bool {

	return func() bool {

		if m.ContainsLoadField("Chat") && cu.Chat == 0 {
			m.AddErrorValidate(errors.New("ERROR ChatUser! Chat ID is equal zero"))
		}
		if m.ContainsLoadField("User") && cu.User == 0 {
			m.AddErrorValidate(errors.New("ERROR ChatUser! User ID is equal zero"))
		}
		if m.ContainsLoadField("SessionHash") && cu.SessionHash == "" {
			m.AddErrorValidate(errors.New("ERROR ChatUser! Session hash is empty"))
		}
		if len(m.GetErrorValidate()) > 0 {
			return false
		}
		return true
	}
}

func (cu *ChatUser) GetChat() int {
	return cu.Chat
}
func (cu *ChatUser) SetChat(chatId int) {
	cu.Chat = chatId
}

func (cu *ChatUser) GetUser() int {
	return cu.User
}
func (cu *ChatUser) SetUser(userId int) {
	cu.User = userId
}

func (cu *ChatUser) GetSessionHash() string {
	return cu.SessionHash
}
func (cu *ChatUser) SetSessionHash(sessionHash string) {
	cu.SessionHash = sessionHash
}

/****
func NewChatUser() *ChatUser {

	return &ChatUser{Model: data.Model{NameModel: "User"}}
}

func (u *ChatUser) SetId(id int64) {
	u.id = id
}
func (u *ChatUser) GetId() int64 {
	return u.id
}

func (u *ChatUser) GetStatus() int {
	return u.status
}
func (u *ChatUser) SetStatus(status int) {
	u.status = status
}

func (u *ChatUser) GetCreateAt() int64 {
	return u.createdAt
}
func (u *ChatUser) SetCreateAt(unixSec int64) {
	u.createdAt = unixSec
}

func (u *ChatUser) SetLogin(login string) {
	u.login = login
}
func (u *ChatUser) GetLogin() string {
	return u.login
}

func (u *ChatUser) GetPassHash() string {
	return u.passHash
}
func (u *ChatUser) SetPassHash(passHash string) {
	u.passHash = passHash
}

func (u *ChatUser) GetUserKey() string {
	return u.userKey
}
func (u *ChatUser) SetUserKey(key string) {
	u.userKey = key
}

func (u *ChatUser) GetIdentKey() string {
	return u.identiKey
}
func (u *ChatUser) SetIdentKey(key string) {
	u.identiKey = key
}

func (u *ChatUser) GetEmail() sql.NullString {
	return u.email
}
func (u *ChatUser) SetEmail(email sql.NullString) {
	u.email = email
}
*/
