package domain

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"goj/configuration"
	"log"
	"net/url"
	"time"
)

var DAO *AppDAO

func initDAO(cfg *configuration.Configuration) *AppDAO {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	connectStr := url.URL{
		Scheme:   "mongodb+srv",
		User:     url.UserPassword(cfg.Mongo.User, cfg.Mongo.Pass),
		Host:     cfg.Mongo.Host,
		Path:     cfg.Mongo.Name,
		RawQuery: cfg.Mongo.Options,
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectStr.String()))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	return &AppDAO{
		client: client,
		name:   cfg.Mongo.Name,
	}
}

func (dao *AppDAO) loadGames() ([]Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	games := dao.client.Database(dao.name).Collection("games")
	cur, err := games.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []Game
	if err = cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (dao *AppDAO) close() {
	if dao.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := dao.client.Disconnect(ctx)
		if err != nil {
			log.Printf("Error when disconneting Mongo %+v", err)
		}
	}
}

func (dao *AppDAO) findGroupById(id int64) (*Group, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	groups := dao.client.Database(dao.name).Collection("groups")
	var result Group
	err := groups.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return nil, false
	}
	return &result, true
}

func (dao AppDAO) findUserById(id int64) (*User, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	users := dao.client.Database(dao.name).Collection("users")
	var result User
	err := users.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return nil, false
	}
	return &result, true
}

func (dao AppDAO) storeUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	users := dao.client.Database(dao.name).Collection("users")
	opts := options.Replace().SetUpsert(true)
	_, err := users.ReplaceOne(ctx, bson.D{{"_id", user.Id}}, user, opts)
	return err
}

func (dao AppDAO) getGameSession(userId int64, gameId string) (*Answer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	answers := dao.client.Database(dao.name).Collection("answers")
	var result Answer
	err := answers.FindOne(ctx, bson.D{{"gameid", gameId}, {"userid", userId}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (dao AppDAO) storeGameSesion(answer *Answer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	answers := dao.client.Database(dao.name).Collection("answers")
	if answer.Id == nil {
		newId := fmt.Sprintf("%s_%d", answer.GameId, answer.UserId)
		answer.Id = &newId
	}
	opts := options.Replace().SetUpsert(true)
	_, err := answers.ReplaceOne(ctx, bson.D{{"_id", answer.Id}}, answer, opts)
	return err
}

type AppDAO struct {
	client *mongo.Client
	name   string
}

func (dao AppDAO) getGameTop(gameId string, limit int) ([]RatingEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	answers := dao.client.Database(dao.name).Collection("answers")
	cur, err := answers.Aggregate(ctx,
		bson.A{
			bson.M{"$match": bson.M{"gameid": gameId}},
			bson.M{"$score": bson.M{"score": -1}},
			bson.M{"$limit": limit},
			bson.M{"$lookup": bson.M{
				"from":         "users",
				"localField":   "userid",
				"foreignField": "_id",
				"as":           "user",
			}},
			bson.M{
				"$replaceRoot": bson.M{
					"newRoot": bson.M{
						"$mergeObjects": bson.A{
							bson.M{"$arrayElemAt": bson.A{"$user", 0}},
							"$$ROOT",
						},
					},
				},
			},
			bson.M{"$project": bson.M{"img": 1, "_id": 0, "name": 1, "lastname": 1, "score": 1, "userid": 1}},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []RatingEntry
	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (dao AppDAO) getAdminUsers() ([]AdminUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	admins := dao.client.Database(dao.name).Collection("admins")
	cur, err := admins.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []AdminUser
	if err = cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (dao AppDAO) getGameById(gameId string) (*Game, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	groups := dao.client.Database(dao.name).Collection("games")
	var result Game
	err := groups.FindOne(ctx, bson.D{{"_id", gameId}}).Decode(&result)
	if err != nil {
		return nil, false
	}
	return &result, true
}
