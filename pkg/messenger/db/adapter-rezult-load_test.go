package db

import (
	"encoding/json"
	"testing"

	data "github.com/alex988334/messenger/pkg/messenger/data"

	mysql_helper "github.com/alex988334/messenger/pkg/messenger/db/mysql-helper"

	"github.com/alex988334/messenger/pkg/messenger/data/chat"
	"github.com/alex988334/messenger/pkg/messenger/data/message"
)

func TestSelectRezultLoad(t *testing.T) {

	t.Log("Check select rezult load to IModel")

	c := chat.NewChat()
	c.Author = 45
	c.Id = 123
	c.CreateAt = "2016-10-06"
	c.Name = "MAD MAX CHAT"
	c.Status = data.CHAT_STATUS_ACTIVE
	c.SetActionType(data.ACTION_FIND)

	c1 := chat.NewChat()
	c1.Author = 45
	c1.Id = 123
	c1.CreateAt = "2016-10-06"
	c1.Name = "MAD MAX CHAT"
	c1.Status = data.CHAT_STATUS_ACTIVE
	c1.SetActionType(data.ACTION_FIND)

	m := message.NewMessage()
	m.Author = 23
	m.ChatId = 123
	m.Date = "2016-08-18"
	m.FileUrl = ""
	m.Id = 746
	m.Message = "MAD MAX HELLO!"
	m.ParrentMessage = 122
	m.Time = "19:46:31"
	m.SetLoadFields()
	m.SetActionType(data.ACTION_FIND)

	m1 := message.NewMessage()
	m1.Author = 20
	m1.ChatId = 120
	m1.Date = "2016-08-18"
	m1.FileUrl = ""
	m1.Id = 743
	m1.Message = "MAD Report!"
	m1.Time = "19:46:00"
	m1.SetActionType(data.ACTION_FIND)
	m1.AddLink(data.NewLink(c.GetNameModel(), "Id", m1.GetNameModel(), "ChatId", data.LINK_WEIGHT_EQUILIBRIUM))

	c.SetMessages(&[]data.IModel{m1})

	a := NewAdapter()
	s := mysql_helper.NewSQLModel(NewDB())
	s.AddSelectFields([]string{"id", "author", "alias", "create_at", "status", ""}, "chat")
	s.AddJoin(mysql_helper.LEFT, "chat", "chat_message", "id", "id_chat")

	a.SetModels(c)
	a.sqlModel = s

	a1 := NewAdapter()
	s1 := mysql_helper.NewSQLModel(NewDB())
	s1.AddSelectFields([]string{"id", "author", "alias", "create_at", "status", ""}, "chat")
	a1.SetModels(c1)
	a1.sqlModel = s1

	required := [][]map[string]data.IModel{
		{
			{
				"Chat":    c,
				"Message": m,
			},
			{
				"Chat":    c,
				"Message": m1,
			},
		},
		{
			{
				"Chat": c,
			},
		},
	}
	params := [][]map[string]interface{}{
		{
			{
				"c.id":         123,
				"c.author":     45,
				"c.alias":      "MAD MAX CHAT",
				"c.create_at":  "2016-10-06",
				"c.status":     data.CHAT_STATUS_ACTIVE,
				"ch.id":        746,
				"ch.id_chat":   123,
				"ch.id_user":   23,
				"ch.parent_id": 122,
				"ch.message":   "MAD MAX HELLO!",
				"ch.file":      "",
				"ch.date":      "2016-08-18",
				"ch.time":      "19:46:31",
			},
			{
				"c.id":        123,
				"c.author":    45,
				"c.alias":     "MAD MAX CHAT",
				"c.create_at": "2016-10-06",
				"c.status":    data.CHAT_STATUS_ACTIVE,
				"ch.id":       743,
				"ch.id_chat":  120,
				"ch.id_user":  20,

				"ch.message": "MAD Report!",
				"ch.file":    "",
				"ch.date":    "2016-08-18",
				"ch.time":    "19:46:00",
			},
		},
		{
			{
				"id":        123,
				"author":    45,
				"alias":     "MAD MAX CHAT",
				"create_at": "2016-10-06",
				"status":    data.CHAT_STATUS_ACTIVE,
			},
		},
	}
	adapters := []Adapter{*a, *a1}

	for i := 0; i < len(params); i++ {

		adapters[i].selectRezult = params[i]
		adapters[i].loadModelFromSQLQuery()

		rez, _ := json.Marshal(adapters[i].modelsRezult)
		req, _ := json.Marshal(required[i])
		if string(rez) != string(req) {

			t.Fatal("Failed Check Load IModel To SQL Model! \nRequired:", string(req),
				"\nrezult:  ", string(rez),
			)
		}
	}
}
