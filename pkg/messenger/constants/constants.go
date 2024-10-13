package constants

const (
	STATUS_ACCEPT = 1
	STATUS_ERROR  = 0

	ERROR_USER_NAME     = 1001
	ERROR_WRITE_BASE    = 1002
	ERROR_SEND_MESSAGE  = 1003
	ERROR_SEND_PARAMETR = 1004

	OP_AUTORIZATE             = 100
	OP_STATUS_MESSAGE         = 101
	OP_NEW_MESSAGE            = 103
	OP_LIST_USERS             = 105
	OP_CREATE_NEW_CHAT        = 106
	OP_WRITEN                 = 107
	OP_SYSTEM                 = 108
	OP_SEARCH_USER            = 109
	OP_GET_CHATS              = 110
	OP_LIST_PREVIOUS_MESSAGES = 111
	OP_LIST_NEXT_MESSAGES     = 112
	OP_EXIT_CHAT              = 113
	OP_REMOVE_USER            = 114
	OP_ADD_USER               = 115
	OP_REMOVE_CHAT            = 116
	OP_BLOCK_USERS            = 117
	OP_UNLOOCK_USERS          = 118
	OP_BLACK_LIST_USERS       = 119
	OP_GET_FILE               = 120
	OP_MY_DATA                = 122
	OP_REGISTRATION           = 123

	KEY_ACTION     = "action"
	KEY_MESSAGE    = "message"
	KEY_CHAT       = "id_chat"
	KEY_ID         = "id"
	KEY_STATUS     = "status"
	KEY_CHAT_NAME  = "chat_name"
	KEY_USERS      = "users"
	KEY_SEARCH     = "search"
	KEY_BLACK_LIST = "blackList"
	KEY_WRITE      = "write"
	KEY_NAME       = "name"
	KEY_AUTHOR     = "author"

	ORDER_AKTIVATE   = 150
	ORDER_DIAKTIVATE = 151
	ORDER_EXECUTED   = 152

	NULL_MESSAGES = 1008
	MESSAGE_ALL   = 1010

	MESSAGE_CREATED    = "created"
	MESSAGE_SENDED     = "sended"
	MESSAGE_DELIVERED  = "delivered"
	MESSAGE_READED     = "readed"
	MESSAGE_BLACK_LIST = "black_list"

	CHAT_ACTIVE      = "active"
	CHAT_DIACTIVATED = "diactivated"
	CHAT_DELETED     = "deleted"

	ROLE_ADMIN        = "admin"
	ROLE_HEAD_MANAGER = "head_manager"
	ROLE_MANAGER      = "manager"
	ROLE_MASTER       = "master"
	ROLE_KLIENT       = "klient"

	USER_LOCK   = 9
	USER_UNLOCK = 10

	LIST_LIMIT_MESSAGES = 25

	CLIENT_CONNECT           = 15
	CLIENT_DISCONNECT        = 16
	ALL_CLIENT_DISCONNECT    = 17
	CLIENT_CONNECT_ALL_CHATS = 18

	PASSHASH_COST = 5
)

func GetDescriptionOperation(operation int) string {

	switch operation {
	case OP_AUTORIZATE:
		return "OP_AUTORIZATE"
	case OP_STATUS_MESSAGE:
		return "OP_STATUS_MESSAGE"
	case OP_NEW_MESSAGE:
		return "OP_NEW_MESSAGE"
	case OP_LIST_USERS:
		return "OP_LIST_USERS"
	case OP_CREATE_NEW_CHAT:
		return "OP_CREATE_NEW_CHAT"
	case OP_WRITEN:
		return "OP_WRITEN"
	case OP_SYSTEM:
		return "OP_SYSTEM"
	case OP_SEARCH_USER:
		return "OP_SEARCH_USER"
	case OP_GET_CHATS:
		return "OP_GET_CHATS"
	case OP_LIST_PREVIOUS_MESSAGES:
		return "OP_LIST_PREVIOUS_MESSAGES"
	case OP_LIST_NEXT_MESSAGES:
		return "OP_LIST_NEXT_MESSAGES"
	case OP_EXIT_CHAT:
		return "OP_EXIT_CHAT"
	case OP_REMOVE_USER:
		return "OP_REMOVE_USER"
	case OP_ADD_USER:
		return "OP_ADD_USER"
	case OP_REMOVE_CHAT:
		return "OP_REMOVE_CHAT"
	case OP_BLOCK_USERS:
		return "OP_BLOCK_USERS"
	case OP_UNLOOCK_USERS:
		return "OP_UNLOOCK_USERS"
	case OP_BLACK_LIST_USERS:
		return "OP_BLACK_LIST_USERS"
	case OP_GET_FILE:
		return "OP_GET_FILE"
	case OP_MY_DATA:
		return "OP_MY_DATA"
	case OP_REGISTRATION:
		return "OP_REGISTRATION"
	default:
		return "Not found operation!!!"
	}
}
