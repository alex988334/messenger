package user

import (
	"errors"
	"regexp"

	cons "github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
)

type User struct {
	data.IModel
	Id             int
	Login          string
	Alias          string
	AuthKey        string
	PassHash       string
	PassResetToken string
	Email          string
	Status         int
	CreateAt       int64
	UpdateAt       int64
	Avatar         string
}

func NewUser() *User {

	u := User{}

	m := data.NewModel(data.MODEL_USER)
	m.SetValidator(u.createValidator(m))
	m.AddConditionField("Id", data.EQUAL)
	m.SetComparisonHandler(u.getComparisonHandler())

	u.IModel = m
	return &u
}

func (u *User) getComparisonHandler() func(data.IModel) bool {

	return func(model data.IModel) bool {

		if m, ok := model.(*User); !ok || u.Id != m.Id || u.Login != m.Login ||
			u.Alias != m.Alias || u.Email != m.Email || u.AuthKey != m.AuthKey ||
			u.Status != m.Status || u.PassHash != m.PassHash || u.CreateAt != m.CreateAt ||
			u.UpdateAt != m.UpdateAt || u.PassResetToken != m.PassResetToken {
			return false
		}

		return true
	}
}

func (u *User) createValidator(m *data.Model) func() bool {

	return func() bool {

		if m.ContainsLoadField("Id") && u.Id == 0 {
			m.AddErrorValidate(errors.New("User ID is equal zero"))
		}
		if m.ContainsLoadField("Login") && len(u.Login) < 8 {
			m.AddErrorValidate(errors.New("Login length less than 8 characters"))
		}
		if m.ContainsLoadField("Alias") && len(u.Alias) < 6 {
			m.AddErrorValidate(errors.New("User alias less than 6 characters"))
		}
		if m.ContainsLoadField("AuthKey") && u.AuthKey == "" {
			m.AddErrorValidate(errors.New("Auth key is empty"))
		}
		if m.ContainsLoadField("PassHash") && len(u.PassHash) < 8 {
			m.AddErrorValidate(errors.New("Password less than 8 characters"))
		}
		if m.ContainsLoadField("PassResetToken") && u.PassResetToken == "" {
			m.AddErrorValidate(errors.New("Password reset token is empty"))
		}
		if m.ContainsLoadField("Email") {
			if match, _ := regexp.MatchString(`@gmail.com\z`, u.Email); !match { //`\w{6,40}@gmail.com\z`
				m.AddErrorValidate(errors.New(
					"User`s email  must be length less 50 characters and it must contains @gmail.com"))
			}
		}
		if m.ContainsLoadField("Status") && u.Status != cons.USER_LOCK && u.Status != cons.USER_UNLOCK {
			m.AddErrorValidate(errors.New("User ID is not valid"))
		}
		if m.ContainsLoadField("CreateAt") && u.CreateAt == 0 {
			m.AddErrorValidate(errors.New("Create at is equal zero"))
		}
		if m.ContainsLoadField("UpdateAt") && u.UpdateAt == 0 {
			m.AddErrorValidate(errors.New("Update at is equal zero"))
		}
		if m.ContainsLoadField("Avatar") && len(u.Avatar) > 100 {
			m.AddErrorValidate(errors.New("Length file name more 100 characters"))
		}

		if len(m.GetErrorValidate()) > 0 {
			return false
		}
		return true
	}
}

func (u *User) SetId(id int) {
	u.Id = id
}
func (u *User) GetId() int {
	return u.Id
}

func (u *User) SetLogin(login string) {
	u.Login = login
}
func (u *User) GetLogin() string {
	return u.Login
}
func (u *User) SetAlias(alias string) {
	u.Alias = alias
}
func (u *User) GetAlias() string {
	return u.Alias
}

func (u *User) SetAuthKey(authKey string) {
	u.AuthKey = authKey
}
func (u *User) GetAuthKey() string {
	return u.AuthKey
}

func (u *User) SetPassHash(passHash string) {
	u.PassHash = passHash
}
func (u *User) GetPassHash() string {
	return u.PassHash
}

func (u *User) SetPassResetToken(passResetToken string) {
	u.PassResetToken = passResetToken
}
func (u *User) GetPassResetToken() string {
	return u.PassResetToken
}

func (u *User) SetEmail(email string) {
	u.Email = email
}
func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) SetUserStatus(status int) {
	u.Status = status
}
func (u *User) GetUserStatus() int {
	return u.Status
}

func (u *User) SetCreateAt(createAt int64) {
	u.CreateAt = createAt
}
func (u *User) GetCreateAt() int64 {
	return u.CreateAt
}

func (u *User) SetUpdateAt(updateAt int64) {
	u.UpdateAt = updateAt
}
func (u *User) GetUpdateAt() int64 {
	return u.UpdateAt
}
