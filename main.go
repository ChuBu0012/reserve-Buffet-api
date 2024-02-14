package main

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

type User struct {
	Name string `json:"name"`
}

var usr *User
var connections []*websocket.Conn
var lock = sync.Mutex{}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/getname", func(c *fiber.Ctx) error {
		return c.JSON(usr)
	})

	app.Post("/setname", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		usr = user
		notifyAllConnections(usr.Name)
		return c.SendString("Update Name Successful!!")
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		lock.Lock()
		connections = append(connections, c)
		lock.Unlock()

		defer func() {
			lock.Lock()
			for i, conn := range connections {
				if conn == c {
					connections = append(connections[:i], connections[i+1:]...)
					break
				}
			}
			lock.Unlock()
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}
			fmt.Printf("Received message: %s\n", msg)

			notifyAllConnections(string(msg))
		}
	}))

	app.Listen(":8080")
}

func notifyAllConnections(msg string) {
	lock.Lock()
	defer lock.Unlock()

	for _, conn := range connections {
		if err := conn.WriteMessage(1, []byte(msg)); err != nil {
			fmt.Println("Error writing message:", err)
		}
	}
}
