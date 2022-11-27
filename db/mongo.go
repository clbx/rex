//MongoDB Wrapper
package db

import (
	"context"
	"log"

	"github.com/clbx/rex/platform"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Holds structures related to the MongoDB
type DB struct {
	client *mongo.Client
}

//Initialize the connection to MongoDB
func InitMongoDB(ctx context.Context, uri string) (*DB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	db := DB{
		client: client,
	}

	return &db, nil
}

//Add a new game to the database, and check for existing entries
func AddGame(db *DB, ctx context.Context, game platform.Game) error {
	rexDatabase := db.client.Database("rex")
	gamesCollection := rexDatabase.Collection("games")
	//Check if file is already in the database.
	filter := bson.D{{"path", game.Path}}
	var readGame platform.Game
	err := gamesCollection.FindOne(ctx, filter).Decode(&readGame)

	//If no document was returned, add the game to the database
	if err == mongo.ErrNoDocuments {
		_, err = gamesCollection.InsertOne(ctx, game)
		log.Printf("[%s][%s] Added to database", game.Platform.Name, game.Name)
		return err
	}

	if err != nil {
		return err
	}

	//If the game was found as its shown in the database, then nothing needs to be done.
	if platform.CompareGames(game, readGame) {
		log.Printf("[%s][%s] Found in database", game.Platform.Name, game.Name)
		return nil

	} else { // If it was not, the old entry should be removed and a new one created
		log.Printf("[%s][%s] Does not match entry found in Database... Overwriting", game.Platform.Name, game.Name)
		_, err := gamesCollection.DeleteOne(ctx, readGame)
		if err != nil {
			log.Fatal(err)
		}
		_, err = gamesCollection.InsertOne(ctx, game)
	}
	return err
}

// Get every game in the database
func GetAllGames(db *DB, ctx context.Context) ([]platform.Game, error) {
	rexDatabase := db.client.Database("rex")
	gamesCollection := rexDatabase.Collection("games")
	cursor, err := gamesCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	allGames := []platform.Game{}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var game platform.Game
		if err = cursor.Decode(&game); err != nil {
			return nil, err
		}
		allGames = append(allGames, game)
	}
	return allGames, nil
}
