package mysql_helper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alex988334/messenger/pkg/messenger/functions"
)

const (
	NOT         = " NOT"
	UPDATE      = "UPDATE "
	SET         = " SET "
	SELECT      = "SELECT "
	AS          = " AS "
	FROM        = " FROM "
	WHERE       = " WHERE "
	EXISTS      = "EXISTS"
	NOT_EXISTS  = "NOT EXISTS"
	ORDER       = " ORDER BY "
	ASC         = " ASC"
	DESC        = " DESC"
	BETWEEN     = " BETWEEN "
	IN          = " IN "
	NOT_IN      = " NOT IN "
	LIMIT       = " LIMIT "
	MIN         = "min"
	MAX         = "max"
	AND         = " AND "
	OR          = " OR "
	DELETE      = "DELETE"
	DUAL        = "DUAL"
	INSERT_INTO = "INSERT INTO "
	VALUES      = " VALUES "
	JOIN        = " JOIN "
	LEFT        = " LEFT"
	RIGHT       = " RIGHT"
	ON          = " ON "
	GROUP_BY    = " GROUP BY "
	ORDER_BY    = " ORDER BY "
	NULL        = " NULL"
	IS          = " IS"
	EQUAL       = "="
	MORE        = ">"
	LESS        = "<"
	MORE_EQUAL  = MORE + EQUAL
	LESS_EQUAL  = LESS + EQUAL
	NOT_EQUAL   = LESS + MORE
	PARAMETR    = "?"
	DELIMETR    = ", "
)

type ExecSQLRezult struct {
	LastId int64
	Total  int64
	Err    error
}

func (r *ExecSQLRezult) GetLastId() int64 {
	return r.LastId
}
func (r *ExecSQLRezult) GetTotalRow() int64 {
	return r.Total
}
func (r *ExecSQLRezult) GetError() error {
	return r.Err
}

type SQLModel struct {
	db           DBInterface
	query        string
	queryParams  []interface{}
	selectRezult []map[string]interface{}

	// Fields in select part
	selectFields  []string
	selectAliases []string

	// Tables in from part
	sqlTables        []string
	sqlTablesAliases []string
	//	fromTables  []string
	//	fromAliases []string

	// Params in join part
	joinTables     []string
	joinAliases    []string
	joinDirection  []string
	joinLinkFields []string
	joinLinkAlias  []string

	// Dependencies in where part
	whereDependenciesFields         []string
	whereDependenciesAliases        []string
	whereDependenciesOperator       []string
	whereDependenciesLogicOperators []string

	// Conditions in where part
	whereFields          []string
	whereAliases         []string
	whereLengthArrFields map[string]int
	whereFieldsOperator  []string
	whereLogicOperators  []string

	// fields in group by part
	groupFields  []string
	groupAliases []string

	// Fields in order part
	orderFileds    []string
	orderAliases   []string
	orderDirection []string

	updSetFields  []string
	updSetAliases []string

	insertFields []string

	deleteSql bool

	// Set limit rows in rezult of query
	limit int

	subFields         []string
	subAliases        []string
	subOperators      []string
	subLogicOperators []string
	subModels         []*SQLModel
}

func NewSQLModel(dbI DBInterface) *SQLModel {
	return &SQLModel{
		db:                   dbI,
		queryParams:          make([]interface{}, 0, 10),
		whereLengthArrFields: make(map[string]int, 5),
	}
}

/*
**

	part of loading data in model

**
*/
func (m *SQLModel) GetDbConnection() DBInterface {
	return m.db
}

func (m *SQLModel) Commit() {
	//m.db.CloseTransact()
}

func (m *SQLModel) CloseDBConnection() {
	m.db.CloseDB()
}

func (m *SQLModel) GetSubParams() []interface{} {

	subParams := make([]interface{}, 0, len(m.subModels)*5)

	for _, obj := range m.subModels {
		subParams = append(subParams, obj.GetQueryParams()...)
	}

	return subParams
}

func (m *SQLModel) GetSubQuries() []string {

	subQuires := make([]string, 0, len(m.subModels))

	for _, obj := range m.subModels {
		subQuires = append(subQuires, obj.GetQuery())
	}

	return subQuires
}

