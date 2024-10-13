package security_module

type IClient interface {
	IsAutorizate() bool
	GetId() int
}

type SecurityHandler struct {
	client *IClient
}

/*
func (s *SecurityHandler) setClient(client *IClient) {
	s.client = client
}*/
/*
func (s *SecurityHandler) CheckRequest(request *map[string]interface{}, client IClient) bool {

	s.client = &client

	data := *request

	if v, ok := data[cons.KEY_ACTION]; !ok || v == nil || v == "" {
		log.Println("ERROR security handler:: operation => ", v)
		return false
	}

	var action int = int(data[cons.KEY_ACTION].(float64))
	if action != cons.OP_SET_USER_NAME && !(*s.client).IsAutorizate() {
		fmt.Println("ERROR security handler:: Client is not autorizated!")
		return false
	}
	data[cons.KEY_AUTHOR] = client.GetId()

	switch action {
	case cons.OP_MY_DATA:
		return true
	case cons.OP_NEW_MESSAGE:
		if v, ok := data[cons.KEY_MESSAGE]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: message =>", v)
			return false
		}
		if v, ok := data[cons.KEY_CHAT]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: id_chat =>", v)
			return false
		}
	case cons.OP_STATUS_MESSAGE:
		if v, ok := data[cons.KEY_ID]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: id =>", v)
			return false
		}
		if v, ok := data[cons.KEY_CHAT]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: id_chat => ", v)
			return false
		}
		if v, ok := data[cons.KEY_STATUS]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: status => ", v)
			return false
		}
		if data[cons.KEY_STATUS] != cons.MESSAGE_READED && data[cons.KEY_STATUS] != cons.MESSAGE_DELIVERED {
			log.Println("ERROR security handler:: status => ", data[cons.KEY_STATUS])
			return false
		}
	case cons.OP_GET_CHATS:
		return true
	case cons.OP_LIST_PREVIOUS_MESSAGES:
		if v, ok := data[cons.KEY_CHAT]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: KEY_CHAT => ", v)
			return false
		}
		if v, ok := data[cons.KEY_ID]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: id => ", v)
			return false
		}
	case cons.OP_BLACK_LIST_USERS:
		return true
	case cons.OP_LIST_USERS:
		if v, ok := data[cons.KEY_ID]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: id => ", v)
			return false
		}
	case cons.OP_REMOVE_CHAT:
		if v, ok := data[cons.KEY_ID]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: id => ", v)
			return false
		}
	case cons.OP_EXIT_CHAT:
		if v, ok := data[cons.KEY_USERS]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: KEY_USERS => ", v)
			return false
		}
		if v, ok := data[cons.KEY_CHAT]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: KEY_CHAT => ", v)
			return false
		}
	case cons.OP_CREATE_NEW_CHAT:
		if v, ok := data[cons.KEY_CHAT_NAME]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: chat_name => ", v)
			return false
		}
		if v, ok := data[cons.KEY_USERS]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: users => ", v)
			return false
		}

	case cons.OP_ADD_USER:
		if v, ok := data[cons.KEY_CHAT]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: id => ", v)
			return false
		}
		if v, ok := data[cons.KEY_USERS]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: users => ", v)
			return false
		}
	case cons.OP_SEARCH_USER:
		if v, ok := data[cons.KEY_SEARCH]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: search => ", v)
			return false
		}
	case cons.OP_REMOVE_USER:
		if v, ok := data[cons.KEY_CHAT]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: id => ", v)
			return false
		}
		if v, ok := data[cons.KEY_USERS]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: users => ", v)
			return false
		}
	case cons.OP_BLOCK_USERS:
		if v, ok := data[cons.KEY_USERS]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: users => ", v)
			return false
		}
		if v, ok := data[cons.KEY_BLACK_LIST]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: blackList => ", v)
			return false
		}
	case cons.OP_UNLOOCK_USERS:
		if v, ok := data[cons.KEY_USERS]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: users => ", v)
			return false
		}
		if v, ok := data[cons.KEY_BLACK_LIST]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: blackList => ", v)
			return false
		}
	case cons.OP_WRITEN:
		if v, ok := data[cons.KEY_CHAT]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: id_chat => ", v)
			return false
		}
		if v, ok := data[cons.KEY_WRITE]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: write => ", v)
			return false
		}
	case cons.OP_SET_USER_NAME:
		if v, ok := data[cons.KEY_NAME]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: name => ", v)
			return false
		}
	case cons.OP_SYSTEM:
		if v, ok := data[cons.KEY_ID]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: id => ", v)
			return false
		}
		if v, ok := data[cons.KEY_STATUS]; !ok || v == nil || v == "" {
			log.Println("ERROR security handler:: status => ", v)
			return false
		}
		status := int(data[cons.KEY_STATUS].(float64))
		if status != cons.ORDER_AKTIVATE && status != cons.ORDER_DIAKTIVATE && status != cons.ORDER_EXECUTED {
			log.Println("ERROR security handler:: status value error =>", status)
			return false
		}
	default:
		log.Println("ERROR security handler:: NOT SUPPORT OPERATION!!!")
		return false
	}

	return true
}
*/
