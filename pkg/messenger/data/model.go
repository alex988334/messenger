package data

import (
	"reflect"

	"github.com/alex988334/messenger/pkg/messenger/functions"
)

type IModel interface {
	GetUserLoadFields() map[string]interface{}
	AddUserLoadField(fieldName string)
	InitUserLoadFields(dataModel IModel)
	SetLoadFields(fields ...string)
	AddLoadFields(fields ...string)
	GetLoadfields() map[string]interface{}
	/*GetChangeFields() *[]string
	SetChangeFields([]string)*/
	GetOperation() int
	SetOperation(int)
	GetNameModel() string
	GetActionType() int
	SetActionType(queryType int)
	Validate() bool

	ContainsLoadField(field string) bool
	//SetConditionsFields(conditionsFields []string, operators []string)
	AddConditionField(field string, operator string)
	GetConditionsFields() []string
	SetGroupModels(groupFields []string)
	AddGroupModels(groupFields ...string)
	GetGroupModels() []string
	SetSortModels(sortFields []string, directions []string)
	AddSortModels(direction string, sortFields ...string)
	GetSortFieldsModels() []string
	GetDirectionSortModels() []string
	GetErrorValidate() []error
	SetErrorValidate(errors []error)
	AddErrorValidate(err error)
	AddLink(link *Link)
	GetLinks() []*Link
	SetArrayFields(fields ...string)
	GetArrayFields() []string
	GetFieldOperator(field string) string

	GetSubOperator() string
	SetSubOperator(subOperator string)
	SetSubField(subField string)
	GetSubField() string

	SetLimitCountModel(limitModels int)
	GetLimitCountModel() int
	ResetLinks()
	SetValidator(func() bool)

	IsEqualModels(IModel) bool
	SetComparisonHandler(func(IModel) bool)

	GetUsers() []IModel
	SetUsers(users *[]IModel)
	GetChatUsers() []IModel
	SetChatUsers(chatUsers *[]IModel)
	GetChats() []IModel
	SetChats(chats *[]IModel)
	GetMessages() []IModel
	SetMessages(messages *[]IModel)
	GetStatus() []IModel
	SetStatus(status *[]IModel)
	GetBlackLists() []IModel
	SetBlackLists(blackLists *[]IModel)
}

type ISystem interface {
	GetClientId() int
}

type Model struct {
	userLoadFields map[string]interface{}
	loadFields     map[string]interface{}
	fieldsOperator map[string]string
	//ChangeFields   []string
	operation int
	nameModel string
	links     []*Link
	queryType int

	comparator       func(IModel) bool
	validator        func() bool
	errorValidate    []error
	conditionsFields []string
	limitModels      int
	sortBy           []string
	sortByDirection  []string
	groupBy          []string
	arrayFields      []string
	subOperator      string
	subField         string

	users      []IModel
	chatUsers  []IModel
	chats      []IModel
	messages   []IModel
	status     []IModel
	blackLists []IModel

	//Name string `json:"name_field"`
}

func NewModel(name string) *Model {

	return &Model{
		nameModel:      name,
		loadFields:     make(map[string]interface{}),
		userLoadFields: make(map[string]interface{}),
		links:          []*Link{},
		fieldsOperator: make(map[string]string),
	}
}

func (m *Model) IsEqualModels(model IModel) bool {
	return m.comparator(model)
}
func (m *Model) SetComparisonHandler(handler func(IModel) bool) {
	m.comparator = handler
}

/*
func EqualIModelByKeysField(model1 IModel, model2 IModel) bool {

	t1 := reflect.TypeOf(model1).Elem()
	t2 := reflect.TypeOf(model2).Elem()

	if t1.Name() != t2.Name() {
		return false
	}

	keyFields := model1.GetConditionsFields()
	keyInd := make([]int, 0, len(keyFields))

	for _, keyField := range keyFields {
		for i := 0; i < t1.NumField(); i++ {
			if t1.Field(i).Name == keyField {
				keyInd = append(keyInd, i)
				break
			}
		}
	}

	v1 := reflect.ValueOf(model1).Elem()
	v2 := reflect.ValueOf(model2).Elem()

	for _, ind := range keyInd {
		if !reflect.DeepEqual(v1.Field(ind), v2.Field(ind)) {
			return false
		}
	}

	return true
}*/

func (m *Model) ResetLinks() {
	m.links = []*Link{}
}
func (m *Model) SetErrorValidate(errors []error) {
	m.errorValidate = errors
}
func (m *Model) AddErrorValidate(err error) {
	m.errorValidate = append(m.errorValidate, err)
}
func (m *Model) GetErrorValidate() []error {

	if len(m.errorValidate) > 0 {
		return m.errorValidate
	} else {
		return []error{}
	}
}

/*
	func (m *Model) GetKeyFields() []string {
		return m.keyFields
	}
*/
func (m *Model) SetGroupModels(groupFields []string) {
	m.groupBy = groupFields
}
func (m *Model) AddGroupModels(groupFields ...string) {
	m.groupBy = append(m.groupBy, groupFields...)
}
func (m *Model) GetGroupModels() []string {
	return m.groupBy
}

func (m *Model) SetSortModels(sortFields []string, directions []string) {
	m.sortBy = sortFields
	m.sortByDirection = directions
}
func (m *Model) AddSortModels(direction string, sortFields ...string) {
	m.sortByDirection = append(m.sortByDirection, functions.GenerateArrayOfvalue(len(sortFields), direction)...)
	m.sortBy = append(m.sortBy, sortFields...)
}
func (m *Model) GetSortFieldsModels() []string {
	return m.sortBy
}
func (m *Model) GetDirectionSortModels() []string {
	return m.sortByDirection
}
func (m *Model) GetSubOperator() string {
	return m.subOperator
}
func (m *Model) SetSubOperator(subOperator string) {
	m.subOperator = subOperator
}
func (m *Model) SetSubField(subField string) {
	m.subField = subField
}
func (m *Model) GetSubField() string {
	return m.subField
}

