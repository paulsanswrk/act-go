package db

import (
	"ACT_GO/db/entities"
	"ACT_GO/utils"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {
	var n int

	DB.Raw("select count(*) from binance.Prices").Scan(&n)
	fmt.Printf("Prices: %v\n", n)
	DB.Raw("select count(*) from app.log").Scan(&n)
	fmt.Printf("Log: %v\n", n)
}

func TestInsertPrice(t *testing.T) {
	var price = &entities.Price{
		Symbol:        "test",
		Time:          0,
		Open:          0.1,
		High:          0,
		Low:           0,
		Close:         0,
		Volume:        0,
		WeightedPrice: 0,
	}

	res := DB.Create(price)
	fmt.Printf("Prices: %v\n", res)

}

func TestInsertLog(t *testing.T) {
	var log = &entities.Log{
		//Id:       0,
		Category: 2,
		Message:  "Message",
		Tag:      "Tag",
		Request:  string(utils.First(json.Marshal([]struct{ X string }{{"qqq"}}))),
		//Response: "",
		//Date:     time.Now(),
	}

	res := DB.Create(log)
	fmt.Printf("Log: %v\n", res)

}

func TestAddUser(t *testing.T) {
	user := &User{Email: "test", Pwd: fmt.Sprintf("%x", sha512.Sum512([]byte("-pwd-111"+pwd_salt)))}

	res := DB.Create(user)
	fmt.Printf("Log: %v\n", res)
}

func TestGetUser(t *testing.T) {
	var user User
	res := DB.Where("email = ?", "test").First(&user)
	fmt.Printf("Log: %v\n", res)
	fmt.Printf("Log: %v\n", user)
}

func TestGetUser2(t *testing.T) {
	var user = GetUser("test", "-pwd-111")
	fmt.Printf("Log: %v\n", user)
}
