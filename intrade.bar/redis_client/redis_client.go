package redis_client

import(
	rs "github.com/go-redis/redis/v8"
)

type c rs.Client

var C c

func InitNew(socket string) *c{
	var client = rs.NewClient(&rs.Options{
		Addr:     socket,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return (*c)(client)
}

