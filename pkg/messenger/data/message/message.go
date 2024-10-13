package message

import (
	"errors"

	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/functions"
)

type Message struct {
	data.IModel
	Id             int64
	ChatId         int
	Author         int
	ParrentMessage int64
	Message        string
	FileUrl        string
	Date           string
	Time           string
}

func NewMessage() *Message {

	mes := &Message{}

	m := data.NewModel(data.MODEL_MESSAGE)
	m.SetValidator(mes.createValidator(m))
	m.AddConditionField("Id", data.EQUAL)
	m.SetComparisonHandler(mes.getComparisonHandler())

	mes.IModel = m

	return mes
}

func (m *Message) getComparisonHandler() func(data.IModel) bool {

	return func(model data.IModel) bool {

		if ms, ok := model.(*Message); !ok || ms.Id != m.Id || ms.ChatId != m.ChatId ||
			ms.Author != m.Author || ms.ParrentMessage != m.ParrentMessage ||
			ms.Message != m.Message || ms.FileUrl != m.FileUrl ||
			ms.Date != m.Date || ms.Time != m.Time {
			return false
		}

		return true
	}
}

func (m *Message) createValidator(mod *data.Model) func() bool {

	return func() bool {
		if mod.ContainsLoadField("Id") && m.Id == 0 {
			mod.AddErrorValidate(errors.New("ERROR Message! Message ID is equal zero"))
		}
		if mod.ContainsLoadField("ChatId") && m.ChatId == 0 {
			mod.AddErrorValidate(errors.New("ERROR Message! Chat ID is equal zero"))
		}
		if mod.ContainsLoadField("Author") && m.Author == 0 {
			mod.AddErrorValidate(errors.New("ERROR Message! Message author id is equal zero"))
		}
		if mod.ContainsLoadField("ParrentMessage") && m.ParrentMessage == 0 {
			mod.AddErrorValidate(errors.New("ERROR Message! Parent message id is equal zero"))
		}
		if mod.ContainsLoadField("Message") && m.Message == "" {
			mod.AddErrorValidate(errors.New("ERROR Message! Message is empty"))
		}
		if mod.ContainsLoadField("FileUrl") && m.FileUrl == "" {
			mod.AddErrorValidate(errors.New("ERROR Message! File Url is empty"))
		}
		if mod.ContainsLoadField("Date") && (m.Date == "" || !functions.IsDate(m.Date)) {
			mod.AddErrorValidate(errors.New("ERROR Message! Message date is not valid"))
		}
		if mod.ContainsLoadField("Time") && (m.Time == "" || !functions.IsTime(m.Time)) {
			mod.AddErrorValidate(errors.New("ERROR Message! Message time is not valid"))
		}
		if len(mod.GetErrorValidate()) > 0 {
			return false
		}
		return true
	}
}

func (m *Message) GetId() int64 {
	return m.Id
}
func (m *Message) SetId(id int64) {
	m.Id = id
}

func (m *Message) GetChatId() int {
	return m.ChatId
}
func (m *Message) SetChatId(chatId int) {
	m.ChatId = chatId
}

func (m *Message) GetAuthor() int {
	return m.Author
}
func (m *Message) SetAuthor(author int) {
	m.Author = author
}

func (m *Message) GetParrentMessage() int64 {
	return m.ParrentMessage
}
func (m *Message) SetParrentMessage(parrentMessage int64) {
	m.ParrentMessage = parrentMessage
}

func (m *Message) GetMessage() string {
	return m.Message
}
func (m *Message) SetMessage(message string) {
	m.Message = message
}

func (m *Message) GetFileUrl() string {
	return m.FileUrl
}
func (m *Message) SetFileUrl(fileUrl string) {
	m.FileUrl = fileUrl
}

func (m *Message) GetDate() string {
	return m.Date
}
func (m *Message) SetDate(date string) {
	m.Date = date
}

func (m *Message) GetTime() string {
	return m.Time
}
func (m *Message) SetTime(time string) {
	m.Time = time
}
