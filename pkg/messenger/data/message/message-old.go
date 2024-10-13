package message

/*****
import (
	"database/sql"
	//	"encoding/json"
	///	"fmt"
	//	"log"
	////	"strconv"
	//	"strings"
	//	"time"
	_ "github.com/alex988334/messenger/pkg/messenger/wss-server"
)

type Message struct {
	id, idChat, idUser  int
	parentId            sql.NullInt64
	message, date, time string
	file                sql.NullString
}

type StatusMessage struct {
	idMessage, idUser, status int
	date, time                string
}

func systemMessage(client *wss_server.Client, request *map[string]interface{}) {

}

/*****
func setStatus(client *Client, request *map[string]interface{}) {
	r := *request
	db := getDB()
	defer db.Close()

	if strings.Compare(fmt.Sprint(r["id"]), "0") == 0 {
		query := "SELECT id_message FROM chat_message_status ms, chat_message m, chat_user u " +
			" WHERE m.id=ms.id_message AND ms.id_user=u.id_user AND status_message <> ? AND u.id_user=? " +
			" AND m.id_chat=? AND ms.id_user <> m.id_user GROUP BY id_message"
		rows, err := db.Query(query, MESSAGE_READED, client.hub.clients[client].idUser, r["id_chat"])
		if err != nil {
			log.Println("ERROR: err => ", err)
			return
		}

		mass := make([]interface{}, 0)
		for rows.Next() {
			err = rows.Scan(mass...)
			if err != nil {
				log.Fatal(err)
			}
		}
		rows.Close()

		query2 := "UPDATE chat_message_status ms, chat_message m, chat_user u SET status_message=? WHERE  m.id=ms.id_message " +
			" AND ms.id_user=u.id_user AND status_message <> ? AND u.id_user=? AND m.id_chat=? AND ms.id_user <> m.id_user"
		res, err := db.Exec(query2, MESSAGE_READED, MESSAGE_READED, client.hub.clients[client].idUser, r["id_chat"])
		if err != nil {
			log.Println("UPDATE chat_message_status ERROR!!!!")
			return
		}

		log.Println("setStatus() 0")
		if i, err := res.RowsAffected(); i > 0 && err == nil {
			securityStatusMessage(mass, client)
		}
	} else {
		// проверить что статус пришел не от автора сообщения
		query2 := "UPDATE chat_message_status ms, chat_message m SET status_message=? WHERE ms.id_message=? " +
			" AND ms.id_user=? AND ms.id_user <> m.id_user"
		res, err := db.Exec(query2, r["status"], r["id"], client.hub.clients[client].idUser)
		if err != nil {
			log.Println("UPDATE chat_message_status ERROR!!!!")
			return
		}

		log.Println("setStatus() id")

		if i, err := res.RowsAffected(); i > 0 && err == nil {
			securityStatusMessage([]interface{}{r["id"]}, client)
		}
	}
}

func securityStatusMessage(idMessages []interface{}, c *Client) {
	db := getDB()
	defer db.Close()
	log.Println("securityStatusMessage()")
	for _, id := range idMessages {
		query := "SELECT send, delivered, readed, total, id_user, id_chat FROM (SELECT COUNT(status_message) as send " +
			" FROM chat_message_status WHERE id_message=? " +
			" AND status_message=\"send\") s, (SELECT COUNT(status_message) as delivered FROM chat_message_status " +
			" WHERE id_message=? AND status_message=\"delivered\") d, (SELECT COUNT(status_message) as readed " +
			" FROM chat_message_status WHERE id_message=? AND status_message=\"readed\") r, " +
			" (SELECT COUNT(status_message) as total FROM chat_message_status WHERE id_message=?) t, chat_message  WHERE id=?"

		rows, err := db.Query(query, id, id, id, id, id)
		if err != nil {
			log.Println("ERROR: err => ", err)
			return
		}
		var total, send, delivered, readed, idAuthor, idChat int
		for rows.Next() {
			err = rows.Scan(&send, &delivered, &readed, &total, &idAuthor, &idChat)
			if err != nil {
				log.Fatal(err)
			}
		}
		rows.Close()

		log.Println("total => ", total, ", send => ", send, ", delivered => ", delivered, ", readed => ", readed,
			", idAuthor => ", idAuthor, ", idChat => ", idChat)
		status := MESSAGE_SEND

		if (delivered + 1) == total {
			status = MESSAGE_DELIVERED
		}
		if (readed + 1) == total {
			status = MESSAGE_READED
		}

		log.Println("status => ", status)
		if strings.Compare(status, MESSAGE_SEND) != 0 {
			query = "UPDATE chat_message_status cms, chat_message cm SET status_message=? WHERE cms.id_message=cm.id AND cms.id_user=cm.id_user AND id_message=?"
			res, err := db.Exec(query, status, id)
			log.Println(query, "; status => ", status, "; id => ", id)
			if err != nil {
				log.Println("UPDATE chat_message_status ERROR!!!!")
				return
			}

			if i, err := res.RowsAffected(); i > 0 && err == nil {
				response := map[string]map[string]map[string]string{
					"messages": {fmt.Sprint(id): {"id": fmt.Sprint(id), "id_chat": strconv.Itoa(idChat), "status": status}},
					"status":   {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_STATUS_MESSAGE)}},
				}
				log.Println(response)
				str, err := json.Marshal(response)
				if err != nil {
					log.Println(err)
					return
				}

				for client, elem := range c.hub.clients {
					log.Println("client => ", client, "; elem => ", elem, "; idAuthor => ", idAuthor)
					if elem.idUser == idAuthor {
						log.Println("TRUE, idAuthor => ", idAuthor)
						client.send <- str
					}
				}
			}
		}
	}
}

func createMessage(m *Message) bool {
	db := getDB()
	defer db.Close()

	//idChat := strconv.Atoi(params["id_chat"].(string))
	//if err != nil { log.Println("ERROR 472: id_chat => err => ", err); return 0, err; }

	idUsers, err := getIdUsersOfChat(m.idChat)
	if err != nil {
		log.Println("ERROR 851: func => getIdUsersOfChat(id), id => ", m.idChat, ", err => ", err)
		return false
	}

	query := "INSERT INTO chat_message(id_chat, id_user, parent_id, message, file, date, time) VALUES (?, ?, ?, ?, ?, ?, CURTIME())"
	res, err := db.Exec(query, m.idChat, m.idUser, m.parentId.Int64, m.message, m.file.String, m.date)
	if err != nil {
		log.Println("ERROR 624: sql => ", err)
		return false
	}

	if id, err := res.LastInsertId(); err != nil {
		log.Println("ERROR 238: id => ", err)
		return false
	} else {
		m.id = int(id)
	}

	for _, val := range idUsers {
		log.Println("createMessage() val =>", val)
		query = "INSERT INTO chat_message_status(id_message, id_user, status_message, date, time) VALUES (?, ?, ?, CURDATE(), CURTIME())"
		res, err = db.Exec(query, m.id, val, MESSAGE_SEND)
		if err != nil {
			log.Println("ERROR 693: sql => ", err)
			return false
		}
		if _, err := res.LastInsertId(); err != nil {
			log.Println("ERROR 235: id => ", err)
			return false
		}
	}

	return true
	//log.Println("RETURN id =>", id)
}

func newMessage(client *Client, request *map[string]interface{}) {
	r := *request
	db := getDB()
	defer db.Close()

	idChat, err := strconv.Atoi(r["id_chat"].(string))
	if err != nil {
		log.Println("ERROR 3761: err =>", err)
		return
	}

	var pId sql.NullInt64
	if id, ok := r["parent_id"]; ok && id != nil {
		pId = sql.NullInt64{
			Int64: int64(r["parent_id"].(float64)), Valid: true,
		}
	} else {
		pId = sql.NullInt64{
			Int64: 0, Valid: false,
		}
	}
	var f sql.NullString
	if id, ok := r["file"]; ok && id != nil {
		f = sql.NullString{
			String: r["file"].(string), Valid: true,
		}
	} else {
		f = sql.NullString{
			String: "", Valid: false,
		}
	}

	m := Message{
		idChat:   idChat,
		idUser:   client.hub.clients[client].idUser,
		parentId: pId,
		message:  r["message"].(string),
		file:     f,
		date:     time.Now().Format("2006-01-02"),
		time:     time.Now().Format("15:04:05"),
	}

	if !createMessage(&m) {
		log.Println("ERROR 9525: error create new message")
		return
	}

	/*	query := "INSERT INTO chat_message(id_chat, id_user, parent_id, message, file, date, time) VALUES (?, ?, ?, ?, ?, CURDATE(), CURTIME())"
		res, err := db.Exec(query, r["id_chat"], client.hub.clients[client].idUser, r["parent_id"], r["message"], r["file"])
		if err != nil { log.Println("ERROR: sql => ", err); return; }
		if id, err = res.LastInsertId(); err != nil {
			log.Println("ERROR 301: id => ", err)
			return
		}

		messages := make(map[string]map[string]string, 0)

		/*query := "SELECT id, id_chat, id_user, parent_id, message, file, date, time FROM chat_message WHERE id=" + strconv.Itoa(m.id)
		rows, err := db.Query(query)
		for rows.Next() {
			err = rows.Scan(&m.id, &m.idChat, &m.idUser, &m.parentId, &m.message, &m.file, &m.date, &m.time)
			if err != nil { log.Fatal(err);	return;	}
			messages[strconv.Itoa(m.idChat)] = map[string]string{"id": strconv.Itoa(m.id), "id_chat": strconv.Itoa(m.idChat),
				"id_user": strconv.Itoa(m.idUser), "parent_id": strconv.Itoa(int(m.parentId.Int64)), "message": m.message,
				"file": m.file.String, "date": m.date, "time": m.time, "autor": client.hub.clients[client].userName, "status_m": MESSAGE_SEND,
			}
		}
		rows.Close()*/

