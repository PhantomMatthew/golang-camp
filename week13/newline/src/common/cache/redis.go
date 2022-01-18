package cache

import (
	"github.com/gomodule/redigo/redis"
)

// Redis Redis
type Redis struct {
	host string
}

// Init Init
func (r *Redis) Init(host string) {
	r.host = host
}

// Exec Exec
func (r *Redis) Exec(db int, handle func(c redis.Conn) (interface{}, error)) (interface{}, error) {
	c, err := redis.Dial("tcp", r.host, redis.DialDatabase(db))
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return handle(c)
}
