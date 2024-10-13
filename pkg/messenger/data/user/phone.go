package user

import (
	"errors"
	"regexp"

	"github.com/alex988334/messenger/pkg/messenger/data"
)

type UserPhone struct {
	data.IModel
	UserId int
	Phone  string
}

func NewUserPhone() *UserPhone {

	p := UserPhone{}

	m := data.NewModel(data.MODEL_USER_PHONE)
	m.SetValidator(p.createValidator(m))
	m.AddConditionField("UserId", data.EQUAL)
	m.AddConditionField("Phone", data.EQUAL)
	m.SetComparisonHandler(p.getComparisonHandler())

	p.IModel = m
	return &p
}

func (u *UserPhone) getComparisonHandler() func(data.IModel) bool {

	return func(model data.IModel) bool {

		if m, ok := model.(*UserPhone); !ok || u.UserId != m.UserId || u.Phone != m.Phone {
			return false
		}

		return true
	}
}

func (u *UserPhone) createValidator(m *data.Model) func() bool {

	return func() bool {

		if m.ContainsLoadField("UserId") && u.UserId == 0 {
			m.AddErrorValidate(errors.New("ERROR UserPhone! User ID is equal zero"))
		}
		if m.ContainsLoadField("Phone") {
			if match, _ := regexp.MatchString(`^\+[0-9]{11}`, u.Phone); !match {
				m.AddErrorValidate(errors.New("ERROR UserPhone! Phone number is not valid"))
			}
		}
		if len(m.GetErrorValidate()) > 0 {
			return false
		}
		return true
	}
}

func (u *UserPhone) GetUserId() int {
	return u.UserId
}
func (u *UserPhone) SetUserId(userId int) {
	u.UserId = userId
}

func (u *UserPhone) GetPhone() string {
	return u.Phone
}
func (u *UserPhone) SetPhone(phone string) {
	u.Phone = phone
}
