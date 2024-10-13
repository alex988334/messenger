package db

import (
	"context"
	"database/sql"
	"log"

	mysql "github.com/alex988334/messenger/pkg/messenger/db/mysql-helper"
	_ "github.com/go-sql-driver/mysql"
)

type runExec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type DB struct {
	conn          *sql.DB
	tx            *sql.Tx
	rollback      bool
	context       context.Context
	contextCancel context.CancelFunc
}

func NewDB() *DB {
	c, err := sql.Open("mysql", UserDB+":"+PasswordDB+"@/"+Database)
	if err != nil {
		log.Fatal("Failed to create new database connection!")
		return nil
	}

	return &DB{conn: c, tx: nil, rollback: false}
}

func (db *DB) SetRollBack(flag bool) {
	db.rollback = flag
}

// считывает из результата запроса значения полей и записывает их в массив карт (поле - значение)
func scanResult(result *[]map[string]interface{} /*keys []string,*/, rows *sql.Rows) {

	columns, e := rows.Columns()
	if e != nil {
		panic("ERROR SQL driver, rezult select query not support columns names!!!")
	}

	for rows.Next() { //	пробегается по всем возвращенным строкам
		mass := make([]interface{}, 0, len(columns)) //	создаем пустой массив интерфейсов, который заполним значениями

		for i := 0; i < len(columns); i++ { //	наполняем его ссылками на пустые интерфейсы, в колличестве всех полей
			var b interface{}
			mass = append(mass, &b)
		}

		err := rows.Scan(mass...) //	автоматом читаем и записываем переменные из результата запроса
		if err != nil {
			panic("ERROR SQL! Rezult do not scan!!!," + e.Error())
		}
		line := make(map[string]interface{}, 0)
		for i, key := range columns {

			line[key] = *mass[i].(*interface{})

			switch line[key].(type) {
			case []byte:
				line[key] = string(line[key].([]byte))
				//	case nil:
				//		line[val] = ""
			}
		}
		*result = append(*result, line)
	}
	//log.Println("result", result)
}

func (db *DB) SelectSQL(query string, params []interface{}) *[]map[string]interface{} {
	var (
		rows *sql.Rows
		err  error
	)

	if params != nil && len(params) > 0 {
		rows, err = db.conn.Query(query, params...)
	} else {
		rows, err = db.conn.Query(query)
	}
	defer rows.Close()

	resp := make([]map[string]interface{}, 0)
	if err != nil {
		log.Println("error sql => 58054521201", err.Error())
		return &resp
	}

	scanResult(&resp /*keys,*/, rows)
	return &resp
}

func decodeResult(res *sql.Result, err *error) *mysql.ExecSQLRezult {
	if *err == nil {
		id, _ := (*res).LastInsertId()
		total, _ := (*res).RowsAffected()
		return &mysql.ExecSQLRezult{id, total, *err}
	} else {
		log.Println("sql exec 45354352432 =>", (*err).Error())
		return &mysql.ExecSQLRezult{0, 0, *err}
	}
}

func (db *DB) ExecSQL(query string, params [][]interface{}) []mysql.ExecSQLRezult {
	var (
		res sql.Result
		err error
		run runExec
	)
	mass := []mysql.ExecSQLRezult{}
	if db.tx != nil {
		run = db.tx
	} else {
		run = db.conn
	}

	if len(params) > 0 {
		for _, v := range params {
			//	fmt.Println("res, err = run.Exec(query, v...) =>", query, "\n", v)
			res, err = run.Exec(query, v...)
			//	fmt.Println("res, =>", res, ", err=>", err)

			if db.tx != nil && err != nil {
				db.rollback = true
				mass = append(mass, *decodeResult(&res, &err))
				break
			}
			mass = append(mass, *decodeResult(&res, &err))
		}
	} else {
		res, err = run.Exec(query)
		if db.tx != nil && err != nil {
			db.rollback = true
		}
		mass = append(mass, *decodeResult(&res, &err))
	}
	return mass
}

func (db *DB) IsRunTransact() bool {

	return db.tx != nil
}

func (db *DB) StartTransact() {

	db.context, db.contextCancel = context.WithCancel(context.TODO())

	tx, err := db.conn.BeginTx(db.context, nil /*&sql.TxOptions{Isolation: sql.LevelSerializable}*/)
	if err != nil {
		log.Fatal("StartTransact()", err)
	}
	db.tx = tx
}

func (db *DB) CloseTransact() {

	if db.tx == nil {
		return
	}

	var err error

	if db.rollback {
		err = db.tx.Rollback()
	} else {
		err = db.tx.Commit()
	}

	if err != nil {
		log.Fatal(err)
	}
	db.tx = nil

	db.contextCancel()
	db.contextCancel = nil
}

func (db *DB) CloseDB() {

	db.CloseTransact()

	err := db.conn.Close()
	if err != nil {
		log.Println("db close 7554757 =>", err)
	}
}
