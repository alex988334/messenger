package business_logic

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"time"

	cons "github.com/alex988334/messenger/pkg/messenger/constants"
	"github.com/alex988334/messenger/pkg/messenger/data"
	"github.com/alex988334/messenger/pkg/messenger/data/chat"
	"github.com/alex988334/messenger/pkg/messenger/data/message"
	"github.com/alex988334/messenger/pkg/messenger/data/user"

	//	data "github.com/alex988334/messenger/pkg/messenger/data"

	//	"github.com/alex988334/messenger/pkg/messenger/data"

	//	chat "github.com/alex988334/messenger/pkg/messenger/data/chat"
	//	user "github.com/alex988334/messenger/pkg/messenger/data/user"

	db "github.com/alex988334/messenger/pkg/messenger/db"
)

type IClient interface {
	IsAutorizate() bool
	SetId(userId int)
	GetId() int
}

/*
	type ConnectingUser struct {
		UserId          int
		ConnectingChats []int
		ConnectedChats  []int
	}
*/
type DateTimeStamp struct {
	unixDate int64
	date     string
	time     string
}

func (dt *DateTimeStamp) GetDate() string {
	return dt.date
}
func (dt *DateTimeStamp) GetTime() string {
	return dt.time
}
func (dt *DateTimeStamp) GetUnix() int64 {
	return dt.unixDate
}

func NewDateTimeStamp() *DateTimeStamp {

	t := time.Now()

	return &DateTimeStamp{
		unixDate: t.Unix(),
		date:     t.Format("2006-01-02"),
		time:     t.Format("15:04:05"),
	}
}

type ActionUserConnections struct {
	UsersId        []int
	ActionChats    []int
	ConnectedChats []int
	Operation      int
}
type BusinessHandler struct {
	httpQuery             map[string]string
	models                []data.IModel
	client                IClient
	adapter               db.Adapter
	opearation            int
	actionUserConnections *ActionUserConnections
	//connectUser      *ConnectingUser
	response         map[string][]data.IModel
	chatsId          []int
	registrationTime *DateTimeStamp
	err              error
}

func NewBusinessHandler(requestParams *map[string]string, clientData IClient, stamp *DateTimeStamp) *BusinessHandler {

	if stamp == nil {
		stamp = NewDateTimeStamp()
	}

	return &BusinessHandler{
		httpQuery:             *requestParams,
		client:                clientData,
		adapter:               *db.NewAdapter(),
		actionUserConnections: nil,
		//connectUser:      nil,
		registrationTime: stamp,
		response:         make(map[string][]data.IModel),
		chatsId:          []int{},
	}
}
func (b *BusinessHandler) TimeRegistration() string {
	return b.registrationTime.time
}
func (b *BusinessHandler) DateRegistration() string {
	return b.registrationTime.date
}

func (b *BusinessHandler) generateListConnection(users []int, actionChats []int, operation int) {

	connectedCh := make([]int, 0, len(users)*10)

	chUs := chat.NewChatUser()
	chUs.SetActionType(data.ACTION_FIND)
	chUs.SetLoadFields("Chat", "User")
	b.models = []data.IModel{chUs}

	for _, userId := range users {

		chUs.User = userId

		if !b.runAdapter() {
			continue
		}

		rez := b.adapter.GetModelRezult()
		chatsId := make([]int, 0, len(rez))

		for i := 0; i < len(rez); i++ {
			chatsId = append(chatsId, rez[i][data.MODEL_CHAT_USER].(*chat.ChatUser).Chat)
		}

		sort.Ints(chatsId)
		chatsId = slices.Compact(chatsId)

		connectedCh = append(connectedCh, chatsId...)
	}

	sort.Ints(connectedCh)
	connectedCh = slices.Compact(connectedCh)
	connectedCh = slices.Clip(connectedCh)

	b.actionUserConnections = &ActionUserConnections{
		UsersId:        users,
		ActionChats:    actionChats,
		ConnectedChats: connectedCh,
		Operation:      operation,
	}
}

func (b *BusinessHandler) GetActionUserConnections() *ActionUserConnections {
	return b.actionUserConnections
}

/*
	func (b *BusinessHandler) GetConnectingUser() *ConnectingUser {
		return b.connectUser
	}
*/
func (b *BusinessHandler) GetOperation() int {
	return b.opearation
}

