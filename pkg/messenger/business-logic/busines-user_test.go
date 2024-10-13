package business_logic

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/data/chat"
	"github.com/alex988334/messenger/pkg/messenger/data/user"
	db "github.com/alex988334/messenger/pkg/messenger/db"
)

func TestRegistration(t *testing.T) {

	t.Log("Generate Registration")

	dt := NewDateTimeStamp()

	reqParam1 := map[string]string{}

	us := user.NewUser()
	us.Login = "ObivanKinobi"
	us.Alias = "DjedayMadYoda"
	us.PassHash = "DjedayMadYodaObivanKinobi"
	us.Email = "DjedayMadYoda@gmail.com"

	required := []string{
		"{\"Status\":{\"Operation\":" + strconv.Itoa(constants.OP_REGISTRATION) +
			",\"Status\":1,\"Message\":\"\"},\"Chat\":null,\"ChatUser\":null," +
			"\"Message\":null,\"MessageStatus\":null,\"User\":[{\"IModel\":{},\"Id\":0,\"Login\":\"" +
			us.Login + "\",\"Alias\":\"" + us.Alias + "\",\"AuthKey\":\"",

		"\",\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"" + us.Email +
			"\",\"Status\":0,\"CreateAt\":" + strconv.Itoa(int(dt.unixDate)) +
			",\"UpdateAt\":" + strconv.Itoa(int(dt.unixDate)) + ",\"Avatar\":\"\"}],\"BlackList\":null,\"UserPhone\":null}",

		"null",
		"[]",
	}

	str, _ := json.Marshal(us)
	reqParam1[data.MODEL_USER] = string(str)

	st := StatusRequest{
		Operation: constants.OP_REGISTRATION,
	}
	str, _ = json.Marshal(st)
	reqParam1[data.MODEL_STATUS] = string(str)

	bh := NewBusinessHandler(&reqParam1, &ex1{id: 0}, dt)

	resp, _, err := bh.ProcessinRequest()
	if err != nil {
		t.Log("ERROR bh.ProcessinRequest() => " + err.Error())
	}

	b := map[string]any{}
	err = json.Unmarshal(resp, &b)
	if err != nil {
		t.Log("ERRROR RESPONSE =>", err)
	}

	conn := db.NewDB()
	defer conn.CloseDB()
	defer conn.ExecSQL("DELETE FROM `user` WHERE username=?", [][]interface{}{{us.Login}})

	rez := *(conn.SelectSQL("SELECT `auth_key` FROM `user` WHERE username=?", []interface{}{us.Login}))
	if len(rez) == 0 {
		t.Fatal("Failed Generate Registration! \nRequired:", required[0]+required[1], "; \nrezult:  ", string(resp))
	}

	auth := rez[0]["auth_key"].(string)
	if required[0]+auth+required[1] != string(resp) {
		t.Fatal("Failed Generate Registration! \nRequired:", required[0]+auth+required[1], "; \nrezult:  ", string(resp))
	}

	enc, _ := json.Marshal(bh.GetActionUserConnections())
	if required[2] != string(enc) {
		t.Fatal("Failed Generate Registration! \nRequired:", required[2], "; \nrezult:  ", string(enc))
	}

	enc, _ = json.Marshal(bh.GetSendChatsId())
	if required[3] != string(enc) {
		t.Fatal("Failed Generate Registration! \nRequired:", required[3], "; \nrezult:  ", string(enc))
	}
}

