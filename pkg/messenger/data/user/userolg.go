package user

/*****
import (
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	//	"math/rand"
	//	_ "math/rand"
	"messenger/service"
	//"time"
	//	"unsafe"
	//wss_server "messenger/units/wss-server"
	//_ "messenger/units/wss-server"
	//	wss_server "messenger/units/wss-server"
	//"messenger/units/wss-server"
)

/*
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}*/

/*type ChatUser struct {
	idChat, idUser int
}*/

/*type BlackList struct {
	blocking, locked int
	date, time       string
}*/
/*
type User struct {
	id, status, createdAt int
	login                 string
	email                 sql.NullString
}
*/
/*****
func GetUserChatsId(client *UserInterface) []interface{} {
	data := (*client).GetDB().SelectSQL("SELECT id_chat FROM chat_users WHERE id_user=?",
		[]string{"id_chat"}, []interface{}{(*client).GetId()})

	return service.ArrayFromField("id_chat", data)
}

func SecurityUser(client *UserInterface) bool {
	if _, ok := (*client).GetQuery()["ik"]; ok {
		if !securityIdentityKey(client) {
			return false
		}
		if !createIdentityKey(client) {
			return false
		}
		return true
	}

	if securityForPassword(client) {
		if !createUserKey(client) {
			return false
		}
		if !createIdentityKey(client) {
			return false
		}
		return true
	}
	return false
}

func UpdateUserKey(client *UserInterface) bool {
	return true
}

func UpdateIdentityKey(client *UserInterface) bool {
	return true
}

func createUserKey(client *UserInterface) bool {
	/*	var param []interface{} = []interface{}{(*client).GetId()}
		res := *(*client).GetDB().SelectSQL("SELECT user_key, id FROM users WHERE id=?", []string{"key", "id"}, param)
		if len(res) > 0 && res[0]["key"] != "" {
			log.Println("ERROR createUserKey 096736463, user key not NULL:", res)
			return true
		}

		arg := RandStringBytesMaskImprSrcUnsafe(255)
		var params [][]interface{} = [][]interface{}{ []interface{}{arg, (*client).GetId()} }
		result := (*client).GetDB().ExecSQL("UPDATE users SET user_key=? WHERE id=?", params)
		if (*result)[0].Err != nil {
			log.Println("ERROR createUserKey 749652641:", (*result)[0].Err.Error())
			return false
		}
		(*client).GetQuery()["uk"] = arg
*/
/****
	return true
}

func createIdentityKey(client *UserInterface) bool {

	lenaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println("ERROR createIdentityKey 48455051:", err)
		return false
	}
	lenaPublicKey := &lenaPrivateKey.PublicKey
	alisaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println("ERROR createIdentityKey 48455051:", err)
		return false
	}
	alisaPublicKey := &alisaPrivateKey.PublicKey
	log.Println("Private Key : ", lenaPrivateKey, "\n\n\n")
	log.Println("Public key ", lenaPublicKey, "\n\n\n")
	log.Println("Private Key : ", alisaPrivateKey, "\n\n\n")
	log.Println("Public key ", alisaPublicKey, "\n\n\n")

	/*arg := RandStringBytesMaskImprSrcUnsafe(255)
	var params [][]interface{} = [][]interface{}{ []interface{}{arg, (*client).GetId()} }
	res := (*client).GetDB().ExecSQL("UPDATE `users` SET `identification_key`=? WHERE id=?", params)
	if (*res)[0].Err != nil {
		log.Println("ERROR createIdentityKey 48455051:", (*res)[0].Err.Error())
		return false
	}
	(*client).GetQuery()["ik"] = arg*/