/*****	idUsers, err := getIdUsersOfChat(m.idChat)
if err != nil {
	log.Println("ERROR 349: func => getIdUsersOfChat(id), id => ", m.idChat, ", err => ", err)
	return
}
/*
	for _, val := range idUsers {
		query = "INSERT INTO chat_message_status(id_message, id_user, status_message, date, time) VALUES (?, ?, ?, ?, ?)"
		res, err = db.Exec(query, m.id, val, MESSAGE_SEND, m.date, m.time)
		if err != nil {
			log.Println("ERROR: sql => ", err);
			return;
		}
		if id, err = res.LastInsertId(); err != nil {
			log.Println("ERROR 301: id => ", err)
			return
		}
	}*/

/***	response := map[string]map[string]map[string]string{
		"messages": {strconv.Itoa(m.idChat): {
			"id": strconv.Itoa(m.id), "id_chat": strconv.Itoa(m.idChat),
			"id_user": strconv.Itoa(m.idUser), "parent_id": strconv.Itoa(int(m.parentId.Int64)), "message": m.message,
			"file": m.file.String, "date": m.date, "time": m.time, "autor": client.hub.clients[client].userName, "status_m": MESSAGE_SEND}},
		"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_INPUT_MESSAGE)}},
	}
	client.hub.sendClients(idUsers, response, nil)
}

func getHistoryMessage(client *Client, r *map[string]interface{}) {
	request := *r
	db := getDB()
	defer db.Close()
	var (
		d          interface{}
		newDate    sql.NullString
		idMessages = make([]string, 0)
	)
	//	ищем ID всех сообщений из чата id_chat со статусом не равным "readed"
	query2 := "SELECT ms.id_message FROM chat_message_status ms LEFT JOIN chat_message m ON m.id=ms.id_message WHERE ms.id_user=" +
		strconv.Itoa(client.hub.clients[client].idUser) +
		" AND m.id_chat=? AND ms.id_user <> m.id_user AND status_message <> \"" + MESSAGE_READED + "\""

	rows, err := db.Query(query2, request["id"])
	if err != nil {
		log.Println("ERROR 65445: SQL => ", query2)
		return
	}
	for rows.Next() {
		var idM sql.NullString
		err = rows.Scan(&idM)
		if err != nil {
			log.Fatal("ERROR 56215:", err)
			return
		}
		idMessages = append(idMessages, idM.String)
	}
	rows.Close()

	//	если сообщения найдены обновляем их статусы до "readed"
	//if len(idMessages) > 0 {
	//str := strings.Join(idMessages, ", ")																		//	создаем строку из массива
	query2 = "UPDATE chat_message_status ms, chat_message cm SET status_message=\"readed\" " +
		" WHERE ms.id_message=cm.id AND ms.id_user=? AND cm.id_chat=?"

	/*query2 := "UPDATE chat_message_status SET status_message=\"" + MESSAGE_READED + "\" WHERE id_user=" +
	strconv.Itoa(client.hub.clients[client].idUser)+ " AND id_message IN (" + str + ")"*/