func (m *SQLModel) SetLimitRows(limitRow int) {
	m.limit = limitRow
}

func (m *SQLModel) GetTablesName() []string {
	return m.sqlTables
}
func (m *SQLModel) GetTablesAlias() []string {
	return m.sqlTablesAliases
}
func (m *SQLModel) GetJoinTables() []string {
	return m.joinTables
}
func (m *SQLModel) GetJoinTableAlias() []string {
	return m.joinAliases
}

func (m *SQLModel) findAlias(table string) string {

	if ind := functions.FindInArray(m.sqlTables, table); ind > -1 {
		if ind < len(m.sqlTablesAliases) {
			return m.sqlTablesAliases[ind]
		}
	}
	if ind := functions.FindInArray(m.joinTables, table); ind > -1 {
		if ind < len(m.joinAliases) {
			return m.joinAliases[ind]
		}
	}

	return ""
}

func (m *SQLModel) generateTableAlias(tableName *string, lengthAlias int) string {

	alias := (*tableName)[0:lengthAlias]

	if functions.FindInArray(m.sqlTablesAliases, alias) > -1 ||
		functions.FindInArray(m.joinAliases, alias) > -1 {

		alias = m.generateTableAlias(tableName, lengthAlias+1)
	}

	return alias
}

func (m *SQLModel) AddSelectFields(fields []string, tableName string) string {

	alias := m.AddSqlTable(tableName)

	aliasFi := functions.GenerateArrayOfvalue(len(fields), alias)

	m.selectFields = append(m.selectFields, fields...)
	m.selectAliases = append(m.selectAliases, aliasFi...)

	return alias
}

func (m *SQLModel) AddSqlTable(table string) string {

	if ind := functions.FindInArray(m.sqlTables, table); ind > -1 {
		return m.sqlTablesAliases[ind]
	}
	if ind := functions.FindInArray(m.joinTables, table); ind > -1 {
		return m.joinAliases[ind]
	}

	alias := ""

	if len(m.insertFields) == 0 && !m.deleteSql {
		alias = m.generateTableAlias(&table, 1)
		m.sqlTablesAliases = append(m.sqlTablesAliases, alias)
	}

	m.sqlTables = append(m.sqlTables, table)

	return alias
}

func (m *SQLModel) AddLinkWhere(firstTable string, secondTable string,
	firstField string, secondField string, operator string, logicOperator string) {

	fInd := functions.FindInArray(m.sqlTables, firstTable)
	sInd := functions.FindInArray(m.sqlTables, secondTable)
	if fInd == -1 || sInd == -1 {
		panic("MySQLHelper ERROR: table " + firstTable + " or " + secondTable + " not found in FROM part!")
	}

	m.whereDependenciesAliases = append(m.whereDependenciesAliases, m.sqlTablesAliases[fInd], m.sqlTablesAliases[sInd])
	m.whereDependenciesFields = append(m.whereDependenciesFields, firstField, secondField)
	m.whereDependenciesOperator = append(m.whereDependenciesOperator, operator)
	m.whereDependenciesLogicOperators = append(m.whereDependenciesLogicOperators, logicOperator)
}

func (m *SQLModel) AddWhereArrField(table string, field string, arrayOperator string, lengthArr int, logicOperator string) {

	alias := m.findAlias(table)
	if alias != "" {
		alias += "."
	}
	m.whereLengthArrFields[alias+field] = lengthArr

	m.AddWhere(table, field, arrayOperator, logicOperator)
}

func (m *SQLModel) AddWhere(table string, field string, operator string, logicOperator string) {

	alias := m.findAlias(table)

	m.whereAliases = append(m.whereAliases, alias)
	m.whereFields = append(m.whereFields, field)
	m.whereFieldsOperator = append(m.whereFieldsOperator, operator)
	m.whereLogicOperators = append(m.whereLogicOperators, logicOperator)
}

func (m *SQLModel) AddJoin(direction string, baseTable string, joinTable string,
	baseField string, joinField string) {

	alias := m.findAlias(baseTable)
	joinAlias := m.generateTableAlias(&joinTable, 1)

	m.joinTables = append(m.joinTables, joinTable)
	m.joinAliases = append(m.joinAliases, joinAlias)
	m.joinDirection = append(m.joinDirection, direction)
	m.joinLinkFields = append(m.joinLinkFields, baseField, joinField)
	m.joinLinkAlias = append(m.joinLinkAlias, alias, joinAlias)
}

