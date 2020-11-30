package database

import (
	"github.com/go-logr/logr"
	"github.com/markuszm/netpol-visualizer/model"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Neo4jClient struct {
	session neo4j.Session
	log     logr.Logger
}

func (client *Neo4jClient) Insert(policies model.Policies) error {
	result, err := client.session.Run("MATCH (x) RETURN (x)", map[string]interface{}{})
	if err != nil {
		client.log.Error(err, "Cipher query failed")
	}

	client.log.Info("Cipher query returned result", "session", client.session, "result", result)

	return nil
}

func NewNeo4jClient(url, username, password string) *Neo4jClient {
	driver := createDriver(url, username, password)
	session, err := driver.NewSession(neo4j.SessionConfig{})
	if err != nil {
		panic(err)
	}
	return &Neo4jClient{session: session}
}

func createDriver(url, username, password string) neo4j.Driver {
	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }

	driver, err := neo4j.NewDriver(url, neo4j.BasicAuth(username, password, ""), configForNeo4j40)
	if err != nil {
		panic(err)
	}

	return driver
}
