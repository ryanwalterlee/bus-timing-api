package src

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	// "log"
	"encoding/json"
)

type BusStopLocation struct {
	VehiclInfo []VehicleInfo `json:"vehicles"` 
}

type VehicleInfo struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
	PlateNumber string `json:"registration_code"`
}


func GetBusLocation(c *gin.Context) {
	busId := c.Query("bus-id")
	url := fmt.Sprintf("https://dummy.uwave.sg/busline/%s", busId)
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject BusStopLocation
	json.Unmarshal(body, &responseObject)

	c.JSON(200, responseObject)
}