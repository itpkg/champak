package jobber

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
	logging "github.com/op/go-logging"
)

//RedisJobber job by readis store
type RedisJobber struct {
	Redis  *redis.Pool     `inject:""`
	Logger *logging.Logger `inject:""`

	Timeout  int
	Handlers map[string]Handler
}

//Register register a job-handler
func (p *RedisJobber) Register(queue string, handler Handler) {
	p.Handlers[p.key(queue)] = handler
}

//Push add a job task
func (p *RedisJobber) Push(queue string, args interface{}) error {
	p.Logger.Infof("push job into %s", queue)
	c := p.Redis.Get()
	defer c.Close()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(args); err != nil {
		return err
	}
	_, err := c.Do("LPUSH", p.key(queue), buf.Bytes())
	return err
}

//Start start to process job
func (p *RedisJobber) Start() error {
	var err error
	for {
		err = p.run()
		if err != nil && err != redis.ErrNil {
			break
		}
	}
	return err
}

func (p *RedisJobber) run() error {
	const stop = ".stop"
	if _, err := os.Stat(stop); err == nil {
		return fmt.Errorf("find file %s, exit", stop)
	}

	if len(p.Handlers) == 0 {
		return errors.New("null handlers")
	}
	c := p.Redis.Get()
	defer c.Close()
	var keys []interface{}
	for k := range p.Handlers {
		keys = append(keys, k)
	}
	keys = append(keys, p.Timeout)
	args, err := redis.ByteSlices(c.Do("BRPOP", keys...))
	if err != nil {
		return err
	}
	queue := string(args[0])
	p.Logger.Infof("get a job from %s", queue)
	return p.Handlers[queue](args[1])
}

func (p *RedisJobber) key(queue string) string {
	return fmt.Sprintf("task://%s", queue)
}
