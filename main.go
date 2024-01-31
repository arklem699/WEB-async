package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	MyToken = "access_token"
)

type DataArray struct {
	Data []DataStruct `json:"data"`
}

type DataStruct struct {
	ID      int `json:"id"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Doctor  string `json:"doctor"`
	ID_appl int    `json:"id_appl"`
	Status  string `json:"status"`
}

type Result struct {
	ID_appl    int `json:"id_appl"`
	ID_appoint    int `json:"id_appoint"`
	Was   bool   `json:"was"`
	Token string `json:"token"`
}

func randomStatus() bool {
	time.Sleep(5 * time.Second)
	return rand.Intn(2) == 0
}

func SendStatus(id_appl int, url string, id_appoint int, data DataStruct) {
	result := randomStatus()

	dataResult := Result{ID_appl: id_appl, ID_appoint: id_appoint, Was: result, Token: MyToken}
	_, err := performPUTRequest(url, dataResult)
	if err != nil {
		fmt.Println("Error sending status:", err)
		return
	}

	fmt.Println("Status sent successfully for id_appl:", id_appl, "and id_appoint:", id_appoint, "-", result)
}

func performPUTRequest(url string, data Result) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending status:", err)
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	r.POST("/was/", func(c *gin.Context) {
		var dataArray DataArray

		if err := c.ShouldBindJSON(&dataArray); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		var wg sync.WaitGroup
		for _, data := range dataArray.Data {
			wg.Add(1)
			go func(data DataStruct) {
				defer wg.Done()
				SendStatus(data.ID_appl, "http://127.0.0.1:8000/appapp/async/put/", data.ID, data)
			}(data)
		}

		wg.Wait()

		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated for all IDs"})
	})

	r.Run("localhost:9000")
}