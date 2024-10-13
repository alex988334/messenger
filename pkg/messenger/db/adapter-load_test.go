package db

import (
	"testing"

	data "github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/functions"

	mysql_helper "github.com/alex988334/messenger/pkg/messenger/db/mysql-helper"

	"github.com/alex988334/messenger/pkg/messenger/data/chat"
	"github.com/alex988334/messenger/pkg/messenger/data/message"
)

func TestLoadDeleteIModelToSQLModel(t *testing.T) {

	t.Log("Check Load Delete IModel To SQL Model")

	required := []struct {
		query  string
		params []interface{}
	}{
		{
			query: "DELETE FROM chat WHERE id=? AND author=(SELECT c.id_chat AS \"c.id_chat\" " +
				"FROM chat_user c WHERE c.client_hash=? AND c.id_user NOT IN (?, ?) " +
				"AND EXISTS(SELECT c.id AS \"c.id\", c.id_chat AS \"c.id_chat\", c.id_user AS \"c.id_user\" " +
				"FROM chat_message c WHERE c.id_user=? AND c.id IN (?, ?) AND c.id_chat NOT IN (?, ?)))",
			params: []interface{}{234, "werwerwer", 12001, 12002, 5677, 356, 357, 7000, 7001},
		},
		{
			query: "DELETE FROM chat_message_status WHERE id_message IN (SELECT c.id AS \"c.id\" " +
				"FROM chat_message c WHERE c.file<>? AND c.id_chat IN (?, ?))",
			params: []interface{}{"QWERT", 3, 4},
		},
	}

	m1 := chat.NewChat()
	m1.SetActionType(data.ACTION_DELETE)
	m1.AddConditionField("Id", data.EQUAL)
	m1.Id = 234

	cu := chat.NewChatUser()
	cu.User = 12001
	cu.SessionHash = "werwerwer"
	cu.SetActionType(data.ACTION_FIND)
	cu.SetLoadFields("Chat")
	cu.SetSubField("Author")
	cu.SetSubOperator(data.EQUAL)
	cu.AddConditionField("User", data.NOT_IN)
	cu.SetArrayFields("User")

	cu1 := chat.NewChatUser()
	cu1.User = 12002
	cu1.SetArrayFields("User")

	cm := message.NewMessage()
	cm.SetArrayFields("Id", "ChatId")
	cm.Id = 356
	cm.Author = 5677
	cm.ChatId = 7000
	cm.AddConditionField("Id", data.IN)
	cm.AddConditionField("ChatId", data.NOT_IN)
	cm.SetActionType(data.ACTION_FIND)
	cm.SetLoadFields("Id", "ChatId", "Author")
	cm.SetSubField("")
	cm.SetSubOperator(data.EXISTS)

	cm1 := message.NewMessage()
	cm1.SetArrayFields("Id", "ChatId")
	cm1.Id = 357
	cm1.ChatId = 7001

	m1.SetChatUsers(&[]data.IModel{cu, cu1})
	cu.SetMessages(&[]data.IModel{cm, cm1})

	ms := message.NewMessageStatus()
	ms.SetActionType(data.ACTION_DELETE)

	mess := message.NewMessage()
	mess.SetActionType(data.ACTION_FIND)
	mess.FileUrl = "QWERT"
	mess.ChatId = 3
	mess.SetArrayFields("ChatId")
	mess.AddConditionField("FileUrl", data.NOT_EQUAL)
	mess.AddConditionField("ChatId", data.IN)
	mess.SetLoadFields("Id")
	mess.SetSubField("MessageId")
	mess.SetSubOperator(data.IN)

	mess1 := message.NewMessage()
	mess1.SetArrayFields("ChatId")
	mess1.ChatId = 4

	ms.SetMessages(&[]data.IModel{mess, mess1})

	params := []data.IModel{
		m1, ms,
	}

	adapter := NewAdapter()
	db := NewDB()

	for i := 0; i < len(params); i++ {

		adapter.sqlModel = mysql_helper.NewSQLModel(db)

		loadModelToSQLModel(params[i], adapter.sqlModel)

		adapter.sqlModel.PrepareQuery()

		if adapter.sqlModel.GetQuery() != required[i].query ||
			!functions.ArraysIsEqual(adapter.sqlModel.GetQueryParams(), required[i].params) {

			t.Fatal("Failed Check Load Delete IModel To SQL Model! \nRequiredSQL:", required[i].query,
				"\nrezultSQL:  ", adapter.sqlModel.GetQuery(),
				"\nRequired params:", required[i].params,
				";\nrezult params:  ", adapter.sqlModel.GetQueryParams(),
			)
		}
	}
}

