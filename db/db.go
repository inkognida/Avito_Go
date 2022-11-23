package db

import (
	"Avito_go/pkg/model"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func createUser(db *gorm.DB) {
	//ADD USER BALANCE
	//data := model.Balance{
	//	Model:       gorm.Model{},
	//	UserId:      uuid.New(),
	//	UserBalance: 100,
	//	Info:        "User balance created",
	//}
	//db.Create(&data)
	//
	////ADD USER ACCOUNT
	//data_ := model.Account{
	//	Model:   gorm.Model{},
	//	UserId:  data.UserId,
	//	ItemID:  uuid.New(),
	//	OrderID: uuid.New(),
	//	Price:   100,
	//	Info:    "",
	//}
	//db.Create(&data_)

	// UPDATE USER
	//db.Model(&Balance{}).Where("user_id = ?", "aboba").Update("user_balance", 200)
}

func CreateTables() *gorm.DB {
	dsn := "user=hardella password=123 dbname=avito sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	// BALANCE
	BalanceInfo := model.Balance{}
	if err := db.AutoMigrate(&BalanceInfo); err != nil {
		log.Fatalln(err)
	}

	// ACCOUNT
	AccountInfo := model.Account{}
	if err := db.AutoMigrate(&AccountInfo); err != nil {
		log.Fatalln(err)
	}

	// RECORDS
	RecordsInfo := model.Records{}
	if err := db.AutoMigrate(&RecordsInfo); err != nil {
		log.Fatalln(err)
	}

	createUser(db)
	return db
}
