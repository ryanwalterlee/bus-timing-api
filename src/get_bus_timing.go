package src

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"math"
	"time"

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

type ForecastFormatted struct {
	ForecastArray []int32
	BusId int32
}

type LastCheckedTime struct {
	Time time.Time
	BusTimingMap map[string]ForecastFormatted
}

var mapCache map[string]LastCheckedTime = make(map[string]LastCheckedTime)


func GetBusTiming(c *gin.Context) {

	busId := c.Query("bus-stop-id")
	lastCheckedTime , entryExists := mapCache[busId]

	// if map entry exist and is not outdated, return it
	if (entryExists && time.Now().Sub(lastCheckedTime.Time) < 60*1000*1000*1000) {
		c.JSON(200, lastCheckedTime.BusTimingMap)

	} else { // else call the api and reupdate map
		apiCall(c, busId)
	}
	
}

func apiCall(c *gin.Context, busId string) {
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

	fmt.Println(responseObject)

	BusStopTimingFormatted := formatBusStopTiming(responseObject)

	fmt.Println(BusStopTimingFormatted)

	m := structToMap(BusStopTimingFormatted)

	// update cache
	mapCache[busId] = LastCheckedTime{time.Now(), m}

	c.JSON(200, m)
}

// place all fields on the same level
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

// convert seconds to minutes with error handling
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

// turn the struct into a map with an array of forecast timings for each bus service
func structToMap(busStopTimingFormattedArray []BusStopTimingFormatted) map[string]ForecastFormatted{
	m := make(map[string]ForecastFormatted)
	for i := 0; i < len(busStopTimingFormattedArray); i++ {
		currBusStop := busStopTimingFormattedArray[i].ShortName
		forecastFormatted, busStopExists := m[currBusStop]
		
		if (busStopExists) {
			listOfTimings := forecastFormatted.ForecastArray
			m[currBusStop] = ForecastFormatted{
				append(listOfTimings, busStopTimingFormattedArray[i].ForecastMinutes),
				busStopTimingFormattedArray[i].BusId,
			}
		} else {
			var newArray []int32
			m[currBusStop] = ForecastFormatted{
				append(newArray, busStopTimingFormattedArray[i].ForecastMinutes),
				busStopTimingFormattedArray[i].BusId,
			}
		}
		
	}
	return m
}