func TestLoadInsertIModelToSQLModel(t *testing.T) {

	t.Log("Check Load Insert IModel To SQL Model")

	required := []struct {
		query  string
		params []interface{}
	}{
		{
			query:  "INSERT INTO chat(author, alias, create_at, status) VALUES (?, ?, ?, ?)",
			params: []interface{}{64, "PROBA", "2015-06-24", data.CHAT_STATUS_ACTIVE},
		},
		{
			query: "INSERT INTO chat(author, alias, create_at, status) SELECT ?, ?, ?, ? FROM DUAL" +
				" WHERE EXISTS(SELECT c.id_chat AS \"c.id_chat\" FROM chat_user c WHERE c.id_user=?) " +
				"AND EXISTS(SELECT c.id AS \"c.id\", c.id_chat AS \"c.id_chat\", " +
				"c.id_user AS \"c.id_user\" FROM chat_message c WHERE c.id=?)",
			params: []interface{}{64, "PROBA", "2015-06-24", data.CHAT_STATUS_ACTIVE, 12345, 356},
		},
	}

	m5 := chat.NewChat()
	m5.SetActionType(data.ACTION_SAVE)
	m5.Author = 64
	m5.CreateAt = "2015-06-24"
	m5.Name = "PROBA"
	m5.Status = data.CHAT_STATUS_ACTIVE

	m1 := chat.NewChat()
	m1.SetActionType(data.ACTION_SAVE)
	m1.Author = 64
	m1.CreateAt = "2015-06-24"
	m1.Name = "PROBA"
	m1.Status = data.CHAT_STATUS_ACTIVE

	cu := chat.NewChatUser()
	cu.User = 12345
	cu.SetActionType(data.ACTION_FIND)
	cu.SetSubOperator(data.EXISTS)
	cu.SetLoadFields("Chat")

	cm := message.NewMessage()
	cm.Id = 356
	cm.SetActionType(data.ACTION_FIND)
	cm.SetSubOperator(data.EXISTS)
	cm.SetLoadFields("Id", "ChatId", "Author")

	m1.SetChatUsers(&[]data.IModel{cu})
	m1.SetMessages(&[]data.IModel{cm})

	params := []data.IModel{
		m5, m1,
	}

	adapter := NewAdapter()
	db := NewDB()

	for i := 0; i < len(params); i++ {

		adapter.sqlModel = mysql_helper.NewSQLModel(db)

		loadModelToSQLModel(params[i], adapter.sqlModel)

		adapter.sqlModel.PrepareQuery()

		if adapter.sqlModel.GetQuery() != required[i].query ||
			!functions.ArraysIsEqual(adapter.sqlModel.GetQueryParams(), required[i].params) {

			t.Fatal("Failed Check Load IModel To SQL Model! \nRequiredSQL:", required[i].query,
				"\nrezultSQL:  ", adapter.sqlModel.GetQuery(),
				"\nRequired params:", required[i].params,
				";\nrezult params:  ", adapter.sqlModel.GetQueryParams(),
			)
		}
	}
}

