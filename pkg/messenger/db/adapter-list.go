package db

var associateTables = map[string]string{
	"User":          "user",
	"ChatUser":      "chat_user",
	"BlackList":     "chat_black_list",
	"Chat":          "chat",
	"Message":       "chat_message",
	"MessageStatus": "chat_message_status",
	"UserPhone":     "user_phone",
}

/*
	var associateModels = map[string]string{
		"user":                "User",
		"chat_user":           "ChatUser",
		"chat_black_list":     "BlackList",
		"chat":                "Chat",
		"chat_message":        "Message",
		"chat_message_status": "MessageStatus",
		"user_phone":          "UserPhone",
	}
*/
var associateFields = map[string]map[string]string{
	"User": {
		"Id":             "id",
		"Login":          "username",
		"Alias":          "alias",
		"AuthKey":        "auth_key",
		"PassHash":       "password_hash",
		"PassResetToken": "password_reset_token",
		"Email":          "email",
		"Status":         "status",
		"CreateAt":       "created_at",
		"UpdateAt":       "updated_at",
		"Avatar":         "avatar",
	},
	"ChatUser": {
		"Chat":        "id_chat",
		"User":        "id_user",
		"SessionHash": "client_hash",
	},
	"BlackList": {
		"User":        "blocking",
		"BlockedUser": "locked",
		"Date":        "date",
		"Time":        "time",
	},
	"Chat": {
		"Id":       "id",
		"Author":   "author",
		"Name":     "alias",
		"CreateAt": "create_at",
		"Status":   "status",
	},
	"Message": {
		"Id":             "id",
		"ChatId":         "id_chat",
		"Author":         "id_user",
		"ParrentMessage": "parent_id",
		"Message":        "message",
		"FileUrl":        "file",
		"Date":           "date",
		"Time":           "time",
	},
	"MessageStatus": {
		"MessageId": "id_message",
		"UserId":    "id_user",
		"Status":    "status_message",
		"Date":      "date",
		"Time":      "time",
	},
	"UserPhone": {
		"UserId": "user_id",
		"Phone":  "phone",
	},
}
