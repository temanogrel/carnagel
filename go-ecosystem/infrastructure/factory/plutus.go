package factory

import (
	"errors"
	"fmt"
	"time"

	"context"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

func NewPlutusConnection(ctx context.Context, consul ecosystem.ConsulClient) (*grpc.ClientConn, error) {

	instances, _, err := consul.API().Catalog().Service("plutus", "", &consulapi.QueryOptions{})
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, errors.New("No plutus service available in consul")
	}

	grpcConfig := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBackoffConfig(grpc.BackoffConfig{
			MaxDelay: 1 * time.Second,
		}),
	}

	// Connected to plutus
	return grpc.DialContext(ctx, fmt.Sprintf("%s:%d", instances[0].ServiceAddress, instances[0].ServicePort), grpcConfig...)
}
