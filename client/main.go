package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type ResponseFromAPI struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var decodedRes ResponseFromAPI
	err = json.NewDecoder(res.Body).Decode(&decodedRes)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = generateTxtFile(f, decodedRes)
	if err != nil {
		panic(err)
	}
}

func generateTxtFile(f *os.File, r ResponseFromAPI) error {
	_, err := f.Write([]byte("Dólar: " + r.Bid))
	if err != nil {
		return errors.New("error to write file")
	}
	return nil
}
