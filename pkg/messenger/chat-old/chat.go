package chat_old

/****
import (
//	"log"
//	"strconv"
//	"strings"
//	"messenger/units/wss-server"
)

type Chat struct {
	id, autor                int
	alias, status, create_at string
}


/****
func createNewChat(client *wss_server.Client, request *map[string]interface{}) {

	r := *request
	db := client.GetDB()
	defer db.CloseDB()

	query := "INSERT INTO chats (author, alias, create_at, status) VALUES (?, ?, CURDATE(), ?)"
	res, err := db.Exec(query, strconv.Itoa(client.hub.clients[client].idUser), r["chat_name"], CHAT_ACTIVE)
	if err != nil {
		log.Println("ERROR: err => ", err)
		return
	}
	var id int64
	id, err = res.LastInsertId()
	if err != nil {
		log.Println("ERROR 301: id => ", err)
		return
	}

	users := r["users"].([]interface{})
	users = append(users, strconv.Itoa(client.hub.clients[client].idUser))
	//	log.Println("users =>", users)
	query = "INSERT INTO chat_user (id_chat, id_user) VALUES (?, ?)"
	for _, v := range users {
		res, err := db.Exec(query, id, v)
		if err != nil {
			log.Println("ERROR: err => ", err)
			return
		}
		if _, err = res.LastInsertId(); err != nil {
			log.Println("ERROR 301: id => ", err)
			return
		}
	}

	idUsers, err := getIdUsersOfChat(int(id))
	if err != nil {
		log.Println("ERROR 546 idUsers, err => ", err)
	}

	response := map[string]map[string]map[string]string{
		"chats": {strconv.Itoa(int(id)): {"id": strconv.Itoa(int(id)), "alias": r["chat_name"].(string),
			"author": strconv.Itoa(client.hub.clients[client].idUser)}},
		"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_CREATE_NEW_CHAT)}},
	}
	client.hub.sendClients(idUsers, response, nil)
}

func getChats(client *Client) {
	db := getDB()
	defer db.Close()

	//	ищем все доступные чаты для клиента
	query := "SELECT id, autor, alias, status FROM chat c, chat_user u  WHERE c.id=u.id_chat AND c.status <> \"" +
		CHAT_DELETED + "\" AND u.id_user=" + strconv.Itoa(client.hub.clients[client].idUser)
	log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		log.Println("QUERY ABORT _333")
	}
	//	создаем карту, куда будем собирать чаты
	chats := make(map[string]map[string]string, 0)
	//	а также создаем массив строк с номерами чатов
	idChats := make([]string, 0)
	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.id, &chat.autor, &chat.alias, &chat.status)
		if err != nil {
			log.Fatal(err)
		}
		idChat := strconv.Itoa(chat.id)
		idChats = append(idChats, idChat)
		chats[idChat] = map[string]string{"id": idChat, "autor": strconv.Itoa(chat.autor),
			"alias": chat.alias, "status": chat.status}
	}
	rows.Close()
	//	преобразуем в строку
	var numberChats string = strings.Join(idChats, ", ")

	//	ищем последнее сообщение для каждого из чатов
	query = "SELECT cm.id, cm.id_chat, cm.id_user, u.username, cm.parent_id, cm.message, cm.file, cm.date, cm.time, cms.status_message " +
		" FROM user u, chat_message cm, chat_message_status cms WHERE cm.id_user=u.id AND cms.id_message=cm.id AND cms.id_user=" +
		strconv.Itoa(client.hub.clients[client].idUser) + " AND (cm.id_chat, cm.date, cm.time) IN " +
		" (SELECT id_chat, date, MAX(time) FROM chat_message c WHERE (id_chat, date) IN (SELECT id_chat, MAX(date) " +
		" FROM chat_message WHERE id_chat IN (" + numberChats + ") GROUP BY id_chat)	GROUP BY id_chat)"
	rows, err = db.Query(query)
	if err != nil {
		log.Println("QUERY ABORT => ", query)
		return
	}
	//	создаем и наполняем карту сообщений
	messages := make(map[string]map[string]string, 0)
	for rows.Next() {
		var m Message
		var username, status string
		err = rows.Scan(&m.id, &m.idChat, &m.idUser, &username, &m.parentId, &m.message, &m.file, &m.date, &m.time, &status)
		if err != nil {
			log.Fatal(err)
		}
		messages[strconv.Itoa(m.idChat)] = map[string]string{"id": strconv.Itoa(m.id), "id_chat": strconv.Itoa(m.idChat),
			"id_user": strconv.Itoa(m.idUser), "username": username, "parent_id": strconv.Itoa(int(m.parentId.Int64)),
			"message": m.message, "file": m.file.String, "date": m.date, "time": m.time, "status": status,
		}
	}
	rows.Close()

	//создаем карту ответа клиенту
	status := map[string]map[string]string{"status": {"operation": strconv.Itoa(OP_GET_CHATS), "status": strconv.Itoa(STATUS_ACCEPT)}}
	response := map[string]map[string]map[string]string{"chats": chats, "messages": messages, "status": status}
	client.hub.sendClients(nil, response, client)
}

func exitChat(client *Client, request *map[string]interface{}) {
	/*	r := *request
		db := getDB()
		defer db.Close()

		idChat, err := strconv.Atoi(r["id"].(string))
		if err != nil {
			log.Println(err)
			return
		}

		idUsers, err := getIdUsersOfChat(idChat)
		if err != nil {
			log.Println("ERROR 409: err =>", err)
			return
		}
		var newAuthor = ""
		log.Println("idUsers =>", idUsers)
		log.Println("id_user =>", r["id_user"])
		flag := securityAuthor(idChat, client.hub.clients[client].idUser)

		if flag {
			if idU, ok := r["id_user"]; ok {
				for _, v := range idUsers {
					if strings.Compare(idU.(string), strconv.Itoa(v)) == 0 {
						log.Println("СОВПАДЕНИЕ: flag =>", flag, "; newAuthor =>", newAuthor)
						break
					}
				}
			}
			if strings.Compare(newAuthor, "") == 0 {
				log.Println("newAuthor == \"\" =>", newAuthor)
				for _, v := range idUsers {
					if v != client.hub.clients[client].idUser {
						newAuthor = strconv.Itoa(v)
						break
					}
				}
			}

			query := "UPDATE chat SET autor=? WHERE id=?"
			log.Println("newAuthor =>", newAuthor, "; r[\"id\"] =>", r["id"], "\n query =>", query)
			_, err = db.Exec(query, newAuthor, r["id"])
			if err != nil {
				log.Println("ERROR 590: sql => ", err)
				return
			}
		}

		query := "DELETE FROM chat_user WHERE id_chat=? AND id_user=?"
		_, err = db.Exec(query, r["id"], client.hub.clients[client].idUser)
		if err != nil {
			log.Println("ERROR: sql => ", err)
			return
		}

		response := map[string]map[string]map[string]string{
			"chats": {strconv.Itoa(idChat): {"id": strconv.Itoa(idChat), "id_user": strconv.Itoa(client.hub.clients[client].idUser)}},
			"users": {strconv.Itoa(client.hub.clients[client].idUser): {"id": strconv.Itoa(client.hub.clients[client].idUser),
				"username": client.hub.clients[client].userName}},
			"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_EXIT_CHAT),
				"message": "Успех"}},
		}
		if flag {
			response["chats"][strconv.Itoa(idChat)]["author"] = newAuthor
		}

		client.hub.sendClients(idUsers, response, nil)*/
