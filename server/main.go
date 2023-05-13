package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type CotacaoEntity struct {
	Bid string
	gorm.Model
}

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	mux := http.NewServeMux()
	db, err := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&CotacaoEntity{})

	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		cot, err := getPrice()
		if err != nil {
			panic(err)
		}
		err = saveToDatabase(db, cot)
		if err != nil {
			panic(errors.New("error to save to database"))
		}
		returnDataToClient(w, cot)
	})

	http.ListenAndServe(":8080", mux)
}

func getPrice() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var decodedRes Cotacao
	err = json.NewDecoder(res.Body).Decode(&decodedRes)
	if err != nil {
		return nil, err
	}

	return &decodedRes, nil
}

func saveToDatabase(db *gorm.DB, c *Cotacao) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	err := db.WithContext(ctx).Model(&CotacaoEntity{}).Create(&CotacaoEntity{Bid: c.USDBRL.Bid}).Error
	if err != nil {
		return err
	}
	return nil
}

func returnDataToClient(w http.ResponseWriter, cot *Cotacao) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(CotacaoResponse{Bid: cot.USDBRL.Bid})
	if err != nil {
		panic(err)
	}
}
