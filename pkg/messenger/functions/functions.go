package functions

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const (
	DELIMETR = ", "
)

func RandomString(length int) string {

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func IsDate(date string) bool {
	_, er := time.Parse("2006-01-02", date)
	return er == nil
}
func IsTime(timee string) bool {
	_, er := time.Parse("15:04:05", timee)
	return er == nil
}

func GetNowTimeStr() string {

	return time.Now().Format("15:04:05")
}

func GetNowDateStr() string {

	return time.Now().Format("2006-01-02")
}

func GetNowDateUnix() int64 {

	return time.Now().Unix()
}

func RemoveFromArray[T string | int | int64 | float32 | float64](array []T, value T) []T {

	rez := make([]T, 0, len(array))

	for i := 0; i < len(array); i++ {
		if array[i] != value {
			rez = append(rez, value)
		}
	}

	return rez
}

func FindInArray[T string | int | int64 | float32 | float64](array []T, value T) int {

	for i, v := range array {
		if v == value {
			return i
		}
	}

	return -1
}

func ArraysIsEqual(first []interface{}, second []interface{}) bool {

	if len(first) != len(second) {
		return false
	}

	f, _ := json.Marshal(first)
	s, _ := json.Marshal(second)

	if string(f) != string(s) {
		fmt.Println("DTATATAT: \nstring(f)->", string(f), "\nstring(s)->", string(s))
		return false
	}

	return true
}

func GenerateStrFromArr(arr []string, delimetr string) string {

	if delimetr == "" {
		delimetr = DELIMETR
	}
	var rezult string

	for i := 0; i < len(arr); i++ {
		if arr[i] == "" {
			continue
		}

		rezult += arr[i] + delimetr
	}

	if rezult == "" {
		return rezult
	}

	return rezult[0 : len(rezult)-len(delimetr)]
}

func FillArray[T string | int | float32 | float64](array *[]T, value T) {

	ar := *array
	for i := 0; i < len(ar); i++ {
		ar[i] = value
	}
}

func GenerateArrayOfvalue[T string | int | float32 | float64](length int, value T) []T {

	array := make([]T, length)
	FillArray(&array, value)

	return array
}

func SetValueToField(fieldName string, field reflect.Value, value interface{}) {

	switch d := value.(type) {
	case int:
		field.SetInt(int64(d))
	case int64:
		field.SetInt(d)
	case string:
		field.SetString(d)
	case float32:
		field.SetFloat(float64(d))
	case float64:
		field.SetFloat(d)
	case nil:
	default:
		//		fmt.Println("value.(type) => ", reflect.TypeOf(value).Name())
		//		fmt.Println("value.(type) => ", reflect.TypeOf(value).Elem().Name())
		panic("ERROR adapter load sql rezult to IModel! Not supported type of field: " + fieldName)
	}
}
