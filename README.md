UUID Support for BSON Serialization in Mongo Go Driver
===

---

### Motivation
While working with a legacy system we are indexing a fairly large amount of
data into a time series collection where the IDs are UUID types.  Our first
instinct was to use strings as they are easier to work around with the go
drivers for MongoDB, but we can cut our index sizes considerably using the
binary representation.  Thus, bson-uuid was born.

### Installation
```
go get -u github.com/vonix-networks/bsonuuid
```

### Deployment
There is a builder for a quick replacement if this is the only change to 
the registry needed:

```
client, err := mongo.Connect(context.Background(),
    options.Client().ApplyURI("uri"),
    options.Client().SetRegistry(bsonuuid.BuildRegistry()))
```

This will enable serialization from and to the standard UUID MongoDB type
(specifically the binary subtype 0x04).  It will also attempt to 
automatically parse strings.
