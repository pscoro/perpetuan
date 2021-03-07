package main

import (
	"context"
	//"context"
	//"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	//"time"
)

func getGrid(name string) (grid [][]gridLocation, err error) {

	gridRes := connect.gridCollection.FindOne(context.Background(), bson.M{"_id": name})
	if gridRes.Err() != nil {
		log.Println("Find error: ", gridRes.Err())
		return nil, gridRes.Err()
	}

	BSONData := struct {
		Setting string           `bson:"_id"`
		Grid    [][]gridLocation `bson:"Grid"`
	}{}

	decodeError := gridRes.Decode(&BSONData)

	if decodeError != nil {
		log.Println("Decode error: ", decodeError)
		return nil, decodeError
	}

	return BSONData.Grid, nil
}

func getSetting(name string) (set setting, err error) {

	setRes := connect.settingCollection.FindOne(context.Background(), bson.M{"_id": name})
	if setRes.Err() != nil {
		log.Println("Find error: ", setRes.Err())
		return setting{}, setRes.Err()
	}

	BSONData := struct {
		name    string  `bson:"_id"`
		Setting setting `bson:"Setting"`
	}{}

	decodeError := setRes.Decode(&BSONData)

	if decodeError != nil {
		log.Println("Decode error: ", decodeError)
		return setting{}, decodeError
	}

	return BSONData.Setting, nil
}

func addGrid(name string, grid [][]gridLocation) error {

	res, err := connect.gridCollection.InsertOne(context.Background(), bson.M{"_id": name, "Grid": grid})
	if err != nil {
		log.Println("Find error: ", err)
		return err
	} else {
		log.Println("Grid for: ", res, " created")
		return nil
	}
}

func addSetting(name string, portals []portal, width int, height int, grid [][]gridLocation) error {
	err := addGrid(name, grid)

	if err != nil {
		log.Fatal(err)
		return err
	} else {
		res, err := connect.settingCollection.InsertOne(context.Background(), bson.M{"_id": name, "Setting": setting{
			key:       name,
			portals:   portals,
			sizeX:     width,
			sizeY:     height,
			imagePath: "./api/assets/img/regions/" + name + ".png",
			grid:      grid,
		}})
		if err != nil {
			log.Println("Insert error: ", err)
			return err
		} else {
			log.Println("Setting for: ", res, " created")
			return nil
		}
	}
}

func updateSettingList(old []setting) ([]setting, bool) {
	cur, err := connect.settingCollection.Find(context.Background(), bson.D{})

	var settingList []setting

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	BSONData := struct {
		name    string  `bson:"_id"`
		Setting setting `bson:"Setting"`
	}{}

	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		err := cur.Decode(&BSONData)
		if err != nil {
			log.Fatal(err)
			return old, false
		}
		// do something with result...
		settingList = append(settingList, BSONData.Setting)
	}
	if err := cur.Err(); err != nil {
		return old, false
	} else {
		return settingList, true
	}

}
