package database

import (
	"book-management-system/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoDB *mongo.Client

// InitMongoDB инициализирует подключение к MongoDB
func InitMongoDB() {
	uri := config.GetEnv("MONGO_URI", "mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}

	// Проверяем подключение
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Ошибка пинга MongoDB: %v", err)
	}

	log.Println("Подключение к MongoDB успешно установлено")

	MongoDB = client
}
