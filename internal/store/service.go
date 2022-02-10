package store

import (
	"context"
	"geoip/pkg/logger"
	"geoip/pkg/model"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	cfg *model.Cfg
	log *logger.Logger
	KV  *KV
	Doc *Doc
}

// Doc holds document databases
type Doc struct {
	Mongo                 *mongo.Client
	usersCollection       *mongo.Collection
	logineventsCollection *mongo.Collection
	deviceIDCollection    *mongo.Collection
	mailMsgCollection     *mongo.Collection
}

// KV redis storage object
type KV struct {
	Redis *redis.Client
}

// New creates a new instance of store
func New(ctx context.Context, cfg *model.Cfg, log *logger.Logger) (*Service, error) {
	s := &Service{
		cfg: cfg,
		log: log,
	}

	if err := s.newKV(ctx); err != nil {
		return nil, err
	}

	if err := s.newDoc(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) newDoc(ctx context.Context) error {
	ctx, _ = context.WithTimeout(ctx, 20*time.Second)

	//	credential := options.Credential{
	//		AuthMechanism:           "",
	//		AuthMechanismProperties: map[string]string{},
	//		AuthSource:              "",
	//		Username:                "eduid_geoip",
	//		Password:                "eduid_geoip_pw",
	//		PasswordSet:             false,
	//	}

	clientOptions := options.Client().ApplyURI(s.cfg.Mongodb.Addr)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	s.Doc = &Doc{
		Mongo: client,
	}

	s.Doc.logineventsCollection = s.Doc.Mongo.Database("eduid_geoip").Collection("loginevents")
	s.Doc.mailMsgCollection = s.Doc.Mongo.Database("eduid_mail").Collection("xxx")

	if err := s.Doc.createLoginEventsIndexes(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) newKV(ctx context.Context) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr: s.cfg.KVStorage.Redis.Addr,
		DB:   s.cfg.KVStorage.Redis.DB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return err
	}

	kv := &KV{
		Redis: redisClient,
	}

	s.KV = kv

	return nil
}
