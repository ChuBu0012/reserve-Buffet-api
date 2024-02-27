package main

import (
	"encoding/json"
	// "fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func MainWebsocket(c *websocket.Conn) {
	lock.Lock()
	connections = append(connections, c)
	lock.Unlock()

	defer func() {
		lock.Lock()
		defer lock.Unlock()

		for i, conn := range connections {
			if conn == c {
				connections = append(connections[:i], connections[i+1:]...)
				break
			}
		}
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
			log.Printf("Write Message Unsuccessful: %v", err)
		}
	}
}

func UpdateRow(c *fiber.Ctx) error {
	tableID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	tableUpdate := new(Table)
	if err := c.BodyParser(tableUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var checkLimit []Table
	db.Where("phone = ?", tableUpdate.Phone).Find(&checkLimit)
	if len(checkLimit) == 2 {
		return c.SendString("Limited to 2 phone numbers per table.")
	}

	tableUpdate.TableID = int8(tableID)
	if tableUpdate.Status == 1 {
		tableUpdate.Code = uuid.New().String()
	} else if tableUpdate.Status == 0 {
		tableUpdate.Code = ""
	}

	result := db.Save(&tableUpdate)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update table.")
	}

	tableUpdateJSON, err := json.Marshal(tableUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to convert object to JSON.")
	}

	NotifyAllConnections(string(tableUpdateJSON))
	return c.SendString("Code: " + tableUpdate.Code)
}

func GetTable(c *fiber.Ctx) error {
	var tables []Table
	result := db.Find(&tables)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve tables.")
	}

	return c.JSON(&tables)
}