func (b *BusinessHandler) ProcessinRequest() ([]byte, []int, error) {

	var ok bool

	b.models, b.err = LoadModels(&b.httpQuery, b.client)

	if b.err != nil {
		fmt.Println("LoadModels LoadModels=>", b.err)
		return b.GenerateResponse()
	}

	b.opearation = b.models[0].GetOperation()

	if !b.client.IsAutorizate() && b.opearation != cons.OP_AUTORIZATE && b.opearation != cons.OP_REGISTRATION {
		b.err = errors.New("ERROR! You must be logged in")
		return b.GenerateResponse()
	}

	b.initUserLoadFields()
	if ok, b.err = b.validateModels(); !ok {

		return b.GenerateResponse()
	}

	b.operationProcessing()

	if b.adapter.Err != nil {
		b.err = b.adapter.Err
	}

	return b.GenerateResponse()
}

func (b *BusinessHandler) initUserLoadFields() {

	for i := 0; i < len(b.models); i++ {
		b.models[i].InitUserLoadFields(b.models[i])
	}
}

func (b *BusinessHandler) AppendModelToResponse(models ...data.IModel) {

	for _, v := range models {
		b.response[v.GetNameModel()] = append(b.response[v.GetNameModel()], v)
	}
}

func (b *BusinessHandler) GetSendChatsId() []int {
	return b.chatsId
}

// load with except copy models from rezult sql select v3.0

func (b *BusinessHandler) loadModelsFromAdapterRezult(modelNames ...string) {

	rezult := b.adapter.GetModelRezult()

	//fmt.Println("b.adapter.GetModelRezult() =>", rezult)

	if len(rezult) == 0 {
		return
	}

	if len(modelNames) == 0 {
		for key, _ := range rezult[0] {
			modelNames = append(modelNames, key)
		}
	}

	if b.response == nil {
		b.response = make(map[string][]data.IModel)
	}

	/*keys := make([]string, 0, len(rezult[0]))
	for k, _ := range rezult[0] {
		keys = append(keys, k)
	}*/

	for _, key := range modelNames {

		var models []data.IModel
		count := len(rezult)

		if v, ok := b.response[key]; ok {
			models = make([]data.IModel, 0, count+len(v))
			models = append(models, v...)
		} else {
			models = make([]data.IModel, 0, count)
		}

		for _, val := range rezult {

			find := false
			for _, v := range models {

				if _, ok := val[key]; ok && val[key].IsEqualModels(v) {

					find = true
				}
			}

			if !find {

				models = append(models, val[key])
			}
		}
		b.response[key] = models
	}

}

/*
// load with except copy models from rezult sql select v2.0
func (b *BusinessHandler) loadModelsFromAdapterRezult() {

	rezult := b.adapter.GetModelRezult()
	if len(rezult) == 0 {
		return
	}

	if b.response == nil {
		b.response = make(map[string][]data.IModel)
	}

	keys := make([]string, 0, len(rezult[0]))
	for k, _ := range rezult[0] {
		keys = append(keys, k)
	}

	for _, key := range keys {

		var models []data.IModel
		count := len(rezult)

		var modelsEncodes []string
		if v, ok := b.response[key]; ok {
			models = make([]data.IModel, 0, count+len(v))
			models = append(models, v...)
			modelsEncodes = make([]string, 0, len(v)+count)
		} else {
			models = make([]data.IModel, 0, count)
			modelsEncodes = make([]string, 0, count)
		}

		for i := 0; i < len(b.response[key]); i++ {
			en, err := json.Marshal(b.response[key][i])
			if err == nil {
				modelsEncodes = append(modelsEncodes, string(en))
			}
		}

		for _, val := range rezult {

			en, _ := json.Marshal(val[key])
			str := string(en)
			if functions.FindInArray(modelsEncodes, str) == -1 {
				modelsEncodes = append(modelsEncodes, str)
				models = append(models, val[key])
			}
			/*		replaced := false
					for i := 0; i < len(models); i++ {

						if data.EqualIModelByKeysField(models[i], val[key]) {
							models[i] = val[key]
							replaced = true
							break
						}
					}
*/

/*	}
		b.response[key] = models
	}
}
*/
/*  load All Models from rezult sql select  v 1.0
func (b *BusinessHandler) loadModelsFromAdapterRezult() {

	rezult := b.adapter.GetModelRezult()
	if len(rezult) == 0 {
		return
	}

	if b.response == nil {
		b.response = make(map[string][]data.IModel)
	}

	keys := make([]string, 0, len(rezult[0]))
	for k, _ := range rezult[0] {
		keys = append(keys, k)
	}

	for _, key := range keys {

		var models []data.IModel
		count := len(rezult)

		if v, ok := b.response[key]; ok {
			models = make([]data.IModel, 0, count+len(v))
			models = append(models, v...)
		} else {
			models = make([]data.IModel, 0, count)
		}

		for _, val := range rezult {
			models = append(models, val[key])
		}

		b.response[key] = models
	}
}*/