/*****}

func removeChat(client *Client, request *map[string]interface{}) {
	/*	r := *request
		db := getDB()
		defer db.Close()

		idChat, err := strconv.Atoi(r["id"].(string))
		if err != nil {
			log.Println(err)
			return
		}

		if !securityAuthor(idChat, client.hub.clients[client].idUser) {
			response := map[string]map[string]map[string]string{
				"status": {"status": {"status": strconv.Itoa(STATUS_ERROR), "operation": strconv.Itoa(OP_ADD_USER),
					"message": "Вы не являетесь автором чата"}},
			}
			client.hub.sendClients(nil, response, client)
			return
		}

		users, err := getIdUsersOfChat(idChat)

		query := "UPDATE chat SET status=? WHERE id=? AND autor=?"
		_, err = db.Exec(query, CHAT_DELETED, idChat, client.hub.clients[client].idUser)
		if err != nil {
			log.Println("ERROR 590: sql => ", err)
			return
		}

		response := map[string]map[string]map[string]string{
			"chats": {strconv.Itoa(idChat): {"id": strconv.Itoa(idChat), "author": strconv.Itoa(client.hub.clients[client].idUser)}},
			"users": {strconv.Itoa(client.hub.clients[client].idUser): {"id": strconv.Itoa(client.hub.clients[client].idUser),
				"username": client.hub.clients[client].userName}},
			"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_REMOVE_CHAT)}},
		}
		client.hub.sendClients(users, response, nil)*/
/***}***/
