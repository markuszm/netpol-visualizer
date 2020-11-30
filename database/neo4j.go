package database

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/markuszm/netpol-visualizer/model"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	v1 "k8s.io/api/networking/v1"
	"strconv"
)

type Neo4jClient struct {
	driver  neo4j.Driver
	session neo4j.Session
	log     logr.Logger
}

func (client *Neo4jClient) Insert(policies model.Policies) error {
	client.log.Info("=> Insert", "policies.len", len(policies), "session", client.session)
	//client.session, _ = client.driver.Session(neo4j.AccessModeWrite)
	for _, policy := range policies {
		queryFmt := `MERGE (from:Pod {namespace: $fromNamespace,name: $fromName}) MERGE (to:Pod {namespace: $toNamespace,name: $toName}) MERGE (from)-[:%s%s]->(to) RETURN from,to`
		var queryEdgeProps string
		if policy.Port != 0 {

			queryEdgeProps = " { port: " + strconv.Itoa(policy.Port) + " }"
		} else {
			queryEdgeProps = ""
		}
		var queryEdgeType string
		if policy.PolicyType == v1.PolicyTypeIngress {
			queryEdgeType = "INGRESS_ALLOWED"
		} else if policy.PolicyType == v1.PolicyTypeEgress {
			queryEdgeType = "EGRESS_ALLOWED"
		}

		queryString := fmt.Sprintf(queryFmt, queryEdgeType, queryEdgeProps)
		result, err := client.session.Run(queryString, map[string]interface{}{
			"fromNamespace": policy.From.Namespace,
			"fromName":      policy.From.Name,
			"toNamespace":   policy.To.Namespace,
			"toName":        policy.To.Name,
		})
		if err != nil {
			client.log.Error(err, "Cipher query failed")
		}

		var records []map[string]interface{}

		for result.Next() {
			record, ok := result.Record().Get("from")
			if !ok {
				keys, _ := result.Keys()
				client.log.Error(nil, "no dummy record found", "keys", keys)
			}

			records = append(records, record.(neo4j.Node).Props())
		}

		client.log.Info("Cipher query returned result", "session", client.session, "result", records)
	}
	return nil
}

func NewNeo4jClient(url, username, password string, logger logr.Logger) Neo4jClient {
	driver := createDriver(url, username, password)
	err := driver.VerifyConnectivity()
	if err != nil {
		panic(err)
	}
	session, err := driver.Session(neo4j.AccessModeWrite)
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
