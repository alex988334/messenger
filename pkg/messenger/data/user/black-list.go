package user

import (
	"errors"

	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/functions"
)

type BlackList struct {
	data.IModel
	User        int
	BlockedUser int
	Date        string
	Time        string
}

func NewBlackList() *BlackList {

	bl := BlackList{}

	m := data.NewModel(data.MODEL_BLACK_LIST)
	m.SetValidator(bl.createValidator(m))
	m.AddConditionField("User", data.EQUAL)
	m.AddConditionField("BlockedUser", data.EQUAL)
	m.SetComparisonHandler(bl.getComparisonHandler())

	bl.IModel = m
	return &bl
}

func (bl *BlackList) getComparisonHandler() func(data.IModel) bool {

	return func(model data.IModel) bool {

		if m, ok := model.(*BlackList); !ok || bl.User != m.User || bl.BlockedUser != m.BlockedUser ||
			bl.Date != m.Date || bl.Time != m.Time {
			return false
		}

		return true
	}
}

func (bl *BlackList) createValidator(m *data.Model) func() bool {

	return func() bool {

		if m.ContainsLoadField("User") && bl.User == 0 {
			m.AddErrorValidate(errors.New("ERROR BlackList! User ID is equal zero"))
		}
		if m.ContainsLoadField("BlockedUser") && bl.BlockedUser == 0 {
			m.AddErrorValidate(errors.New("ERROR BlackList! Blocked user ID is equal zero"))
		}
		if m.ContainsLoadField("Date") && (bl.Date == "" || functions.IsDate(bl.Date)) {
			m.AddErrorValidate(errors.New("ERROR BlackList! Date is not valid"))
		}
		if m.ContainsLoadField("Time") && (bl.Time == "" || functions.IsTime(bl.Time)) {
			m.AddErrorValidate(errors.New("ERROR BlackList! Time is not valid"))
		}
		if len(m.GetErrorValidate()) > 0 {
			return false
		}
		return true
	}
}

func (bl *BlackList) GetUser() int {
	return bl.User
}
func (bl *BlackList) SetUser(userId int) {
	bl.User = userId
}

func (bl *BlackList) GetBlockedUser() int {
	return bl.BlockedUser
}
func (bl *BlackList) SetBlockedUser(blockedUser int) {
	bl.BlockedUser = blockedUser
}

func (bl *BlackList) GetDate() string {
	return bl.Date
}
func (bl *BlackList) SetDate(date string) {
	bl.Date = date
}

func (bl *BlackList) GetTime() string {
	return bl.Time
}
func (bl *BlackList) SetTime(time string) {
	bl.Time = time
}
