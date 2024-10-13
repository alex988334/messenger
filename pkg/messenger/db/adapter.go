package db

import (
	"reflect"
	"sort"

	cons "github.com/alex988334/messenger/pkg/messenger/constants"
	data "github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/functions"

	"github.com/alex988334/messenger/pkg/messenger/data/chat"
	mess "github.com/alex988334/messenger/pkg/messenger/data/message"
	"github.com/alex988334/messenger/pkg/messenger/data/user"

	mysql "github.com/alex988334/messenger/pkg/messenger/db/mysql-helper"
)

type IExecSQLRezult interface {
	GetLastId() int64
	GetTotalRow() int64
	GetError() error
}

type ICreaterModels interface {
	CreateModel(nameModel string) data.IModel
}

type Adapter struct {
	//	generationModels ICreaterModels
	models       []data.IModel
	sqlModel     *mysql.SQLModel
	execRezult   []IExecSQLRezult
	modelsRezult []map[string]data.IModel
	selectRezult []map[string]interface{}
	db           mysql.DBInterface
	Err          error
}

func NewAdapter() *Adapter {

	var d mysql.DBInterface = NewDB()
	a := Adapter{
		sqlModel: mysql.NewSQLModel(d),
	}
	a.db = d

	return &a
}

func (a *Adapter) GetModelRezult() []map[string]data.IModel {
	return a.modelsRezult
}
func (a *Adapter) GetExecRezult() []IExecSQLRezult {
	return a.execRezult
}
func (a *Adapter) GetIdOfLastInsertRow() int64 {
	return a.execRezult[len(a.execRezult)-1].GetLastId()
}

func (a *Adapter) CloseTransact() {
	a.db.CloseTransact()
}

func (a *Adapter) CloseDBConnection() {
	a.db.CloseDB()
	a.db = nil
}

func (a *Adapter) RunQuery() int {

	a.modelsRezult = []map[string]data.IModel{}
	a.execRezult = []IExecSQLRezult{}
	a.selectRezult = []map[string]interface{}{}

	for i := 0; i < len(a.models); i++ {

		a.sqlModel = mysql.NewSQLModel(a.db)

		//	fmt.Println("a.models[i] =>", a.models[i])

		loadModelToSQLModel(a.models[i], a.sqlModel)

		a.sqlModel.PrepareQuery()

		qt := a.models[i].GetActionType()

		//	fmt.Println("a.sqlModel.GetQuery() =>", a.sqlModel.GetQuery())
		//	fmt.Println("a.sqlModel.a.sqlModel.GetQueryParams()() =>", a.sqlModel.GetQueryParams())

		if qt == data.ACTION_FIND {
			return a.SelectQuery(a.models[i])
		} else {
			if a.ExecQuery(qt, &a.models[i]) == cons.STATUS_ERROR {
				return cons.STATUS_ERROR
			}

		}
	}
	a.CloseTransact()

	return cons.STATUS_ACCEPT
}

func generateDataModel(modelName string) data.IModel {

	var m data.IModel

	switch modelName {
	case data.MODEL_MESSAGE:
		m = mess.NewMessage()
	case data.MODEL_STATUS_MESSAGE:
		m = mess.NewMessageStatus()
	case data.MODEL_CHAT:
		m = chat.NewChat()
	case data.MODEL_CHAT_USER:
		m = chat.NewChatUser()
	case data.MODEL_USER:
		m = user.NewUser()
	case data.MODEL_USER_PHONE:
		m = user.NewUserPhone()
	case data.MODEL_BLACK_LIST:
		m = user.NewBlackList()
	default:
		panic("ERROR parsing sql rezult! Type model not found!!! Model name => " + modelName)
	}

	return m
}

func getModelsNames(model data.IModel, namesModels *[]string, firstModel bool) {

	if model.GetActionType() == data.ACTION_FIND {

		if firstModel || (!firstModel && len(model.GetLinks()) > 0) {

			*namesModels = append(*namesModels, model.GetNameModel())

			if m := model.GetUsers(); len(m) > 0 {
				getModelsNames(m[0], namesModels, false)
			}
			if m := model.GetChatUsers(); len(m) > 0 {
				getModelsNames(m[0], namesModels, false)
			}
			if m := model.GetChats(); len(m) > 0 {
				getModelsNames(m[0], namesModels, false)
			}
			if m := model.GetMessages(); len(m) > 0 {
				getModelsNames(m[0], namesModels, false)
			}
			if m := model.GetBlackLists(); len(m) > 0 {
				getModelsNames(m[0], namesModels, false)
			}
			if m := model.GetStatus(); len(m) > 0 {
				getModelsNames(m[0], namesModels, false)
			}
		}
	}
}

