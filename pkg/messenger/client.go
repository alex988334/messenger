// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wss_server

import (
	//	"database/sql"

	//	"database/sql"

	//_ "github.com/go-sql-driver/mysql"
	"fmt"

	"encoding/json"
	"net/http"
	"time"

	bl "github.com/alex988334/messenger/pkg/messenger/business-logic"
	"github.com/alex988334/messenger/pkg/messenger/constants"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the Hub.
type Client struct {
	Hub *Hub
	// The websocket connection.
	Conn *websocket.Conn
	// Buffered channel of outbound messages.
	Send   chan []byte
	IdUser int
	//	DB        db.DBInterface
	HttpQuery map[string]string
}

func (client *Client) SetId(id int) {
	client.IdUser = id
}

func (client *Client) GetId() int {
	return client.IdUser
}

func (client *Client) IsAutorizate() bool {

	if isAutorizate, ok := client.Hub.noHaveChatClients[client]; ok && !isAutorizate {
		return false
	}

	return true
}

// readPump pumps messages from the websocket connection to the Hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.

func (client *Client) readPump() {
	fmt.Println("readPump")

	defer client.closeConnection()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, mess, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Is Unexpected Close Error")
			}
			fmt.Printf("readPump error: %v", err)
			break
		}
		var data map[string]string
		if err := json.Unmarshal(mess, &data); err != nil {
			fmt.Println("НЕ удалось распарсить Message:", mess)
			continue
		}

		client.HttpQuery = data

		busLogic := bl.NewBusinessHandler(&data, client, nil)
		response, idch, err := busLogic.ProcessinRequest()
		operat := busLogic.GetOperation()

		fmt.Println("\nRequest -> complete; operation ->",
			constants.GetDescriptionOperation(operat), "\nResponse ->", string(response))

		if err != nil || isSendOnlyClient(operat, idch) {

			if cn := busLogic.GetActionUserConnections(); cn != nil &&
				cn.Operation == constants.CLIENT_CONNECT_ALL_CHATS {

				client.Hub.autorizate <- struct {
					client  *Client
					idChats []int
				}{
					client:  client,
					idChats: cn.ConnectedChats,
				}
			}

			client.Send <- response
			continue
		} else {
			if cn := busLogic.GetActionUserConnections(); cn != nil &&
				cn.Operation != constants.CLIENT_CONNECT_ALL_CHATS {

				var c bl.ActionUserConnections = *cn
				client.Hub.actionUserConnections <- c
			}

			if len(idch) > 0 {

				client.Hub.broadcast <- struct {
					message string
					idChats []int
				}{message: string(response), idChats: idch}
			}
		}
	}
}

func isSendOnlyClient(operation int, sendedChats []int) bool {

	ok := false

	switch operation {
	case constants.OP_LIST_USERS:
		ok = true
	case constants.OP_GET_CHATS:
		ok = true
	case constants.OP_LIST_NEXT_MESSAGES:
		ok = true
	case constants.OP_LIST_PREVIOUS_MESSAGES:
		ok = true
	case constants.OP_SEARCH_USER:
		ok = true
	case constants.OP_BLOCK_USERS:
		ok = true
	case constants.OP_UNLOOCK_USERS:
		ok = true
	case constants.OP_BLACK_LIST_USERS:
		ok = true
	case constants.OP_AUTORIZATE:
		ok = true
	case constants.OP_REGISTRATION:
		ok = true
	}

	return ok && (len(sendedChats) == 0)
}

func (client *Client) closeConnection() {

	/*
	   var u user.UserInterface = client
	   w := user.GetUserChatsId(&u)

	   	if len(w) > 0 {
	   		client.Hub.unregister <- struct {
	   			client  *Client
	   			idChats []interface{}
	   		}{client: client, idChats: w}
	   	}

	   //   	client.DB.CloseDB()
	   client.Conn.Close()
	*/
}

func (client *Client) autorizationClient() {
	/******
	var u user.UserInterface = client
	w := user.GetUserChatsId(&u)
	if len(w) > 0 {
		client.Hub.autorizate <- struct {
			client  *Client
			idChats []interface{}
		}{client: client, idChats: w}
	}

	/*	_ := message.Message{12, 12, 12, sql.NullInt64{}, "gdfg", "wrer", "ewrwe",
		sql.NullString{"sdfsd", true}}*/

}

func (client *Client) excRequest(data *map[string]interface{}) {
	//var resp
	/******
		switch int((*data)["action"].(float64)) {
		case OP_NEW_MESSAGE:
			newMessage(client, &data)
		case OP_STATUS_MESSAGE:
			setStatus(client, &data)
		case OP_GET_CHATS:
			getChats(client)
		case OP_GET_HISTORY_MESSAGE:
			getHistoryMessage(client, &data)
		case OP_BLACK_LIST_USERS:
			user.BlackList(client)
		case OP_LIST_USERS:
			listUsers(client, &data)
		case OP_REMOVE_CHAT:
			removeChat(client, &data)
		case OP_EXIT_CHAT:
			exitChat(client, &data)
		case OP_CREATE_NEW_CHAT:
			createNewChat(client, &data)
		case OP_ADD_USER:
			addUserInChat(client, &data)
		case OP_SEARCH_USER:
			searchUser(client, &data)
		case OP_REMOVE_USER:
			removeUsersFromChat(client, &data)
		case OP_BLOCK_USERS:
			blockUsers(client, &data)
		case OP_UNLOOCK_USERS:
			unlockUsers(client, &data)
		case OP_WRITEN:
			userWrite(client, &data)
		case OP_SET_USER_NAME:
			setName(client, &data)
		case OP_SYSTEM:
			systemMessage(client, &data)
		case OP_MY_DATA:
			getMyData(client)
		default:
			fmt.Println("NOT SUPPORT OPERATION!!!")
		}
	*****/
}

// writePump перекачивает сообщения из концентратора в соединение websocket.
// writePump pumps messages from the Hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (client *Client) writePump() {

	/*	defer func() {
			if r := recover(); r != nil {
				fmt.Println("SERVER RECOVERED writePump() =>", r)
			}
		}()
	*/

	fmt.Println("START writePump()")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		fmt.Println("STOP writePump()")
		ticker.Stop()
		client.Conn.Close()
	}()

	newline := []byte{'\n'}
	//	space   := []byte{' '}

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The Hub closed the channel.
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, responseWriter http.ResponseWriter, request *http.Request) {

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(responseWriter, request, nil) //	создается соединение WS
	if err != nil {                                             //	проверяем на ошибки
		fmt.Println("upgrader.Upgrade() error =>", err)
		return
	}

	conn.SetReadLimit(10240000)
	client := &Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		IdUser: 0,
		//	DB: nil
	} //	создаем структуру клиента на базе соединения WS
	client.Hub.register <- client //	отправляем готового клиента в канал регистрации хаба

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump() //	запускаем поток писания клиента
	go client.readPump()  //	запускаем поток чтения клиента

	/*defer func() {
		if r := recover(); r != nil {
			fmt.Println("SERVER RECOVERED =>", r)
		}
	}()*/
}