func TestMyData(t *testing.T) {

	t.Log("Generate My Data")
	//	prepare test data
	conn := db.NewDB()

	dt := NewDateTimeStamp()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Kruger\", \"Mad Kruger\", \"$2y$1ddd.fhh23233Kg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVvbmxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	phones := [][]interface{}{
		{clientId, "+70632345678"},
		{clientId, "+70632345679"},
		{clientId, "+70632345680"},
	}

	for i := 0; i < len(phones); i++ {
		conn.ExecSQL("INSERT INTO `user_phone`(`user_id`, `phone`) VALUES (?, ?)", [][]interface{}{phones[i]})
	}

	defer conn.CloseDB()
	defer conn.ExecSQL("DELETE FROM `user` WHERE id=?", [][]interface{}{{clientId}})
	defer func() {
		for i := 0; i < len(phones); i++ {
			conn.ExecSQL("DELETE FROM `user_phone` WHERE user_id=? AND phone=?",
				[][]interface{}{phones[i]})
		}
	}()

	// Test start
	required := []string{
		"{\"Status\":{\"Operation\":122,\"Status\":1,\"Message\":\"\"}," +
			"\"Chat\":null,\"ChatUser\":null,\"Message\":null,\"MessageStatus\":null," +
			"\"User\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(clientId) + ",\"Login\":\"Kruger\"," +
			"\"Alias\":\"Mad Kruger\",\"AuthKey\":\"\",\"PassHash\":\"\",\"PassResetToken\":\"\"," +
			"\"Email\":\"\",\"Status\":10,\"CreateAt\":1540361413,\"UpdateAt\":1563035897,\"Avatar\":\"\"}],\"BlackList\":null," +
			"\"UserPhone\":[{\"IModel\":{},\"UserId\":" + strconv.Itoa(clientId) +
			",\"Phone\":\"+70632345678\"},{\"IModel\":{},\"UserId\":" + strconv.Itoa(clientId) +
			",\"Phone\":\"+70632345679\"},{\"IModel\":{},\"UserId\":" + strconv.Itoa(clientId) +
			",\"Phone\":\"+70632345680\"}]}",
		"null",
		"[]",
	}

	reqParam1 := map[string]string{}

	st := StatusRequest{
		Operation: constants.OP_MY_DATA,
	}
	str, _ := json.Marshal(st)
	reqParam1[data.MODEL_STATUS] = string(str)

	bh := NewBusinessHandler(&reqParam1, &ex1{id: clientId}, dt)

	resp, _, err := bh.ProcessinRequest()
	if err != nil {
		t.Log("ERROR bh.ProcessinRequest() => " + err.Error())
	}

	b := map[string]any{}
	err = json.Unmarshal(resp, &b)
	if err != nil {
		t.Log("ERRROR RESPONSE =>", err)
	}

	if required[0] != string(resp) {
		t.Fatal("Failed Generate My Data! \nRequired:", required[0], "; \nrezult:  ", string(resp))
	}

	enc, _ := json.Marshal(bh.GetActionUserConnections())
	if required[1] != string(enc) {
		t.Fatal("Failed Generate My Data! \nRequired:", required[1], "; \nrezult:  ", string(enc))
	}

	enc, _ = json.Marshal(bh.GetSendChatsId())
	if required[2] != string(enc) {
		t.Fatal("Failed Generate My Data! \nRequired:", required[2], "; \nrezult:  ", string(enc))
	}
}

