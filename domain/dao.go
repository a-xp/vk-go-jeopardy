package domain

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
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
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectStr.String()).SetMonitor(cmdMonitor))
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

func (dao *AppDAO) loadGameHeaders() ([]*GameHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	games := dao.client.Database(dao.name).Collection("games")
	opts := options.Find().SetProjection(bson.M{"_id": 1, "name": 1, "active": 1, "new": 1})
	cur, err := games.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []*GameHeader
	if err = cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (dao *AppDAO) loadActiveGames() ([]*Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	games := dao.client.Database(dao.name).Collection("games")
	cur, err := games.Find(ctx, bson.D{{"active", true}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []*Game
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

func (dao AppDAO) getGameSession(userId int64, gameId *string) (*Answer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	answers := dao.client.Database(dao.name).Collection("answers")
	var result Answer
	err := answers.FindOne(ctx, bson.D{{"gameId", gameId}, {"userId", userId}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (dao AppDAO) storeGameSession(answer *Answer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	answers := dao.client.Database(dao.name).Collection("answers")
	if answer.Id == nil {
		newId := fmt.Sprintf("%s_%d", *answer.GameId, answer.UserId)
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

func (dao AppDAO) getUserRating(gameId *string, userId int64) *RatingEntry {
	answer, err := dao.getGameSession(userId, gameId)
	if err != nil || answer == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := dao.client.Database(dao.name).Collection("answers").CountDocuments(ctx,
		bson.M{"$or": bson.A{
			bson.M{"gameId": gameId, "score": bson.M{"$gt": answer.Score}},
			bson.M{"gameId": gameId, "score": answer.Score, "completeTime": bson.M{"$lt": answer.CompleteTime}}}})
	if err != nil {
		return nil
	}
	return &RatingEntry{
		Pos:    int(count) + 1,
		UserId: answer.UserId,
		Score:  answer.Score,
	}
}

func (dao AppDAO) getGameTop(gameId *string, limit int) ([]*RatingEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	answers := dao.client.Database(dao.name).Collection("answers")
	cur, err := answers.Aggregate(ctx,
		bson.A{
			bson.M{"$match": bson.M{"gameId": gameId}},
			bson.M{"$sort": bson.M{"score": -1, "completeTime": 1}},
			bson.M{"$limit": limit},
			bson.M{"$lookup": bson.M{
				"from":         "users",
				"localField":   "userId",
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
			bson.M{"$project": bson.M{"img": 1, "_id": 0, "name": 1, "lastname": 1, "score": 1, "userId": 1}},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []*RatingEntry
	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	for i, _ := range result {
		result[i].Pos = i + 1
	}
	return result, nil
}

func (dao AppDAO) getAdminUsers() ([]*AdminUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	admins := dao.client.Database(dao.name).Collection("admins")
	cur, err := admins.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []*AdminUser
	if err = cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (dao AppDAO) getGameById(gameId *string) (*Game, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	groups := dao.client.Database(dao.name).Collection("games")
	oId, err := primitive.ObjectIDFromHex(*gameId)
	if err != nil {
		return nil, false
	}
	var result Game
	err = groups.FindOne(ctx, bson.D{{"_id", oId}}).Decode(&result)
	if err != nil {
		return nil, false
	}
	return &result, true
}

func (dao AppDAO) removeAdmin(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	admins := dao.client.Database(dao.name).Collection("admins")
	_, err := admins.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (dao AppDAO) addAdmin(user *AdminUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	admins := dao.client.Database(dao.name).Collection("admins")
	opts := options.Replace().SetUpsert(true)
	_, err := admins.ReplaceOne(ctx, bson.M{"_id": user.Id}, user, opts)
	return err
}

func (dao AppDAO) getGroups() ([]*Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	groups := dao.client.Database(dao.name).Collection("groups")
	cur, err := groups.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []*Group
	if err = cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (dao AppDAO) removeGroup(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	groups := dao.client.Database(dao.name).Collection("groups")
	_, err := groups.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (dao AppDAO) addGroup(group *Group) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	groups := dao.client.Database(dao.name).Collection("groups")
	opts := options.Replace().SetUpsert(true)
	_, err := groups.ReplaceOne(ctx, bson.M{"_id": group.Id}, group, opts)
	return err
}

func (dao AppDAO) storeGame(game *Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	games := dao.client.Database(dao.name).Collection("games")
	var err error
	if game.Id != nil {
		oId, err := primitive.ObjectIDFromHex(*game.Id)
		if err != nil {
			return err
		}
		game.Id = nil
		_, err = games.ReplaceOne(ctx, bson.M{"_id": oId}, game)
	} else {
		_, err = games.InsertOne(ctx, game)
	}
	return err
}