func (m *SQLModel) AddGroupBy(fields []string, table string) {

	alias := m.findAlias(table)
	if alias == "" {
		panic("ERROR MySQLHelper! Not found table " + table + " in FROM or JOIN parts!!!")
	}

	aliases := functions.GenerateArrayOfvalue(len(fields), alias)
	m.groupFields = append(m.groupFields, fields...)
	m.groupAliases = append(m.groupAliases, aliases...)
}

func (m *SQLModel) AddOrderBy(fields []string, table string, directions []string) {

	alias := m.findAlias(table)
	if alias == "" {
		panic("ERROR MySQLHelper! Not found table " + table + " in FROM or JOIN parts!!!")
	}

	aliases := functions.GenerateArrayOfvalue(len(fields), alias)

	m.orderFileds = append(m.orderFileds, fields...)
	m.orderAliases = append(m.orderAliases, aliases...)
	m.orderDirection = append(m.orderDirection, directions...)
}

func (m *SQLModel) AddSetFields(fields []string, table string) {

	alias := m.findAlias(table)
	if alias == "" {
		panic("ERROR MySQLHelper! Not found table " + table + " in Tables or JOIN parts!!!")
	}

	aliases := functions.GenerateArrayOfvalue(len(fields), alias)

	m.updSetAliases = append(m.updSetAliases, aliases...)
	m.updSetFields = append(m.updSetFields, fields...)
}

func (m *SQLModel) AddInsertFields(fields []string) {

	if len(fields) == 0 {
		panic("ERROR MySQLHelper! Insert fields count equal zero!!!")
	}

	m.insertFields = fields
}

func (m *SQLModel) InitDeleteQuery() {

	m.deleteSql = true
}

func (m *SQLModel) AddParamsQuery(param ...interface{}) {

	m.queryParams = append(m.queryParams, param...)
}

func (m *SQLModel) AddSubQuery(subModel *SQLModel, table string, whereField string,
	fieldOperator string, logicOperator string) {

	m.subAliases = append(m.subAliases, m.findAlias(table))
	m.subFields = append(m.subFields, whereField)
	m.subOperators = append(m.subOperators, fieldOperator)
	m.subLogicOperators = append(m.subLogicOperators, logicOperator)
	m.subModels = append(m.subModels, subModel)
}

func (m *SQLModel) ExecSQL(params [][]interface{}) []ExecSQLRezult {

	if m.query == "" {
		m.PrepareQuery()
	}

	p := params
	if len(p) == 0 {
		p = [][]interface{}{m.queryParams}
	}

	if !m.db.IsRunTransact() {
		m.db.StartTransact()
	}

	return m.db.ExecSQL(m.query, p)
}

func (m *SQLModel) GetQueryParams() []interface{} {
	return m.queryParams
}

func (m *SQLModel) GetQuery() string {
	return m.query
}

/*
**

	part SQL generation

**
*/
func (m *SQLModel) Select() *[]map[string]interface{} {

	if m.query == "" {
		m.PrepareQuery()
	}
	//	fmt.Println("m.GetQuery()=>", m.GetQuery())
	//	fmt.Println("m.GetQueryParams()=>", m.GetQueryParams())

	m.selectRezult = *m.db.SelectSQL(m.query, m.queryParams)

	return &m.selectRezult
}

func PrepareSelectQuery(m *SQLModel) {

	m.query = generateSelect(m)
	m.query += generateFrom(m)

	if len(m.joinTables) > 0 {
		m.query += generateJoin(m)
	}

	if len(m.whereDependenciesFields) > 0 {
		m.query += generateWhereDependencies(m)
	}

	if len(m.whereFields) > 0 {
		m.query += generateWhereConditions(m)
	}

	if len(m.subModels) > 0 {
		m.query += generateSubSelectQuries(m)
	}

	if len(m.groupFields) > 0 {
		m.query += generateGroupBy(m)
	}

	if len(m.orderFileds) > 0 {
		m.query += generateOrderBy(m)
	}

	if m.limit > 0 {
		m.query += generateLimit(m)
	}

}

