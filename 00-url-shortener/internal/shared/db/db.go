package db

import (
	"fmt"

	"github.com/gocql/gocql"
)

func New(clusters []string, keyspace string) (*gocql.Session, error) {
	var cluster = gocql.NewCluster(clusters[0])
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum // Just for development, not for production

	var session, err = cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return session, nil
}
