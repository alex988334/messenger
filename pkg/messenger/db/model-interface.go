package db

/*****
import (
	"strconv"
)

const (
	MONO        = true
	POLY        = false
	UPDATE      = "UPDATE "
	SET         = " SET "
	SELECT      = "SELECT "
	FROM        = " FROM "
	WHERE       = " WHERE "
	ORDER       = " ORDER BY "
	ASC         = " ASC "
	DESC        = " DESC "
	BETWEEN     = " BETWEEN "
	IN          = " IN "
	LIMIT       = " LIMIT "
	MIN         = "min"
	MAX         = "max"
	AND         = " AND "
	OR          = " OR "
	NOT_IN      = " NOT IN "
	DELETE      = "DELETE"
	INSERT_INTO = "INSERT INTO "
	VALUES      = " VALUES "
)

type Model interface {
	//	SetLoadingAttr(attr []string) Model
	//	LoadAttr(request *map[string]interface{}) Model
	SetAttr(attr string, value interface{}) Model
	//	GetAttr(key string) (interface{}, bool)
	AddSelect(fields []string) Model
	AddOrder(fields []string) Model
	SetLimit(limit int) Model
	AndWhere(operat string, field string, params []interface{}) Model
	OrWhere(operat string, field string, params []interface{}) Model
	FindOne() (Model, bool)
	FindAll() ([]Model, bool)
	FindBySQL(query string, fields []string, params []interface{}) []map[string]interface{}
	Save() bool
	Delete() bool
	Insert() bool

	// Validate() bool
	// SetRules(rules *map[string]map[string]interface{})
}

type where struct {
	fields string
	params []interface{}
}

type models struct {
	attr        map[string]interface{}
	loadingAttr []string
	tableName   string
	fieldsKeys  []string
	rules       map[string]map[string]interface{}
	primaryKeys []string
	sqlSelect   []string
	sqlWhere    where
	sqlOrder    string
	sqlLimit    int
	DB
}

func CreateModel(table string, fieldKeys []string, primKeys []string, connection *DB) Model {
	var m Model = &models{
		attr:        make(map[string]interface{}),
		loadingAttr: []string{},
		tableName:   table,
		fieldsKeys:  fieldKeys,
		primaryKeys: primKeys,
		DB:          *connection,
		sqlSelect:   []string{},
		sqlLimit:    0,
		sqlWhere: where{
			fields: "",
			params: []interface{}{},
		},
		sqlOrder: "",
	}
	return m
}

func (m *models) SetRules(rules *map[string]map[string]interface{}) {

}

func (m *models) SetAttr(attr string, value interface{}) Model {
	if securityAttr(m.fieldsKeys, attr) {
		m.attr[attr] = value
	}
	return m
}

func (m *models) GetAttr(key string) (interface{}, bool) {
	if val, ok := m.attr[key]; ok {
		return val, true
	}
	return "", false
}

func (m *models) SetLimit(limit int) Model {
	m.sqlLimit = limit
	return m
}

func (m *models) AddSelect(fields []string) Model {
	m.sqlSelect = fields
	return m
}

func (m *models) AddOrder(fields []string) Model {
	switch fields[1] {
	case ASC, DESC:
	default:
		panic("ERROR 774443447: Not support value for ORDER BY, value = " + fields[1])
	}
	if m.sqlOrder == "" {
		m.sqlOrder += (fields[0] + fields[1])
	} else {
		m.sqlOrder += (", " + fields[0] + fields[1])
	}

	return m
}

func generateWhere(andOr string, operat string, field string, params []interface{}) (string, []interface{}) {
	wh := ""
	par := []interface{}{}
	switch operat {
	case "=", "<", ">", "<=", ">=", "<>":
		wh = wh + andOr + field + operat + "?"
		par = append(par, params[0])
	case BETWEEN:
		wh = wh + andOr + field + BETWEEN + "?" + AND + "?"
		par = append(par, params[0], params[1])
	case IN, NOT_IN:
		wh = wh + andOr + field + operat + "("
		for i := 0; i < len(params); i++ {
			wh += "?, "
			par = append(par, params[i])
		}
		wh = wh[0 : len(wh)-2]
		wh += ")"
	}
	return wh, par
}

func appendWhere(unit *where, andOr string, newWhere string, params []interface{}) {

	if (*unit).fields == "" {
		(*unit).fields += newWhere[len(andOr):len(newWhere)]
	} else {
		(*unit).fields += newWhere
	}
	for _, v := range params {
		(*unit).params = append((*unit).params, v)
	}
}

func (m *models) AndWhere(operat string, field string, params []interface{}) Model {

	wh, par := generateWhere(AND, operat, field, params)
	appendWhere(&m.sqlWhere, OR, wh, par)
	return m
}

func (m *models) OrWhere(operat string, field string, params []interface{}) Model {

	wh, par := generateWhere(OR, operat, field, params)
	appendWhere(&m.sqlWhere, OR, wh, par)
	return m
}

func (m *models) createSelect() string {
	if len(m.sqlSelect) == 0 {
		m.sqlSelect = m.fieldsKeys
	}
	sel := ""
	for _, v := range m.sqlSelect {
		sel += (v + ", ")
	}
	sel = sel[0 : len(sel)-2]

	return SELECT + sel
}

func (m *models) createFrom() string {
	return FROM + m.tableName
}

func (m *models) createWhere() string {
	str := WHERE
	if len(m.attr) > 0 {
		m.sqlWhere.params = []interface{}{}
		for k, v := range m.attr {
			str += (k + "=?" + AND)
			m.sqlWhere.params = append(m.sqlWhere.params, v)
		}
		str = str[0 : len(str)-len(AND)]
	} else if m.sqlWhere.fields != "" {
		str += m.sqlWhere.fields
	} else {
		panic("ERROR 886356228: in Model not set WHERE")
	}

	return str
}

func (m *models) createOrder() string {
	if m.sqlOrder != "" {
		return ORDER + m.sqlOrder
	}
	return ""
}

func (m *models) createLimit() string {
	if m.sqlLimit > 0 {
		return LIMIT + strconv.Itoa(m.sqlLimit)
	}
	return ""
}

func securityAttr(fields []string, key string) bool {
	for _, v := range fields {
		if v == key {
			return true
		}
	}
	return false
}

/*
Security model attributes. If model has no valid attributing name, it remove.
*/

