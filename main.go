package main

import (
	// "fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Name string `json:"name"`
}

type Table struct {
	TableID int8   `json:"tableid" gorm:"primaryKey"`
	Phone   string `json:"phone" gorm:"size:255"`
	Status  uint8  `json:"status"`
	Code    string `json:"code" gorm:"size:255"`
	StartTime string `json:"startTime" gorm:"size:255"`
	EndTime string `json:"endTime" gorm:"size:255"`
}

var connections []*websocket.Conn
var lock = sync.Mutex{}

var db *gorm.DB

func main() {
	app := fiber.New()
	app.Use(cors.New())

	dsn := "root:1234@tcp(localhost:3306)/testgocon?parseTime=true"
	dbs, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = dbs

	if err != nil {
		panic(err)
	}
	app.Put("/updatestate/:id",UpdateRow)
	app.Get("/gettables",GetTable)

	app.Get("/ws", websocket.New(MainWebsocket))

	app.Listen(":8080")
}


