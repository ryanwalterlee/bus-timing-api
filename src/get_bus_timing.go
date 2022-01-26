package src

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	// "log"
	"encoding/json"
)

func GetBusTiming(c *gin.Context) {
	busId := c.Query("bus-id")
	url := fmt.Sprintf("https://dummy.uwave.sg/busstop/%s", busId)
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	type BusName struct {
		Name string `json:"short_name"`
	}

	type Forecast struct {
		ForecastSeconds float32 `json:"forecast_seconds"`
		BusId           int32   `json:"rv_id"`
		BusName BusName `json:"route"`
	}

	type BusStopTiming struct {
		Forecast []Forecast `json:"forecast"`
	}

	var responseObject BusStopTiming
	json.Unmarshal(body, &responseObject)

	c.JSON(200, gin.H{"forecast": responseObject.Forecast})
}