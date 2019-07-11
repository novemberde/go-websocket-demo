package api

import (
	"log"
	"math/rand"
)

func createRoom() {
	roomID := createRoomID(16)
	log.Println(roomID)

	// TODO: Be sure creating a unique chat room id.
}

func createRoomID(strLen int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	b := make([]byte, strLen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
