package api

import "time"

// Room Room table on DynamoDB
// partitionkey is stuctured like "ROOM|room_id", "ROOM|room_id"
type Room struct {
	RoomID    string    `dynamo:"room_id"`    // Hash key, a.k.a. partition key(pk)
	CreatedAt time.Time `dynamo:"created_at"` // Range key, a.k.a. sort key(sk)
}

// Message Message table on DynamoDB
type Message struct {
	RoomID    string    `dynamo:"room_id"`    // pk
	CreatedAt time.Time `dynamo:"created_at"` // sk
	SenderID  int       `dynamo:"sender_id"`  // user_id
	Content   string    `dynamo:"content"`    // Message content. string
	// ArticleID int       `dynamo:"article_id"` // referenced article. nullable
	// ImageURLs []string  `dynamo:"image_urls"`
}

// UserRoom UserRoom table on DynamoDB
// Get the room list of certain user.
type UserRoom struct {
	UserID    int       `dynamo:"user_id"`    // pk
	CreatedAt time.Time `dynamo:"created_at"` // sk
	Rooms     []string  `dynamo:"rooms"`      // Lists(Sorted list)
	ArticleID int       `dynamo:"article_id"`
}

// Connection Manage user connections
// TODO: Apply ttl on expired_at
type Connection struct {
	UserID       int       `dynamo:"user_id"`       // pk
	ConnectionID string    `dynamo:"connection_id"` // sk
	CreatedAt    time.Time `dynamo:"created_at"`
	ExpiredAt    time.Time `dynamo:"expired_at"`
}
