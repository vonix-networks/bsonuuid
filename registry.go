package bsonuuid

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

func BuildRegistry() *bson.Registry {
	registry := bson.NewRegistry()
	registry.RegisterTypeEncoder(TypeUUID, bson.ValueEncoderFunc(UUIDEncodeValue))
	registry.RegisterTypeDecoder(TypeUUID, bson.ValueDecoderFunc(UUIDDecodeValue))
	return registry
}