func PrepareInsertQuery(m *SQLModel) {

	if len(m.sqlTables) == 0 || m.sqlTables[0] == "" {
		panic("ERROR MySQL! Table name is empty!!!")
	}

	m.query = INSERT_INTO + createListTables(m.sqlTables, m.sqlTablesAliases) + generateInsertFields(m)

	if len(m.subModels) > 0 {
		m.query += generateSubInsertQuries(m)
	} else {
		m.query += VALUES + "(" +
			functions.GenerateStrFromArr(
				functions.GenerateArrayOfvalue(len(m.insertFields), "?"),
				DELIMETR,
			) + ")"
	}
}

func PrepareDeleteQuery(m *SQLModel) {

	//fmt.Println("m.query (PrepareDeleteQuery)=>", m.query)
	if len(m.sqlTables) == 0 || m.sqlTables[0] == "" {
		panic("ERROR MySQL! Table name is empty!!!")
	}

	m.query = DELETE + FROM + m.sqlTables[0] + generateWhereConditions(m)

	if len(m.subModels) > 0 {
		m.query += generateSubSelectQuries(m)
	}
	//	fmt.Println("m.query (PrepareDeleteQuery)=>", m.query)
}

func generateExcludedObjects(m *SQLModel) (excludedTables []string, excludedAliases []string) {

	excludedTables = make([]string, 0, len(m.sqlTables))
	excludedAliases = make([]string, 0, len(m.sqlTables))

	if len(m.selectAliases) > 0 {

		for i := 0; i < len(m.selectAliases); i++ {

			if functions.FindInArray(excludedAliases, m.selectAliases[i]) == -1 {
				excludedAliases = append(excludedAliases, m.selectAliases[i])
			}
		}
		for i := 0; i < len(excludedTables); i++ {
			excludedTables[i] = m.sqlTables[functions.FindInArray(m.sqlTablesAliases, excludedAliases[i])]
		}
	}
	return excludedTables, excludedAliases
}

func PrepareUpdateQuery(m *SQLModel) {

	if len(m.sqlTables) == 0 {
		panic("ERROR MySQL! Table name is empty!!!")
	}

	m.query = UPDATE + createListTables(m.sqlTables, m.sqlTablesAliases) +
		SET + generateUpdateSet(m)

	if len(m.whereDependenciesFields) > 0 {
		m.query += generateWhereDependencies(m)
	}

	if len(m.whereFields) > 0 {
		m.query += generateWhereConditions(m)
	}

	if len(m.subModels) > 0 {
		m.query += generateSubSelectQuries(m)
	}
}
func (m *SQLModel) prepareSubModels() {

	for _, v := range m.subModels {
		v.prepareSubModels()
		v.prepareSQLModel()
	}
}

func (m *SQLModel) prepareSQLModel() {

	//fmt.Println("m.query (prepareSQLModel)=>", m.query)
	//	fmt.Println("*SQLModel (prepareSQLModel)=>", m)

	if len(m.selectFields) > 0 && len(m.updSetFields) == 0 &&
		len(m.insertFields) == 0 && !m.deleteSql {

		PrepareSelectQuery(m)
	} else if len(m.insertFields) > 0 {
		PrepareInsertQuery(m)
	} else if len(m.updSetFields) > 0 {
		PrepareUpdateQuery(m)
	} else if m.deleteSql {
		PrepareDeleteQuery(m)
	} else {
		panic("ERROR MySQL! SQL model is empty, data not loaded!!!")
	}
	//fmt.Println("m.query (prepareSQLModel)<=", m.query)
	//	fmt.Println("*SQLModel (prepareSQLModel)<=", m)
	//	fmt.Println(0)
}

func (m *SQLModel) PrepareQuery() {

	m.prepareSubModels()

	m.prepareSQLModel()

	//	fmt.Println("END")
}

func generateUpdateSet(model *SQLModel) string {

	rezult := ""

	for i := 0; i < len(model.updSetFields); i++ {

		if i < len(model.updSetAliases) && model.updSetAliases[i] != "" {
			rezult += model.updSetAliases[i] + "."
		}

		rezult += model.updSetFields[i] + EQUAL + PARAMETR + DELIMETR
	}

	return rezult[0 : len(rezult)-len(DELIMETR)]
}

