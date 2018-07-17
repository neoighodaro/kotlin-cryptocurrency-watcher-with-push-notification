package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"time"

	"errors"

	"../../notification"
)

// CoinPrice represents a single coin resource
type CoinPrice map[string]interface{}

// Device represents a single device resource
type Device struct {
	ID     int64  `json:"id"`
	UUID   string `json:"uuid"`
	BTCMin int64  `json:"btc_min"`
	BTCMax int64  `json:"btc_max"`
	ETHMin int64  `json:"eth_min"`
	ETHMax int64  `json:"eth_max"`
}

// Devices represents a collection of Devices
type Devices struct {
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
func GetCoinPrices(simulate bool) (CoinPrice, error) {
	coinPrice := make(CoinPrice)
	currencies := [2]string{"ETH", "BTC"}

	for _, curr := range currencies {
		if simulate == true {
			min := 1000.0
			max := 15000.0
			price, _ := big.NewFloat(min + rand.Float64()*(max-min)).SetPrec(8).Float64()
			coinPrice[curr] = map[string]interface{}{"USD": price}
			continue
		}

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
			coinPrice[curr] = map[string]interface{}{"USD": price.(float64)}
		}
	}

	return coinPrice, nil
}

func minMaxQuery(curr string) string {
	return `(` + curr + `_min > 0 AND ` + curr + `_min > ?) OR (` + curr + `_max > 0 AND ` + curr + `_max < ?)`
}

// NotifyDevicesOfPriceChange returns the devices that are within the range
func NotifyDevicesOfPriceChange(db *sql.DB, prices CoinPrice) (Devices, error) {
	devices := Devices{}

	for currency, price := range prices {
		pricing := price.(map[string]interface{})
		rows, err := db.Query("SELECT * FROM devices WHERE "+minMaxQuery(currency), pricing["USD"], pricing["USD"])
		if err != nil {
			return devices, err
		}

		defer rows.Close()

		for rows.Next() {
			device := Device{}
			err = rows.Scan(&device.ID, &device.UUID, &device.BTCMin, &device.BTCMax, &device.ETHMin, &device.ETHMax)
			if err != nil {
				return devices, err
			}

			devices.Devices = append(devices.Devices, device)

			notification.SendNotification(currency, pricing["USD"].(float64), device.UUID)
		}
	}

	return devices, nil
}
