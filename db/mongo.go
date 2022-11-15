package db

//MongoDB wrapper
import (
	"context"
	"log"

	"github.com/clbx/rex/platform"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

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

func AddGame(db *DB, ctx context.Context, game platform.Game) error {
	rexDatabase := db.client.Database("rex")
	gamesCollection := rexDatabase.Collection("games")
	//Check if file is already in the database.

	filter := bson.D{{"path", game.Path}}
	var readGame platform.Game
	err := gamesCollection.FindOne(ctx, filter).Decode(&readGame)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", readGame)
	if readGame == game {
		log.Printf("%s found in database", game.Name)
		return nil
	}
	_, err = gamesCollection.InsertOne(ctx, game)
	log.Printf("Added %s to database", game.Name)
	return err
}

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