func TestLoadUpdateIModelToSQLModel(t *testing.T) {

	t.Log("Check Load Update IModel To SQL Model")

	required := []struct {
		query  string
		params []interface{}
	}{
		{
			query: "UPDATE chat c SET c.author=?, c.create_at=? WHERE c.id=? AND EXISTS(" +
				"SELECT c.id_chat AS \"c.id_chat\" FROM chat_user c WHERE c.client_hash=? AND c.id_user IN (?, ?))",
			params: []interface{}{64, "2015-06-24", 14, "sadsfds", 12345, 223451},
		},
	}

	m1 := chat.NewChat()
	m1.SetActionType(data.ACTION_UPDATE)
	m1.Author = 64
	m1.CreateAt = "2015-06-24"
	m1.Id = 14

	cu := chat.NewChatUser()
	cu.User = 12345
	cu.SessionHash = "sadsfds"
	cu.SetActionType(data.ACTION_FIND)
	cu.SetSubOperator(data.EXISTS)
	cu.SetLoadFields("Chat")
	cu.AddConditionField("User", data.IN)
	cu.SetArrayFields("User")

	cu1 := chat.NewChatUser()
	cu1.User = 223451
	cu.SetArrayFields("User")

	//	cu.AddLink(data.NewLink(m1.GetNameModel(), "Id", cu.GetNameModel(), "Chat", data.LINK_WEIGHT_PARENT_MORE))
	m1.SetChatUsers(&[]data.IModel{cu, cu1})

	params := []data.IModel{
		m1,
	}

	adapter := NewAdapter()
	db := NewDB()

	for i := 0; i < len(params); i++ {

		adapter.sqlModel = mysql_helper.NewSQLModel(db)

		loadModelToSQLModel(params[i], adapter.sqlModel)

		adapter.sqlModel.PrepareQuery()

		if adapter.sqlModel.GetQuery() != required[i].query ||
			!functions.ArraysIsEqual(adapter.sqlModel.GetQueryParams(), required[i].params) {

			t.Fatal("Failed Check Load IModel To SQL Model! \nRequiredSQL:", required[i].query,
				"\nrezultSQL:  ", adapter.sqlModel.GetQuery(),
				"\nRequired params:", required[i].params,
				";\nrezult params:  ", adapter.sqlModel.GetQueryParams(),
			)
		}
	}
}

