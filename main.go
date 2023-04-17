package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HashTable struct {
	data [][][]string
}

func NewHashTable(size int) *HashTable {
	return &HashTable{
		data: make([][][]string, size),
	}
}

func (h *HashTable) _hash(key string) int {
	hash := 0
	for i := 0; i < len(key); i++ {
		hash = (hash + int(key[i])*i) % len(h.data)
	}
	return hash
}

func (h *HashTable) set(key string, value string) [][][]string {
	address := h._hash(key)
	if h.data[address] == nil {
		h.data[address] = make([][]string, 0)
	}
	h.data[address] = append(h.data[address], []string{key, value})
	return h.data
}

func (h *HashTable) get(key string) string {
	address := h._hash(key)
	currentBucket := h.data[address]
	if currentBucket != nil {
		for i := 0; i < len(currentBucket); i++ {
			if currentBucket[i][0] == key {
				return currentBucket[i][1]
			}
		}
	}
	return ""
}

var cache HashTable = *NewHashTable(1000)

func getMongoClient() (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://haniif02aj:uQz7VLcY70SoS4JA@test.dkwiku6.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %s", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %s", err)
	}
	return client, nil
}

///////////////////

type Room struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Description  string             `bson:"description"`
	Price        float64            `bson:"price"`
	IsAvailable  bool               `bson:"isAvailable"`
	NumberOfBeds int                `bson:"numberOfBeds"`
}

func createRoom(c *gin.Context) {
	client, err := getMongoClient()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	collection := client.Database("hotel_booking").Collection("rooms")

	var room Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := collection.InsertOne(context.Background(), &room)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	room.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, room)
}

func getRooms(c *gin.Context) {
	client, err := getMongoClient()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	collection := client.Database("hotel_booking").Collection("rooms")

	var rooms []Room
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := cursor.All(context.Background(), &rooms); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func getRoomByID(c *gin.Context) {

	client, err := getMongoClient()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	collection := client.Database("hotel_booking").Collection("rooms")

	roomID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var room Room
	err = collection.FindOne(context.Background(), bson.M{"_id": roomID}).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, room)
}

func updateRoom(c *gin.Context) {
	client, err := getMongoClient()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	collection := client.Database("hotel_booking").Collection("rooms")
	roomID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var room Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": roomID}, bson.M{"$set": &room})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteRoom(c *gin.Context) {
	client, err := getMongoClient()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	collection := client.Database("hotel_booking").Collection("rooms")
	roomID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": roomID})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

//////////////////////////////////////////

type Reservation struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	CustomerName string             `bson:"customerName"`
	RoomId       string             `bson:"roomId"`
	StayingDays  int                `bson:"days"`
	StartDate    primitive.DateTime `bson:"startDate"`
	TotalPrice   float64            `bson:"totalPrice"`
}

func createReservation(c *gin.Context) {
	client, err := getMongoClient()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	collection := client.Database("hotel_booking").Collection("reservations")

	var reservation Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := collection.InsertOne(context.Background(), &reservation)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	reservation.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, reservation)
}

func getReservations(c *gin.Context) {
	client, err := getMongoClient()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	collection := client.Database("hotel_booking").Collection("reservations")

	var reservations []Reservation
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := cursor.All(context.Background(), &reservations); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, reservations)
}

func getReservationByID(c *gin.Context) {

	client, err := getMongoClient()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	collection := client.Database("hotel_booking").Collection("reservations")

	reservationID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var reservation Reservation
	err = collection.FindOne(context.Background(), bson.M{"_id": reservationID}).Decode(&reservation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func updateReservation(c *gin.Context) {
	client, err := getMongoClient()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	collection := client.Database("hotel_booking").Collection("reservations")
	reservationID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var reservation Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": reservationID}, bson.M{"$set": &reservation})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteReservation(c *gin.Context) {
	client, err := getMongoClient()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	collection := client.Database("hotel_booking").Collection("reservations")
	reservationID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": reservationID})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

///////////////////////////////

func main() {
	r := gin.Default()

	r.POST("/rooms", createRoom)
	r.GET("/rooms", getRooms)
	r.GET("/rooms/:id", getRoomByID)
	r.PUT("/rooms/:id", updateRoom)
	r.DELETE("/rooms/:id", deleteRoom)

	r.POST("/reservations", createReservation)
	r.GET("/reservations", getReservations)
	r.GET("/reservations/:id", getReservationByID)
	r.PUT("/reservations/:id", updateReservation)
	r.DELETE("/reservations/:id", deleteReservation)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
