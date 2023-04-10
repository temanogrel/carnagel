package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

func NewRedisConnection(consul *api.Client) (*redis.Client, error) {

	services, _, err := consul.Catalog().Service("redis", "", &api.QueryOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query for redis service")
	}

	if len(services) == 0 {
		return nil, errors.New("No redis service available")
	}

	client := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     fmt.Sprintf("%s:%d", services[0].ServiceAddress, services[0].ServicePort),
		PoolSize: 20,
	})

	if err := client.Ping().Err(); err != nil {
		return client, err
	}

	return client, err
}
