package db

import (
	"testing"
	"time"

	mysql "github.com/alex988334/messenger/pkg/messenger/db/mysql-helper"
)

func TestConnectionDB(t *testing.T) {

	t.Log("Test create DB connection")

	db := NewDB()
	if db == nil {
		t.Fatal("Failed to create new database connection!")
	}
	db.CloseDB()
}

func TestTransactDB(t *testing.T) {

	t.Log("Test create transact")

	db := NewDB()
	db.StartTransact()
	if db.tx == nil {
		t.Fatal("Failed to create transact!")
	}
	db.rollback = true
	db.CloseTransact()
	db.CloseDB()
}

func TestInsertUserDB(t *testing.T) {

	t.Log("Test Insert into user")

	db := NewDB()
	db.ExecSQL("DELETE FROM `user` WHERE username=\"\" AND alias=\"Ivanych\"", [][]interface{}{})

	db.StartTransact()

	now := time.Now().Unix()
	var f []mysql.ExecSQLRezult = db.ExecSQL("INSERT INTO `user`(`username`, `alias`, `auth_key`, "+
		"`password_hash`, `created_at`, `updated_at`) "+
		"VALUES (?, ?, ?, ?, ?, ?)",
		[][]interface{}{
			{"Ivan", "Ivanych", "fsfgdjfsdndssdff", "rfdgdsagfds", now, now},
		})

	if err := f[0].Err; err != nil {
		t.Fatal("Failed insert into user!", err.Error())
	}

	db.rollback = true
	db.CloseTransact()
	db.CloseDB()
}

func TestInsertUserPhoneDB(t *testing.T) {

	t.Log("Test Insert into User Phone")

	db := NewDB()
	db.ExecSQL("DELETE FROM `user_phone` WHERE user_id=?", [][]interface{}{{36}})

	db.StartTransact()

	var f []mysql.ExecSQLRezult = db.ExecSQL("INSERT INTO `user_phone`(`user_id`, `phone`) VALUES (?, ?)",
		[][]interface{}{
			{36, "79633640116"},
		})

	if err := f[0].Err; err != nil {
		t.Fatal("Failed insert into User Phone!", err.Error())
	}

	db.rollback = true
	db.CloseTransact()
	db.CloseDB()
}

func TestInsertChatsDB(t *testing.T) {

	t.Log("Test Insert into chats")

	db := NewDB()
	db.ExecSQL("DELETE FROM `chat` WHERE author=9 AND alias=\"Wer\"", [][]interface{}{})

	db.StartTransact()

	var f []mysql.ExecSQLRezult = db.ExecSQL("INSERT INTO `chat`"+
		"(`author`, `alias`, `create_at`, `status`)"+
		" VALUES (?, ?, ?, ?)",
		[][]interface{}{
			{9, "Wer", "2016-02-16", "active"},
		},
	)

	if err := f[0].Err; err != nil {
		t.Fatal("Failed insert into chats!", err.Error())
	}

	db.rollback = true
	db.CloseTransact()
	db.CloseDB()
}

func TestInsertChatsUserDB(t *testing.T) {

	t.Log("Test Insert into chat user")

	db := NewDB()
	db.ExecSQL("DELETE FROM `chat_user` WHERE id_chat=1 AND id_user=2", [][]interface{}{})

	db.StartTransact()
	var f []mysql.ExecSQLRezult = db.ExecSQL("INSERT INTO `chat_user`"+
		"(`id_chat`, `id_user`) VALUES (?, ?)",
		[][]interface{}{
			{1, 2},
		},
	)

	if err := f[0].Err; err != nil {
		t.Fatal("Failed insert into chat user!", err.Error())
	}

	db.rollback = true
	db.CloseTransact()
	db.CloseDB()
}

func TestInsertChatsBlackListDB(t *testing.T) {

	t.Log("Test Insert into chat black list")

	db := NewDB()
	db.ExecSQL("DELETE FROM `chat_black_list` WHERE blocking=1 AND locked=2", [][]interface{}{})

	db.StartTransact()
	var f []mysql.ExecSQLRezult = db.ExecSQL("INSERT INTO `chat_black_list`"+
		"(`blocking`, `locked`, `date`, `time`) VALUES (?, ?, ?, ?)",
		[][]interface{}{
			{1, 2, "2016-02-19", "12:59:59"},
		},
	)

	if err := f[0].Err; err != nil {
		t.Fatal("Failed insert into chat black list!", err.Error())
	}

	db.rollback = true
	db.CloseTransact()
	db.CloseDB()
}

func TestInsertChatMessageDB(t *testing.T) {

	t.Log("Test Insert into chat message")

	db := NewDB()
	db.ExecSQL("DELETE FROM `chat_message` WHERE id_chat=1 AND id_user=2 AND message=\"WER\"", [][]interface{}{})

	db.StartTransact()
	var f []mysql.ExecSQLRezult = db.ExecSQL("INSERT INTO `chat_message`"+
		"(`id_chat`, `id_user`, `message`, `date`, `time`) "+
		"VALUES (?, ?, ?, ?, ?)",
		[][]interface{}{
			{1, 2, "WER", "2018-05-23", "15:45:23"},
		},
	)

	if err := f[0].Err; err != nil {
		t.Fatal("Failed insert into chat message!", err.Error())
	}

	db.CloseTransact()
	db.CloseDB()
}

func TestInsertChatMessageStatusDB(t *testing.T) {

	t.Log("Test Insert into chat message status")

	db := NewDB()
	db.ExecSQL("DELETE FROM `chat_message_status` WHERE id_message=1 AND id_user=2", [][]interface{}{})

	db.StartTransact()
	var f []mysql.ExecSQLRezult = db.ExecSQL("INSERT INTO `chat_message_status`"+
		"(`id_message`, `id_user`, `status_message`, `date`, `time`) "+
		"VALUES (?, ?, ?, ?, ?)",
		[][]interface{}{
			{1, 2, "created", "2016-05-30", "18:53:45"},
		},
	)

	if err := f[0].Err; err != nil {
		t.Fatal("Failed insert into chat message status!", err.Error())
	}

	db.CloseTransact()
	db.CloseDB()
}
