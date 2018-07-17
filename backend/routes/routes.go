package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"../database/model"

	"github.com/labstack/echo"
)

var postedSettings map[string]string

func formValue(c echo.Context, key string) (string, error) {
	if postedSettings == nil {
		if err := c.Bind(&postedSettings); err != nil {
			return "", err
		}
	}

	return postedSettings[key], nil
}

func getCoinValueFromRequest(key string, c echo.Context) (int64, error) {
	value, _ := formValue(c, key)
	if value != "" {
		setting, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return setting, nil
		}
	}

	return 0, errors.New("Invalid or empty key for: " + key)
}

// SaveDeviceSettings saves the device settings
func SaveDeviceSettings(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid, _ := formValue(c, "uuid")
		field := make(map[string]int64)

		if btcmin, err := getCoinValueFromRequest("minBTC", c); err == nil {
			field["btc_min"] = btcmin
		}

		if btcmax, err := getCoinValueFromRequest("maxBTC", c); err == nil {
			field["btc_max"] = btcmax
		}

		if ethmin, err := getCoinValueFromRequest("minETH", c); err == nil {
			field["eth_min"] = ethmin
		}

		if ethmax, err := getCoinValueFromRequest("maxETH", c); err == nil {
			field["eth_max"] = ethmax
		}

		defer func() { postedSettings = nil }()

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
		prices, err := model.GetCoinPrices(true)
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

		devices, err := model.NotifyDevicesOfPriceChange(db, prices)
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
