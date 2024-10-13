package data

const (
	MODEL_USER           = "User"
	MODEL_CHAT_USER      = "ChatUser"
	MODEL_CHAT           = "Chat"
	MODEL_MESSAGE        = "Message"
	MODEL_STATUS_MESSAGE = "MessageStatus"
	MODEL_BLACK_LIST     = "BlackList"
	MODEL_USER_PHONE     = "UserPhone"
	MODEL_STATUS         = "Status"
	MODEL_SYSTEM         = "System"

	ACTION_FIND   = 200
	ACTION_UPDATE = 201
	ACTION_SAVE   = 202
	ACTION_DELETE = 203

	LINK_WEIGHT_EQUILIBRIUM  = 0
	LINK_WEIGHT_PARENT_MORE  = 1
	LINK_WEIGHT_CURRENT_MORE = 2

	CHAT_STATUS_DELETED = "deleted"
	CHAT_STATUS_ACTIVE  = "active"

	EQUAL      = "="
	MORE       = ">"
	LESS       = "<"
	MORE_EQUAL = MORE + EQUAL
	LESS_EQUAL = LESS + EQUAL
	NOT_EQUAL  = LESS + MORE
	IN         = " IN "
	NOT_IN     = " NOT IN "

	DIRECTION_ASC  = " ASC"
	DIRECTION_DESC = " DESC"

	EXISTS     = "EXISTS"
	NOT_EXISTS = "NOT EXISTS"
)
