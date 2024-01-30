package main

import (
	"math/rand"
	"time"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

const (
	MyToken = "access_token"
)

// Ваша структура с булевской переменной
type Result struct {
	Was bool `json:"was"`
	Token string `json:"token"`
}

// Функция для генерации случайного статуса
func randomStatus() bool {
	time.Sleep(5 * time.Second) // Задержка на 5 секунд
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 0
}

// Функция для отправки статуса в отдельной горутине
func SendStatus(id string, url string) {
	// Выполнение расчётов с randomStatus
	result := randomStatus()

	// Отправка PUT-запроса к основному серверу
	data := Result{Was: result, Token: MyToken}
	_, err := performPUTRequest(url, data)
	if err != nil {
		fmt.Println("Error sending status:", err)
		return
	}

	fmt.Println("Status sent successfully for id:", id)
}

func performPUTRequest(url string, data Result) (*http.Response, error) {
	// Сериализация структуры в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Создание PUT-запроса
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}

func main() {
	// Создание роутера Gin
	r := gin.Default()

	// Добавление middleware для обработки CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Укажите адрес вашего клиента
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	// Обработчик POST-запроса для set_status
	r.POST("/was/:id/", func(c *gin.Context) {
		// Получение значения "id" из параметра запроса
		id := c.Param("id")

		// Запуск горутины для отправки статуса
		go SendStatus(id, "http://127.0.0.1:8000/application/" + id + "/async/put/")

		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated"})
	})

	// Запуск сервера
	r.Run("localhost:9000")
}