/*****
	return true
}

func CreatePasswordHash(password string) string {

	//	не реализован
	return password
}

func securityForPassword(client *UserInterface) bool {
	var (
		login, pass string
		ok          bool
	)
	if login, ok = (*client).GetQuery()["log"].(string); !ok {
		log.Println("ERROR security user log 68343763")
		return false
	}
	if pass, ok = (*client).GetQuery()["pass"].(string); !ok {
		log.Println("ERROR security user pass 465234766")
		return false
	}
	res := *(*client).GetDB().SelectSQL("SELECT id, password_hash FROM users WHERE login=?",
		[]string{"id", "password"}, []interface{}{login})
	if len(res) > 0 {
		if id, ok := res[0]["id"].(int64); ok {
			if CreatePasswordHash(pass) == res[0]["password"].(string) {
				(*client).SetId(id)
				return true
			} else {
				log.Println("ERROR security user pass 7846678")
			}
		}
	}

	log.Println("ERROR security user pass 83236234")
	return false
}

func securityIdentityKey(client *UserInterface) bool {
	res := *(*client).GetDB().SelectSQL("SELECT id FROM users WHERE identification_key=?",
		[]string{"id"}, []interface{}{(*client).GetQuery()["ik"].(string)})
	if len(res) > 0 {
		if id, ok := res[0]["id"]; ok {
			(*client).SetId(id.(int64))
			return true
		}
	}
	log.Println("ERROR security user ik 575562673657")
	return false
}

func BlackList(client *wss_server.Client) {
	//	db :=
	//	defer db.Close()

	/*query := "SELECT username, id_klient AS id, CONCAT(familiya, ' ', imya, ' ', otchestvo) AS fio FROM klient k " +
	" JOIN user u ON  k.id_klient=u.id WHERE u.id IN (SELECT locked FROM chat_black_list WHERE blocking=?) " +
	" UNION SELECT username, id_master AS id, CONCAT(familiya, ' ', imya, ' ', otchestvo) AS fio " +
	" FROM master m JOIN user u ON m.id_master=u.id WHERE u.id IN (SELECT locked FROM chat_black_list " +
	" WHERE blocking=?) UNION SELECT username, id_manager AS id, CONCAT(familiya, ' ', imya, ' ', otchestvo) " +
	" AS fio FROM manager mg JOIN user u ON mg.id_manager=u.id WHERE u.id IN (SELECT locked " +
	" FROM chat_black_list WHERE blocking=?)"
*/

/****
	rows, err := db.Query(query, strconv.Itoa(client.hub.clients[client].idUser), strconv.Itoa(client.hub.clients[client].idUser),
		strconv.Itoa(client.hub.clients[client].idUser))
	if err != nil {
		log.Println("ERROR: err => ", err)
		return
	}
	users := make(map[string]map[string]string)
	for rows.Next() {
		var uname, fio string
		var id int
		err = rows.Scan(&uname, &id, &fio)
		if err != nil {
			log.Fatal(err)
		}
		users[strconv.Itoa(id)] = map[string]string{"id": strconv.Itoa(id), "username": uname, "fio": fio}
	}
	rows.Close()

	response := map[string]map[string]map[string]string{
		"users":  users,
		"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_BLACK_LIST_USERS)}},
	}
	client.hub.sendClients(nil, response, client)
}

func RemoveUsersFromChat(client *Client, request *map[string]interface{}) {
	r := *request
	db := getDB()
	defer db.Close()

	idChat, err := strconv.Atoi(r["id_chat"].(string))
	if err != nil {
		log.Println(err)
		return
	}

	if !securityAuthor(idChat, client.hub.clients[client].idUser) {
		response := map[string]map[string]map[string]string{
			"status": {"status": {"status": strconv.Itoa(STATUS_ERROR), "operation": strconv.Itoa(OP_REMOVE_USER),
				"message": "Либо вы не являетесь автором чата, либо чат был заблокирован"}},
		}
		client.hub.sendClients(nil, response, client)
		return
	}

	var users = r["users"].([]interface{})
	for _, v := range users {
		query := "DELETE FROM chat_user WHERE id_chat=? AND id_user IN (?)"
		_, err = db.Exec(query, r["id_chat"], v)
		if err != nil {
			log.Println("ERROR: sql => ", err)
			return
		}
	}

	response := map[string]map[string]map[string]string{
		"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_REMOVE_USER),
			"message": "Успех"}},
	}
	client.hub.sendClients(nil, response, client)
}

func securityAuthor(idChat int, idUser int) bool {
	db := getDB()
	defer db.Close()

	query := "SELECT status FROM chat WHERE autor=? AND id=?"
	rows, err := db.Query(query, idUser, idChat)
	if err != nil {
		log.Println("ERROR: err => ", err)
		return false
	}
	var status string
	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			log.Fatal(err)
		}
	}
	rows.Close()

	if strings.Compare(status, CHAT_ACTIVE) == 0 {
		return true
	} else {
		return false
	}
}

func AddUserInChat(client *Client, request *map[string]interface{}) {
	/*r := *request
	db := getDB()
	defer db.Close()

	log.Println("BLOCK 1")

	idChat, err := strconv.Atoi(r["id_chat"].(string))
	if err != nil {
		log.Println("ERROR 865863:", err)
		return
	}

	log.Println("BLOCK 2")

	if !securityAuthor(idChat, client.hub.clients[client].idUser) {
		response := map[string]map[string]map[string]string{
			"status": {"status": {"status": strconv.Itoa(STATUS_ERROR), "operation": strconv.Itoa(OP_ADD_USER),
				"message": "Либо вы не являетесь автором чата, либо чат был заблокирован"}},
		}
		client.hub.sendClients(nil, response, client)
		return
	}

	log.Println("BLOCK 3")

	var users []interface{} = r["users"].([]interface{})
	for _, v := range users {
		query := "SELECT id_chat FROM chat_user WHERE id_chat=? AND id_user=?"
		rows, err := db.Query(query, r["id_chat"], v)
		if err != nil {
			log.Println("ERROR 5521556: err => ", err)
			return
		}
		var (
			id   int
			flag = true
		)
		for rows.Next() {
			err = rows.Scan(&id)
			if err != nil {
				log.Fatal("ERROR 56325632:", err)
			}
			flag = false
		}
		rows.Close()
		if flag {
			query = "INSERT INTO chat_user(id_chat, id_user) VALUES (?, ?)"
			_, err = db.Exec(query, r["id_chat"], v)
			if err != nil {
				log.Println("ERROR 75215662: sql => ", err)
				return
			}
		}
	}

	response := map[string]map[string]map[string]string{
		"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_ADD_USER),
			"message": "Успех"}},
	}
	client.hub.sendClients(nil, response, client)*/