func generateInsertFields(model *SQLModel) string {

	rezult := "("

	for i := 0; i < len(model.insertFields); i++ {

		rezult += model.insertFields[i] + DELIMETR
	}

	return rezult[0:len(rezult)-len(DELIMETR)] + ")"
}

func generateSelect(model *SQLModel) string {

	return SELECT + createSelect(model.selectFields, model.selectAliases)
}

func generateFrom(model *SQLModel) string {

	return FROM + createListTables(model.sqlTables, model.sqlTablesAliases)
}

func generateJoin(model *SQLModel) string {

	join := ""

	for i := 0; i < len(model.joinTables); i++ {

		join += createJoin(model.joinDirection[i], model.joinLinkAlias[i*2], model.joinLinkFields[i*2],
			model.joinTables[i], model.joinLinkAlias[i*2+1], model.joinLinkFields[i*2+1])
	}

	return join
}

func generateWhereDependencies(m *SQLModel) string {

	depend := ""

	if !strings.Contains(m.query, WHERE) {
		depend = WHERE
	}

	for i := 0; i < len(m.whereDependenciesOperator); i++ {

		operat := m.whereDependenciesOperator[i]
		if operat == "" {
			operat = EQUAL
		}
		depend += createWhereDependencies(m.whereDependenciesAliases[i*2], m.whereDependenciesFields[i*2],
			m.whereDependenciesAliases[i*2+1], m.whereDependenciesFields[i*2+1], operat)

		if i < len(m.whereDependenciesOperator)-1 {

			logicOperat := AND
			if i < len(m.whereDependenciesLogicOperators) && logicOperat != "" {
				logicOperat = m.whereDependenciesLogicOperators[i]
			}

			depend += logicOperat
		}
	}

	return depend
}

func generateWhereConditions(model *SQLModel) string {

	cond := createCondition(model.whereFields, model.whereAliases,
		model.whereFieldsOperator, model.whereLogicOperators,
		&model.whereLengthArrFields)

	if cond != "" {
		if !strings.Contains(model.query, WHERE) {
			cond = WHERE + cond
		} else {
			cond = AND + cond
		}
	}
	return cond
}

func generateSubSelectQuries(model *SQLModel) string {

	sub := ""
	existWhere := strings.Contains(model.query, WHERE)
	if !existWhere {
		sub = WHERE
	}

	sub += generateSubQuries(model, existWhere)

	return sub
}

func generateSubInsertQuries(model *SQLModel) string {

	sub := " " + SELECT + functions.GenerateStrFromArr(
		functions.GenerateArrayOfvalue(len(model.insertFields), "?"),
		DELIMETR,
	) + FROM + DUAL + WHERE

	sub += generateSubQuries(model, false)

	return sub
}

func generateSubQuries(model *SQLModel, whereExists bool) string {

	////	fmt.Println("m.query (generateSubQuries)=>", model.query)
	//	fmt.Println("*SQLModel (generateSubQuries)=>", model)
	sub := ""
	/*	sub := " " + SELECT + functions.GenerateStrFromArr(
			functions.GenerateArrayOfvalue(len(model.insertFields), "?"),
			DELIMETR,
		) + FROM + DUAL + WHERE

		/*if !strings.Contains(model.query, WHERE) {
			sub = WHERE
		}*/

	for i := 0; i < len(model.subModels); i++ {

		quer := model.subModels[i].GetQuery()

		if quer == "" {
			fmt.Println("Пропущен "+strconv.Itoa(i)+", model ", model.subModels[i].sqlTables)
			continue
		}

		if i > 0 || whereExists {
			sub += model.subLogicOperators[i]
		}

		//	fmt.Println("model.subOperators[i] =>", model.subOperators[i])

		switch model.subOperators[i] {
		case EXISTS:
			sub += EXISTS + "(" + quer + ")"
		case NOT_EXISTS:
			sub += NOT_EXISTS + "(" + quer + ")"
		case IN:
			subF := model.subFields[i]
			sub += subF + IN + "(" + quer + ")"
		case NOT_IN:
			subF := model.subFields[i]
			sub += subF + NOT_IN + "(" + quer + ")"
		case EQUAL:
			subF := model.subFields[i]
			sub += subF + EQUAL + "(" + quer + ")"
			/*
						TODO
				case IN:
					/*if model.subAliases[i] != "" {
						sub += model.subAliases[i] + "."
					}
					sub += model.subFields[i] + IN + "(" + model.subModels[i].GetQuery() + ")"
					model.queryParams = append(model.queryParams, model.subModels[i].GetQueryParams()...)*/
		default:
			panic("ERROR MySQL! Unsupported operator for sub query, operator \"" + model.subOperators[i] + "\"")
		}

		model.queryParams = append(model.queryParams, model.subModels[i].GetQueryParams()...)
		//	fmt.Println("m.query (generateSubQuries)=>", model.query, " sub=>", sub)

	}

	return sub
}

