package app

import (
	"API/internal/config"
	"API/internal/database"
	"API/internal/handlers"
	"API/internal/kafka"
	"API/internal/repository"
	"API/internal/services"
	"log"
)

type DIContainer struct {
	DB *database.Database

	UserService               *services.UserService
	UserHandler               *handlers.UserHandler
	SegmentService            *services.SegmentService
	SegmentHandler            *handlers.SegmentHandler
	UserSegmentService        *services.UserSegmentService
	UserSegmentHandler        *handlers.UserSegmentHandler
	UserSegmentHistoryService *services.UserSegmentHistoryService
	UserSegmentHistoryHandler *handlers.UserSegmentHistoryHandler

	KafkaProducer *kafka.Producer
	KafkaConsumer *kafka.Consumer
}

func InitDI(cfg config.AppConfig) *DIContainer {
	db := initDatabase(cfg)

	producer, consumer := initKafka(cfg)

	userRepo, segmentRepo, userSegmentRepo, userSegmentHistoryRepo := initRepositories(db)

	userService, segmentService, userSegmentService, userSegmentHistoryService := initServices(
		userRepo,
		segmentRepo,
		userSegmentRepo,
		userSegmentHistoryRepo,
		producer,
	)

	userHandler, segmentHandler, userSegmentHandler, userSegmentHistoryHandler := initHandlers(
		userService,
		segmentService,
		userSegmentService,
		userSegmentHistoryService,
	)

	return &DIContainer{
		DB:                        db,
		UserService:               userService,
		UserHandler:               userHandler,
		SegmentService:            segmentService,
		SegmentHandler:            segmentHandler,
		UserSegmentService:        userSegmentService,
		UserSegmentHandler:        userSegmentHandler,
		UserSegmentHistoryService: userSegmentHistoryService,
		UserSegmentHistoryHandler: userSegmentHistoryHandler,
		KafkaProducer:             producer,
		KafkaConsumer:             consumer,
	}
}

func initDatabase(cfg config.AppConfig) *database.Database {
	db, err := database.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	return db
}

func initKafka(cfg config.AppConfig) (*kafka.Producer, *kafka.Consumer) {
	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal("Could not initialize Kafka producer: ", err)
	}

	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal("Could not initialize Kafka consumer: ", err)
	}

	return producer, consumer
}

func StartTTLConsumer(consumer *kafka.Consumer, service *services.UserSegmentService) {
	go consumer.ConsumeMessages("segment_expiry", service.ProcessTTLExpiryMessage)
}

func initRepositories(db *database.Database) (
	repository.UserRepository,
	repository.SegmentRepository,
	repository.UserSegmentRepository,
	repository.UserSegmentHistoryRepository,
) {
	userRepo := repository.NewUserRepository(db.DB)
	segmentRepo := repository.NewSegmentRepository(db.DB)
	userSegmentHistoryRepo := repository.NewUserSegmentHistoryRepository(db.DB)
	userSegmentRepo := repository.NewUserSegmentRepository(db.DB, userRepo, segmentRepo, userSegmentHistoryRepo)

	return userRepo, segmentRepo, userSegmentRepo, userSegmentHistoryRepo
}

func initServices(
	userRepo repository.UserRepository,
	segmentRepo repository.SegmentRepository,
	userSegmentRepo repository.UserSegmentRepository,
	userSegmentHistoryRepo repository.UserSegmentHistoryRepository,
	producer *kafka.Producer,
) (
	*services.UserService,
	*services.SegmentService,
	*services.UserSegmentService,
	*services.UserSegmentHistoryService,
) {
	userService := services.NewUserService(userRepo)
	segmentService := services.NewSegmentService(segmentRepo)
	userSegmentService := services.NewUserSegmentService(userSegmentRepo, userSegmentHistoryRepo, producer)
	userSegmentHistoryService := services.NewUserSegmentHistoryService(userSegmentHistoryRepo)

	return userService, segmentService, userSegmentService, userSegmentHistoryService
}

func initHandlers(
	userService *services.UserService,
	segmentService *services.SegmentService,
	userSegmentService *services.UserSegmentService,
	userSegmentHistoryService *services.UserSegmentHistoryService,
) (
	*handlers.UserHandler,
	*handlers.SegmentHandler,
	*handlers.UserSegmentHandler,
	*handlers.UserSegmentHistoryHandler,
) {
	userHandler := handlers.NewUserHandler(userService)
	segmentHandler := handlers.NewSegmentHandler(segmentService)
	userSegmentHandler := handlers.NewUserSegmentHandler(userSegmentService)
	userSegmentHistoryHandler := handlers.NewUserSegmentHistoryHandler(userSegmentHistoryService)

	return userHandler, segmentHandler, userSegmentHandler, userSegmentHistoryHandler
}
