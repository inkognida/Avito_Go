package api

import (
	"Avito_go/pkg/model"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func makeOrder(db *gorm.DB, acc model.AccountRequest, w http.ResponseWriter) interface{} {
	var user model.Balance

	userExists := db.Where("user_id = ?", acc.UserId).Find(&user)
	if userExists.RowsAffected == 0 { //no double user check
		w.WriteHeader(http.StatusBadRequest)
		return model.AccountRespondError{
			UserId:  acc.UserId,
			ItemID:  acc.ItemID,
			OrderID: acc.OrderID,
			Info:    "No such user",
		}
	} else if acc.Price > user.UserBalance {
		w.WriteHeader(http.StatusBadRequest)
		return model.AccountRespondError{
			UserId:  acc.UserId,
			ItemID:  acc.ItemID,
			OrderID: acc.OrderID,
			Info:    "Not enough money to buy it",
		}
	} else if acc.DeletedAt.Valid {
		w.WriteHeader(http.StatusBadRequest)
		return model.AccountRespondError{
			UserId:  acc.UserId,
			ItemID:  acc.ItemID,
			OrderID: acc.OrderID,
			Info:    "Deleted user",
		}
	} else {
		db.Model(&model.Balance{}).Where("user_id = ?", user.UserId).Update("user_balance", user.UserBalance-acc.Price)
		w.WriteHeader(http.StatusOK)
		return model.AccountRespond{
			UserId:      acc.UserId,
			ItemID:      acc.ItemID,
			OrderID:     acc.OrderID,
			UserBalance: user.UserBalance - acc.Price,
		}
	}
}

func AccountOrder(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Not allowed method for doing an order", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		accountOrder := &model.AccountRequest{}
		err := json.NewDecoder(r.Body).Decode(accountOrder)
		if err != nil {
			log.Fatalln(err)
		}
		json.NewEncoder(w).Encode(makeOrder(db, *accountOrder, w))
	}
	return http.HandlerFunc(fn)
}
