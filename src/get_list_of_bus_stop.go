package src

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	// "log"
	"encoding/json"
)

type BusStopInfo struct {
	BusStopName string `json:"name"`
	BusId       int32   `json:"id"`
}

var listOfBusStopNames []BusStopInfo

func GetListOfBusStop(c *gin.Context) {

	if len(listOfBusStopNames) > 0 {
		c.JSON(200, gin.H{"listOfBusStopNames": listOfBusStopNames})
	} else {
		listOfBusStopId := []string{"378204", "383050", "378202", "383049", "382998", "378237", "378233", "378230",
			"378229", "378228", "378227", "382995", "378224", "378226", "383010", "383009",
			"383006", "383004", "378234", "383003", "378222", "383048", "378203", "382999",
			"378225", "383014", "383013", "383011", "377906", "383018", "383015", "378207"}

		for i := 0; i < len(listOfBusStopId); i++ {
			busId := listOfBusStopId[i]
			url := fmt.Sprintf("https://dummy.uwave.sg/busstop/%s", busId)
			response, err := http.Get(url)

			if err != nil {
				fmt.Print(err.Error())
			}

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Print(err.Error())
			}

			type Response struct {
				Name string `json:"name"`
				Id   int32  `json:"id"`
			}

			var responseObject Response
			json.Unmarshal(body, &responseObject)

			listOfBusStopNames = append(listOfBusStopNames, BusStopInfo{responseObject.Name, responseObject.Id})

		}
		c.JSON(200, gin.H{"listOfBusStopNames": listOfBusStopNames})
	}
}