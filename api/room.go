package api

import (
	"log"
	"math/rand"
	"time"
)

var (
	rooms = []Room{}
)

func joinRoom(c *Client, room *Room) {
	room.Users = append(room.Users, c)
}

func createRoom() *Room {
	roomID := getRandomString(24)
	log.Println(roomID)
	room := Room{
		RoomID:    roomID,
		CreatedAt: time.Now(),
	}
	rooms = append(rooms, room)

	log.Println(len(rooms))
	// TODO: Be sure creating a unique chat room id. Check the id on database

	return &room
}

func getRandomString(strLen int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	b := make([]byte, strLen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
