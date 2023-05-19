package models

import (
	"encoding/json"
	"strconv"
)

const (
	Snapshot = "Snapshot"
	Update   = "Update"
)

type Item struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

func (item *Item) UnmarshalJSON(data []byte) error {
	type Alias struct {
		Price string `json:"price"`
		Size  string `json:"size"`
	}

	temp := Alias{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	price, err := strconv.ParseFloat(temp.Price, 64)
	if err != nil {
		return err
	}

	size, err := strconv.ParseFloat(temp.Size, 64)
	if err != nil {
		return err
	}

	item.Price = price
	item.Size = size
	return nil
}

func (item Item) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Price string `json:"price"`
		Size  string `json:"size"`
	}

	temp := Alias{}
	temp.Price = strconv.FormatFloat(item.Price, 'f', -1, 64)
	temp.Size = strconv.FormatFloat(item.Size, 'f', -1, 64)

	return json.Marshal(temp)
}

type Data struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Bids []Item `json:"bids"`
	Asks []Item `json:"asks"`
}
