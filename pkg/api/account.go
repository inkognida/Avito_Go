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

	userExists := db.Where("user_id = ?", acc.AccountOrder.UserId).Find(&user)
	if userExists.RowsAffected == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return model.AccountRespondError{
			UserId:  acc.AccountOrder.UserId,
			ItemID:  acc.AccountOrder.ItemID,
			OrderID: acc.AccountOrder.OrderID,
			Info:    "No such user",
		}
	} else if acc.AccountOrder.Price > user.UserBalance {
		w.WriteHeader(http.StatusBadRequest)
		return model.AccountRespondError{
			UserId:  acc.AccountOrder.UserId,
			ItemID:  acc.AccountOrder.ItemID,
			OrderID: acc.AccountOrder.OrderID,
			Info:    "Not enough money to buy it",
		}
	} else if acc.AccountOrder.DeletedAt.Valid {
		w.WriteHeader(http.StatusBadRequest)
		return model.AccountRespondError{
			UserId:  acc.AccountOrder.UserId,
			ItemID:  acc.AccountOrder.ItemID,
			OrderID: acc.AccountOrder.OrderID,
			Info:    "Deleted user",
		}
	} else {
		db.Model(&model.Balance{}).Where("user_id = ?", user.UserId).Update("user_balance", user.UserBalance-acc.AccountOrder.Price)
		w.WriteHeader(http.StatusOK)
		return model.AccountRespond{
			UserId:      acc.AccountOrder.UserId,
			ItemID:      acc.AccountOrder.ItemID,
			OrderID:     acc.AccountOrder.OrderID,
			UserBalance: user.UserBalance - acc.AccountOrder.Price,
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