/****
}

func SearchUser(client *Client, request *map[string]interface{}) {
	r := *request
	db := getDB()
	defer db.Close()

	var myMap map[string]interface{} = r["search"].(map[string]interface{})
	where := ""
	whereP := ""
	flag := false
	var param = make([]interface{}, 0)
	var paramP = make([]interface{}, 0)
	var p = make([]interface{}, 0)

	for k, v := range myMap {
		val := v.(string)
		if strings.Compare(val, "") != 0 && strings.Compare(k, "phone") != 0 {
			where = where + k + " LIKE ? AND "
			param = append(param, "%"+val+"%")
			paramP = append(paramP, "%"+val+"%")
			flag = true
		}
	}
	whereP = where + whereP

	if k, ok := myMap["phone"]; ok && k.(string) != "" {
		where = where + " phone LIKE ? AND "
		whereP = whereP + " phone1 LIKE ? OR phone2 LIKE ? OR phone3 LIKE ? AND "
		param = append(param, "%"+k.(string)+"%")
		for i := 0; i < 3; i++ {
			paramP = append(paramP, "%"+k.(string)+"%")
		}
		flag = true
	}

	if flag {
		where = " WHERE " + where
		whereP = " WHERE " + whereP
		where = where[:len([]rune(where))-4]
		whereP = whereP[:len([]rune(whereP))-4]
		for i := 0; i < 2; i++ {
			for _, v := range param {
				p = append(p, v)
			}
		}
		for _, v := range paramP {
			p = append(p, v)
		}
	}

	query := "SELECT username, id_klient AS id, CONCAT(familiya, ' ', imya, ' ', otchestvo) AS fio FROM klient k \n" +
		" JOIN user u ON k.id_klient=u.id " + where +
		" \n UNION SELECT username, id_master AS id, CONCAT(familiya, ' ', imya, ' ', otchestvo) AS fio FROM master m \n" +
		" JOIN user u ON m.id_master=u.id " + where +
		" \n UNION SELECT username, id_manager AS id, CONCAT(familiya, ' ', imya, ' ', otchestvo) AS fio FROM manager mg \n" +
		" JOIN user u ON mg.id_manager=u.id " + whereP
	//	log.Println(query, "; params =>", p)
	var rows *sql.Rows
	var err error
	if flag {
		rows, err = db.Query(query, p...)
	} else {
		rows, err = db.Query(query)
	}
	/*query := "SELECT username, id_klient AS id, CONCAT(familiya, ' ', imya, ' ', otchestvo) AS fio FROM klient k " +
			" JOIN user u ON k.id_klient=u.id " +
	//	"WHERE username LIKE ? AND phone LIKE ? AND familiya LIKE ?  AND imya LIKE ? AND otchestvo LIKE ? " +
			" UNION SELECT username, id_master AS id, CONCAT(familiya, ' ', imya, ' ', otchestvo) AS fio FROM master m " +
			" JOIN user u ON m.id_master=u.id " +
	//	"WHERE username LIKE ? AND phone LIKE ? AND familiya LIKE ?  AND imya LIKE ? AND otchestvo LIKE ? " +
			" UNION SELECT username, id_manager AS id, CONCAT(familiya, ' ', imya, ' ', otchestvo) AS fio FROM manager mg " +
			" JOIN user u ON mg.id_manager=u.id " +
		//"WHERE username LIKE ? AND phone1 LIKE ? AND phone2 LIKE ? AND phone3 LIKE ? AND familiya LIKE ? AND imya LIKE ? AND otchestvo LIKE ?"

	rows, err := db.Query(query, "%" + myMap["username"].(string) + "%", "%" + myMap["phone"].(string) + "%", "%" + myMap["familiya"].(string) + "%", "%" + myMap["imya"].(string) + "%", "%" + myMap["otchestvo"].(string) + "%",
			"%" + myMap["username"].(string) + "%", "%" + myMap["phone"].(string) + "%", "%" + myMap["familiya"].(string) + "%", "%" + myMap["imya"].(string) + "%", "%" + myMap["otchestvo"].(string) + "%",
			"%" + myMap["username"].(string) + "%", "%" + myMap["phone"].(string) + "%", "%" + myMap["phone"].(string) + "%", "%" + myMap["phone"].(string) + "%", "%" + myMap["familiya"].(string) + "%", "%" + myMap["imya"].(string) + "%", "%" + myMap["otchestvo"].(string) + "%",
	)*/
