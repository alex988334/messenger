package user

import (
	"testing"

	"github.com/alex988334/messenger/pkg/messenger/data"
)

/*
	type User struct {
		data.IModel
		Id             int
		Login          string
		Alias          string
		AuthKey        string
		PassHash       string
		PassResetToken string
		Email          string
		Status         int
		CreateAt       int64
		UpdateAt       int64
	}
*/
func TestCompareUserModels(t *testing.T) {

	t.Log("Generate Compare User Models")
	required := []bool{
		true,
		true,
		true,
		true,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
	}
	params := [][][]interface{}{
		{
			{140, "Grand", "Torino", "kice,fiei", "safdsdsewe", "fdgfbgha", "fsdsdf", 5, 2324, 5343},
			{140, "Grand", "Torino", "kice,fiei", "safdsdsewe", "fdgfbgha", "fsdsdf", 5, 2324, 5343},
		},
		{
			{0, "Grand", "Torino", "kice,fiei", "safdsdsewe", "fdgfbgha", "fsdsdf", 0, 0, 0},
			{0, "Grand", "Torino", "kice,fiei", "safdsdsewe", "fdgfbgha", "fsdsdf", 0, 0, 0},
		},
		{
			{140, "", "", "", "", "", "", 5, 2324, 5343},
			{140, "", "", "", "", "", "", 5, 2324, 5343},
		},
		{
			{0, "", "", "", "", "", "", 0, 0, 0},
			{0, "", "", "", "", "", "", 0, 0, 0},
		},
		// false
		{
			{0, "", "", "", "", "", "", 0, 0, 0},
			{2442, "", "", "", "", "", "", 0, 0, 0},
		},
		{
			{0, "dsfdsf", "", "", "", "", "", 0, 0, 0},
			{0, "", "", "", "", "", "", 0, 0, 0},
		},
		{
			{0, "", "", "", "", "", "", 0, 0, 0},
			{0, "", "fdsgdg", "", "", "", "", 0, 0, 0},
		},
		{
			{0, "", "", "sdfdsfsdf", "", "", "", 0, 0, 0},
			{0, "", "", "", "", "", "", 0, 0, 0},
		},
		{
			{0, "", "", "", "", "", "", 0, 0, 0},
			{0, "", "", "", "asdasfsd", "", "", 0, 0, 0},
		},
		{
			{0, "", "", "", "", "sdfhjf", "", 0, 0, 0},
			{0, "", "", "", "", "", "", 0, 0, 0},
		},
		{
			{0, "", "", "", "", "", "", 0, 0, 0},
			{0, "", "", "", "", "", "asdsad", 0, 0, 0},
		},
		{
			{0, "", "", "", "", "", "", 4324, 0, 0},
			{0, "", "", "", "", "", "", 0, 0, 0},
		},
		{
			{0, "", "", "", "", "", "", 0, 0, 0},
			{0, "", "", "", "", "", "", 0, 34324, 0},
		},
		{
			{0, "", "", "", "", "", "", 0, 0, 34543},
			{0, "", "", "", "", "", "", 0, 0, 0},
		},
		{
			{}, {},
		},
	}

	for ind, test := range params {

		m := NewUser()
		if len(test[0]) > 0 {
			m.Id = test[0][0].(int)
			m.Login = test[0][1].(string)
			m.Alias = test[0][2].(string)
			m.AuthKey = test[0][3].(string)
			m.PassHash = test[0][4].(string)
			m.PassResetToken = test[0][5].(string)
			m.Email = test[0][6].(string)
			m.Status = test[0][7].(int)
			m.CreateAt = int64(test[0][8].(int))
			m.UpdateAt = int64(test[0][9].(int))
		}

		m1 := NewUser()
		if len(test[1]) > 0 {
			m1.Id = test[1][0].(int)
			m1.Login = test[1][1].(string)
			m1.Alias = test[1][2].(string)
			m1.AuthKey = test[1][3].(string)
			m1.PassHash = test[1][4].(string)
			m1.PassResetToken = test[1][5].(string)
			m1.Email = test[1][6].(string)
			m1.Status = test[1][7].(int)
			m1.CreateAt = int64(test[1][8].(int))
			m1.UpdateAt = int64(test[1][9].(int))
		}
		var im1 data.IModel = m1
		if ind == len(params)-1 {
			im1 = NewUserPhone()
		}

		if required[ind] != im1.IsEqualModels(m) {
			t.Fatal("Failed Generate List Messages! \nRequired:", required[ind], "; \nmodel1: ", m, "; \nmodel2: ", im1)
		} else {
			//	fmt.Println("index =>", strconv.Itoa(ind), ", true")
		}
	}
}