/****
func (m *models) securityModel() {
	rAttr := []string{}
	for k := range m.attr {
		if !securityAttr(m.fieldsKeys, k) {
			rAttr = append(rAttr, k)
		}
	}
	for _, v := range rAttr {
		delete(m.attr, v)
	}
}

func (m *models) generateSelectQuery() string {

	return m.createSelect() + m.createFrom() + m.createWhere() +
		m.createOrder() + m.createLimit()
}

/*
if attr length > 0 - find for attr without sql params,
else find sql params
*/

/*****
func (m *models) FindOne() (Model, bool) {
	m.securityModel()
	m.SetLimit(1)
	query := m.generateSelectQuery()

	res := *(m.DB.SelectSQL(query, m.sqlSelect, m.sqlWhere.params))
	if len(res) > 0 {
		m.attr = res[0]
		m.resetSql()
		return m, true
	} else {
		return nil, false
	}
}

func (m *models) resetSql() {
	m.sqlLimit = 0
	m.sqlOrder = ""
	m.sqlWhere = where{fields: "", params: []interface{}{}}
	m.sqlSelect = nil
}

/*
if attr length > 0 - find for attr without sql params,
else find sql params
*/

/*****
func (m *models) FindAll() ([]Model, bool) {
	m.securityModel()
	modelss := []Model{}
	ok := false
	query := m.generateSelectQuery()
	res := *(m.DB.SelectSQL(query, m.sqlSelect, m.sqlWhere.params))

	for _, v := range res {
		ok = true
		var mod Model = CreateModel(m.tableName, m.fieldsKeys, m.primaryKeys, &m.DB)
		for k, val := range v {
			mod.SetAttr(k, val)
		}
		modelss = append(modelss, mod)
	}
	m.resetSql()
	return modelss, ok
}

/*
Find all attr of models in map
*/

