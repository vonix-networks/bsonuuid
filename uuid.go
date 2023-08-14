package bsonuuid

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

var TypeUUID = reflect.TypeOf(uuid.UUID{})

// UUIDEncodeValue attempts to marshal an uuid.UUID type into a MongoDB binary uuid type.
func UUIDEncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != TypeUUID {
		return bsoncodec.ValueEncoderError{Name: "UUIDEncodeValue", Types: []reflect.Type{TypeUUID}, Received: val}
	}

	/// Zero val should emit null
	if val.IsZero() {
		return vw.WriteNull()
	}

	// MarshalBinary always returns nil error, ignoring
	b, _ := val.Interface().(uuid.UUID).MarshalBinary()
	return vw.WriteBinaryWithSubtype(b, bson.TypeBinaryUUID)
}

// UUIDDecodeValue attempts to unmarshal a string or a binary with subtype 0x00 (generic) or 0x04 (uuid) into an
// uuid.UUID type.
func UUIDDecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || val.Type() != TypeUUID {
		return bsoncodec.ValueDecoderError{Name: "UUIDDecodeValue", Types: []reflect.Type{TypeUUID}, Received: val}
	}

	var err error

	switch vrType := vr.Type(); vrType {
	case bson.TypeString:
		str, err := vr.ReadString()
		if err == nil {
			v, err := uuid.Parse(str)
			if err == nil {
				val.Set(reflect.ValueOf(v))
			}
		}
	case bson.TypeNull:
		err = vr.ReadNull()
	case bson.TypeUndefined:
		err = vr.ReadUndefined()
	case bson.TypeBinary:
		b, bType, err := vr.ReadBinary()
		if err == nil {
			switch bType {
			// If stored as a generic type (as was the default if you were to store with the default configuration on
			// go.mongodb.org/mongo-driver@1.12.1) and we're asking for an uuid type then attempt to decode
			case bson.TypeBinaryGeneric, bson.TypeBinaryUUID:
				v, err := uuid.FromBytes(b)
				if err == nil {
					val.Set(reflect.ValueOf(v))
				}
			default:
				return fmt.Errorf("cannot decode binary type %v into a uuid", bType)
			}
		}
	default:
		return fmt.Errorf("cannot decode %v into a uuid", vrType)
	}

	return err
}