/*func (m *Model) SetConditionsFields(conditionsFields []string, operators []string) {

	m.conditionsFields = conditionsFields

	for i := 0; i < len(conditionsFields); i++ {
		if i < len(operators) {
			m.fieldsOperator[conditionsFields[i]] = operators[i]
		} else {
			m.fieldsOperator[conditionsFields[i]] = EQUAL
		}
	}
}*/

func (m *Model) AddConditionField(field string, operator string) {

	if functions.FindInArray(m.conditionsFields, field) == -1 {
		m.conditionsFields = append(m.conditionsFields, field)
	}
	m.fieldsOperator[field] = operator
}
func (m *Model) GetConditionsFields() []string {
	return m.conditionsFields
}
func (m *Model) GetFieldOperator(field string) string {

	if val, ok := m.fieldsOperator[field]; ok {
		return val
	} else {
		return EQUAL
	}
}

func (m *Model) SetValidator(validator func() bool) {
	m.validator = validator
}

func (m *Model) InitUserLoadFields(dataModel IModel) {

	t := reflect.TypeOf(dataModel).Elem()
	for i := 0; i < t.NumField(); i++ {

		typ := t.Field(i).Type.Name()
		if ((typ == "int" || typ == "int64") && reflect.ValueOf(dataModel).Elem().Field(i).Int() > 0) ||
			(typ == "string" && reflect.ValueOf(dataModel).Elem().Field(i).String() != "") ||
			((typ == "float32" || typ == "float64") && reflect.ValueOf(dataModel).Elem().Field(i).Float() > 0) {
			m.AddUserLoadField(t.Field(i).Name)
		}
	}
}

func (m *Model) AddUserLoadField(fieldName string) {
	m.userLoadFields[fieldName] = nil
}

func (m *Model) ContainsLoadField(field string) bool {
	if _, ok := m.userLoadFields[field]; ok {
		return true
	} else {
		return false
	}
}

func (m *Model) Validate() bool {
	return m.validator()
}

func (m *Model) GetUsers() []IModel {
	return m.users
}
func (m *Model) SetUsers(users *[]IModel) {
	m.users = *users
}
func (m *Model) GetChatUsers() []IModel {
	return m.chatUsers
}
func (m *Model) SetChatUsers(chatUsers *[]IModel) {
	m.chatUsers = *chatUsers
}
func (m *Model) GetChats() []IModel {
	return m.chats
}
func (m *Model) SetChats(chats *[]IModel) {
	m.chats = *chats
}
func (m *Model) GetMessages() []IModel {
	return m.messages
}
func (m *Model) SetMessages(messages *[]IModel) {
	m.messages = *messages
}
func (m *Model) GetStatus() []IModel {
	return m.status
}
func (m *Model) SetStatus(status *[]IModel) {
	m.status = *status
}
func (m *Model) GetBlackLists() []IModel {
	return m.blackLists
}
func (m *Model) SetBlackLists(blackLists *[]IModel) {
	m.blackLists = *blackLists
}

func (m *Model) GetNameModel() string {
	return m.nameModel
}

func (m *Model) AddLink(link *Link) {
	m.links = append(m.links, link)
}
func (m *Model) GetLinks() []*Link {
	return m.links
}

func (m *Model) SetArrayFields(fields ...string) {
	m.arrayFields = fields
}
func (m *Model) GetArrayFields() []string {
	return m.arrayFields
}
func (m *Model) SetLimitCountModel(limitModels int) {
	m.limitModels = limitModels
}
func (m *Model) GetLimitCountModel() int {
	return m.limitModels
}

/*
	func (m *Model) SetUserLoadFields(fields map[string]interface{}) {
		m.UserLoadFields = fields
	}
*/

func (m *Model) AddLoadFields(fields ...string) {

	for _, v := range fields {
		m.loadFields[v] = nil
	}
}

func (m *Model) SetLoadFields(fields ...string) {

	for k, _ := range m.loadFields {
		delete(m.loadFields, k)
	}

	for _, v := range fields {
		m.loadFields[v] = nil
	}
}
func (m *Model) GetLoadfields() map[string]interface{} {
	return m.loadFields
}
func (m *Model) GetUserLoadFields() map[string]interface{} {
	return m.userLoadFields
}

/*
func (m *Model) SetChangeFields(changeFields []string) {
	m.ChangeFields = changeFields
}
func (m *Model) GetChangeFields() *[]string {
	return &m.ChangeFields
}*/

func (m *Model) SetOperation(operation int) {
	m.operation = operation
}
func (m *Model) GetOperation() int {
	return m.operation
}

func (m *Model) GetActionType() int {
	return m.queryType
}
func (m *Model) SetActionType(queryType int) {
	m.queryType = queryType
}

/*
func main() {
	m := Model{Name: "имя_объекта"}
	field, _ := reflect.TypeOf(m).FieldByName("Name")
	val := reflect.ValueOf(m)
	fmt.Println(field, "=>", val.FieldByName("Name"))

	reflect.ValueOf(&m).Elem().FieldByName("Name").SetString("Строка была изменена")

	fmt.Println(wss.OP_ADD_USER)
}*/
