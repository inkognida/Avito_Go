package api

import (
	"Avito_go/pkg/model"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func MakeRecords(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Not allowed method for making records", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		record := &model.RecordsRequest{}
		err := json.NewDecoder(r.Body).Decode(record)

		if err != nil {
			log.Fatalln(err)
		}

		var order model.Account
		var records []model.Records

		orderExists := db.Where("order_id = ?", record.OrderID).Find(&order)
		if orderExists.RowsAffected == 1 && !order.DeletedAt.Valid {
			recordExists := db.Where("order_id = ?", record.OrderID).Find(&records)
			if recordExists.RowsAffected != 0 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(model.RecordsRespond{
					Status: false,
					Info:   "Bad record",
				})
				return
			}
			db.Model(&model.Records{}).Create(record)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(model.RecordsRespond{
				Status: true,
				Info:   "Record done",
			})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.RecordsRespond{
				Status: false,
				Info:   "Bad record",
			})
		}
	}
	return http.HandlerFunc(fn)
}

func GetRecords(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Not allowed method for getting records", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)

		item := &model.RecordsRequest{}
		err := json.NewDecoder(r.Body).Decode(item)

		if err != nil {
			log.Fatalln(err)
		}

		var records []model.Records
		db.Where("item_id = ?", item.ItemID).Find(&records)
		json.NewEncoder(w).Encode(records)
	}
	return http.HandlerFunc(fn)
}
