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

type Table struct {
	TableID   int8   `json:"tableid" gorm:"primaryKey"`
	Phone     string `json:"phone" gorm:"size:255"`
	Status    uint8  `json:"status"`
	Code      string `json:"code" gorm:"size:255"`
	StartTime string `json:"startTime" gorm:"size:255"`
	EndTime   string `json:"endTime" gorm:"size:255"`
}

var connections []*websocket.Conn
var lock = sync.Mutex{}

var db *gorm.DB

func main() {
	app := fiber.New()
	app.Use(cors.New())

	dsn := "ur1xpked8rzvzaan:wTUXKsx9Wr05B9NSF874@tcp(bovpwi5ezwdtd5jhp1ab-mysql.services.clever-cloud.com:3306)/bovpwi5ezwdtd5jhp1ab?parseTime=true"
	dbs, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = dbs

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Table{})

	app.Put("/updatestate/:id", UpdateRow)
	app.Get("/gettables", GetTable)

	app.Get("/ws", websocket.New(MainWebsocket))

	app.Listen(":8080")
}