/****
	if err != nil {
		log.Println("ERROR: err => ", err)
		return
	}
	var (
		uname, fio sql.NullString
		id         int
		resMap     = make(map[string]map[string]string)
	)
	for rows.Next() {
		err = rows.Scan(&uname, &id, &fio)
		if err != nil {
			log.Fatal(err)
		}
		resMap[strconv.Itoa(id)] = map[string]string{"id": strconv.Itoa(id), "username": uname.String, "fio": fio.String}
	}
	rows.Close()

	response := map[string]map[string]map[string]string{
		"users":  resMap,
		"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_SEARCH_USER)}},
	}
	client.hub.sendClients(nil, response, client)
}

func unlockUsers(client *Client, request *map[string]interface{}) {
	r := *request
	var idUsers []interface{} = r["users"].([]interface{})

	if len(idUsers) <= 0 {
		log.Println("Массив пользователей пуст  len(idUsers) <= 0")
		return
	}

	db := getDB()
	defer db.Close()
	query := "DELETE FROM chat_black_list WHERE blocking=" + strconv.Itoa(client.hub.clients[client].idUser) +
		" AND locked=?"
	for _, v := range idUsers {
		_, err := db.Exec(query, v)
		if err != nil {
			log.Println("ERROR: sql => ", err)
			return
		}
	}

	mass := make(map[string]map[string]string)
	for _, v := range idUsers {
		mass[v.(string)] = nil
	}
	response := map[string]map[string]map[string]string{
		"users": mass,
		"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_UNLOOCK_USERS),
			"blackList": fmt.Sprint(r["blackList"])}},
	}

	client.hub.sendClients(nil, response, client)

}

func BlockUsers(client *Client, request *map[string]interface{}) {
	r := *request
	var idUsers []interface{} = r["users"].([]interface{})
	var idUs []string

	/*if err != nil {
		fmt.Println("ERROR 321 => ", err)
		return
	}*/
