package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoClient() (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://haniif02aj:uQz7VLcY70SoS4JA@test.dkwiku6.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
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

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
