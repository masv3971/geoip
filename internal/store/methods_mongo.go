package store

import (
	"context"
	"geoip/pkg/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func (s *Doc) createLoginEventsIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "userid", Value: bsonx.String("text")}},
		},
		{
			Keys: bsonx.Doc{{Key: "timestamp", Value: bsonx.Int32(1)}},
		},
	}

	opts := options.CreateIndexes().SetMaxTime(20 * time.Second)

	_, err := s.loginEventsCollection.Indexes().CreateMany(ctx, indexes, opts)
	if err != nil {
		return err
	}

	return nil
}

// AddLoginEvent return id, or error
func (s *Doc) AddLoginEvent(ctx context.Context, loginEvent *model.LoginEvent) (interface{}, error) {
	res, err := s.loginEventsCollection.InsertOne(ctx, loginEvent)
	if err != nil {
		return "", err
	}
	return res.InsertedID, nil
}

// GetLoginEvents return all loginEvents for eppn and/or deviceID, or error
func (s *Doc) GetLoginEvents(ctx context.Context, eppn string) (model.LoginEvents, error) {
	filter := bson.M{
		"eppn_hashed": eppn,
	}

	loginEvents := model.LoginEvents{}

	curser, err := s.loginEventsCollection.Find(ctx, filter, &options.FindOptions{})
	if err != nil {
		return nil, err
	}

	defer curser.Close(ctx)
	for curser.Next(ctx) {
		loginEvent := &model.LoginEvent{}
		if err := curser.Decode(loginEvent); err != nil {
			return nil, err
		}
		loginEvents = append(loginEvents, loginEvent)
	}

	return loginEvents, nil
}

// GetLoginEventsAll return map[eppnHashed]LoginEvents
func (s *Doc) GetLoginEventsAll(ctx context.Context) (map[string]model.LoginEvents, error) {
	loginEvents := map[string]model.LoginEvents{}

	curser, err := s.loginEventsCollection.Find(ctx, bson.M{}, &options.FindOptions{})
	if err != nil {
		return nil, err
	}

	defer curser.Close(ctx)
	for curser.Next(ctx) {
		loginEvent := &model.LoginEvent{}
		if err := curser.Decode(loginEvent); err != nil {
			return nil, err
		}
		loginEvents[loginEvent.EppnHashed] = append(loginEvents[loginEvent.EppnHashed], loginEvent)
	}
	if err := curser.Err(); err != nil {
		return nil, err
	}

	return loginEvents, nil
}

//GetLatestLoginEvent return the latest loginEvent associated with user, else return error
func (s *Doc) GetLatestLoginEvent(ctx context.Context, eppn string) (*model.LoginEvent, error) {
	filter := bson.M{
		"eppn": eppn,
	}

	opts := options.FindOne()
	opts.SetSort(bson.M{"timestamp": -1})

	loginEvent := &model.LoginEvent{}
	if err := s.loginEventsCollection.FindOne(ctx, filter, opts).Decode(loginEvent); err != nil {
		return nil, err
	}

	return loginEvent, nil
}

// RemoveLoginEventForUser removes one loginevents corresponding to id
func (s *Doc) RemoveLoginEventForUser(ctx context.Context, eppn string) error {
	filter := bson.M{
		"eppn": eppn,
	}
	_, err := s.loginEventsCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// IsDeviceIDNew return true if the deviceID is new, else false
func (s *Doc) IsDeviceIDNew(ctx context.Context) bool {
	return false
}

func (s *Doc) ping(ctx context.Context) error {
	return s.Mongo.Ping(ctx, readpref.Primary())
}
