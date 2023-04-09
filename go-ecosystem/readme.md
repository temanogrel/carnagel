Go-ecosystem
============

The purpose of this repository is to provide factories for infrastructure and domain services and systems.
Without having to duplicate all the factories and configurations in each project but relying on consul.


# Factories

Factories allow standard access to our shared services

 - `NewLogger(application, hostname string, consul *api.Client) (logrus.FieldLogger, error)`
 - `NewRedisConnection(consul *api.Client) (*redis.Client, error)`
 - `NewMinervaConnection(consul *api.Client) (*grpc.ClientConn, error)`
 - `NewInfinityConnection(consul *api.Client) (*grpc.ClientConn, error)`

# Middleware

Http middleware that can be useful