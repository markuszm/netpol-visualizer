package database

import (
	"github.com/go-logr/logr"
	"github.com/markuszm/netpol-visualizer/model"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Neo4jClient struct {
	driver  neo4j.Driver
	session neo4j.Session
	log     logr.Logger
}

func (client *Neo4jClient) Insert(policies model.Policies) error {
	client.log.Info("=> Insert", "policies", policies, "session", client.session)
	//client.session, _ = client.driver.Session(neo4j.AccessModeWrite)
	//defer client.session.Close()
	result, err := client.session.Run("MATCH (x) RETURN (x)", map[string]interface{}{})
	if err != nil {
		client.log.Error(err, "Cipher query failed")
	}

	var records []map[string]interface{}

	for result.Next() {
		record, ok := result.Record().Get("x")
		if !ok {
			keys, _ := result.Keys()
			client.log.Error(nil, "no dummy record found", "keys", keys)
		}

		records = append(records, record.(neo4j.Node).Props())
	}

	client.log.Info("Cipher query returned result", "session", client.session, "result", records)
	return nil
}

func NewNeo4jClient(url, username, password string, logger logr.Logger) Neo4jClient {
	driver := createDriver(url, username, password)
	err := driver.VerifyConnectivity()
	if err != nil {
		panic(err)
	}
	session, err := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		panic(err)
	}
	return Neo4jClient{driver: driver, session: session, log: logger}
}

func (client *Neo4jClient) Destroy() {
	err := client.session.Close()
	if err != nil {
		client.log.Error(err, "Failed to close session")
	}
	err = client.driver.Close()
	if err != nil {
		client.log.Error(err, "Failed to close driver")
	}
}

func createDriver(url, username, password string) neo4j.Driver {
	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }

	driver, err := neo4j.NewDriver(url, neo4j.BasicAuth(username, password, ""), configForNeo4j40)
	if err != nil {
		panic(err)
	}

	return driver
}