/*****
func (m *models) FindBySQL(query string, fields []string, params []interface{}) []map[string]interface{} {

	res := *(m.DB.SelectSQL(query, fields, params))
	return res
}

/*func (m models) LoadAttributes(request *map[string]interface{}) {
}
*/
/*func (m models) SetLoadingAttr(attr []string) {
	m.attr = attr
}*/

/****
func (m *models) generateUpdateQuery() string {
	str, param := m.createSet()
	query := m.createUpdate() + str + m.createExecWhere()
	m.sqlWhere.params = append(param, m.sqlWhere.params...)
	return query
}

/*
func (m *models) generateDeleteQuery() string {

}
*/
/*****
func (m *models) createUpdate() string {
	return UPDATE + m.tableName
}

func (m *models) createExecWhere() string {
	str := WHERE
	par := []interface{}{}
	for _, v := range m.primaryKeys {
		str = str + v + "=?" + AND
		par = append(par, m.attr[v])
	}
	m.sqlWhere.params = par
	return str[0 : len(str)-len(AND)]
}

func (m *models) createSet() (string, []interface{}) {
	fields := ""
	params := []interface{}{}
	for k, v := range m.attr {
		flag := true
		for i := 0; i < len(m.primaryKeys); i++ {
			if m.primaryKeys[i] == k {
				flag = false
				break
			}
		}
		if flag {
			params = append(params, v)
			fields += (k + "=?, ")
		}
	}
	return SET + fields[0:len(fields)-2], params
}

func (m *models) securityPrimKeys() {
	for _, val := range m.primaryKeys {
		if _, ok := m.attr[val]; !ok {
			panic("ERROR 766253547: in model not setted attribute for primary keys")
		}
	}
}

func (m *models) securitySave() {

	m.securityPrimKeys()
	if m.sqlWhere.fields != "" || len(m.sqlWhere.params) > 0 {
		panic("ERROR 0895364583: in model setted WHERE")
	}
	if len(m.attr) <= len(m.primaryKeys) {
		panic("ERROR 253437675: not setted attributed for model")
	}
}

func (m *models) Save() bool {
	m.securitySave()
	query := m.generateUpdateQuery()
	res := *(m.DB.ExecSQL(query, [][]interface{}{m.sqlWhere.params}))
	if res[0].Total > 0 {
		return true
	}
	return false
}

func (m *models) createInsert() string {
	return INSERT_INTO + m.tableName
}

func (m *models) createInsertFields() string {
	fields := "("
	values := "("
	params := []interface{}{}
	for k, v := range m.attr {
		fields += k + ", "
		params = append(params, v)
		values += "?, "
	}
	m.sqlWhere.params = params
	return fields[0:len(fields)-2] + ")" + VALUES + values[0:len(values)-2] + ")"
}

func (m *models) securityInsert() {
	m.securityModel()
	if len(m.attr) == 0 {
		panic("ERROR 6734576868: not setted attributed for model")
	}
}

func (m *models) Insert() bool {
	m.securityInsert()
	query := m.createInsert() + m.createInsertFields()
	res := *(m.DB.ExecSQL(query, [][]interface{}{m.sqlWhere.params}))
	if res[0].Total > 0 {
		return true
	}
	return false
}

func (m *models) createDelete() string {
	return DELETE + FROM + m.tableName
}

func (m *models) Delete() bool {
	m.securityPrimKeys()
	query := m.createDelete() + m.createExecWhere()
	res := *(m.DB.ExecSQL(query, [][]interface{}{m.sqlWhere.params}))
	if res[0].Total > 0 {
		return true
	}
	return false
}
****/
