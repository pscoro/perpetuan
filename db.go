package main

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type dbConnection struct {
	settingCollection *mongo.Collection
	gridCollection    *mongo.Collection
	userCollection	  *mongo.Collection
	ctx               context.Context
	cancel            context.CancelFunc
}

func initDB() dbConnection {
	var connection dbConnection
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	usersDB, err := mongo.Connect(connection.ctx, options.Client().ApplyURI(
		os.Getenv("USER_DB_URI"),
	))
	if err != nil {
		log.Fatal(err)
	}

	connection.ctx, connection.cancel = context.WithTimeout(context.Background(), 100*time.Minute)
	defer connection.cancel()
	regionsDB, err := mongo.Connect(connection.ctx, options.Client().ApplyURI(
		os.Getenv("USER_DB_URI"),
	))
	if err != nil {
		log.Fatal(err)
	}

	connection.settingCollection = regionsDB.Database("regiondata").Collection("regions")
	connection.gridCollection = regionsDB.Database("regiondata").Collection("gridlocations")
	connection.userCollection = usersDB.Database("userdata").Collection("users")

	log.Println("DB connected")

	return connection
}

type userData struct {
	Username string
	Email string
	Password string
	Player player
}
