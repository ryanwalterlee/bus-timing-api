package src

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"math"

	"github.com/gin-gonic/gin"

	// "log"
	"encoding/json"
)

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

type BusStopTimingFormatted struct {
	ForecastMinutes int32 
	BusId           int32
	ShortName string
}


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

	var responseObject BusStopTiming
	json.Unmarshal(body, &responseObject)

	BusStopTimingFormatted := formatBusStopTiming(responseObject)

	m := structToMap(BusStopTimingFormatted)

	c.JSON(200, m)
}

func formatBusStopTiming(responseObject BusStopTiming) []BusStopTimingFormatted{
	
	var forecastArray []BusStopTimingFormatted
	for i := 0; i < len(responseObject.Forecast); i++ {
		forecast := responseObject.Forecast[i]
		currBusStopTimingFormatted := BusStopTimingFormatted{
			ForecastMinutes: secondsToMinutes(forecast.ForecastSeconds),
			BusId: forecast.BusId,
			ShortName: forecast.BusName.Name,
		}
		forecastArray = append(forecastArray, currBusStopTimingFormatted)
	}
	return forecastArray
}

func secondsToMinutes(seconds float32) int32{
	seconds64 := float64(seconds) / 60
	var minutes int32
	if (seconds <= 0) {
		minutes = 0
	} else {
		minutes = int32(math.Ceil(seconds64))
	}
	return minutes
}

func structToMap(busStopTimingFormattedArray []BusStopTimingFormatted) map[string][]int32{
	m := make(map[string][]int32)
	for i := 0; i < len(busStopTimingFormattedArray); i++ {
		currBusStop := busStopTimingFormattedArray[i].ShortName
		listOfTimings, busStopExists := m[currBusStop]
		if (busStopExists) {
			m[currBusStop] = append(listOfTimings, busStopTimingFormattedArray[i].ForecastMinutes)
		} else {
			var newArray []int32
			m[currBusStop] = append(newArray, busStopTimingFormattedArray[i].ForecastMinutes)
		}
		
	}
	return m
}