package bsonuuid

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsonrw/bsonrwtest"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"reflect"
	"testing"
)

func TestUUIDEncodeValue(t *testing.T) {
	type UuidTest struct {
		Id uuid.UUID `bson:"_id"`
	}

	t.Run("encode", func(t *testing.T) {
		var testCases = []struct {
			name    string
			value   reflect.Value
			invoked bsonrwtest.Invoked
			err     string
		}{
			{
				name:    "success",
				value:   reflect.ValueOf(uuid.MustParse("7b68db73-a514-460e-900a-b3f47bbc7eaa")),
				invoked: bsonrwtest.WriteBinaryWithSubtype,
				err:     "",
			}, {
				name:    "wrong type",
				value:   reflect.ValueOf("wrong"),
				invoked: bsonrwtest.Nothing,
				err:     "UUIDEncodeValue can only encode valid uuid.UUID, but got string",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				vw := &bsonrwtest.ValueReaderWriter{}

				reg := BuildRegistry()
				enc := bsoncodec.ValueEncoderFunc(UUIDEncodeValue)
				err := enc.EncodeValue(bsoncodec.EncodeContext{Registry: reg}, vw, tc.value)

				assert.Equal(t, tc.invoked, vw.Invoked)
				if tc.err != "" {
					assert.EqualError(t, err, tc.err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	type UuidPtrTest struct {
		Id *uuid.UUID `bson:"_id"`
	}

	t.Run("marshal", func(t *testing.T) {
		reg := BuildRegistry()

		buf := new(bytes.Buffer)
		vw, err := bsonrw.NewBSONValueWriter(buf)
		assert.NoError(t, err)
		enc, err := bson.NewEncoder(vw)
		assert.NoError(t, err)
		_ = enc.SetRegistry(reg)

		value := &UuidTest{
			Id: uuid.MustParse("0F7D33CE-AF9F-4DDE-A4EC-ED630EAE74E2"),
		}

		err = enc.Encode(value)
		assert.NoError(t, err)

		binary, _ := value.Id.MarshalBinary()
		doc := buildDocument(bsoncore.AppendBinaryElement(nil, "_id", bson.TypeBinaryUUID, binary))
		assert.Equal(t, doc, buf.Bytes())
	})

	t.Run("marshal zero value", func(t *testing.T) {
		reg := BuildRegistry()

		buf := new(bytes.Buffer)
		vw, err := bsonrw.NewBSONValueWriter(buf)
		assert.NoError(t, err)
		enc, err := bson.NewEncoder(vw)
		assert.NoError(t, err)
		_ = enc.SetRegistry(reg)

		value := &UuidTest{
			Id: uuid.UUID{},
		}

		err = enc.Encode(value)
		assert.NoError(t, err)

		doc := buildDocument(bsoncore.AppendNullElement(nil, "_id"))
		assert.Equal(t, doc, buf.Bytes())
	})

	t.Run("marshal pointer", func(t *testing.T) {
		reg := BuildRegistry()

		buf := new(bytes.Buffer)
		vw, err := bsonrw.NewBSONValueWriter(buf)
		assert.NoError(t, err)
		enc, err := bson.NewEncoder(vw)
		assert.NoError(t, err)
		_ = enc.SetRegistry(reg)

		v := uuid.MustParse("0F7D33CE-AF9F-4DDE-A4EC-ED630EAE74E2")
		value := &UuidPtrTest{
			Id: &v,
		}

		err = enc.Encode(value)
		assert.NoError(t, err)

		binary, _ := value.Id.MarshalBinary()
		doc := buildDocument(bsoncore.AppendBinaryElement(nil, "_id", bson.TypeBinaryUUID, binary))
		assert.Equal(t, doc, buf.Bytes())
	})

	t.Run("marshal nil pointer", func(t *testing.T) {
		reg := BuildRegistry()

		buf := new(bytes.Buffer)
		vw, err := bsonrw.NewBSONValueWriter(buf)
		assert.NoError(t, err)
		enc, err := bson.NewEncoder(vw)
		assert.NoError(t, err)
		_ = enc.SetRegistry(reg)

		value := &UuidPtrTest{
			Id: nil,
		}

		err = enc.Encode(value)
		assert.NoError(t, err)

		doc := buildDocument(bsoncore.AppendNullElement(nil, "_id"))
		assert.Equal(t, doc, buf.Bytes())
	})
}

func TestUUIDDecodeValue(t *testing.T) {
	type UuidTest struct {
		Id uuid.UUID `bson:"_id"`
	}

	t.Run("decoder", func(t *testing.T) {
		var testCases = []struct {
			name     string
			expected interface{}
			b        []byte
			err      error
		}{
			{
				"decode from string",
				UuidTest{Id: uuid.MustParse("57CCFF46-A71D-4019-AF2D-70C4BBAF28B9")},
				buildDocument(bsoncore.AppendStringElement(nil, "_id", "57CCFF46-A71D-4019-AF2D-70C4BBAF28B9")),
				nil,
			}, {
				"decode from binary uuid type",
				UuidTest{Id: uuid.MustParse("FFFFFFFF-FFFF-FFFF-FFFF-000000000000")},
				buildDocument(bsoncore.AppendBinaryElement(nil, "_id", bson.TypeBinaryUUID, []byte{
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0})),
				nil,
			}, {
				"decode from binary generic type",
				UuidTest{Id: uuid.MustParse("00000000-0000-0000-0000-FFFFFFFFFFFF")},
				buildDocument(bsoncore.AppendBinaryElement(nil, "_id", bson.TypeBinaryGeneric, []byte{
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255})),
				nil,
			}, {
				"undefined returns zero value",
				UuidTest{},
				buildDocument(bsoncore.AppendUndefinedElement(nil, "_id")),
				nil,
			}, {
				"null returns zero value",
				UuidTest{},
				buildDocument(bsoncore.AppendNullElement(nil, "_id")),
				nil,
			}, {
				"error on bad type",
				UuidTest{},
				buildDocument(bsoncore.AppendInt32Element(nil, "_id", 42)),
				errors.New("error decoding key _id: cannot decode 32-bit integer into a uuid"),
			}, {
				"error on bad binary type",
				UuidTest{},
				buildDocument(bsoncore.AppendBinaryElement(nil, "_id", bson.TypeBinaryBinaryOld, []byte{
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0})),
				errors.New("error decoding key _id: cannot decode binary type 2 into a uuid"),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				vr := bsonrw.NewBSONDocumentReader(tc.b)
				dec, err := bson.NewDecoder(vr)
				assert.NoError(t, err)

				reg := BuildRegistry()
				_ = dec.SetRegistry(reg)

				var result UuidTest
				err = dec.Decode(&result)
				if tc.err == nil {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.err.Error())
				}
				assert.Equal(t, tc.expected, result)
			})
		}
	})

	type UuidPtrTest struct {
		Id *uuid.UUID `bson:"_id"`
	}

	t.Run("unmarshal pointer", func(t *testing.T) {
		b := buildDocument(bsoncore.AppendBinaryElement(nil, "_id", bson.TypeBinaryUUID, []byte{
			255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0}))
		vr := bsonrw.NewBSONDocumentReader(b)
		dec, err := bson.NewDecoder(vr)
		assert.NoError(t, err)

		reg := BuildRegistry()
		_ = dec.SetRegistry(reg)

		var result UuidPtrTest
		err = dec.Decode(&result)
		assert.NoError(t, err)
		v := uuid.MustParse("FFFFFFFF-FFFF-FFFF-FFFF-000000000000")
		expected := UuidPtrTest{Id: &v}
		assert.Equal(t, expected, result)
	})

	t.Run("unmarshal null pointer", func(t *testing.T) {
		b := buildDocument(bsoncore.AppendNullElement(nil, "_id"))
		vr := bsonrw.NewBSONDocumentReader(b)
		dec, err := bson.NewDecoder(vr)
		assert.NoError(t, err)

		reg := BuildRegistry()
		_ = dec.SetRegistry(reg)

		var result UuidPtrTest
		err = dec.Decode(&result)
		assert.NoError(t, err)
		expected := UuidPtrTest{Id: nil}
		assert.Equal(t, expected, result)
	})
}

// buildDocument inserts elems inside of a document.
func buildDocument(elems []byte) []byte {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = append(doc, elems...)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
	return doc
}
