package business_logic

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/data/chat"
	db "github.com/alex988334/messenger/pkg/messenger/db"
	"github.com/alex988334/messenger/pkg/messenger/functions"
)

type ex1 struct {
	id int
}

func (p *ex1) IsAutorizate() bool {
	return true
}
func (p *ex1) GetId() int {
	return p.id
}
func (p *ex1) SetId(id int) {
	p.id = id
}

func TestExitChat(t *testing.T) {

	t.Log("Generate Exit Chat")
	//	prepare test data
	conn := db.NewDB()

	dt := NewDateTimeStamp()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Blizard\", \"Blizard\", \"$2y$1ddd.fhhdda/1111Kg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVvbmxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Rigel\", \"Rigel\", \"$2y$10$s8kkkktGjsFjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)

	chatId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(clientId)+", \"MAD Blizard\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)
	chatSupId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(userId)+", \"MAD Rigel\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

	conn.ExecSQL("INSERT INTO `chat_user`(`id_chat`, `id_user`) VALUES (?, ?)",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatSupId, clientId},
			{chatSupId, userId},
		},
	)

	defer conn.CloseDB()
	defer conn.ExecSQL("DELETE FROM `user` WHERE id=?",
		[][]interface{}{
			{clientId},
			{userId},
		},
	)

	defer conn.ExecSQL("DELETE FROM `chat` WHERE id=?",
		[][]interface{}{
			{chatId},
			{chatSupId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_user` WHERE id_chat=? AND id_user=?",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatSupId, clientId},
			{chatSupId, userId},
		},
	)

	// Test start
	required := []string{
		"{\"Status\":{\"Operation\":113,\"Status\":1,\"Message\":\"\"}," +
			"\"Chat\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(chatId) + ",\"Author\":" +
			strconv.Itoa(userId) + ",\"Name\":\"\",\"CreateAt\":\"\",\"Status\":\"\"}]," +
			"\"ChatUser\":[{\"IModel\":{},\"Chat\":" + strconv.Itoa(chatId) + ",\"User\":" +
			strconv.Itoa(clientId) + ",\"SessionHash\":\"\"}],\"Message\":null," +
			"\"MessageStatus\":null,\"User\":null,\"BlackList\":null,\"UserPhone\":null}",

		"{\"UsersId\":[" + strconv.Itoa(clientId) + "],\"ActionChats\":[" +
			strconv.Itoa(chatId) + "],\"ConnectedChats\":[" + strconv.Itoa(chatSupId) + "],\"Operation\":16}",
		"[" + strconv.Itoa(chatId) + "]",
	}

	reqParam := map[string]string{}

	ch := chat.NewChatUser()
	ch.User = userId
	ch.Chat = chatId

	str, _ := json.Marshal(ch)
	reqParam[data.MODEL_CHAT_USER] = string(str)

	st := StatusRequest{
		Operation: constants.OP_EXIT_CHAT,
	}
	str, _ = json.Marshal(st)
	reqParam[data.MODEL_STATUS] = string(str)

	bh := NewBusinessHandler(&reqParam, &ex1{id: clientId}, dt)

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
		t.Fatal("Failed Generate Exit Chat! \nRequired:", required[0], "; \nrezult:  ", string(resp))
	}

	//fmt.Println("bh.GetConnectingUser() =>", res)
	enc, _ := json.Marshal(bh.GetActionUserConnections())
	if required[1] != string(enc) {
		t.Fatal("Failed Generate Exit Chat! \nRequired:", required[1], "; \nrezult:  ", string(enc))
	}

	enc, _ = json.Marshal(bh.GetSendChatsId())
	if required[2] != string(enc) {
		t.Fatal("Failed Generate Exit Chat! \nRequired:", required[2], "; \nrezult:  ", string(enc))
	}
}

func TestListChats(t *testing.T) {

	t.Log("Generate List Chats")
	//	prepare test data
	conn := db.NewDB()

	dt := NewDateTimeStamp()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Turuk\", \"Turuk\", \"$2y$10$cKabhu.fhhdda/1111Kg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVvbmxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Murzik\", \"Murzik\", \"$2y$10$s8222265Cws2wFGjsFjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
	userSubId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Hichkok\", \"Hichkok\", \"$2y$1e$kfglr5a55552wFej4Fjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)

	chatId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(clientId)+", \"MAD Turuk\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)
	chatSupId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(userId)+", \"MAD Murzik\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

	messageChatsId := []int{chatId, chatSupId, chatId, chatSupId, chatId, chatSupId, chatId,
		chatSupId, chatId, chatSupId, chatId, chatSupId, chatId, chatSupId, chatSupId}
	messagesId := [15]int64{}
	messUsersId := []int{clientId, userId, userSubId, clientId, userId, userSubId, clientId,
		userId, userSubId, clientId, clientId, userId, userSubId, clientId, userId}
	messText := []string{
		"1111111111", "2222222222", "33333333333", "44444444444", "55555555555", "66666666666",
		"7777777777", "8888888888", "99999999999", "101010101010", "12121212212", "13131313131",
		"Last", "151515151515", "Last Sup",
	}

	conn.ExecSQL("INSERT INTO `chat_user`(`id_chat`, `id_user`) VALUES (?, ?)",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatId, userSubId},
			{chatSupId, clientId},
			{chatSupId, userId},
			{chatSupId, userSubId},
		},
	)
	for i := 0; i < len(messagesId); i++ {
		messagesId[i] = int64(conn.ExecSQL(
			"INSERT INTO `chat_message`(`id_chat`, `id_user`, `message`, `date`, `time`) VALUES (?, ?, ?, ?, ?)",
			[][]interface{}{
				{messageChatsId[i], messUsersId[i], messText[i], dt.date, dt.time},
			},
		)[0].LastId)
	}

	messStatus := [][]string{
		{constants.MESSAGE_DELIVERED, constants.MESSAGE_CREATED, constants.MESSAGE_READED},
		{constants.MESSAGE_READED, constants.MESSAGE_CREATED, constants.MESSAGE_READED},
		{constants.MESSAGE_READED, constants.MESSAGE_DELIVERED, constants.MESSAGE_READED},
		{constants.MESSAGE_READED, constants.MESSAGE_CREATED, constants.MESSAGE_DELIVERED},
		{constants.MESSAGE_CREATED, constants.MESSAGE_CREATED, constants.MESSAGE_DELIVERED},
		{constants.MESSAGE_DELIVERED, constants.MESSAGE_READED, constants.MESSAGE_READED},
		{constants.MESSAGE_READED, constants.MESSAGE_READED, constants.MESSAGE_READED},
		{constants.MESSAGE_DELIVERED, constants.MESSAGE_READED, constants.MESSAGE_DELIVERED},
		{constants.MESSAGE_DELIVERED, constants.MESSAGE_CREATED, constants.MESSAGE_READED},
		{constants.MESSAGE_READED, constants.MESSAGE_READED, constants.MESSAGE_READED},
		{constants.MESSAGE_DELIVERED, constants.MESSAGE_READED, constants.MESSAGE_READED},
		{constants.MESSAGE_CREATED, constants.MESSAGE_CREATED, constants.MESSAGE_READED},
		{constants.MESSAGE_DELIVERED, constants.MESSAGE_CREATED, constants.MESSAGE_DELIVERED},
		{constants.MESSAGE_DELIVERED, constants.MESSAGE_DELIVERED, constants.MESSAGE_DELIVERED},
		{constants.MESSAGE_DELIVERED, constants.MESSAGE_CREATED, constants.MESSAGE_READED},
	}

	for ind, mId := range messagesId {
		conn.ExecSQL("INSERT INTO `chat_message_status`(`id_message`, `id_user`, `status_message`, `date`, `time`) VALUES (?,?,?,?,?)",
			[][]interface{}{
				{mId, clientId, messStatus[ind][0], dt.date, dt.time},
				{mId, userId, messStatus[ind][1], dt.date, dt.time},
				{mId, userSubId, messStatus[ind][2], dt.date, dt.time},
			},
		)
	}

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
			{chatSupId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_user` WHERE id_chat=? AND id_user=?",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatId, userSubId},
			{chatSupId, clientId},
			{chatSupId, userId},
			{chatSupId, userSubId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_message` WHERE id=?",
		[][]interface{}{
			{messagesId[0]}, {messagesId[1]}, {messagesId[2]}, {messagesId[3]}, {messagesId[4]}, {messagesId[5]},
			{messagesId[6]}, {messagesId[7]}, {messagesId[8]}, {messagesId[9]}, {messagesId[10]}, {messagesId[11]},
			{messagesId[12]}, {messagesId[13]}, {messagesId[14]},
		},
	)
	defer func() {
		for _, val := range messagesId {
			conn.ExecSQL("DELETE FROM `chat_message_status` WHERE id_message=? AND id_user=?",
				[][]interface{}{
					{val, clientId},
					{val, userId},
					{val, userSubId},
				},
			)
		}
	}()

	// Test start
	start := time.Now().UnixMilli()

	required := []string{
		"{\"Status\":{\"Operation\":" + strconv.Itoa(constants.OP_GET_CHATS) + ",\"Status\":1,\"Message\":\"\"}," +
			"\"Chat\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(chatId) + ",\"Author\":" + strconv.Itoa(clientId) +
			",\"Name\":\"MAD Turuk\",\"CreateAt\":\"2016-08-18\",\"Status\":\"active\"},{\"IModel\":{},\"Id\":" +
			strconv.Itoa(chatSupId) + ",\"Author\":" + strconv.Itoa(userId) + ",\"Name\":\"MAD Murzik\"," +
			"\"CreateAt\":\"2016-08-18\",\"Status\":\"active\"}],\"ChatUser\":null," +
			"\"Message\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(int(messagesId[12])) +
			",\"ChatId\":" + strconv.Itoa(chatId) + ",\"Author\":" + strconv.Itoa(userSubId) +
			",\"ParrentMessage\":0,\"Message\":\"" + messText[12] + "\",\"FileUrl\":\"\",\"Date\":\"" + dt.date +
			"\",\"Time\":\"" + dt.time + "\"},{\"IModel\":{},\"Id\":" + strconv.Itoa(int(messagesId[14])) +
			",\"ChatId\":" + strconv.Itoa(chatSupId) + ",\"Author\":" + strconv.Itoa(userId) +
			",\"ParrentMessage\":0,\"Message\":\"" + messText[14] + "\",\"FileUrl\":\"\",\"Date\":\"" + dt.date +
			"\",\"Time\":\"" + dt.time + "\"}]," +
			"\"MessageStatus\":[{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[12])) +
			",\"UserId\":0,\"Status\":\"" + constants.MESSAGE_DELIVERED + "\",\"Date\":\"\",\"Time\":\"\"}," +
			"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[14])) +
			",\"UserId\":0,\"Status\":\"" + constants.MESSAGE_DELIVERED + "\",\"Date\":\"\",\"Time\":\"\"}]," +
			"\"User\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(userSubId) + ",\"Login\":\"\",\"Alias\":\"Hichkok\"," +
			"\"AuthKey\":\"\",\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0," +
			"\"CreateAt\":0,\"UpdateAt\":0,\"Avatar\":\"\"},{\"IModel\":{},\"Id\":" + strconv.Itoa(userId) +
			",\"Login\":\"\",\"Alias\":\"Murzik\",\"AuthKey\":\"\",\"PassHash\":\"\",\"PassResetToken\":\"\"," +
			"\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0,\"Avatar\":\"\"}],\"BlackList\":null,\"UserPhone\":null}",
		"[]",
	}

	reqParam := map[string]string{}

	ch := chat.NewChatUser()
	ch.User = clientId

	str, _ := json.Marshal(ch)
	reqParam[data.MODEL_CHAT_USER] = string(str)

	st := StatusRequest{
		Operation: constants.OP_GET_CHATS,
	}
	str, _ = json.Marshal(st)
	reqParam[data.MODEL_STATUS] = string(str)

	bh := NewBusinessHandler(&reqParam, &ex1{id: clientId}, dt)

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
		t.Fatal("Failed Generate List Chats! \nRequired:", required[0], "; \nrezult:  ", string(resp))
	}

	enc, _ := json.Marshal(bh.GetSendChatsId())
	if required[1] != string(enc) {
		t.Fatal("Failed Generate List Chats! \nRequired:", required[1], "; \nrezult:  ", string(enc))
	}
	t.Log("Lead time test ", float32(time.Now().UnixMilli()-start)/float32(1000), " seconds")
}

func TestCreateNewChat(t *testing.T) {

	t.Log("Generate New Chat")
	//	prepare test data
	conn := db.NewDB()

	chatName := "Mad Mark"
	dt := NewDateTimeStamp()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Mark\", \"Mark\", \"$2y$10$YKxieu.fJ0nda/ElFyKg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVzeHxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)

	rez := *(conn.SelectSQL("SELECT id FROM `chat` ORDER BY id DESC limit 1", []interface{}{}))

	chatId := 1
	if len(rez) > 0 {
		chatId = int(rez[0]["id"].(int64)) + 1
	}

	conn.ExecSQL("ALTER TABLE chat AUTO_INCREMENT="+strconv.Itoa(chatId), [][]interface{}{})

	defer conn.CloseDB()
	defer conn.ExecSQL("DELETE FROM `user` WHERE id=?",
		[][]interface{}{
			{clientId},
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
		},
	)

	// Test start
	required := []string{
		"{\"Status\":{\"Operation\":106,\"Status\":1,\"Message\":\"\"}," +
			"\"Chat\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(chatId) + ",\"Author\":" + strconv.Itoa(clientId) +
			",\"Name\":\"Mad Mark\",\"CreateAt\":\"" + functions.GetNowDateStr() + "\"," +
			"\"Status\":\"active\"}],\"ChatUser\":[{\"IModel\":{},\"Chat\":" + strconv.Itoa(chatId) + ",\"User\":" +
			strconv.Itoa(clientId) + ",\"SessionHash\":\"\"}],\"Message\":null,\"MessageStatus\":null,\"User\":null," +
			"\"BlackList\":null,\"UserPhone\":null}",

		"{\"UsersId\":[" + strconv.Itoa(clientId) + "],\"ActionChats\":[" + strconv.Itoa(chatId) +
			"],\"ConnectedChats\":[" + strconv.Itoa(chatId) + "],\"Operation\":15}",
		"[" + strconv.Itoa(chatId) + "]",
	}

	reqParam := map[string]string{}

	ch := chat.NewChat()
	ch.Author = clientId
	ch.Name = chatName
	ch.CreateAt = dt.date
	ch.Status = data.CHAT_STATUS_ACTIVE

	str, _ := json.Marshal(ch)
	reqParam[data.MODEL_CHAT] = string(str)

	st := StatusRequest{
		Operation: constants.OP_CREATE_NEW_CHAT,
	}
	str, _ = json.Marshal(st)
	reqParam[data.MODEL_STATUS] = string(str)

	bh := NewBusinessHandler(&reqParam, &ex1{id: clientId}, dt)

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
		t.Fatal("Failed Generate New Chat! \nRequired:", required[0], "; \nrezult:  ", string(resp))
	}

	//fmt.Println("bh.GetConnectingUser() =>", res)
	enc, _ := json.Marshal(bh.GetActionUserConnections())
	if required[1] != string(enc) {
		t.Fatal("Failed Generate New Chat! \nRequired:", required[1], "; \nrezult:  ", string(enc))
	}

	enc, _ = json.Marshal(bh.GetSendChatsId())
	if required[2] != string(enc) {
		t.Fatal("Failed Generate New Chat! \nRequired:", required[2], "; \nrezult:  ", string(enc))
	}
}

func TestRemoveChat(t *testing.T) {

	t.Log("Generate Remove Chat")
	//	prepare test data
	conn := db.NewDB()

	dt := NewDateTimeStamp()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Hazard\", \"Hazard\", \"$2y$10$ggghu.fhhdda/1111Kg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVvbmxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Stark\", \"Stark\", \"$2y$10$s82222ttttGjsFjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)

	chatId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(clientId)+", \"MAD Hazard\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)
	chatSupId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(userId)+", \"MAD Stark\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

	messageChatsId := []int{chatId, chatId, chatId}
	messagesId := [3]int64{}
	messUsersId := []int{clientId, userId, userId}
	messText := []string{
		"1111111111", "2222222222", "33333333333",
	}

	conn.ExecSQL("INSERT INTO `chat_user`(`id_chat`, `id_user`) VALUES (?, ?)",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatSupId, clientId},
			{chatSupId, userId},
		},
	)
	for i := 0; i < len(messagesId); i++ {
		messagesId[i] = int64(conn.ExecSQL(
			"INSERT INTO `chat_message`(`id_chat`, `id_user`, `message`, `date`, `time`) VALUES (?, ?, ?, ?, ?)",
			[][]interface{}{
				{messageChatsId[i], messUsersId[i], messText[i], dt.date, dt.time},
			},
		)[0].LastId)
	}

	for i := 0; i < len(messagesId); i++ {

		conn.ExecSQL(
			"INSERT INTO `chat_message_status`(`id_message`, `id_user`, `status_message`, `date`, `time`) VALUES (?, ?, ?, ?, ?)",
			[][]interface{}{
				{messagesId[i], clientId, "created", dt.date, dt.time},
			},
		)
		conn.ExecSQL(
			"INSERT INTO `chat_message_status`(`id_message`, `id_user`, `status_message`, `date`, `time`) VALUES (?, ?, ?, ?, ?)",
			[][]interface{}{
				{messagesId[i], userId, "created", dt.date, dt.time},
			},
		)
	}

	defer conn.CloseDB()
	defer conn.ExecSQL("DELETE FROM `user` WHERE id=?",
		[][]interface{}{
			{clientId},
			{userId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat` WHERE id=?",
		[][]interface{}{
			{chatId},
			{chatSupId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_user` WHERE id_chat=? AND id_user=?",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatSupId, clientId},
			{chatSupId, userId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_message` WHERE id=?",
		[][]interface{}{
			{messagesId[0]}, {messagesId[1]}, {messagesId[2]},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_message_status` WHERE id_message=?",
		[][]interface{}{
			{messagesId[0]}, {messagesId[1]}, {messagesId[2]},
		},
	)

	// Test start
	required := []string{
		"{\"Status\":{\"Operation\":116,\"Status\":1,\"Message\":\"\"}," +
			"\"Chat\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(chatId) + ",\"Author\":" +
			strconv.Itoa(clientId) + ",\"Name\":\"\",\"CreateAt\":\"\",\"Status\":\"\"}]," +
			"\"ChatUser\":null,\"Message\":null,\"MessageStatus\":null,\"User\":null,\"BlackList\":null,\"UserPhone\":null}",

		"{\"UsersId\":[" + strconv.Itoa(clientId) + "," + strconv.Itoa(userId) + "],\"ActionChats\":[" +
			strconv.Itoa(chatId) + "],\"ConnectedChats\":[" + strconv.Itoa(chatSupId) + "],\"Operation\":17}",
		"[" + strconv.Itoa(chatId) + "]",
	}

	reqParam := map[string]string{}

	ch := chat.NewChat()
	ch.Id = chatId
	ch.Author = clientId

	str, _ := json.Marshal(ch)
	reqParam[data.MODEL_CHAT] = string(str)

	st := StatusRequest{
		Operation: constants.OP_REMOVE_CHAT,
	}
	str, _ = json.Marshal(st)
	reqParam[data.MODEL_STATUS] = string(str)

	bh := NewBusinessHandler(&reqParam, &ex1{id: clientId}, dt)

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
		t.Fatal("Failed Generate Remove Chat! \nRequired:", required[0], "; \nrezult:  ", string(resp))
	}

	//fmt.Println("bh.GetConnectingUser() =>", res)
	enc, _ := json.Marshal(bh.GetActionUserConnections())
	if required[1] != string(enc) {
		t.Fatal("Failed Generate Remove Chat! \nRequired:", required[1], "; \nrezult:  ", string(enc))
	}

	enc, _ = json.Marshal(bh.GetSendChatsId())
	if required[2] != string(enc) {
		t.Fatal("Failed Generate Remove Chat! \nRequired:", required[2], "; \nrezult:  ", string(enc))
	}

	rez := conn.SelectSQL("SELECT `id` FROM `chat` WHERE id=?", []interface{}{chatId})
	if len(*rez) > 0 {
		t.Fatal("Failed Remove Chat! Chat is finded: " + strconv.Itoa(chatId))
	}

	p := [][]interface{}{
		{chatId, clientId},
		{chatId, userId},
	}

	for i := 0; i < len(p); i++ {

		rez = conn.SelectSQL("SELECT `id_chat` FROM `chat_user` WHERE id_chat=? AND id_user=?", p[i])
		if len(*rez) > 0 {
			t.Fatal("Failed Remove Chat! ChatUser is finded: chat " + strconv.Itoa(p[i][0].(int)) +
				" user " + strconv.Itoa(p[i][0].(int)))
		}
	}

	for i := 0; i < len(messagesId); i++ {

		rez = conn.SelectSQL("SELECT `id` FROM `chat_message` WHERE id=?", []interface{}{messagesId[i]})
		if len(*rez) > 0 {
			t.Fatal("Failed Remove Chat! ChatMessage is finded: message " + strconv.Itoa(int(messagesId[i])))
		}
	}

	for i := 0; i < len(messagesId); i++ {

		rez = conn.SelectSQL("SELECT `id_message` FROM `chat_message_status` WHERE id_message=?", []interface{}{messagesId[i]})
		if len(*rez) > 0 {
			t.Fatal("Failed Remove Chat! MessageStatus is finded: message " + strconv.Itoa(int(messagesId[i])))
		}
	}
}

/*
func TestCreateNewChat(t *testing.T) {

		t.Log("Generate New Chat")
		//	prepare test data
		conn := db.NewDB()

		chatName := "Mad Mark"
		dt := NewDateTimeStamp()

		clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
			"VALUES (\"Mark\", \"Mark\", \"$2y$10$YKxieu.fJ0nda/ElFyKg0uk5R\", "+
			"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVzeHxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)

		rez := *(conn.SelectSQL("SELECT id FROM `chat` ORDER BY id DESC limit 1", []interface{}{}))

		chatId := 1
		if len(rez) > 0 {
			chatId = int(rez[0]["id"].(int64)) + 1
		}

		conn.ExecSQL("ALTER TABLE chat AUTO_INCREMENT="+strconv.Itoa(chatId), [][]interface{}{})

		defer conn.CloseDB()
		defer conn.ExecSQL("DELETE FROM `user` WHERE id=?",
			[][]interface{}{
				{clientId},
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
			},
		)

		// Test start
		required := []string{
			"{\"Chat\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(chatId) + ",\"Author\":" + strconv.Itoa(clientId) +
				",\"Name\":\"Mad Mark\",\"CreateAt\":\"" + functions.GetNowDateStr() + "\"," +
				"\"Status\":\"active\"}],\"ChatUser\":[{\"IModel\":{},\"Chat\":" + strconv.Itoa(chatId) + ",\"User\":" +
				strconv.Itoa(clientId) + ",\"SessionHash\":\"\"}]," +
				"\"Status\":{\"Operation\":106,\"Status\":1,\"Message\":\"\"}}",

			"{\"UserId\":" + strconv.Itoa(clientId) + ",\"ConnectingChats\":[" + strconv.Itoa(chatId) +
				"],\"ConnectedChats\":[" + strconv.Itoa(chatId) + "]}",
			"[" + strconv.Itoa(chatId) + "]",
		}

		reqParam := map[string][]byte{}

		ch := chat.NewChat()
		ch.Author = clientId
		ch.Name = chatName
		ch.CreateAt = dt.date
		ch.Status = data.CHAT_STATUS_ACTIVE

		str, _ := json.Marshal(ch)
		reqParam[data.MODEL_CHAT] = str

		st := StatusRequest{
			Operation: constants.OP_CREATE_NEW_CHAT,
		}
		str, _ = json.Marshal(st)
		reqParam[data.MODEL_STATUS] = str

		bh := NewBusinessHandler(&reqParam, &ex1{id: clientId}, dt)

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
			t.Fatal("Failed Generate New Chat! \nRequired:", required[0], "; \nrezult:  ", string(resp))
		}

		//fmt.Println("bh.GetConnectingUser() =>", res)
		enc, _ := json.Marshal(bh.GetActionUserConnections())
		if required[1] != string(enc) {
			t.Fatal("Failed Generate New Chat! \nRequired:", required[1], "; \nrezult:  ", string(enc))
		}

		enc, _ = json.Marshal(bh.GetSendChatsId())
		if required[2] != string(enc) {
			t.Fatal("Failed Generate New Chat! \nRequired:", required[2], "; \nrezult:  ", string(enc))
		}
	}
*/
func TestAddUserInChat(t *testing.T) {

	t.Log("Generate Add User In Chat")
	//	prepare test data
	conn := db.NewDB()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Mikola\", \"Mikolay\", \"$2y$10$YKadjf.fJ0nda/ElFyKg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVzeHxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Djon\", \"Djon\", \"$2y$10$2hllr5a5Cwe2wFGj4Fjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
	//fmt.Println("clientId =>", clientId, "userId =>", userId)

	chatId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(clientId)+", \"MAD Mikola\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

	chatSuppId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(userId)+", \"MAD Djon\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)
	//	fmt.Println("chatId =>", chatId, "chatSuppId =>", chatSuppId)

	conn.ExecSQL("INSERT INTO `chat_user`(`id_chat`, `id_user`) VALUES (?, ?)",
		[][]interface{}{
			{chatId, clientId},
			{chatSuppId, userId},
		},
	)
	defer conn.CloseDB()
	defer conn.ExecSQL("DELETE FROM `user` WHERE id=?",
		[][]interface{}{
			{userId},
			{clientId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat` WHERE id=?",
		[][]interface{}{
			{chatId},
			{chatSuppId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_user` WHERE id_chat=? AND id_user=?",
		[][]interface{}{
			{chatId, clientId},
			{chatSuppId, userId},
			{chatId, userId},
		},
	)

	// Test start
	required := []string{
		"{\"Status\":{\"Operation\":115,\"Status\":1,\"Message\":\"\"},\"Chat\":null," +
			"\"ChatUser\":[{\"IModel\":{},\"Chat\":" + strconv.Itoa(chatId) + ",\"User\":" + strconv.Itoa(userId) +
			",\"SessionHash\":\"\"}],\"Message\":null,\"MessageStatus\":null,\"User\":[{\"IModel\":{}," +
			"\"Id\":" + strconv.Itoa(userId) + ",\"Login\":\"\",\"Alias\":\"Djon\",\"AuthKey\":\"\",\"PassHash\":\"\"," +
			"\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0,\"Avatar\":\"\"}]," +
			"\"BlackList\":null,\"UserPhone\":null}",
		"{\"UsersId\":[" + strconv.Itoa(userId) + "],\"ActionChats\":[" + strconv.Itoa(chatId) +
			"],\"ConnectedChats\":[" + strconv.Itoa(chatId) + "," + strconv.Itoa(chatSuppId) + "],\"Operation\":15}",
		"[" + strconv.Itoa(chatId) + "]",
	}

	reqParam := map[string]string{}

	chUs := chat.NewChatUser()
	chUs.User = userId
	chUs.Chat = chatId

	str, _ := json.Marshal(chUs)
	reqParam[data.MODEL_CHAT_USER] = string(str)

	st := StatusRequest{
		Operation: constants.OP_ADD_USER,
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
		t.Fatal("Failed Add User In Chat! \nRequired:", required[0], "; \nrezult:  ", string(resp))
	}

	enc, _ := json.Marshal(*bh.GetActionUserConnections())
	if required[1] != string(enc) {
		t.Fatal("Failed Add User In Chat! \nRequired:", required[1], "; \nrezult:  ", string(enc))
	}

	enc, _ = json.Marshal(bh.GetSendChatsId())
	if required[2] != string(enc) {
		t.Fatal("Failed Add User In Chat! \nRequired:", required[2], "; \nrezult:  ", string(enc))
	}
}

/*
func TestAddUserInChat(t *testing.T) {

		t.Log("Generate Add User In Chat")
		//	prepare test data
		conn := db.NewDB()

		clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
			"VALUES (\"Mikola\", \"Mikolay\", \"$2y$10$YKadjf.fJ0nda/ElFyKg0uk5R\", "+
			"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVzeHxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
		userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
			"`created_at`, `updated_at`) VALUES (\"Djon\", \"Djon\", \"$2y$10$2hllr5a5Cwe2wFGj4Fjcx.Xzh\","+
			" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
		//fmt.Println("clientId =>", clientId, "userId =>", userId)

		chatId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
			"VALUES ("+strconv.Itoa(clientId)+", \"MAD Mikola\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

		chatSuppId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
			"VALUES ("+strconv.Itoa(userId)+", \"MAD Djon\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)
		//	fmt.Println("chatId =>", chatId, "chatSuppId =>", chatSuppId)

		conn.ExecSQL("INSERT INTO `chat_user`(`id_chat`, `id_user`) VALUES (?, ?)",
			[][]interface{}{
				{chatId, clientId},
				{chatSuppId, userId},
			},
		)
		defer conn.CloseDB()
		defer conn.ExecSQL("DELETE FROM `user` WHERE id=?",
			[][]interface{}{
				{userId},
				{clientId},
			},
		)
		defer conn.ExecSQL("DELETE FROM `chat` WHERE id=?",
			[][]interface{}{
				{chatId},
				{chatSuppId},
			},
		)
		defer conn.ExecSQL("DELETE FROM `chat_user` WHERE id_chat=? AND id_user=?",
			[][]interface{}{
				{chatId, clientId},
				{chatSuppId, userId},
				{chatId, userId},
			},
		)

		// Test start
		required := []string{
			"{\"ChatUser\":[{\"IModel\":{},\"Chat\":" + strconv.Itoa(chatId) + ",\"User\":" + strconv.Itoa(userId) +
				",\"SessionHash\":\"\"}]," +
				"\"Status\":{\"Operation\":115,\"Status\":1,\"Message\":\"\"},\"User\":[{\"IModel\":{}," +
				"\"Id\":" + strconv.Itoa(userId) + ",\"Login\":\"\",\"Alias\":\"Djon\",\"AuthKey\":\"\",\"PassHash\":\"\"," +
				"\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0}]}",
			"{\"UserId\":" + strconv.Itoa(userId) + ",\"ConnectingChats\":[" + strconv.Itoa(chatId) +
				"],\"ConnectedChats\":[" + strconv.Itoa(chatId) + "," + strconv.Itoa(chatSuppId) + "]}",
			"[" + strconv.Itoa(chatId) + "]",
		}

		reqParam := map[string][]byte{}

		chUs := chat.NewChatUser()
		chUs.User = userId
		chUs.Chat = chatId

		str, _ := json.Marshal(chUs)
		reqParam[data.MODEL_CHAT_USER] = str

		st := StatusRequest{
			Operation: constants.OP_ADD_USER,
		}
		str, _ = json.Marshal(st)
		reqParam[data.MODEL_STATUS] = str

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
			t.Fatal("Failed Add User In Chat! \nRequired:", required[0], "; \nrezult:  ", string(resp))
		}

		enc, _ := json.Marshal(*bh.GetActionUserConnections())
		if required[1] != string(enc) {
			t.Fatal("Failed Add User In Chat! \nRequired:", required[1], "; \nrezult:  ", string(enc))
		}

		enc, _ = json.Marshal(bh.GetSendChatsId())
		if required[2] != string(enc) {
			t.Fatal("Failed Add User In Chat! \nRequired:", required[2], "; \nrezult:  ", string(enc))
		}
	}
*/
func TestRemoveUserFromChat(t *testing.T) {

	t.Log("Generate Remove User From Chat")
	//	prepare test data
	conn := db.NewDB()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Major\", \"Major\", \"$2y$10rrrjf.fJ0nda/ElFyKg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVzeHxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Snork\", \"Snork\", \"$2y$10$2hllrrrre2wFGj4Fjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
	userSupId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Lun\", \"Lun\", \"rrr2hllr5a5Cwe2wFGj4Fjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)

	chatId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(clientId)+", \"MAD Major\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

	chatSuppId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(userId)+", \"MAD Snork\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

	conn.ExecSQL("INSERT INTO `chat_user`(`id_chat`, `id_user`) VALUES (?, ?)",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatId, userSupId},
			{chatSuppId, userId},
			{chatSuppId, clientId},
			{chatSuppId, userSupId},
		},
	)
	defer conn.CloseDB()
	defer conn.ExecSQL("DELETE FROM `user` WHERE id=?",
		[][]interface{}{
			{userId},
			{clientId},
			{userSupId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat` WHERE id=?",
		[][]interface{}{
			{chatId},
			{chatSuppId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_user` WHERE id_chat=? AND id_user=?",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatId, userSupId},
			{chatSuppId, userId},
			{chatSuppId, clientId},
			{chatSuppId, userSupId},
		},
	)

	// Test start
	required := []string{
		"{\"Status\":{\"Operation\":114,\"Status\":1,\"Message\":\"\"},\"Chat\":null," +
			"\"ChatUser\":[{\"IModel\":{},\"Chat\":" + strconv.Itoa(chatId) + ",\"User\":" + strconv.Itoa(userId) +
			",\"SessionHash\":\"\"}],\"Message\":null,\"MessageStatus\":null,\"User\":null,\"BlackList\":null,\"UserPhone\":null}",
		"{\"UsersId\":[" + strconv.Itoa(userId) + "],\"ActionChats\":[" + strconv.Itoa(chatId) +
			"],\"ConnectedChats\":[" + strconv.Itoa(chatSuppId) + "],\"Operation\":16}",
		"[" + strconv.Itoa(chatId) + "]",
	}

	reqParam := map[string]string{}

	chUs := chat.NewChatUser()
	chUs.User = userId
	chUs.Chat = chatId

	str, _ := json.Marshal(chUs)
	reqParam[data.MODEL_CHAT_USER] = string(str)

	st := StatusRequest{
		Operation: constants.OP_REMOVE_USER,
	}
	str, _ = json.Marshal(st)
	reqParam[data.MODEL_STATUS] = string(str)

	bh := NewBusinessHandler(&reqParam, &ex1{id: userId}, nil)

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
		t.Fatal("Failed Remove User From Chat! \nRequired:", required[0], "; \nrezult:  ", string(resp))
	}

	enc, _ := json.Marshal(*bh.GetActionUserConnections())
	if required[1] != string(enc) {
		t.Fatal("Failed Remove User From Chat! \nRequired:", required[1], "; \nrezult:  ", string(enc))
	}

	enc, _ = json.Marshal(bh.GetSendChatsId())
	if required[2] != string(enc) {
		t.Fatal("Failed Remove User From Chat! \nRequired:", required[2], "; \nrezult:  ", string(enc))
	}
}
