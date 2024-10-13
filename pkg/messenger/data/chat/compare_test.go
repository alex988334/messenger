package chat

import (
	"testing"

	"github.com/alex988334/messenger/pkg/messenger/data"
)

/*
	type ChatUser struct {
		data.IModel
		Chat        int
		User        int
		SessionHash string
	}
*/
func TestCompareChatUserModels(t *testing.T) {

	t.Log("Generate Compare Chat User Models")
	required := []bool{
		true,
		true,
		true,
		true,
		true,
		false,
		false,
		false,
		false,
		false,
	}
	params := [][][]interface{}{
		{
			{10, 110, "Version"}, {10, 110, "Version"},
		},
		{
			{99, 0, "Version"}, {99, 0, "Version"},
		},
		{
			{0, 220, "Version"}, {0, 220, "Version"},
		},
		{
			{0, 220, ""}, {0, 220, ""},
		},
		{
			{0, 0, ""}, {0, 0, ""},
		},
		{
			{10, 110, "Version"}, {10, 0, "Version"},
		},
		{
			{10, 110, "Version"}, {10, 110, ""},
		},
		{
			{0, 0, "v"}, {0, 0, ""},
		},
		{
			{110, 2, "vrwere"}, {150, 3, "fdfdsfdsfdsf"},
		},
		{
			{}, {},
		},
	}

	for ind, test := range params {

		m := NewChatUser()
		if len(test[0]) > 0 {
			m.Chat = test[0][0].(int)
			m.User = test[0][1].(int)
			m.SessionHash = test[0][2].(string)
		}

		m1 := NewChatUser()
		if len(test[1]) > 0 {
			m1.Chat = test[1][0].(int)
			m1.User = test[1][1].(int)
			m1.SessionHash = test[1][2].(string)
		}
		var im1 data.IModel = m1
		if ind == len(params)-1 {
			im1 = NewChat()
		}

		if required[ind] != im1.IsEqualModels(m) {
			t.Fatal("Failed Generate List Messages! \nRequired:", required[ind], "; \nmodel1: ", m, "; \nmodel2: ", im1)
		} else {
			//	fmt.Println("index =>", strconv.Itoa(ind), ", true")
		}
	}
}

/*
	type Chat struct {
		data.IModel
		Id       int
		Author   int
		Name     string
		CreateAt string
		Status   string
	}
*/
func TestCompareChatModels(t *testing.T) {

	t.Log("Generate Compare Chat Models")
	required := []bool{
		true,
		true,
		true,
		true,
		true,
		false,
		false,
		false,
		false,
		false,
	}
	params := [][][]interface{}{
		{
			{10, 110, "Version", "2016-09-08", data.CHAT_STATUS_ACTIVE}, {10, 110, "Version", "2016-09-08", data.CHAT_STATUS_ACTIVE},
		},
		{
			{99, 0, "Version", "2016-09-08", data.CHAT_STATUS_ACTIVE}, {99, 0, "Version", "2016-09-08", data.CHAT_STATUS_ACTIVE},
		},
		{
			{0, 0, "", "", data.CHAT_STATUS_ACTIVE}, {0, 0, "", "", data.CHAT_STATUS_ACTIVE},
		},
		{
			{0, 220, "", "", ""}, {0, 220, "", "", ""},
		},
		{
			{0, 0, "", "", ""}, {0, 0, "", "", ""},
		},
		{
			{10, 110, "Version", "2016-09-08", data.CHAT_STATUS_ACTIVE}, {10, 0, "Version", "2016-09-08", data.CHAT_STATUS_DELETED},
		},
		{
			{0, 0, "", "", ""}, {10, 0, "", "", ""},
		},
		{
			{0, 0, "", "2016-09-08", ""}, {0, 0, "", "", ""},
		},
		{
			{110, 2, "vrwere", "2016-09-08", data.CHAT_STATUS_DELETED}, {150, 3, "fdfdsfdsfdsf", "2016-09-08", data.CHAT_STATUS_ACTIVE},
		},
		{
			{}, {},
		},
	}

	m := NewChat()
	m1 := NewChat()
	for ind, test := range params {

		if len(test[0]) > 0 {
			m.Id = test[0][0].(int)
			m.Author = test[0][1].(int)
			m.Name = test[0][2].(string)
			m.CreateAt = test[0][3].(string)
			m.Status = test[0][4].(string)
		} else {
			m = NewChat()
		}

		if len(test[1]) > 0 {
			m1.Id = test[1][0].(int)
			m1.Author = test[1][1].(int)
			m1.Name = test[1][2].(string)
			m1.CreateAt = test[1][3].(string)
			m1.Status = test[1][4].(string)
		} else {
			m1 = NewChat()
		}
		var im1 data.IModel = m1
		if ind == len(params)-1 {
			im1 = NewChatUser()
		}

		if required[ind] != im1.IsEqualModels(m) {
			t.Fatal("Failed Generate List Messages! \nRequired:", required[ind], "; \nmodel1: ", m, "; \nmodel2: ", im1)
		} else {
			//	fmt.Println("index =>", strconv.Itoa(ind), ", true")
		}
	}
}
