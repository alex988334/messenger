package business_logic

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/data/message"
	db "github.com/alex988334/messenger/pkg/messenger/db"
)

func TestListMessages(t *testing.T) {

	t.Log("Generate List Messages")
	//	prepare test data
	conn := db.NewDB()

	dt := NewDateTimeStamp()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Leonard\", \"Leonard\", \"$2y$10$cbbbb.fhhdda/1111Kg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVvbmxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Sheldon\", \"Sheldon\", \"$2y$10$s822226nnnnnFGjsFjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
	userSubId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Hovard\", \"Hovard\", \"$2y$1e$kfeee5a55552wFej4Fjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)

	chatId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(clientId)+", \"MAD Leonard\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

	messageChatsId := []int{chatId, chatId, chatId, chatId, chatId, chatId, chatId,
		chatId, chatId, chatId, chatId, chatId, chatId, chatId, chatId, chatId, chatId, chatId, chatId, chatId, chatId, chatId,
		chatId, chatId, chatId, chatId, chatId, chatId, chatId, chatId}

	messagesId := make([]int64, len(messageChatsId))
	messUsersId := []int{clientId, userId, userSubId, clientId, userId, userSubId, clientId,
		userId, userSubId, clientId, clientId, userId, userSubId, clientId, userId,
		clientId, userId, userSubId, clientId, userId, userSubId, clientId,
		userId, userSubId, clientId, clientId, userId, userSubId, clientId, userId,
	}

	messText := []string{
		"1111111111", "2222222222", "33333333333", "44444444444", "55555555555", "66666666666",
		"7777777777", "8888888888", "99999999999", "101010101010", "12121212212", "13131313131",
		"Last", "151515151515", "Last Sup",
		"1111111111", "2222222222", "33333333333", "44444444444", "55555555555", "66666666666",
		"7777777777", "8888888888", "99999999999", "101010101010", "12121212212", "13131313131",
		"Last", "151515151515", "Last Sup",
	}

	conn.ExecSQL("INSERT INTO `chat_user`(`id_chat`, `id_user`) VALUES (?, ?)",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatId, userSubId},
		},
	)
	for i := 0; i < len(messagesId); i++ {
		messagesId[i] = conn.ExecSQL(
			"INSERT INTO `chat_message`(`id_chat`, `id_user`, `message`, `date`, `time`) VALUES (?, ?, ?, ?, ?)",
			[][]interface{}{
				{messageChatsId[i], messUsersId[i], messText[i], dt.date, dt.time},
			},
		)[0].LastId
	}
	parentId := []int64{0, 0, 0, messagesId[0], 0, messagesId[1], 0, 0, 0,
		messagesId[8], messagesId[9], 0, 0, 0, messagesId[12],
		0, 0, 0, messagesId[0], 0, messagesId[1], 0, 0, 0,
		messagesId[8], messagesId[9], 0, 0, 0, messagesId[12],
	}

	for i := 0; i < len(messagesId); i++ {
		conn.ExecSQL("UPDATE `chat_message` SET `parent_id`=? WHERE `id`=?",
			[][]interface{}{
				{parentId[i], messagesId[i]},
			},
		)
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
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_user` WHERE id_chat=? AND id_user=?",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatId, userSubId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_message` WHERE id=?",
		[][]interface{}{
			{messagesId[0]}, {messagesId[1]}, {messagesId[2]}, {messagesId[3]}, {messagesId[4]}, {messagesId[5]},
			{messagesId[6]}, {messagesId[7]}, {messagesId[8]}, {messagesId[9]}, {messagesId[10]}, {messagesId[11]},
			{messagesId[12]}, {messagesId[13]}, {messagesId[14]},
			{messagesId[15]}, {messagesId[16]}, {messagesId[17]}, {messagesId[18]}, {messagesId[19]}, {messagesId[20]},
			{messagesId[21]}, {messagesId[22]}, {messagesId[23]}, {messagesId[24]}, {messagesId[25]}, {messagesId[26]},
			{messagesId[27]}, {messagesId[28]}, {messagesId[29]},
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

	str := "{\"Status\":{\"Operation\":" + strconv.Itoa(constants.OP_LIST_PREVIOUS_MESSAGES) +
		",\"Status\":1,\"Message\":\"\"},\"Chat\":null,\"ChatUser\":null,\"Message\":["
	for i := len(messagesId) - 1 - 1; i > len(messagesId)-constants.LIST_LIMIT_MESSAGES-1-1; i-- {
		str += "{\"IModel\":{},\"Id\":" + strconv.Itoa(int(messagesId[i])) + ",\"ChatId\":" + strconv.Itoa(chatId) +
			",\"Author\":" + strconv.Itoa(messUsersId[i]) + ",\"ParrentMessage\":" + strconv.Itoa(int(parentId[i])) +
			",\"Message\":\"" + messText[i] + "\",\"FileUrl\":\"\",\"Date\":\"" + dt.date + "\",\"Time\":\"" + dt.time +
			"\"},"
	}
	str += "{\"IModel\":{},\"Id\":" + strconv.Itoa(int(messagesId[0])) + ",\"ChatId\":" + strconv.Itoa(chatId) +
		",\"Author\":" + strconv.Itoa(messUsersId[0]) + ",\"ParrentMessage\":" + strconv.Itoa(int(parentId[0])) +
		",\"Message\":\"" + messText[0] + "\",\"FileUrl\":\"\",\"Date\":\"" + dt.date + "\",\"Time\":\"" + dt.time +
		"\"},{\"IModel\":{},\"Id\":" + strconv.Itoa(int(messagesId[1])) + ",\"ChatId\":" + strconv.Itoa(chatId) +
		",\"Author\":" + strconv.Itoa(messUsersId[1]) + ",\"ParrentMessage\":" + strconv.Itoa(int(parentId[1])) +
		",\"Message\":\"" + messText[1] + "\",\"FileUrl\":\"\",\"Date\":\"" + dt.date + "\",\"Time\":\"" + dt.time +
		"\"}"
	str += "]," +
		"\"MessageStatus\":[" +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[4])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[4])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[5])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[5])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[6])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[6])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[7])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[7])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[8])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[8])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[9])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[10])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[10])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[10])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[11])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[11])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[11])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[12])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[12])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[13])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[13])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[14])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[14])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[14])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[15])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[15])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[16])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[16])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[17])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[18])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[18])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[19])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[19])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[19])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[20])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[21])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[21])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[22])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[22])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[23])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[23])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[24])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[0])) + ",\"UserId\":0,\"Status\":\"created\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[0])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[1])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[1])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[2])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[3])) + ",\"UserId\":0,\"Status\":\"delivered\",\"Date\":\"\",\"Time\":\"\"}," +
		"{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messagesId[3])) + ",\"UserId\":0,\"Status\":\"readed\",\"Date\":\"\",\"Time\":\"\"}" +
		"]," +
		"\"User\":[" +
		"{\"IModel\":{},\"Id\":" + strconv.Itoa(clientId) + ",\"Login\":\"\",\"Alias\":\"Leonard\",\"AuthKey\":\"\"," +
		"\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0}," +
		"{\"IModel\":{},\"Id\":" + strconv.Itoa(userSubId) + ",\"Login\":\"\",\"Alias\":\"Hovard\",\"AuthKey\":\"\"," +
		"\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0},{\"IModel\":{}," +
		"\"Id\":" + strconv.Itoa(userId) + ",\"Login\":\"\",\"Alias\":\"Sheldon\",\"AuthKey\":\"\"," +
		"\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0}]," +
		"\"BlackList\":null,\"UserPhone\":null}"

	str1 := "{\"Message\":["
	for i := 1; i < len(messagesId)-(len(messagesId)-constants.LIST_LIMIT_MESSAGES)+1; i++ {
		str1 += "{\"IModel\":{},\"Id\":" + strconv.Itoa(int(messagesId[i])) + ",\"ChatId\":" + strconv.Itoa(chatId) +
			",\"Author\":" + strconv.Itoa(messUsersId[i]) + ",\"ParrentMessage\":" + strconv.Itoa(int(parentId[i])) +
			",\"Message\":\"" + messText[i] + "\",\"FileUrl\":\"\",\"Date\":\"" + dt.date + "\",\"Time\":\"" + dt.time +
			"\"},"
	}
	str1 += "{\"IModel\":{},\"Id\":" + strconv.Itoa(int(messagesId[0])) + ",\"ChatId\":" + strconv.Itoa(chatId) +
		",\"Author\":" + strconv.Itoa(messUsersId[0]) + ",\"ParrentMessage\":" + strconv.Itoa(int(parentId[0])) +
		",\"Message\":\"" + messText[0] + "\",\"FileUrl\":\"\",\"Date\":\"" + dt.date + "\",\"Time\":\"" + dt.time +
		"\"}"
	str1 += "],\"Status\":{\"Operation\":" + strconv.Itoa(constants.OP_LIST_NEXT_MESSAGES) +
		",\"Status\":1,\"Message\":\"\"},\"User\":[{\"IModel\":{}," +
		"\"Id\":" + strconv.Itoa(userId) + ",\"Login\":\"\",\"Alias\":\"Sheldon\",\"AuthKey\":\"\"," +
		"\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0}," +
		"{\"IModel\":{},\"Id\":" + strconv.Itoa(userSubId) + ",\"Login\":\"\",\"Alias\":\"Hovard\",\"AuthKey\":\"\"," +
		"\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0}," +
		"{\"IModel\":{},\"Id\":" + strconv.Itoa(clientId) + ",\"Login\":\"\",\"Alias\":\"Leonard\",\"AuthKey\":\"\"," +
		"\"PassHash\":\"\",\"PassResetToken\":\"\",\"Email\":\"\",\"Status\":0,\"CreateAt\":0,\"UpdateAt\":0}]}"

	required := [][]string{
		{
			str,
			"[]",
		},
		/*	{
			str1,
			"[]",
		},*/
	}

	params := [][]int64{
		{
			constants.OP_LIST_PREVIOUS_MESSAGES, messagesId[len(messagesId)-1],
		},
		{
			constants.OP_LIST_NEXT_MESSAGES, messagesId[0],
		},
	}

	for key, testt := range params {

		reqParam := map[string]string{}

		m := message.NewMessage()
		m.Id = testt[1]
		str1, _ := json.Marshal(m)
		reqParam[data.MODEL_MESSAGE] = string(str1)

		stat := StatusRequest{
			Operation: int(testt[0]),
		}
		str1, _ = json.Marshal(stat)
		reqParam[data.MODEL_STATUS] = string(str1)

		bh := NewBusinessHandler(&reqParam, &ex1{id: clientId}, dt)

		resp, _, err := bh.ProcessinRequest()
		if err != nil {
			t.Log("ERROR bh.ProcessinRequest() => " + err.Error())
		}

		b := map[string]any{}
		err = json.Unmarshal(resp, &b)
		if err != nil {
			t.Log("ERRROR JSON MARSHALING RESPONSE =>", err)
		}

		if required[key][0] != string(resp) {
			t.Fatal("Failed Generate List Messages! \nRequired:", required[key][0], "; \nrezult:  ", string(resp))
		}

		enc, _ := json.Marshal(bh.GetSendChatsId())
		if required[key][1] != string(enc) {
			t.Fatal("Failed Generate List Messages! \nRequired:", required[key][1], "; \nrezult:  ", string(enc))
		}
	}
	t.Log("Lead time test ", float32(time.Now().UnixMilli()-start)/float32(1000), " seconds")
}