func TestLoadSelectIModelToSQLModel(t *testing.T) {

	t.Log("Check Load Select IModel To SQL Model")

	required := []struct {
		query  string
		params []interface{}
	}{
		{
			query: "SELECT c.alias AS \"c.alias\", c.author AS \"c.author\", c.create_at AS \"c.create_at\", c.id AS \"c.id\"," +
				" c.status AS \"c.status\" FROM chat c WHERE c.id=? AND c.author=? AND c.create_at=?",
			params: []interface{}{14, 64, "2015-06-24"},
		},
		{
			query: "SELECT c.alias AS \"c.alias\", c.author AS \"c.author\", c.create_at AS \"c.create_at\", c.id AS \"c.id\", " +
				"c.status AS \"c.status\", ch.client_hash AS \"ch.client_hash\", ch.id_chat AS \"ch.id_chat\", ch.id_user AS \"ch.id_user\"" +
				" FROM chat c RIGHT JOIN chat_user ch ON c.id=ch.id_chat " +
				"WHERE c.id=? AND c.author=? AND c.create_at=? AND ch.id_user IN (?, ?)",
			params: []interface{}{14, 64, "2015-06-24", 12345, 12346},
		},
		{
			query: "SELECT c.alias AS \"c.alias\", c.author AS \"c.author\", c.create_at AS \"c.create_at\", c.id AS \"c.id\", " +
				"c.status AS \"c.status\", ch.client_hash AS \"ch.client_hash\", ch.id_chat AS \"ch.id_chat\", ch.id_user AS \"ch.id_user\"" +
				" FROM chat c RIGHT JOIN chat_user ch ON c.id=ch.id_chat WHERE c.id=? AND c.author=? " +
				"AND c.create_at=? AND ch.id_chat=? AND ch.id_user IN (?, ?)",
			params: []interface{}{14, 64, "2015-06-24", 90202, 123, 12345},
		},
		{
			query: "SELECT c.alias AS \"c.alias\", c.author AS \"c.author\", c.create_at AS \"c.create_at\", " +
				"c.id AS \"c.id\", c.status AS \"c.status\" FROM chat c WHERE c.id=? AND c.author=? AND c.create_at=? " +
				"AND EXISTS(SELECT c.client_hash AS \"c.client_hash\", c.id_chat AS \"c.id_chat\", " +
				"c.id_user AS \"c.id_user\" FROM chat_user c WHERE c.id_chat=? AND c.id_user=? " +
				"AND EXISTS(SELECT c.client_hash AS \"c.client_hash\", c.id_chat AS \"c.id_chat\", " +
				"c.id_user AS \"c.id_user\" FROM chat_user c WHERE c.id_user=?))",
			params: []interface{}{14, 64, "2015-06-24", 7777, 888, 99999},
		},
	}
	/* 4 query */
	m3 := chat.NewChat()
	m3.SetActionType(data.ACTION_FIND)
	m3.Author = 64
	m3.CreateAt = "2015-06-24"
	m3.Id = 14
	m3.SetLoadFields("Id", "Author", "Name", "CreateAt", "Status")

	cu3 := chat.NewChatUser()
	cu3.SetActionType(data.ACTION_FIND)
	cu3.Chat = 7777
	cu3.User = 888
	cu3.SetSubOperator(data.EXISTS)
	cu3.SetLoadFields("Chat", "User", "SessionHash")

	cu4 := chat.NewChatUser()
	cu4.User = 99999
	cu4.SetActionType(data.ACTION_FIND)
	cu4.SetSubOperator(data.EXISTS)
	cu4.SetLoadFields("Chat", "User", "SessionHash")

	cu3.SetChats(&[]data.IModel{cu4})
	m3.SetChatUsers(&[]data.IModel{cu3})

	/* 1 query */
	m := chat.NewChat()
	m.SetActionType(data.ACTION_FIND)
	m.Author = 64
	m.CreateAt = "2015-06-24"
	m.Id = 14
	m.SetLoadFields("Id", "Author", "Name", "CreateAt", "Status")

	/* 2 query */
	m1 := chat.NewChat()
	m1.SetActionType(data.ACTION_FIND)
	m1.Author = 64
	m1.CreateAt = "2015-06-24"
	m1.Id = 14
	m1.SetLoadFields("Id", "Author", "Name", "CreateAt", "Status")
	cu := chat.NewChatUser()
	cu.User = 12345
	cu.SetArrayFields("User")
	cu.AddConditionField("User", data.IN)
	cu.SetActionType(data.ACTION_FIND)
	cu.SetLoadFields("Chat", "User", "SessionHash")
	cu.AddLink(data.NewLink(m.GetNameModel(), "Id", cu.GetNameModel(), "Chat", data.LINK_WEIGHT_PARENT_MORE))

	cu5 := chat.NewChatUser()
	cu5.User = 12346
	cu.SetArrayFields("User")
	m1.SetChatUsers(&[]data.IModel{cu, cu5})

	/* 3 query */
	m2 := chat.NewChat()
	m2.SetActionType(data.ACTION_FIND)
	m2.Author = 64
	m2.CreateAt = "2015-06-24"
	m2.Id = 14
	m2.SetLoadFields("Id", "Author", "Name", "CreateAt", "Status")

	cu1 := chat.NewChatUser()
	cu1.SetArrayFields("User")
	cu1.AddConditionField("User", data.IN)
	cu1.Chat = 90202
	cu1.User = 123
	cu1.SetActionType(data.ACTION_FIND)
	cu1.SetLoadFields("Chat", "User", "SessionHash")
	cu1.AddLink(data.NewLink(m.GetNameModel(), "Id", cu.GetNameModel(), "Chat", data.LINK_WEIGHT_PARENT_MORE))

	cu2 := chat.NewChatUser()
	cu2.SetArrayFields("User")
	cu2.User = 12345

	m2.SetChatUsers(&[]data.IModel{cu1, cu2})

	params := []data.IModel{
		m, m1, m2, m3,
	}

	adapter := NewAdapter()
	db := NewDB()

	for i := 0; i < len(params); i++ {

		adapter.sqlModel = mysql_helper.NewSQLModel(db)

		loadModelToSQLModel(params[i], adapter.sqlModel)

		adapter.sqlModel.PrepareQuery()

		if adapter.sqlModel.GetQuery() != required[i].query ||
			!functions.ArraysIsEqual(adapter.sqlModel.GetQueryParams(), required[i].params) {

			t.Fatal("Failed Check Load Select IModel To SQL Model! \nRequiredSQL:", required[i].query,
				"\nrezultSQL:  ", adapter.sqlModel.GetQuery(),
				"\nRequired params:", required[i].params,
				";\nrezult params:  ", adapter.sqlModel.GetQueryParams(),
			)
		}
	}
}
