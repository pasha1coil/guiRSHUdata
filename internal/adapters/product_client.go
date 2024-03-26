package adapters

import (
	"demofine/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func ProductClient(from time.Time, to time.Time) models.ProductCalendarInfo {
	client := fiber.AcquireClient()
	to = to.Add(-24 * time.Hour)
	fromStr := from.Format("02012006")
	toStr := to.Format("02012006")

	period := fromStr + "-" + toStr
	fmt.Println(period)
	agent := client.Get(models.ProductUrl + models.Token + "/" + models.Country + "/" + period + "/json")

	statusCode, resBody, errs := agent.Bytes()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
		return models.ProductCalendarInfo{}
	}

	if statusCode != fiber.StatusOK {
		errorMessage := fmt.Sprintf("received an incorrect response from the product calendar service: %d", statusCode)
		fmt.Println(errorMessage)
		return models.ProductCalendarInfo{}
	}

	var calendar models.ProductCalendarInfo
	if err := json.Unmarshal(resBody, &calendar); err != nil {
		fmt.Println(err.Error())
		return models.ProductCalendarInfo{}
	}

	return calendar
}