func TestStatusMessageSendResponse(t *testing.T) {

	t.Log("Generate Status Message")

	//	prepare test data
	conn := db.NewDB()

	mess := "Hello Barbarian!"
	dt := NewDateTimeStamp()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Barbarian\", \"Barbarian\", \"$2y$10$cKabhu.fhhdda/ElFyKg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hD0xVvbmxI5w2ldZ8I8XCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Madyar\", \"Madyar\", \"$2y$10$s8Ql2465Cws2wFGjsFjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
	userSubId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Matroskin\", \"Matroskin\", \"$2y$1e$kfglr5a5ewe2wFej4Fjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)

	chatId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(clientId)+", \"MAD Barbarian\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

	conn.ExecSQL("INSERT INTO `chat_user`(`id_chat`, `id_user`) VALUES (?, ?)",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
			{chatId, userSubId},
		},
	)
	messageSupId := int64(conn.ExecSQL("INSERT INTO `chat_message`(`id_chat`, `id_user`, `message`, `date`, `time`) VALUES (?, ?, ?, ?, ?)",
		[][]interface{}{{chatId, clientId, mess, dt.date, dt.time}},
	)[0].LastId)
	messageId := int64(conn.ExecSQL("INSERT INTO `chat_message`(`id_chat`, `id_user`, `message`, `date`, `time`) VALUES (?, ?, ?, ?, ?)",
		[][]interface{}{{chatId, userId, mess, dt.date, dt.time}},
	)[0].LastId)

	conn.ExecSQL("INSERT INTO `chat_message_status`(`id_message`, `id_user`, `status_message`, `date`, `time`) VALUES (?,?,?,?,?)",
		[][]interface{}{
			{messageSupId, clientId, constants.MESSAGE_CREATED, dt.date, dt.time},
			{messageSupId, userId, constants.MESSAGE_DELIVERED, dt.date, dt.time},
			{messageSupId, userSubId, constants.MESSAGE_READED, dt.date, dt.time},
			{messageId, clientId, constants.MESSAGE_CREATED, dt.date, dt.time},
			{messageId, userId, constants.MESSAGE_DELIVERED, dt.date, dt.time},
			{messageId, userSubId, constants.MESSAGE_READED, dt.date, dt.time},
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
	defer conn.ExecSQL("DELETE FROM `chat_message` WHERE id=?",
		[][]interface{}{
			{messageId},
			{messageSupId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_message_status` WHERE id_message=? AND id_user=?",
		[][]interface{}{
			{messageSupId, clientId},
			{messageSupId, userId},
			{messageSupId, userSubId},
			{messageId, clientId},
			{messageId, userId},
			{messageId, userSubId},
		},
	)

	// Test start

	arrStatuses := []string{constants.MESSAGE_DELIVERED, constants.MESSAGE_READED}
	required := [][]string{
		{
			"{\"Status\":{\"Operation\":101,\"Status\":1,\"Message\":\"\"},\"Chat\":null,\"ChatUser\":null," +
				"\"Message\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(int(messageId)) + ",\"ChatId\":" + strconv.Itoa(int(chatId)) +
				",\"Author\":0,\"ParrentMessage\":0," +
				"\"Message\":\"\",\"FileUrl\":\"\",\"Date\":\"\",\"Time\":\"\"}],\"MessageStatus\":" +
				"[{\"IModel\":{},\"MessageId\":" + strconv.Itoa(int(messageId)) + ",\"UserId\":0,\"Status\":\"" + arrStatuses[0] +
				"\",\"Date\":\"" + dt.date + "\"," +
				"\"Time\":\"" + dt.time + "\"}],\"User\":null,\"BlackList\":null,\"UserPhone\":null}",
			"[" + strconv.Itoa(chatId) + "]",
		},
		{
			"{\"Status\":{\"Operation\":101,\"Status\":1,\"Message\":\"\"},\"Chat\":null,\"ChatUser\":null,\"Message\":null," +
				"\"MessageStatus\":null,\"User\":null,\"BlackList\":null,\"UserPhone\":null}",
			"[]",
		},
	}

	for i := 0; i < len(arrStatuses); i++ {

		reqParam := map[string]string{}

		m := message.NewMessageStatus()
		m.UserId = clientId
		m.MessageId = messageId
		m.Status = arrStatuses[i]
		m.Date = dt.date
		m.Time = dt.time

		str, _ := json.Marshal(m)
		reqParam[data.MODEL_STATUS_MESSAGE] = string(str)

		st := StatusRequest{
			Operation: constants.OP_STATUS_MESSAGE,
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

		if required[i][0] != string(resp) {
			t.Fatal("Failed Generate Status Message! \nRequired:", required[i][0], "; \nrezult:  ", string(resp))
		}

		enc, _ := json.Marshal(bh.GetSendChatsId())
		if required[i][1] != string(enc) {
			t.Fatal("Failed Generate Status Message! \nRequired:", required[i][1], "; \nrezult:  ", string(enc))
		}

		/*	conn.ExecSQL("UPDATE `chat_message_status` SET `status_message`=? WHERE `id_message`=? AND `id_user`=?", [][]interface{}{
			{constants.MESSAGE_CREATED, messageSupId, clientId},
			{constants.MESSAGE_CREATED, messageId, clientId},
		})*/
	}
}

func TestNewMessage(t *testing.T) {

	t.Log("Generate New Message")
	//	prepare test data
	conn := db.NewDB()

	mess := "Hello Maykl!"
	messSup := "Hello Graf"

	dt := NewDateTimeStamp()

	clientId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, `created_at`, `updated_at`) "+
		"VALUES (\"Graf\", \"Graf\", \"$2y$10$YKj6hu.fJ0jjj/ElFyKg0uk5R\", "+
		"\"$2y$10$c.a2SAgBl.Ey2BpcG96dQO1TuaB3hj0xVzeHxI5w2ldZ8IjXCjyM2\", 1540361413, 1563035897)", [][]interface{}{})[0].LastId)
	userId := int(conn.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, `password_hash`, "+
		"`created_at`, `updated_at`) VALUES (\"Maykl\", \"Maykl\", \"$2y$10$68QljjjjCwe2wFGj4Fjcx.Xzh\","+
		" \"$2y$10$gwYeSu1aoCR2NTBPYcU.QuP0h6Lu/udtWaQGffAmBIif57EEQcBPK\", 1540361414, 1540361414)", [][]interface{}{})[0].LastId)
	//fmt.Println("clientId =>", clientId, "userId =>", userId)

	chatId := int(conn.ExecSQL("INSERT INTO `chat`(`author`, `alias`, `create_at`, `status`) "+
		"VALUES ("+strconv.Itoa(clientId)+", \"MAD Graf\", \"2016-08-18\", \"active\")", [][]interface{}{})[0].LastId)

	conn.ExecSQL("INSERT INTO `chat_user`(`id_chat`, `id_user`) VALUES (?, ?)",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
		},
	)
	rez := *(conn.SelectSQL("SELECT `id` FROM `chat_message` ORDER BY id DESC LIMIT 1", []interface{}{}))

	messageSupId := 1
	if len(rez) > 0 {
		messageSupId = int(rez[0]["id"].(uint64)) + 1
	}

	conn.ExecSQL("ALTER TABLE chat_message AUTO_INCREMENT="+strconv.Itoa(messageSupId), [][]interface{}{})

	conn.ExecSQL("INSERT INTO `chat_message`(`id_chat`, `id_user`, `message`, `date`, `time`) VALUES (?, ?, ?, ?, ?)",
		[][]interface{}{
			{chatId, userId, messSup, dt.date, dt.time},
		},
	)
	messageId := messageSupId + 1

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
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_user` WHERE id_chat=? AND id_user=?",
		[][]interface{}{
			{chatId, clientId},
			{chatId, userId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_message` WHERE id=?",
		[][]interface{}{
			{messageId},
			{messageSupId},
		},
	)
	defer conn.ExecSQL("DELETE FROM `chat_message_status` WHERE id_message=? AND id_user=?",
		[][]interface{}{
			{messageId, clientId},
			{messageId, userId},
		},
	)

	// Test start
	required := []string{
		"{\"Status\":{\"Operation\":103,\"Status\":1,\"Message\":\"\"},\"Chat\":null,\"ChatUser\":null," +
			"\"Message\":[{\"IModel\":{},\"Id\":" + strconv.Itoa(messageId) + ",\"ChatId\":" +
			strconv.Itoa(chatId) + ",\"Author\":" + strconv.Itoa(clientId) + ",\"ParrentMessage\":" +
			strconv.Itoa(messageSupId) + ",\"Message\":\"" + mess + "\",\"FileUrl\":\"\",\"Date\":\"" + dt.date +
			"\",\"Time\":\"" + dt.time + "\"}],\"MessageStatus\":null,\"User\":null,\"BlackList\":null,\"UserPhone\":null}",
		"[" + strconv.Itoa(chatId) + "]",
	}

	reqParam := map[string]string{}

	m := message.NewMessage()
	m.Author = clientId
	m.ChatId = chatId
	m.Message = mess
	m.Date = dt.date
	m.Time = dt.time
	m.ParrentMessage = int64(messageSupId)

	str, _ := json.Marshal(m)
	reqParam[data.MODEL_MESSAGE] = string(str)

	st := StatusRequest{
		Operation: constants.OP_NEW_MESSAGE,
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
		t.Fatal("Failed Generate New Message! \nRequired:", required[0], "; \nrezult:  ", string(resp))
	}

	enc, _ := json.Marshal(bh.GetSendChatsId())
	if required[1] != string(enc) {
		t.Fatal("Failed Generate New Message! \nRequired:", required[1], "; \nrezult:  ", string(enc))
	}
}
