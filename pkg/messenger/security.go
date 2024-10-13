package wss_server

func securityRequest(client *Client) bool {
	/*****data := client.Query
	if v, ok := data["action"]; !ok || v == nil || v == "" {
		log.Println("ERROR: operation => ", v)
		return false
	}
	if int(data["action"].(float64)) != OP_SET_USER_NAME {
		if client.Hub.clients[client].conn != true {
			log.Println("ERROR: client not autorizated!!!")
			return false
		}
	}
	switch int(data["action"].(float64)) {
	case OP_MY_DATA:
		return true
	case OP_NEW_MESSAGE:
		if v, ok := data["message"]; !ok || v == nil || v == "" {
			log.Println("ERROR: message => ", v)
			return false
		}
		if v, ok := data["id_chat"]; !ok || v == nil || v == "" {
			log.Println("ERROR: id_chat => ", v)
			return false
		}
	case OP_STATUS_MESSAGE:
		if v, ok := data["id"]; !ok || v == nil || v == "" {
			log.Println("ERROR: id => ", v)
			return false
		}
		if v, ok := data["id_chat"]; !ok || v == nil || v == "" {
			log.Println("ERROR: id_chat => ", v)
			return false
		}
		if v, ok := data["status"]; !ok || v == nil || v == "" {
			log.Println("ERROR: status => ", v)
			return false
		}
		str := fmt.Sprint(data["status"])
		if strings.Compare(str, MESSAGE_READED) != 0 && strings.Compare(str, MESSAGE_DELIVERED) != 0 {
			log.Println("ERROR: status => ", str)
			return false
		}
	case OP_GET_CHATS:
		return true
	case OP_GET_HISTORY_MESSAGE:
		if v, ok := data["id"]; !ok || v == nil || v == "" {
			log.Println("ERROR: id => ", v)
			return false
		}
	case OP_BLACK_LIST_USERS:
		return true
	case OP_LIST_USERS:
		if v, ok := data["id"]; !ok || v == nil || v == "" {
			log.Println("ERROR: id => ", v)
			return false
		}
	case OP_REMOVE_CHAT:
		if v, ok := data["id"]; !ok || v == nil || v == "" {
			log.Println("ERROR: id => ", v)
			return false
		}
	case OP_EXIT_CHAT:
		if v, ok := data["id"]; !ok || v == nil || v == "" {
			log.Println("ERROR: id => ", v)
			return false
		}
	case OP_CREATE_NEW_CHAT:
		if v, ok := data["chat_name"]; !ok || v == nil || v == "" {
			log.Println("ERROR: chat_name => ", v)
			return false
		}
		if v, ok := data["users"]; !ok || v == nil || v == "" {
			log.Println("ERROR: users => ", v)
			return false
		}
	case OP_ADD_USER:
		if v, ok := data["id_chat"]; !ok || v == nil || v == "" {
			log.Println("ERROR: id => ", v)
			return false
		}
		if v, ok := data["users"]; !ok || v == nil || v == "" {
			log.Println("ERROR: users => ", v)
			return false
		}
	case OP_SEARCH_USER:
		if v, ok := data["search"]; !ok || v == nil || v == "" {
			log.Println("ERROR: search => ", v)
			return false
		}
	case OP_REMOVE_USER:
		if v, ok := data["id_chat"]; !ok || v == nil || v == "" {
			log.Println("ERROR: id => ", v)
			return false
		}
		if v, ok := data["users"]; !ok || v == nil || v == "" {
			log.Println("ERROR: users => ", v)
			return false
		}
	case OP_BLOCK_USERS:
		if v, ok := data["users"]; !ok || v == nil || v == "" {
			log.Println("ERROR: users => ", v)
			return false
		}
		if v, ok := data["blackList"]; !ok || v == nil || v == "" {
			log.Println("ERROR: blackList => ", v)
			return false
		}
	case OP_UNLOOCK_USERS:
		if v, ok := data["users"]; !ok || v == nil || v == "" {
			log.Println("ERROR: users => ", v)
			return false
		}
		if v, ok := data["blackList"]; !ok || v == nil || v == "" {
			log.Println("ERROR: blackList => ", v)
			return false
		}
	case OP_WRITEN:
		if v, ok := data["id_chat"]; !ok || v == nil || v == "" {
			log.Println("ERROR: id_chat => ", v)
			return false
		}
		if v, ok := data["write"]; !ok || v == nil || v == "" {
			log.Println("ERROR: write => ", v)
			return false
		}
	case OP_SET_USER_NAME:
		if v, ok := data["name"]; !ok || v == nil || v == "" {
			log.Println("ERROR: name => ", v)
			return false
		}
	case OP_SYSTEM:
	/***	if v, ok := data["id"]; !ok || v == nil || v == "" {
		log.Println("ERROR: id => ", v)
		return false
	}
	if v, ok := data["status"]; !ok || v == nil || v == "" {
		log.Println("ERROR: status => ", v)
		return false
	}
	status := int(data["status"].(float64))
	/*if err != nil {
		log.Println("ERROR: status not int => ", data["status"])
		return false;
	}*/
	/***	if status != ORDER_AKTIVATE && status != ORDER_DIAKTIVATE && status != ORDER_EXECUTED {
		log.Println("ERROR: status value error =>", status)
		return false
	}***/
/*****	default:
		//	log.Println("NOT SUPPORT OPERATION!!!")
		return false
	}*****/
	return true
}