/*****
	if len(idUsers) <= 0 {
		log.Println("Массив пользователей пуст  len(idUsers) <= 0")
		return
	}
	for _, v := range idUsers {
		idUs = append(idUs, v.(string))
	}

	db := getDB()
	defer db.Close()
	str1 := strings.Join(idUs, ", ")
	query := "SELECT locked FROM chat_black_list WHERE blocking=" + strconv.Itoa(client.hub.clients[client].idUser) +
		" AND locked IN (?)"
	rows, err := db.Query(query, str1)
	if err != nil {
		log.Println("ERROR: err => ", err)
		return
	}

	var id []int
	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			log.Fatal(err)
		}
		id = append(id, i)
	}
	rows.Close()

	for _, v := range idUs {
		var flag = true
		for _, val := range id {
			if d, err := strconv.Atoi(v); err == nil && d == val {
				flag = false
			}
		}
		if flag {
			query := "INSERT INTO chat_black_list (blocking, locked, date, time) VALUES (" +
				strconv.Itoa(client.hub.clients[client].idUser) + ", ?, CURDATE(), CURTIME())"
			_, err := db.Exec(query, v)
			if err != nil {
				log.Println("ERROR: sql => ", err)
				return
			}
		}
	}

	mass := make(map[string]map[string]string)
	for _, v := range idUs {
		mass[v] = map[string]string{}
	}

	response := map[string]map[string]map[string]string{
		"users": mass,
		"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_BLOCK_USERS),
			"blackList": fmt.Sprint(r["blackList"])}},
	}
	client.hub.sendClients(nil, response, client)

}

func ListUsers(client *Client, request *map[string]interface{}) {

	r := *request
	db := getDB()
	defer db.Close()
	//	выбираем данные всех пользователей чата
	query := "SELECT id_user AS user_id, autor, username, CONCAT(IF(client_fio IS NULL, '', client_fio), IF(manager_fio IS NULL, " +
		" '', manager_fio),  IF(master_fio IS NULL, '', master_fio)) AS fio  FROM (SELECT res.*, " +
		" CONCAT(k.familiya, ' ', k.imya, ' ', k.otchestvo) AS client_fio FROM (SELECT mas.*, CONCAT(mg.familiya, " +
		" ' ', mg.imya, ' ', mg.otchestvo) AS manager_fio FROM (SELECT id_user, autor, username, " +
		" CONCAT(m.familiya, ' ', m.imya, ' ', m.otchestvo) AS master_fio  FROM chat_user cu, chat c, user u " +
		" LEFT JOIN master m ON m.id_master=u.id WHERE cu.id_chat=c.id AND u.id=cu.id_user AND id_chat=?) mas " +
		" LEFT JOIN manager mg ON mas.id_user=mg.id_manager) res LEFT JOIN klient k ON res.id_user=k.id_klient) o"
	rows, err := db.Query(query, r["id"])
	if err != nil {
		log.Println("ERROR: err => ", err)
		return
	}

	var (
		user_id, autor int
		username       string
		fio            sql.NullString
	)
	users := map[string]map[string]string{}
	for rows.Next() {
		err = rows.Scan(&user_id, &autor, &username, &fio)
		if err != nil {
			log.Fatal(err)
		}
		users[strconv.Itoa(user_id)] = map[string]string{"user_id": strconv.Itoa(user_id), "autor": strconv.Itoa(autor),
			"username": username, "fio": fio.String, "connected": "false"}
	}
	rows.Close()

	//	выбираем всех пользователей внесенных в черный список данным user
	blackList := map[string]map[string]string{}
	query = "SELECT locked FROM chat_black_list WHERE blocking=" + strconv.Itoa(client.hub.clients[client].idUser)
	rows, err = db.Query(query)
	if err != nil {
		log.Println("ERROR: err => ", err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&user_id)
		if err != nil {
			log.Fatal(err)
		}
		blackList[strconv.Itoa(user_id)] = map[string]string{}
	}
	rows.Close()

	for _, val := range client.hub.clients {
		if val.conn {
			if id, ok := users[strconv.Itoa(val.idUser)]; ok {
				id["connected"] = "true"
			}
		}
	}

	response := map[string]map[string]map[string]string{
		"users": users, "black_list": blackList,
		"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_LIST_USERS)}},
	}

	client.hub.sendClients(nil, response, client)
}

func UserWrite(client *Client, request *map[string]interface{}) {
	r := *request
	idChat, err := strconv.Atoi(fmt.Sprint(r["id_chat"]))
	if err != nil {
		log.Println(err)
		return
	}
	idUsers, err := getIdUsersOfChat(idChat)

	response := map[string]map[string]map[string]string{
		"users": {strconv.Itoa(client.hub.clients[client].idUser): {"id_user": strconv.Itoa(client.hub.clients[client].idUser),
			"username": client.hub.clients[client].userName, "id_chat": fmt.Sprint(r["id_chat"]), "write": fmt.Sprint(r["write"])}},
		"status": {"status": {"status": strconv.Itoa(STATUS_ACCEPT), "operation": strconv.Itoa(OP_WRITEN)}},
	}
	client.hub.sendClients(idUsers, response, nil)
}

func getIdUsersOfChat(idChat int) ([]int, error) {
	db := getDB()
	defer db.Close()

	var idUsers = []int{}
	query := "SELECT id_user FROM chat_user cu, chat c WHERE c.id=cu.id_chat AND c.status=\"" + CHAT_ACTIVE +
		"\" AND id_chat=" + strconv.Itoa(idChat)
	log.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		idUsers = append(idUsers, id)
	}
	rows.Close()
	log.Println("idUsers 1 => ", idUsers)
	return idUsers, nil
}

func setName(client *Client, request *map[string]interface{}) {
	data := *request
	db := getDB()
	defer db.Close()
	query := "SELECT id, username, status, imei FROM user WHERE username=? LIMIT 1"
	rows, err := db.Query(query, data["name"])
	if err != nil {
		log.Println("ERROR 9218: sql => ", err)
		return
	}

	for rows.Next() {
		var hit User
		err = rows.Scan(&hit.id, &hit.username, &hit.status, &hit.imei)
		if err != nil {
			log.Fatal(err)
			return
		}
		/*	if hit.id == 0 {
			//	client.hub.unregister <- client
				log.Println("id не найден, login => ", data["name"])
				return
			}*/