func TestSearchUser(t *testing.T) {

	t.Log("Generate Search User")
	//	prepare test data
	conn := db.NewDB()

	dt := NewDateTimeStamp()

	alias1 := "NilsAlias"
	alias2 := "Hans32"
	alias3 := "Hans Mad"

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Nils\", \""+alias1+"\", \"$2y$1ddd.fhh23233Kg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVvbmxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Hans\", \""+alias2+"\", \"$2y$10$s8k54545cx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
	userSupId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Hans43243\", \""+alias3+"\", \"$2y$1338kk34343Fjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
	phones := [][]interface{}{
		{clientId, "+70632345678"},
		{clientId, "+70632345679"},
		{clientId, "+70632345680"},
		{userId, "+70632345681"},
		{userId, "+70632345682"},
		{userId, "+70632345683"},
		{userSupId, "+70632345634"},
		{userSupId, "+70632345645"},
		{userSupId, "+70632345625"},
	}

	for i := 0; i < len(phones); i++ {
		conn.ExecSQL("INSERT INTO `user_phone`(`user_id`, `phone`) VALUES (?, ?)", [][]interface{}{phones[i]})
	}

	defer conn.CloseDB()
	defer conn.ExecSQL("DELETE FROM `user` WHERE id=?",
		[][]interface{}{
			{clientId},
			{userId},
			{userSupId},
		},
	)
	defer func() {
		for i := 0; i < len(phones); i++ {
			conn.ExecSQL("DELETE FROM `user_phone` WHERE user_id=? AND phone=?",
				[][]interface{}{phones[i]})
		}
	}()

	// Test start
	required := [][]string{
		{
			"{\"Status\":{\"Operation\":" + strconv.Itoa(constants.OP_SEARCH_USER) + ",\"Status\":1,\"Message\":\"\"}," +
				"\"Chat\":null,\"ChatUser\":null,\"Message\":null,\"MessageStatus\":null," +
				"\"User\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(userId) + ",\"Login\":\"\",\"Alias\":\"" +
				alias2 + "\",\"AuthKey\":\"\",\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"\"," +
				"\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0,\"Avatar\":\"\"}],\"BlackList\":null,\"UserPhone\":null}",
			"null",
			"[]",
		},
		{
			"{\"Status\":{\"Operation\":" + strconv.Itoa(constants.OP_SEARCH_USER) + ",\"Status\":1,\"Message\":\"\"}," +
				"\"Chat\":null,\"ChatUser\":null,\"Message\":null,\"MessageStatus\":null,\"User\":null,\"BlackList\":null," +
				"\"UserPhone\":[{\"IModel\":{},\"UserId\":" + strconv.Itoa(userSupId) + ",\"Phone\":\"" +
				phones[8][1].(string) + "\"}]}",
			"null",
			"[]",
		},
	}

	reqParam1 := map[string]string{}

	ch := user.NewUser()
	ch.Alias = alias2

	str, _ := json.Marshal(ch)
	reqParam1[data.MODEL_USER] = string(str)

	st := StatusRequest{
		Operation: constants.OP_SEARCH_USER,
	}
	str, _ = json.Marshal(st)
	reqParam1[data.MODEL_STATUS] = string(str)

	reqParam2 := map[string]string{}

	up := user.NewUserPhone()
	up.Phone = phones[8][1].(string)

	str1, _ := json.Marshal(up)
	reqParam2[data.MODEL_USER_PHONE] = string(str1)

	st1 := StatusRequest{
		Operation: constants.OP_SEARCH_USER,
	}
	str1, _ = json.Marshal(st1)
	reqParam2[data.MODEL_STATUS] = string(str1)

	params := []map[string]string{reqParam1, reqParam2}

	for i := 0; i < len(params); i++ {

		bh := NewBusinessHandler(&params[i], &ex1{id: clientId}, dt)

		resp, _, err := bh.ProcessinRequest()
		if err != nil {
			t.Log("ERROR bh.ProcessinRequest() => " + err.Error())
		}

		b := map[string]any{}
		err = json.Unmarshal(resp, &b)
		if err != nil {
			t.Log("ERRROR RESPONSE =>", err)
		}

		if required[i][0] != string(resp) {
			t.Fatal("Failed Generate Search User! \nRequired:", required[i][0], "; \nrezult:  ", string(resp))
		}

		//fmt.Println("bh.GetConnectingUser() =>", res)
		enc, _ := json.Marshal(bh.GetActionUserConnections())
		if required[i][1] != string(enc) {
			t.Fatal("Failed Generate Search User! \nRequired:", required[i][1], "; \nrezult:  ", string(enc))
		}

		enc, _ = json.Marshal(bh.GetSendChatsId())
		if required[i][2] != string(enc) {
			t.Fatal("Failed Generate Search User! \nRequired:", required[i][2], "; \nrezult:  ", string(enc))
		}
	}
}

