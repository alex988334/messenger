package message

import (
	"errors"

	cons "github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/functions"
)

type MessageStatus struct {
	data.IModel
	MessageId int64
	UserId    int
	Status    string
	Date      string
	Time      string
}

func NewMessageStatus() *MessageStatus {

	ms := MessageStatus{}

	m := data.NewModel(data.MODEL_STATUS_MESSAGE)
	m.SetValidator(ms.createValidator(m))
	m.AddConditionField("MessageId", data.EQUAL)
	m.AddConditionField("UserId", data.EQUAL)
	m.SetComparisonHandler(ms.getComparisonHandler())

	ms.IModel = m
	return &ms
}

func (ms *MessageStatus) getComparisonHandler() func(data.IModel) bool {

	return func(model data.IModel) bool {

		if m, ok := model.(*MessageStatus); !ok || ms.MessageId != m.MessageId || ms.UserId != m.UserId ||
			ms.Status != m.Status || ms.Date != m.Date || ms.Time != m.Time {
			return false
		}

		return true
	}
}

func (ms *MessageStatus) createValidator(m *data.Model) func() bool {

	return func() bool {
		if m.ContainsLoadField("MessageId") && ms.MessageId == 0 {
			m.AddErrorValidate(errors.New("ERROR MessageStatus! Message ID is equal zero"))
		}
		if m.ContainsLoadField("UserId") && ms.UserId == 0 {
			m.AddErrorValidate(errors.New("ERROR MessageStatus! User ID is equal zero"))
		}
		if m.ContainsLoadField("Status") && ms.Status != cons.MESSAGE_DELIVERED &&
			ms.Status != cons.MESSAGE_SENDED && ms.Status != cons.MESSAGE_READED && ms.Status != cons.MESSAGE_CREATED {
			m.AddErrorValidate(errors.New("ERROR MessageStatus! Status message is not valid"))
		}
		if m.ContainsLoadField("Date") && (ms.Date == "" || !functions.IsDate(ms.Date)) {
			m.AddErrorValidate(errors.New("ERROR MessageStatus! Message date is not valid"))
		}
		if m.ContainsLoadField("Time") && (ms.Time == "" || !functions.IsTime(ms.Time)) {
			m.AddErrorValidate(errors.New("ERROR MessageStatus! Message time is not valid"))
		}
		if len(m.GetErrorValidate()) > 0 {
			return false
		}
		return true
	}
}

func (ms *MessageStatus) GetMessageId() int64 {
	return ms.MessageId
}
func (ms *MessageStatus) SetMessageId(messageId int64) {
	ms.MessageId = messageId
}

func (ms *MessageStatus) GetUserId() int {
	return ms.UserId
}
func (ms *MessageStatus) SetUserId(userId int) {
	ms.UserId = userId
}

func (ms *MessageStatus) GetStatusMessage() string {
	return ms.Status
}
func (ms *MessageStatus) SetStatusMessage(status string) {
	ms.Status = status
}

func (ms *MessageStatus) GetDate() string {
	return ms.Date
}
func (ms *MessageStatus) SetDate(date string) {
	ms.Date = date
}

func (ms *MessageStatus) GetTime() string {
	return ms.Time
}
func (ms *MessageStatus) SetTime(time string) {
	ms.Time = time
}
