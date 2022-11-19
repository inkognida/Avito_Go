package main

import (
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Balance struct {
	gorm.Model
	UserId      string
	UserBalance int
}

type Account struct {
	gorm.Model
	UserId  string
	ItemID  string
	OrderID string
	Price   int
}

type Records struct {
	gorm.Model
	UserID  string
	ItemID  string
	OrderID string
	Income  int
}

func createTables() {
	dsn := "user=hardella password=123 dbname=avito sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	// BALANCE
	BalanceInfo := Balance{}
	if err := db.AutoMigrate(&BalanceInfo); err != nil {
		log.Fatalln(err)
	}

	// ACCOUNT
	AccountInfo := Account{}
	if err := db.AutoMigrate(&AccountInfo); err != nil {
		log.Fatalln(err)
	}

	// RECORDS
	RecordsInfo := Records{}
	if err := db.AutoMigrate(&RecordsInfo); err != nil {
		log.Fatalln(err)
	}

	// ADD USER
	//data := Balance{
	//	Model:       gorm.Model{},
	//	UserId:      "aboba",
	//	UserBalance: 100,
	//}
	//db.Create(&data)

	// UPDATE USER
	//db.Model(&Balance{}).Where("user_id = ?", "aboba").Update("user_balance", 200)
}

func main() {
	createTables()
	//router := httprouter.New()
}
