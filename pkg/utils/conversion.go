package utils

import (
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConvertUUIDToString преобразует uuid.UUID в строку
func ConvertUUIDToString(id uuid.UUID) string {
	return id.String()
}

// ConvertStringToUUID преобразует строку в uuid.UUID
func ConvertStringToUUID(idStr string) (uuid.UUID, error) {
	if idStr == "" {
		return uuid.UUID{}, errors.New("пустой UUID")
	}
	return uuid.Parse(idStr)
}

func ConvertStringToObjectID(idStr string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(idStr)
}