func generateGroupBy(model *SQLModel) string {

	return createGroupBy(model.groupFields, model.groupAliases)
}

func generateOrderBy(model *SQLModel) string {

	return ORDER_BY + createOrderBy(model.orderFileds, model.orderAliases, model.orderDirection)
}

func generateLimit(model *SQLModel) string {

	lim := model.limit
	if lim == 0 {
		lim = 1
	}
	return LIMIT + strconv.Itoa(lim)
}

func createSelect(fields []string, aliasTable []string) string {

	sel := ""

	for i := 0; i < len(fields); i++ {

		if len(fields[i]) == 0 {
			continue
		}

		aliasExist := len(aliasTable) > 0 && aliasTable[i] != ""
		if aliasExist {
			sel += aliasTable[i] + "."
		}
		sel += fields[i]

		if aliasExist {
			sel += AS + "\"" + aliasTable[i] + "." + fields[i] + "\""
		}
		sel += DELIMETR
	}

	return sel[0 : len(sel)-len(DELIMETR)]
}

func createListTables(tablesNames []string, tableAlias []string) string {

	from := ""

	if len(tablesNames) == 0 {
		panic("ERROR MySQL! Table name is empty!!!")
	}

	if len(tablesNames) == 0 || (len(tablesNames) > 1 && len(tablesNames) != len(tableAlias)) {
		panic("ERROR MySQL! Length table alias not equal length table alias!")
	}

	var alias string
	for i := 0; i < len(tablesNames); i++ {
		if tablesNames[i] == "" {
			continue
		}

		alias = ""
		if i < len(tableAlias) || (i == 0 && len(tableAlias) == 1 && tableAlias[i] != "") {
			alias = " " + tableAlias[i]
		}

		from += tablesNames[i] + alias + DELIMETR
	}

	if len(from) == 0 || from == FROM {
		fmt.Println("tableName =>", tablesNames, "\n", "tableAlias =>", tableAlias)
		panic("ERROR MySQL! All name tables is empty!!!")
	}

	return from[0 : len(from)-len(DELIMETR)]
}

func createCondition(fields []string, tableAlias []string, operators []string,
	logicOperators []string, arrLenFields *map[string]int) string {

	cond := ""
	for i := 0; i < len(fields); i++ {

		if fields[i] == "" {
			continue
		}

		if i < len(logicOperators) && ((i == 0 && logicOperators[i] == OR) ||
			(i > 0 && logicOperators[i-1] != OR && logicOperators[i] == OR)) {
			cond += "("
		}

		if i < len(tableAlias) && tableAlias[i] != "" {
			cond += tableAlias[i] + "."
		}
		cond += fields[i]

		switch operators[i] {
		case NOT + NULL:
			cond += IS + NOT + NULL
		case NULL:
			cond += IS + NULL
		case IN:
			alias := tableAlias[i]
			if alias != "" {
				alias += "."
			}
			cond += IN + "(" + functions.GenerateStrFromArr(
				functions.GenerateArrayOfvalue((*arrLenFields)[alias+fields[i]], "?"), DELIMETR,
			) + ")"
		case NOT + IN:
			alias := tableAlias[i]
			if alias != "" {
				alias += "."
			}
			cond += NOT + IN + "(" + functions.GenerateStrFromArr(
				functions.GenerateArrayOfvalue((*arrLenFields)[alias+fields[i]], "?"), DELIMETR,
			) + ")"
		case "":
			panic("ERROR MySQL! Operation is not support!!!")
		default:
			cond += operators[i] + PARAMETR
		}

		if i < len(fields)-1 {
			if len(logicOperators) > 0 {
				if i > 0 && logicOperators[i-1] == OR && logicOperators[i] != OR {
					cond += ")"
				}
				cond += logicOperators[i]
			}
		} else {
			if len(logicOperators) > 0 && i > 0 && logicOperators[i-1] == OR {
				cond += ")"
			}
		}
	}

	return cond
}

