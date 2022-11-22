package main

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Balance struct {
	gorm.Model
	UserId      string `json:"user_id"`
	UserBalance int    `json:"user_balance"`
}

type Account struct {
	gorm.Model
	UserId  string `json:"user_id"`
	ItemID  string `json:"item_id"`
	OrderID string `json:"order_id"`
	Price   int    `json:"price"`
}

type Records struct {
	gorm.Model
	UserID  string
	ItemID  string
	OrderID string
	Income  int
}

func createTables() *gorm.DB {
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
	return db
}

func balanceUpdate(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Not allowed method for checking your balance", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		income := &Balance{}
		err := json.NewDecoder(r.Body).Decode(income)
		if err != nil {
			log.Fatalln(err)
		}
		if income.UserBalance >= 0 {
			var user Balance
			db.Where("user_id = ?", income.UserId).First(&user)
			db.Model(&Balance{}).Where("user_id = ?", income.UserId).Update("user_balance", income.UserBalance+user.UserBalance)
			log.Println("Balance updated")
			w.WriteHeader(http.StatusOK)
		} else {
			log.Println("You have to add positive income")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	return http.HandlerFunc(fn)
}

//func accountOrder(db *gorm.DB) http.HandlerFunc {
//	fn := func(w http.ResponseWriter, r *http.Request) {
//		if r.Method != http.MethodPost {
//			http.Error(w, "Not allowed method for checking your balance", http.StatusMethodNotAllowed)
//			return
//		}
//		w.Header().Set("Content-Type", "application/json")
//		accountOrder := &Account{}
//		err := json.NewDecoder(r.Body).Decode(accountOrder)
//		if err != nil {
//			log.Fatalln(err)
//		}
//
//		if accountOrder.UserBalance >= 0 {
//			var user Balance
//			db.Where("user_id = ?", income.UserId).First(&user)
//			db.Model(&Balance{}).Where("user_id = ?", income.UserId).Update("user_balance", income.UserBalance+user.UserBalance)
//			log.Println("Balance updated")
//			w.WriteHeader(http.StatusOK)
//		} else {
//			log.Println("You have to add positive income")
//			w.WriteHeader(http.StatusBadRequest)
//		}
//	}
//	return http.HandlerFunc(fn)
//}

//func server(db *gorm.DB) {
//	balance()
//	log.Println("Server is working") // pseudo test
//	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
//		log.Println("ListenAndServe error")
//		panic(err)
//	}
//}

func main() {
	db := createTables()
	http.HandleFunc("/balance", balanceUpdate(db))
	http.HandleFunc("/order", accountOrder(db))
	log.Println("Server is working") // pseudo test
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		log.Println("ListenAndServe error")
		panic(err)
	}
	//server(db)
}
