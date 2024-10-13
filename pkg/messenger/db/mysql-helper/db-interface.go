package mysql_helper

/*
DBInterface interface
*/
type DBInterface interface {
	Transact
	SimpleDB
}

type Transact interface {
	IsRunTransact() bool
	StartTransact()
	CloseTransact()
	SetRollBack(flag bool)
}

type SimpleDB interface {
	SelectSQL(query string, params []interface{}) *[]map[string]interface{}
	ExecSQL(query string, params [][]interface{}) []ExecSQLRezult
	CloseDB()
}
