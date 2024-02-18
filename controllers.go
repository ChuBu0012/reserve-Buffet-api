package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func MainWebsocket(c *websocket.Conn) {
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
			return
		}

		NotifyAllConnections(string(msg))
	}
}

func NotifyAllConnections(msg string) {
	lock.Lock()
	defer lock.Unlock()

	for _, conn := range connections {
		if err := conn.WriteMessage(1, []byte(msg)); err != nil {
		}
	}
}

func UpdateRow(c *fiber.Ctx) error {
	tableId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	tableUpdate := new(Table)

	if err := c.BodyParser(tableUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	tableUpdate.TableID = int8(tableId)

	result := db.Save(&tableUpdate)
	if result.Error != nil {
		log.Fatalf("Update book failed %v", result.Error)
	}
	tableUpdateJSON, err := json.Marshal(tableUpdate)
	if err != nil {
		log.Fatalf("Failed to convert object to JSON: %v", err)
	}

	// NotifyAllConnections("Row with ID-" + strconv.Itoa(tableId) + "-has been updated.")
	NotifyAllConnections(string(tableUpdateJSON))

	return c.SendString("Update row successfully")
}
func GetTable(c *fiber.Ctx) error {
	var tables []Table
	
	result := db.Find(&tables)
	if result.Error != nil {
		log.Fatalf("Update book failed %v", result.Error)
	}

	return c.JSON(&tables)

}