func (a *Adapter) loadModelFromSQLQuery() {

	a.modelsRezult = make([]map[string]data.IModel, len(a.selectRezult))
	namesModels := []string{}

	for _, v := range a.models {

		if v.GetActionType() == data.ACTION_FIND {
			getModelsNames(v, &namesModels, true)
		}
	}

	for c := 0; c < len(a.selectRezult); c++ {

		models := make(map[string]data.IModel, len(namesModels))

		for i := 0; i < len(namesModels); i++ {

			nameModel := namesModels[i]
			tableName := associateTables[nameModel]

			alias := ""
			ind := functions.FindInArray(a.sqlModel.GetTablesName(), tableName)
			if ind == -1 {
				ind = functions.FindInArray(a.sqlModel.GetJoinTables(), tableName)
				alias = a.sqlModel.GetJoinTableAlias()[ind]
			} else {
				alias = a.sqlModel.GetTablesAlias()[ind]
			}

			m := generateDataModel(nameModel)
			isLoaded := false

			t := reflect.TypeOf(m).Elem()
			v := reflect.ValueOf(m).Elem()

			for i := 0; i < t.NumField(); i++ {

				fieldName := t.Field(i).Name
				sqlField := associateFields[nameModel][fieldName]

				var (
					ok bool
					d  interface{}
				)
				d, ok = a.selectRezult[c][sqlField]
				if !ok {
					d, ok = a.selectRezult[c][alias+"."+sqlField]
				}
				if ok {
					functions.SetValueToField(fieldName, v.Field(i), d)
					isLoaded = true
				}
			}
			if isLoaded {
				models[nameModel] = m
			}
		}
		a.modelsRezult[c] = models
	}
}

func (a *Adapter) ExecQuery(queryType int, model *data.IModel) int {

	rezults := a.sqlModel.ExecSQL([][]interface{}{})

	for i := 0; i < len(rezults); i++ {

		a.execRezult = append(a.execRezult, &rezults[i])

		if rezults[i].Err != nil {
			a.Err = rezults[i].Err
			a.db.SetRollBack(true)
			a.db.CloseTransact()
			return cons.STATUS_ERROR
		}
	}

	return cons.STATUS_ACCEPT
}

func (a *Adapter) SelectQuery(model data.IModel) int {

	a.selectRezult = *a.sqlModel.Select()

	if len(a.selectRezult) == 0 {
		return cons.STATUS_ACCEPT
	}

	a.loadModelFromSQLQuery()

	return cons.STATUS_ACCEPT
}

func (a *Adapter) SetModels(models ...data.IModel) {
	a.models = models
}

func getFieldValue(fieldType string, value reflect.Value) (interface{}, bool) {

	if fieldType == "string" && value.String() != "" {
		return interface{}(value.String()), true
	}
	if (fieldType == "int" || fieldType == "int64") && value.Int() > 0 {
		return interface{}(value.Int()), true
	}
	if (fieldType == "float32" || fieldType == "float64") && value.Float() > 0 {
		return interface{}(value.Float()), true
	}

	return nil, false
}

func loadFieldsToSelectPart(model data.IModel, sqlModel *mysql.SQLModel) {

	fields := generateLoadFields(model.GetLoadfields(), model.GetNameModel())
	sqlModel.AddSelectFields(fields, associateTables[model.GetNameModel()])
}

func loadFieldsToWherePart(model data.IModel, sqlModel *mysql.SQLModel) {

	t := reflect.TypeOf(model).Elem()
	v := reflect.ValueOf(model).Elem()
	countF := t.NumField()
	arrayFields := model.GetArrayFields()

	whereFields := make([]string, 0, countF)
	params := make([]interface{}, 0, countF)

	for i := 0; i < countF; i++ {

		if functions.FindInArray(arrayFields, t.Field(i).Name) == -1 {
			if val, ok := getFieldValue(t.Field(i).Type.Name(), v.Field(i)); ok {
				whereFields = append(whereFields, t.Field(i).Name)
				params = append(params, val)
			}
		}
	}

	for i := 0; i < len(whereFields); i++ {
		sqlModel.AddWhere(associateTables[model.GetNameModel()],
			associateFields[model.GetNameModel()][whereFields[i]],
			model.GetFieldOperator(whereFields[i]), mysql.AND)
		sqlModel.AddParamsQuery(params[i])
	}
}

func loadModelToOrderPart(model data.IModel, sqlModel *mysql.SQLModel) {

	fields := model.GetSortFieldsModels()

	if len(fields) == 0 {
		return
	}

	mName := model.GetNameModel()
	sortAlias := make([]string, len(fields))
	for i := 0; i < len(fields); i++ {
		sortAlias[i] = associateFields[mName][fields[i]]
	}

	sqlModel.AddOrderBy(sortAlias, associateTables[mName],
		model.GetDirectionSortModels(),
	)
}