/****
	if imei, ok := data["code"]; ok {
		if strings.Compare(imei.(string), hit.imei.String) != 0 {
			log.Println("imei не совпали")
			return
		}
	}
	response := make(map[string]map[string]map[string]string, 0)
	response["status"] = map[string]map[string]string{"status": {"status": fmt.Sprint(STATUS_ACCEPT), "operation": fmt.Sprint(OP_SET_USER_NAME)}}
	response["users"] = map[string]map[string]string{strconv.Itoa(hit.id): {"id": strconv.Itoa(hit.id), "username": hit.username}}
	log.Println("RETURN => ", response)

	str, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}
	client.hub.clients[client] = &user.ClientData{idUser: hit.id, conn: true, userName: hit.username}
	client.send <- str

}
rows.Close()

/*if client.hub.clients[client].conn {
	query := "SELECT id_message FROM chat_message_status WHERE id_user=? AND status_message <> \"readed\" LIMIT 1"
	rows, err := db.Query(query, client.hub.clients[client].idUser)
	if err != nil { log.Println("ERROR 3462: sql => ", err); return;	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil  { log.Fatal(err); return; }

		response := map[string]map[string]map[string]string{
			"status": { "status": {"status": fmt.Sprint(STATUS_ACCEPT), "operation": fmt.Sprint(OP_HAVE_MESSAGE)}}}
		client.hub.sendClients(nil, response, client)
	}
	rows.Close()
}*/

/*	query = "SELECT id_chat FROM `chat_user` cu, chat c WHERE c.id=cu.id_chat AND c.status=\"" + CHAT_ACTIVE +
			"\" AND id_user=" + strconv.Itoa(hit.id)
	rows, err = db.Query(query)
	if err != nil { log.Println("ERROR: sql => ", err); return;	}
	var idChats = make([]int, 0)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil { log.Fatal(err);	return;	}
		idChats = append(idChats, id)
	}
	rows.Close()

	//client.hub.broadcast[1] = make(chan []byte)

	for index := range idChats {
		if _, ok := client.hub.chatUsers[idChats[index]]; !ok {
			qw := make([]*Client, 0)
			qw = append(qw, client)
			client.hub.chatUsers[idChats[index]] = qw
		//	client.hub.broadcast[idChats[index]] = make(chan []byte)
		} else {
			client.hub.chatUsers[idChats[index]] = append(client.hub.chatUsers[idChats[index]], client)
		}
	}
	// */

/****}*/