func createModel(modelName string) data.IModel {

	switch modelName {
	case data.MODEL_CHAT:
		return chat.NewChat()
	case data.MODEL_MESSAGE:
		return message.NewMessage()
	case data.MODEL_CHAT_USER:
		return chat.NewChatUser()
	case data.MODEL_STATUS_MESSAGE:
		return message.NewMessageStatus()
	case data.MODEL_USER:
		return user.NewUser()
	case data.MODEL_USER_PHONE:
		return user.NewUserPhone()
	case data.MODEL_BLACK_LIST:
		return user.NewBlackList()
	default:
		panic("Not supported type model!")
	}
}

func initArrayFieldCondition(model data.IModel, field string, operator string) {

	if operator != data.IN && operator != data.NOT_IN {
		operator = data.IN
	}
	model.SetArrayFields(field)
	model.AddConditionField(field, operator)
}

/*
**
func initArrayFieldCondition(modelName string, operator string, fields []string, values [][]interface{}) []data.IModel {

		data := make([]data.IModel, len(values))

		for i := 0; i < len(data); i++ {

			model := createModel(modelName)
			vals := values[i]

			t := reflect.TypeOf(model).Elem()
			v := reflect.ValueOf(model).Elem()

			for ind, field := range fields {

				for k:=0; k < t.NumField(); k++ {

					if t.Field(k).Name == field {

						v.Field(k)
						break
					}
				}
			}
		}
			model.SetArrayFields("User")
			model.AddConditionField("User", data.IN)

		return data
	}
*/
func initConditionOperatorOfSubQuery(model data.IModel, operator string, loadFields ...string) {

	model.SetActionType(data.ACTION_FIND)
	/* in new code
	model.SetSubOperator(operator)
	*/

	model.SetSubOperator(operator)
	/***/
	model.SetLoadFields(loadFields...)
}
func initExistConditionOfSubQuery(model data.IModel, operator string, loadFields ...string) {

	model.SetActionType(data.ACTION_FIND)
	model.SetSubOperator(operator)
	model.SetLoadFields(loadFields...)
}

func (b *BusinessHandler) validateModels() (bool, error) {

	for i := 0; i < len(b.models); i++ {
		if !b.models[i].Validate() {
			return false, b.models[i].GetErrorValidate()[0]
		}
	}

	return true, nil
}

func (b *BusinessHandler) runAdapter() bool {

	b.adapter.SetModels(b.models...)

	return b.adapter.RunQuery() != cons.STATUS_ERROR
}

func (b *BusinessHandler) operationProcessing() {

	switch b.models[0].GetOperation() {
	case cons.OP_STATUS_MESSAGE:
		b.processStatusMessage()
	case cons.OP_NEW_MESSAGE:
		b.processNewMessage()
	case cons.OP_LIST_USERS:
		b.processListUsers()
	case cons.OP_CREATE_NEW_CHAT:
		b.processNewChat()
	case cons.OP_WRITEN:
	case cons.OP_SYSTEM:
	case cons.OP_SEARCH_USER:
		b.processSearchUser()
	case cons.OP_GET_CHATS:
		b.processListChats()
	case cons.OP_LIST_NEXT_MESSAGES:
		b.processListMessge()
	case cons.OP_LIST_PREVIOUS_MESSAGES:
		b.processListMessge()
	case cons.OP_EXIT_CHAT:
		b.processExitChat()
	case cons.OP_REMOVE_USER:
		b.processRemoveUserFromChat()
	case cons.OP_ADD_USER:
		b.processAddUserInChat()
	case cons.OP_REMOVE_CHAT:
		b.processRemoveChat()
	case cons.OP_BLOCK_USERS:
	case cons.OP_UNLOOCK_USERS:
	case cons.OP_BLACK_LIST_USERS:
	case cons.OP_GET_FILE:
	case cons.OP_AUTORIZATE:
		b.processAutorizate()
	case cons.OP_MY_DATA:
		b.processMyData()
	case cons.OP_REGISTRATION:
		b.processRegistration()
	default:

	}
}

/*
func prepareStatusMessage(m *mess.MessageStatus) {

	now := time.Now()
	m.SetDate(now.Format("2016-01-01"))
	m.SetTime(now.Format("15:04:01"))

	//	loadFields := []string{"messageId", "userId", "status", "date", "time"}

}*/