/****	_, err = db.Exec(query2, client.hub.clients[client].idUser, request["id"])
	if err != nil {
		log.Println("UPDATE chat_message_status ERROR!!!!")
		return
	}
	//}

	//	проверяем что запрашиваются сообщения на предыдущую дату от указанной, либо на последнюю дату
	if date, ok := request["date"]; ok {
		//	ищем крайнюю дату
		query := "SELECT date FROM chat_message c WHERE id_chat=? AND date < ? ORDER BY date DESC LIMIT 1"
		rows, err := db.Query(query, request["id"], date)
		if err != nil {
			log.Println("ERROR: ", err, " SQL => ", query)
			return
		}
		for rows.Next() {
			err = rows.Scan(&newDate)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		rows.Close()

		//	проверяем что дата существует (есть еще сообщения в истории)
		if newDate.String != "" {
			d = newDate.String
			query2 = "SELECT cm.id AS id, id_chat, cm.id_user, parent_id, message, file, cm.date, cm.time, " +
				" u.username AS autor, status_message FROM chat_message cm, user u, chat_message_status cms " +
				" WHERE cm.id_user=u.id AND cms.id_message=cm.id AND cms.id_user=cm.id_user AND id_chat=? AND cm.date=? " +
				" ORDER BY cm.time DESC"
		} else {
			//	если сообщений в истории больше нет, оповещаем клиента
			response := map[string]map[string]map[string]string{"messages": {}, "status": {"status": {
				"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_GET_HISTORY_MESSAGE),
				"loadingData": fmt.Sprint(request["loadingData"])}},
			}

			client.hub.sendClients(nil, response, client)
			log.Println("БОЛЬШЕ нет сообщений")
			return
		}
	} else {
		//	работает в случае если запрос без указания конкретной даты
		query2 = "SELECT cm.id AS id, id_chat, cm.id_user, parent_id, message, file, cm.date, cm.time, u.username AS autor, status_message FROM chat_message cm, user u, chat_message_status cms WHERE cm.id_user=u.id AND cms.id_message=cm.id AND cms.id_user=cm.id_user AND id_chat=? AND cm.date IN (SELECT MAX(date) FROM chat_message WHERE id_chat=?) ORDER BY cm.time ASC"
		d = fmt.Sprint(request["id"])
	}

	//	ищем и получаем список сообщений
	rows, err = db.Query(query2, request["id"], d)
	if err != nil {
		log.Println("ERROR: ", err, " ; SQL => ", query2)
		return
	}
	log.Println("ERROR:  SQL => ", query2, "  PARAMS: id => ", request["id"], " , d => ", d)
	messages := make(map[string]map[string]string, 0)
	for rows.Next() {
		var (
			mes               Message
			username, statusM string
		)
		err = rows.Scan(&mes.id, &mes.idChat, &mes.idUser, &mes.parentId, &mes.message, &mes.file, &mes.date, &mes.time, &username, &statusM)
		if err != nil {
			log.Fatal(err)
			return
		}
		messages[strconv.Itoa(mes.id)] = map[string]string{"id": strconv.Itoa(mes.id), "id_chat": strconv.Itoa(mes.idChat),
			"id_user": strconv.Itoa(mes.idUser), "parent_id": strconv.Itoa(int(mes.parentId.Int64)), "message": mes.message,
			"file": mes.file.String, "date": mes.date, "time": mes.time, "autor": username, "status_m": statusM,
		}
	}
	rows.Close()

	//	генерируем стандартную карту ответа клиенту
	response := map[string]map[string]map[string]string{"messages": messages}
	response["status"] = map[string]map[string]string{"status": {"status": strconv.Itoa(STATUS_ACCEPT),
		"operation": strconv.Itoa(OP_GET_HISTORY_MESSAGE), "loadingData": fmt.Sprint(request["loadingData"])}}
	client.hub.sendClients(nil, response, client)

	//	запускаем проверку статусов сообщений для их авторов, т.к. вполне возможно что пользователь изменивший
	//	статусы был крайним и следущий на изменение статусов идет автор
	mass := make([]interface{}, 0)
	for _, v := range idMessages {
		mass = append(mass, v)
	}
	securityStatusMessage(mass, client)
/*}
 ****/
