package handler

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

const bucketName = "grpc"

func init() {
	db1, err := bolt.Open("./grpc.db", 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		panic(err)
	}
	db = db1
}

func marshal(v interface{}) []byte {
	var dt, err = json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return dt
}

func unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func bucket(tx *bolt.Tx, name string) *bolt.Bucket {
	_, _ = tx.CreateBucketIfNotExists([]byte(name))
	return tx.Bucket([]byte(name))
}

func set(b *bolt.DB, name string, key string, val interface{}) error {
	return b.Update(func(tx *bolt.Tx) error {
		var bkt = bucket(tx, name)
		return bkt.Put([]byte(key), marshal(val))
	})

}

func del(b *bolt.DB, name string, key string) error {
	return b.Update(func(tx *bolt.Tx) error {
		var bkt = bucket(tx, name)
		return bkt.Delete([]byte(key))
	})
}

func has(b *bolt.DB, name string, key string) bool {
	return b.View(func(tx *bolt.Tx) error {
		var bkt = bucket(tx, name)
		var val = bkt.Get([]byte(key))
		if val != nil {
			return nil
		}
		return errors.New("not found")
	}) == nil
}

func get(b *bolt.DB, name string, key string, v interface{}) error {
	return b.View(func(tx *bolt.Tx) error {
		var bkt = bucket(tx, name)
		var val = bkt.Get([]byte(key))
		return unmarshal(val, v)
	})
}

func list(b *bolt.DB, name string, fn interface{}) error {
	return b.View(func(tx *bolt.Tx) error {
		var bkt = bucket(tx, name)
		return bkt.ForEach(func(k, v []byte) error {
			if k == nil || v == nil {
				return nil
			}

			var mthIn = reflect.New(reflect.TypeOf(fn).In(0).Elem())
			ret := reflect.ValueOf(unmarshal).Call([]reflect.Value{reflect.ValueOf(v), mthIn})
			if !ret[0].IsNil() {
				return ret[0].Interface().(error)
			}

			reflect.ValueOf(fn).Call([]reflect.Value{mthIn})
			return nil
		})
	})
}
