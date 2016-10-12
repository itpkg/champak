package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

//RedisStore redis cache store
type RedisStore struct {
	Redis *redis.Pool `inject:""`
}

//Flush clear cache items
func (p *RedisStore) Flush() error {
	c := p.Redis.Get()
	defer c.Close()
	keys, err := redis.Values(c.Do("KEYS", p.key("*")))
	if err == nil && len(keys) > 0 {
		_, err = c.Do("DEL", keys...)
	}
	return err
}

//Keys list cache items
func (p *RedisStore) Keys() ([]string, error) {
	c := p.Redis.Get()
	defer c.Close()
	return redis.Strings(c.Do("KEYS", p.key("*")))
}

//Set cache item
func (p *RedisStore) Set(key string, val interface{}, ttl uint) error {
	c := p.Redis.Get()
	defer c.Close()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return err
	}
	_, err := c.Do("SET", p.key(key), buf.Bytes(), "EX", ttl)
	return err
}

//Get get from cache
func (p *RedisStore) Get(key string, val interface{}) error {
	c := p.Redis.Get()
	defer c.Close()
	bys, err := redis.Bytes(c.Do("GET", p.key(key)))
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(bys)
	return dec.Decode(val)
}

func (p *RedisStore) key(k string) string {
	return fmt.Sprintf("cache://%s", k)
}