func loadModelToGroupPart(model data.IModel, sqlModel *mysql.SQLModel) {

	fields := model.GetGroupModels()
	mName := model.GetNameModel()

	for i := 0; i < len(fields); i++ {
		fields[i] = associateFields[mName][fields[i]]
	}

	sqlModel.AddGroupBy(fields, associateTables[mName])
}

func loadLinksModelToJoinPart(model data.IModel, sqlModel *mysql.SQLModel) {

	links := model.GetLinks()

	for i := 0; i < len(links); i++ {
		if links[i] == nil {
			continue
		}
		link := *links[i]

		direction := ""

		if link.Weight == data.LINK_WEIGHT_PARENT_MORE {
			direction = mysql.RIGHT
		} else if link.Weight == data.LINK_WEIGHT_CURRENT_MORE {
			direction = mysql.LEFT
		}

		sqlModel.AddJoin(
			direction, associateTables[link.ParentModel], associateTables[link.CurrentModel],
			associateFields[link.ParentModel][link.ParentKey],
			associateFields[link.CurrentModel][link.CurrentKey],
		)
	}
}

func loadLinksModelToWherePart(model data.IModel, sqlModel *mysql.SQLModel) {

	links := model.GetLinks()

	for i := 0; i < len(links); i++ {
		if links[i] == nil {
			continue
		}
		link := *links[i]

		sqlModel.AddLinkWhere(associateTables[link.ParentModel], associateTables[link.CurrentModel],
			associateFields[link.ParentModel][link.ParentKey], associateFields[link.CurrentModel][link.CurrentKey],
			data.EQUAL, mysql.AND)

	}
}

func loadSelectSqlModel(model data.IModel, sqlModel *mysql.SQLModel) {

	loadLinksModelToJoinPart(model, sqlModel)
	loadFieldsToSelectPart(model, sqlModel)
	loadFieldsToWherePart(model, sqlModel)
	loadModelToOrderPart(model, sqlModel)
	loadModelToGroupPart(model, sqlModel)
	loadLimitModel(model, sqlModel)
}

func loadLimitModel(model data.IModel, sqlModel *mysql.SQLModel) {

	if count := model.GetLimitCountModel(); count > 0 {
		sqlModel.SetLimitRows(count)
	}
}

func loadSubModel(parentModel data.IModel, model data.IModel, parent *mysql.SQLModel) *mysql.SQLModel {

	subModel := mysql.NewSQLModel(parent.GetDbConnection())
	loadModelToSQLModel(model, subModel)

	operator := model.GetSubOperator()
	/*	if operator == "" {
		operator = mysql.EXISTS
	}*/
	pName := parentModel.GetNameModel()
	pField := model.GetSubField()
	whereField := associateFields[pName][pField]
	assT := associateTables[pName]
	parent.AddSubQuery(subModel, assT, whereField, operator, mysql.AND)

	return subModel
}

func loadInsertSqlModel(model data.IModel, sqlModel *mysql.SQLModel) {

	tName := associateTables[model.GetNameModel()]

	s := getSeparatedFields(model)

	sqlModel.AddInsertFields(s.ActionFields)
	sqlModel.AddSqlTable(tName)
	sqlModel.AddParamsQuery(s.ActionValues...)

	for i := 0; i < len(model.GetLinks()); i++ {
		loadLinksModelToWherePart(model, sqlModel)
	}
}

type Separator struct {
	SearchFields  []string
	SearchAliases []string
	SearchValues  []interface{}
	ActionFields  []string
	ActionValues  []interface{}
}

func getSeparatedFields(model data.IModel) *Separator {

	findFields := model.GetConditionsFields()

	t := reflect.TypeOf(model).Elem()
	v := reflect.ValueOf(model).Elem()

	mName := model.GetNameModel()
	count := t.NumField()

	s := &Separator{
		SearchFields:  make([]string, 0, count),
		SearchAliases: make([]string, 0, count),
		ActionFields:  make([]string, 0, count),
		SearchValues:  make([]interface{}, 0, count),
		ActionValues:  make([]interface{}, 0, count),
	}

	for i := 0; i < t.NumField(); i++ {

		fieldName := t.Field(i).Name
		if val, ok := getFieldValue(t.Field(i).Type.Name(), v.Field(i)); ok {

			if functions.FindInArray(findFields, fieldName) == -1 || model.GetActionType() == data.ACTION_SAVE {
				s.ActionFields = append(s.ActionFields, associateFields[mName][fieldName])
				s.ActionValues = append(s.ActionValues, val)
			} else {
				s.SearchFields = append(s.SearchFields, associateFields[mName][fieldName])
				s.SearchAliases = append(s.SearchAliases, fieldName)
				s.SearchValues = append(s.SearchValues, val)
			}
		}
	}

	return s
}

