// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wss_server

import (
	"log"
	//	"github.com/gorilla/websocket"
	//	"log"
	//	"security_user"
	bl "github.com/alex988334/messenger/pkg/messenger/business-logic"
	"github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/functions"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

type Hub struct {
	//	Массив зарегистрированных клиентов. ключом служит номер чата
	// ключи: номер чата, номер пользователя, подключенный клиент пользователя
	chats map[int]map[int]map[*Client]bool

	// хранилище подключенных клиентов не входящих в чаты
	noHaveChatClients map[*Client]bool

	// добавление или удаление соединений в пул чатов
	actionUserConnections chan bl.ActionUserConnections

	// Входящие сообщения для чатов.
	broadcast chan struct {
		message string
		idChats []int
	}
	// Регистрация запросов от клиентов.
	autorizate chan struct {
		client  *Client
		idChats []int
	}
	register chan *Client
	// Отмена регистрации запросов от клиентов.
	unregister chan struct {
		client  *Client
		idChats []int
	}
	//	Канал для остановки хаба
	stopServer chan bool
	//	переменная остановки работы хаба
	running bool
}

func NewHub() *Hub {
	return &Hub{
		chats:             make(map[int]map[int]map[*Client]bool),
		noHaveChatClients: make(map[*Client]bool),
		autorizate: make(chan struct {
			client  *Client
			idChats []int
		}),
		register: make(chan *Client),
		unregister: make(chan struct {
			client  *Client
			idChats []int
		}),
		broadcast: make(chan struct {
			message string
			idChats []int
		}),
		actionUserConnections: make(chan bl.ActionUserConnections),
		stopServer:            make(chan bool),
		running:               false,
	}
}

/*
func (h *Hub) getStatistics() map[string]interface{} {
	data := make(map[string]interface{})
	if h.running {
		data["Hub status"] = "active"
	} else {
		data["Hub status"] = "stoped"
	}
	data["Connected not register"] = len(h.unautorizateClients)
	data["Connected to chats"] = len(h.chats)
	return data
}
*/

// если найдет клиента в хранилище неавторизованных пользователей то вернет true
func (h *Hub) isAutorizate(client *Client) bool {

	if isAutoriz, ok := h.noHaveChatClients[client]; ok && !isAutoriz {
		return false
	}

	return true
}

func (h *Hub) sendMessage(message string, idChats []int) {
	for _, v := range idChats {
		for key, clients := range h.chats[v] {
			for client, _ := range clients {
				select {
				case client.Send <- []byte(message):
				default:
					close(client.Send)
					delete(h.chats[v][key], client)
				}
			}
		}
	}
}

func (h *Hub) autorizateClient(client *Client, idChats []int) {

	if len(idChats) > 0 && idChats[0] > 0 {
		delete(h.noHaveChatClients, client)
	} else {
		h.noHaveChatClients[client] = true
	}

	for _, id := range idChats {

		if id == 0 {
			continue
		}
		if _, ok := h.chats[id]; !ok {
			h.chats[id] = make(map[int]map[*Client]bool)
		}
		if _, ok := h.chats[id][client.IdUser]; !ok {
			h.chats[id][client.IdUser] = make(map[*Client]bool)
		}
		if _, ok := h.chats[id][client.IdUser][client]; !ok {
			h.chats[id][client.IdUser][client] = true
		}
	}
}

func (h *Hub) unautorizateClient(client *Client, idChats []int) {

	h.noHaveChatClients[client] = false

	for _, v := range idChats {
		id := v
		if _, ok := h.chats[id]; !ok {
			continue
		}
		if _, ok := h.chats[id][client.IdUser]; !ok {
			continue
		}

		for conn, _ := range h.chats[id][client.IdUser] {
			h.noHaveChatClients[conn] = false
		}

		delete(h.chats[id][client.IdUser], client)

		if len(h.chats[id][client.IdUser]) == 0 {
			delete(h.chats[id], client.IdUser)
		}
	}
}

func (h *Hub) appendUserConnFromChats(userId int, connectingChats []int, connectedChats []int) {

	findedClients := h.getCollectionConnections(userId, connectedChats)

	for _, id := range connectingChats {

		if _, ok := h.chats[id]; !ok {
			h.chats[id] = map[int]map[*Client]bool{}
		}
		if _, ok := h.chats[id][userId]; !ok {
			h.chats[id][userId] = map[*Client]bool{}
		}

		for conn, _ := range findedClients {
			h.chats[id][userId][conn] = true
		}
	}
}

func (h *Hub) appendUserConnFromNoChats(userId int, connectingChats []int) {

	for conn, isAutorizated := range h.noHaveChatClients {

		if !isAutorizated || conn.GetId() != userId {
			continue
		}

		for _, connecting := range connectingChats {

			if _, ok := h.chats[connecting]; !ok {
				h.chats[connecting] = map[int]map[*Client]bool{}
			}
			if _, ok := h.chats[connecting][userId]; !ok {
				h.chats[connecting][userId] = map[*Client]bool{}
			}

			h.chats[connecting][userId][conn] = true
		}
		delete(h.noHaveChatClients, conn)
	}
}

func (h *Hub) getCollectionConnections(userId int, connectedChats []int) map[*Client]bool {

	clients := make(map[*Client]bool)

	for _, connected := range connectedChats {

		if _, ok := h.chats[connected]; !ok {
			continue
		}
		if _, ok := h.chats[connected][userId]; !ok {
			continue
		}

		for conn, _ := range h.chats[connected][userId] {
			clients[conn] = true
		}
	}

	return clients
}

func (h *Hub) appendUserToChats(userId int, connectingChats []int, connectedChats []int) {

	if len(connectingChats) > 0 && connectingChats[0] > 0 {

		h.appendUserConnFromChats(userId, connectingChats, connectedChats)
		h.appendUserConnFromNoChats(userId, connectingChats)
	}
}

func (h *Hub) removeChat(usersId []int, disconnectingChats []int, connectedChats []int) {

	for i := 0; i < len(usersId); i++ {
		h.removeUserFromChats(usersId[i], disconnectingChats, connectedChats)
	}
}

func (h *Hub) removeUserFromChats(userId int, disconnectingChats []int, connectedChats []int) {

	for _, v := range disconnectingChats {
		connectedChats = functions.RemoveFromArray(connectedChats, v)
	}

	allSavedConns := h.getCollectionConnections(userId, connectedChats)

	for _, id := range disconnectingChats {
		if users, ok := h.chats[id]; ok {
			if conns, ok := users[userId]; ok {
				for c, _ := range conns {
					if _, ok = allSavedConns[c]; !ok {
						h.noHaveChatClients[c] = true
					}
					delete(h.chats[id][userId], c)
				}

			}
			delete(h.chats[id], userId)
		}
		if len(h.chats[id]) == 0 {
			delete(h.chats, id)
		}
	}
}

func (h *Hub) stopHub() {

	h.running = false
	for k, v := range h.chats {
		for ke, va := range v {
			for key, _ := range va {
				key.Conn.Close()
				close(key.Send)
				delete(h.chats[k][ke], key)
			}
		}
	}

	for key, _ := range h.noHaveChatClients {
		key.Conn.Close()
		close(key.Send)
		delete(h.noHaveChatClients, key)
	}
}

func (h *Hub) registerClient(client *Client) {

	h.noHaveChatClients[client] = false
}

func (h *Hub) Run() {
	h.running = true
	log.Println("Hub running")

	for h.running {
		select {
		case data := <-h.autorizate:
			h.autorizateClient(data.client, data.idChats)
		case client := <-h.register: //	при поступлении данных в канал регистрации
			h.registerClient(client)
		case data := <-h.unregister: //	при поступлении данных в канал отмены регистрации
			h.unautorizateClient(data.client, data.idChats)
		case <-h.stopServer:
			h.stopHub()
		case unit := <-h.broadcast:
			h.sendMessage(unit.message, unit.idChats)
		case actUserConn := <-h.actionUserConnections:
			if actUserConn.Operation == constants.CLIENT_CONNECT {
				h.appendUserToChats(actUserConn.UsersId[0], actUserConn.ActionChats, actUserConn.ConnectedChats)
			}
			if actUserConn.Operation == constants.CLIENT_DISCONNECT {
				h.removeUserFromChats(actUserConn.UsersId[0], actUserConn.ActionChats, actUserConn.ConnectedChats)
			}
			if actUserConn.Operation == constants.ALL_CLIENT_DISCONNECT {
				h.removeChat(actUserConn.UsersId, actUserConn.ActionChats, actUserConn.ConnectedChats)
			}
		}
	}
	log.Println("Hub stop")
}