func TestListChatUsers(t *testing.T) {

	t.Log("Take List Chat Users")
	//	prepare test data
	conn := db.NewDB()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Ahiles\", \"Ahiles\", \"$2y$10$YKattttfJ0nda/ElFyKg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVzeHxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Marmaduk\", \"Marmaduk\", \"$2y$1$6rrrr5a5Cwe2wFGj4Fjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
	userSubId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Hulk\", \"Hulk\", \"$2y$10$eeeer5a5Cwe2wFGj4Fjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
	//fmt.Println("clientId =>", clientId, "userId =>", userId)

	chatId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(clientId)+", \"MAD Ahiles\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

	conn.ExecSQL("INSERT INTO `chat_user`(`id_chat`, `id_user`) VALUES (?, ?)",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatId, userSubId},
		},
	)

	defer conn.CloseDB()
	defer conn.ExecSQL("DELETE FROM `user` WHERE id=?",
		[][]interface{}{
			{clientId},
			{userId},
			{userSubId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat` WHERE id=?",
		[][]interface{}{
			{chatId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_user` WHERE id_chat=? AND id_user=?",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatId, userSubId},
		},
	)

	// Test start
	required := []string{
		"{\"Status\":{\"Operation\":105,\"Status\":1,\"Message\":\"\"},\"Chat\":null," +
			"\"ChatUser\":[{\"IModel\":{},\"Chat\":" + strconv.Itoa(chatId) + ",\"User\":" +
			strconv.Itoa(clientId) + ",\"SessionHash\":\"\"},{\"IModel\":{},\"Chat\":" + strconv.Itoa(chatId) +
			",\"User\":" + strconv.Itoa(userId) + ",\"SessionHash\":\"\"},{\"IModel\":{},\"Chat\":" +
			strconv.Itoa(chatId) + ",\"User\":" + strconv.Itoa(userSubId) + ",\"SessionHash\":\"\"}]," +
			"\"Message\":null,\"MessageStatus\":null,\"User\":[{\"IModel\":{},\"Id\":" +
			strconv.Itoa(clientId) + ",\"Login\":\"\",\"Alias\":\"Ahiles\",\"AuthKey\":\"\",\"PassHash\":\"\"," +
			"\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0,\"Avatar\":\"\"}," +
			"{\"IModel\":{},\"Id\":" + strconv.Itoa(userId) + ",\"Login\":\"\",\"Alias\":\"Marmaduk\",\"AuthKey\":\"\"," +
			"\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0," +
			"\"UpdateAt\":0,\"Avatar\":\"\"},{\"IModel\":{},\"Id\":" + strconv.Itoa(userSubId) +
			",\"Login\":\"\",\"Alias\":\"Hulk\",\"AuthKey\":\"\"" +
			",\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0,\"Avatar\":\"\"}]," +
			"\"BlackList\":null,\"UserPhone\":null}",
		"[]",
	}

	reqParam := map[string]string{}

	m := chat.NewChatUser()
	m.User = clientId
	m.Chat = chatId

	str, _ := json.Marshal(m)
	reqParam[data.MODEL_CHAT_USER] = string(str)

	st := StatusRequest{
		Operation: constants.OP_LIST_USERS,
	}
	str, _ = json.Marshal(st)
	reqParam[data.MODEL_STATUS] = string(str)

	bh := NewBusinessHandler(&reqParam, &ex1{id: clientId}, nil)

	resp, _, err := bh.ProcessinRequest()
	if err != nil {
		t.Log("ERROR bh.ProcessinRequest() => " + err.Error())
	}

	b := map[string]any{}
	err = json.Unmarshal(resp, &b)
	if err != nil {
		t.Log("ERRROR RESPONSE =>", err)
	}

	if required[0] != string(resp) {
		t.Fatal("Failed List Chat Users! \nRequired:", required[0], "; \nrezult:  ", string(resp))
	}

	enc, _ := json.Marshal(bh.GetSendChatsId())
	if required[1] != string(enc) {
		t.Fatal("Failed List Chat Users! \nRequired:", required[1], "; \nrezult:  ", string(enc))
	}
}