func loadUpdateSqlModel(model data.IModel, sqlModel *mysql.SQLModel) {

	tName := associateTables[model.GetNameModel()]

	s := getSeparatedFields(model)

	sqlModel.AddSqlTable(tName)
	sqlModel.AddSetFields(s.ActionFields, tName)
	sqlModel.AddParamsQuery(s.ActionValues...)

	for i := 0; i < len(model.GetLinks()); i++ {
		loadLinksModelToWherePart(model, sqlModel)
	}

	for i := 0; i < len(s.SearchFields); i++ {
		sqlModel.AddWhere(tName,
			s.SearchFields[i], model.GetFieldOperator(s.SearchAliases[i]), mysql.AND)
	}
	sqlModel.AddParamsQuery(s.SearchValues...)
}

func loadDeleteSqlModel(model data.IModel, sqlModel *mysql.SQLModel) {

	tName := associateTables[model.GetNameModel()]

	sqlModel.InitDeleteQuery()
	sqlModel.AddSqlTable(tName)

	s := getSeparatedFields(model)

	for i := 0; i < len(s.SearchFields); i++ {
		sqlModel.AddWhere(tName,
			s.SearchFields[i], model.GetFieldOperator(s.SearchAliases[i]), mysql.AND)
	}
	sqlModel.AddParamsQuery(s.SearchValues...)
}

func loadModelToSQLModel(model data.IModel, sqlModel *mysql.SQLModel) {

	action := model.GetActionType()

	if action == data.ACTION_FIND {
		loadSelectSqlModel(model, sqlModel)
	}

	if action == data.ACTION_UPDATE {
		loadUpdateSqlModel(model, sqlModel)
	}

	if action == data.ACTION_SAVE {
		loadInsertSqlModel(model, sqlModel)
	}

	if action == data.ACTION_DELETE {
		loadDeleteSqlModel(model, sqlModel)
	}

	models := [][]data.IModel{
		model.GetUsers(), model.GetChatUsers(), model.GetChats(),
		model.GetMessages(), model.GetStatus(), model.GetBlackLists(),
	}

	for i := 0; i < len(models); i++ {
		if len(models[i]) > 0 {
			rangeModels(model, models[i], sqlModel)
		}
	}
}

func rangeModels(parent data.IModel, models []data.IModel, sqlModel *mysql.SQLModel) {

	sqlM := sqlModel
	if len(models) > 0 {
		if len(models[0].GetLinks()) > 0 {
			loadModelToSQLModel(models[0], sqlModel)
		} else {
			//for i := 0; i < len(models); i++ {
			if models[0].GetSubOperator() != "" {
				sqlM = loadSubModel(parent, models[0], sqlModel)
			}
			//	}
		}
	}

	if len(models) > 0 && len(models[0].GetArrayFields()) > 0 {
		loadArrayFields(models, sqlM)
	}
}

func loadArrayFields(models []data.IModel, sqlModel *mysql.SQLModel) {

	arrayfields := models[0].GetArrayFields()
	if len(arrayfields) == 0 {
		return
	}

	var params [][]interface{} = make([][]interface{}, len(arrayfields))

	for i := 0; i < len(arrayfields); i++ {
		params[i] = make([]interface{}, len(models), len(models))
	}

	for i := 0; i < len(models); i++ {

		m := models[i]

		t := reflect.TypeOf(m).Elem()
		v := reflect.ValueOf(m).Elem()

		for k := 0; k < len(arrayfields); k++ {

			findfield := arrayfields[k]
			for m := 0; m < t.NumField(); m++ {

				if findfield == t.Field(m).Name {

					fType := t.Field(m).Type.Name()

					if fType == "string" || fType == "int" || fType == "int64" || fType == "float32" || fType == "float64" {
						params[k][i] = v.Field(m).Interface()
						break
					}
				}
			}
		}
	}

	mName := models[0].GetNameModel()

	for i := 0; i < len(arrayfields); i++ {

		sqlModel.AddWhereArrField(associateTables[mName],
			associateFields[mName][arrayfields[i]],
			models[0].GetFieldOperator(arrayfields[i]),
			len(params[i]), mysql.AND)

		sqlModel.AddParamsQuery(params[i]...)
	}
}

/*
func ExecQuery(models *[]data.IModel) bool {

}*/

func generateLoadFields(loadFields map[string]interface{}, model string) []string {

	f := make([]string, len(loadFields))
	i := 0
	for k, _ := range loadFields {
		f[i] = associateFields[model][k]
		i++
	}

	sort.Strings(f)

	return f
}
