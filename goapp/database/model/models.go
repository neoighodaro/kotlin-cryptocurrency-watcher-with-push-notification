package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"errors"
)

// CoinPrice represents a single coin resource
type CoinPrice map[string]float64

// Device represents a single device resource
type Device struct {
	ID     int64  `json:"id"`
	UUID   string `json:"uuid"`
	BTCMin int64  `json:"btc_min"`
	BTCMax int64  `json:"btc_max"`
	ETHMin int64  `json:"eth_min"`
	ETHMax int64  `json:"eth_max"`
}

// DeviceCollection represents a collection of Devices
type DeviceCollection struct {
	Devices []Device `json:"items"`
}

// CreateSettings creates a new device and saves it to the db
func CreateSettings(db *sql.DB, uuid string) (Device, error) {
	device := Device{UUID: uuid, BTCMin: 0, BTCMax: 0, ETHMin: 0, ETHMax: 0}
	stmt, err := db.Prepare("INSERT INTO devices (uuid, btc_min, btc_max, eth_min, eth_max) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return device, err
	}

	res, err := stmt.Exec(device.UUID, device.BTCMin, device.BTCMax, device.ETHMin, device.ETHMax)
	if err != nil {
		return device, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return device, err
	}

	device.ID = lastID

	return device, nil
}

// GetSettings fetches the settings for a single user from the db
func GetSettings(db *sql.DB, uuid string) (Device, error) {
	device := Device{}

	if len(uuid) <= 0 {
		return device, errors.New("Invalid device UUID")
	}

	err := db.QueryRow("SELECT * FROM devices WHERE uuid=?", uuid).Scan(
		&device.ID,
		&device.UUID,
		&device.BTCMin,
		&device.BTCMax,
		&device.ETHMin,
		&device.ETHMax)
	if err != nil {
		return CreateSettings(db, uuid)
	}

	return device, nil
}

// SaveSettings saves the devices settings
func SaveSettings(db *sql.DB, uuid string, field map[string]int64) (Device, error) {
	device, err := GetSettings(db, uuid)
	if err != nil {
		return Device{}, err
	}

	if btcmin, isset := field["btc_min"]; isset {
		device.BTCMin = btcmin
	}

	if btcmax, isset := field["btc_max"]; isset {
		device.BTCMax = btcmax
	}

	if ethmin, isset := field["eth_min"]; isset {
		device.ETHMin = ethmin
	}

	if ethmax, isset := field["eth_max"]; isset {
		device.ETHMax = ethmax
	}

	stmt, err := db.Prepare("UPDATE devices SET btc_min = ?, btc_max = ?, eth_min = ?, eth_max = ? WHERE uuid = ?")
	if err != nil {
		return Device{}, err
	}

	_, err = stmt.Exec(device.BTCMin, device.BTCMax, device.ETHMin, device.ETHMax, device.UUID)
	if err != nil {
		return Device{}, err
	}

	return device, nil
}

// GetCoinPrices gets the current coin prices
func GetCoinPrices() (CoinPrice, error) {
	coinPrice := make(CoinPrice)
	currencies := [2]string{"ETH", "BTC"}

	for _, curr := range currencies {
		url := fmt.Sprintf("https://min-api.cryptocompare.com/data/pricehistorical?fsym=%s&tsyms=USD&ts=%d", curr, time.Now().Unix())
		res, err := http.Get(url)
		if err != nil {
			return coinPrice, err
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return coinPrice, err
		}

		var f interface{}
		err = json.Unmarshal([]byte(body), &f)
		if err != nil {
			return coinPrice, err
		}

		priceMap := f.(map[string]interface{})[curr]
		for _, price := range priceMap.(map[string]interface{}) {
			coinPrice[curr] = price.(float64)
		}
	}

	return coinPrice, nil
}
