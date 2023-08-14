package bsonuuid

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

func BuildRegistry() *bsoncodec.Registry {
	registry := bson.NewRegistry()
	registry.RegisterTypeEncoder(TypeUUID, bsoncodec.ValueEncoderFunc(UUIDEncodeValue))
	registry.RegisterTypeDecoder(TypeUUID, bsoncodec.ValueDecoderFunc(UUIDDecodeValue))
	return registry
}
