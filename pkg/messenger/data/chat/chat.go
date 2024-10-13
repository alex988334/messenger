package chat

import (
	"errors"

	cons "github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/functions"
)

type Chat struct {
	data.IModel
	Id       int
	Author   int
	Name     string
	CreateAt string
	Status   string
}

func NewChat() *Chat {

	c := Chat{}

	m := data.NewModel(data.MODEL_CHAT)
	m.SetValidator(c.createValidator(m))
	m.AddConditionField("Id", data.EQUAL)
	m.SetComparisonHandler(c.getComparisonHandler())

	c.IModel = m
	return &c
}

func (c *Chat) getComparisonHandler() func(data.IModel) bool {

	return func(model data.IModel) bool {

		if m, ok := model.(*Chat); !ok || c.Id != m.Id || c.Author != m.Author ||
			c.Name != m.Name || c.Status != m.Status || c.CreateAt != m.CreateAt {
			return false
		}

		return true
	}
}

func (c *Chat) createValidator(m *data.Model) func() bool {

	return func() bool {
		if m.ContainsLoadField("Id") && c.Id == 0 {
			m.AddErrorValidate(errors.New("ERROR Chat! Chat ID is equal zero"))
		}
		if m.ContainsLoadField("Author") && c.Author == 0 {
			m.AddErrorValidate(errors.New("ERROR Chat! Chat author ID is equal zero"))
		}
		if m.ContainsLoadField("Name") && c.Name == "" {
			m.AddErrorValidate(errors.New("ERROR Chat! Chat name is empty"))
		}
		if m.ContainsLoadField("CreateAt") && (c.CreateAt == "" || !functions.IsDate(c.CreateAt)) {
			m.AddErrorValidate(errors.New("ERROR Chat! Chat date created is not valid"))
		}
		if m.ContainsLoadField("Status") && c.Status != cons.CHAT_ACTIVE &&
			c.Status != cons.CHAT_DELETED && c.Status != cons.CHAT_DIACTIVATED {
			m.AddErrorValidate(errors.New("ERROR Chat! Chat status is not valid"))
		}
		if len(m.GetErrorValidate()) > 0 {
			return false
		}
		return true
	}
}

func (c *Chat) GetId() int {
	return c.Id
}
func (c *Chat) SetId(id int) {
	c.Id = id
}

func (c *Chat) GetAuthor() int {
	return c.Author
}
func (c *Chat) SetAuthor(authorId int) {
	c.Author = authorId
}

func (c *Chat) GetName() string {
	return c.Name
}
func (c *Chat) SetName(name string) {
	c.Name = name
}

func (c *Chat) GetCreateAt() string {
	return c.CreateAt
}
func (c *Chat) SetCreateAt(createAt string) {
	c.CreateAt = createAt
}

func (c *Chat) GetChatStatus() string {
	return c.Status
}
func (c *Chat) SetChatStatus(status string) {
	c.Status = status
}