func createUnionWhere(wheres []string, operators string) string {

	if wheres == nil || len(wheres) == 0 || &operators == nil || len(operators) == 0 {

		fmt.Println("wheres =>", wheres, "\n", "operators =>", operators)
		panic("ERROR MySQL! Conditions where or operator is empty!!!")
	}

	where := ""
	for i := 0; i < len(wheres); i++ {

		if &wheres[i] == nil || len(wheres[i]) == 0 {
			continue
		}
		where += wheres[i]
		if i < len(wheres)-1 {
			where += operators
		}
	}

	if len(where) == 0 {

		fmt.Println("wheres =>", wheres)
		panic("ERROR MySQL! Each field is empty!!!")
	}

	return where
}

func createJoin(direction string, baseTableAlias string, baseField string,
	joinTable string, joinAlias string, joinField string) string {

	if baseTableAlias == "" || baseField == "" || joinTable == "" || joinField == "" {

		/*	fmt.Println("baseTableAlias =>", baseTableAlias, "\n", "baseField =>", baseField, "\n",
			"joinTable =>", joinTable, "\n", "joinField =>", joinField)*/
		panic("ERROR MySQL! Parametrs are not valid!!!")
	}

	join := JOIN

	if direction == LEFT || direction == RIGHT {
		join = direction + join
	}

	var jAlias string
	if joinAlias == "" {
		jAlias = joinTable
		join += joinTable + ON
	} else {
		jAlias = joinAlias
		join += joinTable + " " + jAlias + ON
	}

	join += baseTableAlias + "." + baseField + "=" + jAlias + "." + joinField

	return join
}

func createWhereDependencies(firstAliasTable string, firstField string,
	secondAliasTable string, secondField string, operator string) string {

	if len(firstAliasTable) == 0 || len(firstField) == 0 || len(secondAliasTable) == 0 || len(secondField) == 0 {

		fmt.Println("firstAliasTable =>", firstAliasTable, "\n", "firstField =>", firstField, "\n",
			"secondAliasTable =>", secondAliasTable, "\n", "secondField =>", secondField)
		panic("ERROR MySQL! Parametrs are not valid!!!")
	}

	return firstAliasTable + "." + firstField + operator + secondAliasTable + "." + secondField
}

func createGroupBy(fields []string, tableAliases []string) string {

	if len(fields) == 0 {

		fmt.Println("fields =>", fields)
		panic("ERROR MySQL! Fields are empty!!!")
	}

	if len(tableAliases) != len(fields) {

		fmt.Println("fields =>", fields, "tableAliases =>", tableAliases)
		panic("ERROR MySQL! Count fields and count tables aliases are not equal!!!")
	}

	group := GROUP_BY

	for i, val := range fields {

		if val == "" {
			continue
		}

		alias := ""
		if i < len(tableAliases) {
			alias = tableAliases[i]
		}

		group += alias + "." + val
		if i < len(fields)-1 {
			group += ", "
		}
	}

	if group == GROUP_BY {
		fmt.Println("fields =>", fields)
		panic("ERROR MySQL! All fields are empty!!!")
	}

	return group
}

func createOrderBy(fields []string, tableAliases []string, directions []string) string {

	if len(fields) == 0 {
		fmt.Println("fields =>", fields)
		panic("ERROR MySQL! Fields are empty!!!")
	}

	order := ""
	for i := 0; i < len(fields); i++ {

		if fields[i] == "" {
			continue
		}

		alias := ""
		if i < len(tableAliases[i]) {
			alias = tableAliases[i] + "."
		}
		order += alias + fields[i] + directions[i] + DELIMETR
	}

	return order[0 : len(order)-len(DELIMETR)]
}
