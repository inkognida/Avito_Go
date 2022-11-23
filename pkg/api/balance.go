package api

import (
	"Avito_go/pkg/model"
	_ "Avito_go/pkg/model"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func makeUpdate(db *gorm.DB, income model.BalanceChangeRequest, w http.ResponseWriter) interface{} {
	var user model.Balance

	userExists := db.Where("user_id = ?", income.UserId).Find(&user)
	if userExists.RowsAffected == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return model.BalanceChangeError{
			UserId: income.UserId,
			Info:   "No such user",
		}
	} else if income.Income < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return model.BalanceChangeError{
			UserId: income.UserId,
			Info:   "You need to enter positive income",
		}
	} else if user.DeletedAt.Valid {
		w.WriteHeader(http.StatusBadRequest)
		return model.BalanceChangeError{
			UserId: income.UserId,
			Info:   "Deleted user",
		}
	} else {
		db.Model(&model.Balance{}).Where("user_id = ?", income.UserId).Update("user_balance", income.Income+user.UserBalance)
		w.WriteHeader(http.StatusOK)
		return model.BalanceChangeRespond{
			UserId:      user.UserId,
			UserBalance: user.UserBalance + income.Income,
			Income:      income.Income,
			Info:        "Balance updated",
		}
	}
}

func UpdateBalance(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Not allowed method for updating your balance", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		income := &model.BalanceChangeRequest{}

		err := json.NewDecoder(r.Body).Decode(income)
		if err != nil {
			log.Fatalln(err)
		}

		json.NewEncoder(w).Encode(makeUpdate(db, *income, w))
	}
	return http.HandlerFunc(fn)
}

func GetBalance(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Not allowed method for checking your balance", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		checkBalance := &model.BalanceRequest{}
		err := json.NewDecoder(r.Body).Decode(checkBalance)
		if err != nil {
			log.Fatalln(err)
		}
		var user model.Balance

		userExists := db.Where("user_id = ?", checkBalance.UserId).Find(&user)
		if userExists.RowsAffected != 0 && !user.DeletedAt.Valid {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(model.BalanceRespond{
				UserId:      user.UserId,
				UserBalance: user.UserBalance,
				Info:        "Balance",
			})
		} else if userExists.RowsAffected == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.BalanceError{
				UserId: checkBalance.UserId,
				Info:   "No such user",
			})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.BalanceError{
				UserId: checkBalance.UserId,
				Info:   "Deleted user",
			})
		}
	}
	return http.HandlerFunc(fn)
}
