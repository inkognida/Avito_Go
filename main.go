package main

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Balance struct {
	gorm.Model
	UserId      string `json:"user_id"`
	UserBalance int    `json:"user_balance"`
	Info        string `json:"info"`
}

type Account struct {
	gorm.Model
	UserId      string `json:"user_id"`
	ItemID      string `json:"item_id"`
	OrderID     string `json:"order_id"`
	Price       int    `json:"price"`
	UserBalance int    `json:"user_balance"`
	Info        string `json:"info"`
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

	//ADD USER BALANCE
	//data := Balance{
	//	Model:       gorm.Model{},
	//	UserId:      "aboba",
	//	UserBalance: 100,
	//}
	//db.Create(&data)
	//
	// ADD USER ACCOUNT
	//data_ := Account{
	//	Model:       gorm.Model{},
	//	UserId:      "aboba",
	//	ItemID:      "123",
	//	OrderID:     "123",
	//	Price:       100,
	//	UserBalance: data.UserBalance,
	//	Info:        "Order",
	//}
	//db.Create(&data_)

	// UPDATE USER
	//db.Model(&Balance{}).Where("user_id = ?", "aboba").Update("user_balance", 200)
	return db
}

func makeUpdate(db *gorm.DB, income Balance) Balance {
	var user Balance
	var msg string

	userExists := db.Where("user_id = ?", income.UserId).Find(&user)
	if income.UserBalance > 0 && userExists.RowsAffected != 0 {
		db.Model(&Balance{}).Where("user_id = ?", income.UserId).Update("user_balance", income.UserBalance+user.UserBalance)
		msg = "Balance updated with " + strconv.Itoa(income.UserBalance)
		log.Println("Balance updated")
	} else if userExists.RowsAffected == 0 {
		log.Println("Wrong user (doesn't exist)")
		msg = "This user doesn't exist"
	} else {
		msg = "You can't add negative sum"
		log.Println("Wrong income ( < 0 )")
	}
	return Balance{
		Model:       gorm.Model{},
		UserId:      user.UserId,
		UserBalance: income.UserBalance + user.UserBalance,
		Info:        msg,
	}
}

func balanceUpdate(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Not allowed method for checking your balance", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		income := &Balance{}
		err := json.NewDecoder(r.Body).Decode(income)
		if err != nil {
			log.Fatalln(err)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(makeUpdate(db, *income))
	}
	return http.HandlerFunc(fn)
}

func makeOrder(db *gorm.DB, accountOrder Account) Account {
	var user Balance
	var msg string
	userExists := db.Where("user_id = ?", accountOrder.UserId).Find(&user)

	if accountOrder.Price > user.UserBalance {
		msg = "You don't have enough money to order it, we have to delete your order"
		db.Model(&Account{}).Where("user_id = ?", accountOrder.OrderID).Delete(&accountOrder)
		log.Println("Wrong order (not enough money)")
	} else if userExists.RowsAffected != 0 && accountOrder.Price <= user.UserBalance && user.DeletedAt.Valid == true {
		db.Model(&Balance{}).Where("user_id = ?", user.UserId).Update("user_balance", user.UserBalance-accountOrder.Price)
		msg = "Order is ready"
		log.Println("Order accepted")
	} else {
		log.Println("Wrong user (doesn't exist)")
		msg = "This user doesn't exist"
	}
	return Account{
		Model:       gorm.Model{},
		UserId:      user.UserId,
		ItemID:      accountOrder.ItemID,
		OrderID:     accountOrder.OrderID,
		Price:       accountOrder.Price,
		UserBalance: user.UserBalance - accountOrder.Price,
		Info:        msg,
	}
}

func accountOrder(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Not allowed method for checking your balance", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		accountOrder := &Account{}
		err := json.NewDecoder(r.Body).Decode(accountOrder)
		if err != nil {
			log.Fatalln(err)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(makeOrder(db, *accountOrder))
	}
	return http.HandlerFunc(fn)
}

func main() {
	db := createTables()
	http.HandleFunc("/balance", balanceUpdate(db))
	http.HandleFunc("/order", accountOrder(db))
	log.Println("Server is working") // pseudo test
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		log.Println("ListenAndServe error")
		panic(err)
	}
}
