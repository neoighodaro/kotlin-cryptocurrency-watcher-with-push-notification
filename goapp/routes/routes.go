package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"../database/model"

	"github.com/labstack/echo"
)

func createInt64(x int64) *int64 {
	return &x
}

func getCoinValueFromRequest(key string, c echo.Context) (int64, error) {
	value := c.FormValue(key)
	if value != "" {
		setting, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return setting, nil
		}
	}

	return 0, errors.New("Invalid or empty key for: " + key)
}

// GetDeviceSetting fetches the device setting
func GetDeviceSetting(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		device, err := model.GetSettings(db, c.QueryParam("uuid"))
		if err != nil {
			return c.JSON(http.StatusNotFound, err)
		}

		return c.JSON(http.StatusOK, device)
	}
}

// SaveDeviceSettings saves the device settings
func SaveDeviceSettings(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.FormValue("uuid")
		field := make(map[string]int64)

		if btcmin, err := getCoinValueFromRequest("btc_min", c); err == nil {
			field["btc_min"] = btcmin
		}

		if btcmax, err := getCoinValueFromRequest("btc_max", c); err == nil {
			field["btc_max"] = btcmax
		}

		if ethmin, err := getCoinValueFromRequest("eth_min", c); err == nil {
			field["eth_min"] = ethmin
		}

		if ethmax, err := getCoinValueFromRequest("eth_max", c); err == nil {
			field["eth_max"] = ethmax
		}

		device, err := model.SaveSettings(db, uuid, field)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, device)
	}
}

// GetPrices returns the coin prices
func GetPrices() echo.HandlerFunc {
	return func(c echo.Context) error {
		prices, err := model.GetCoinPrices(false)
		if err != nil {
			return c.JSON(http.StatusBadGateway, err)
		}

		return c.JSON(http.StatusOK, prices)
	}
}

// SimulatePriceChanges simulates the prices changes
func SimulatePriceChanges(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		prices, err := model.GetCoinPrices(true)
		if err != nil {
			panic(err)
		}

		devices, err := model.FindDevicesToBeNotified(db, prices)
		if err != nil {
			panic(err)
		}

		resp := map[string]interface{}{
			"prices":  prices,
			"devices": devices,
			"status":  "success",
		}

		return c.JSON(http.StatusOK, resp)
	}
}
