package database

import "github.com/markuszm/netpol-visualizer/model"
import neo4j_driver "github.com/neo4j/neo4j-go-driver/neo4j"

type Client interface {
	Insert(policies model.Policies) error
}

type Neo4j struct {
	driver neo4j_driver.Driver
}

func CreateNeo4j(url, username, password string) *Neo4j {
	driver := initDB(url, username, password)
	return &Neo4j{driver: driver}
}

func initDB(url, username, password string) neo4j_driver.Driver {
	configForNeo4j40 := func(conf *neo4j_driver.Config) { conf.Encrypted = false }

	driver, err := neo4j_driver.NewDriver(url, neo4j_driver.BasicAuth(username, password, ""), configForNeo4j40)
	if err != nil {
		panic(err)
	}

	return driver
}

//func (r *Neo4j) exec(query string, args map[string]interface{}) (int64, error) {
//	result, execErr := r.conn.ExecNeo(query, args)
//	if execErr != nil {
//		return -1, execErr
//	}
//	rowsAffected, metaDataErr := result.RowsAffected()
//	if metaDataErr != nil {
//		return 0, nil
//	}
//	return rowsAffected, nil
//}

//func InsertDependency(neo4JDatabase Database, dep model.Dependency) error {
//	_, insertErr := neo4JDatabase.Exec(`
//					MERGE (p1:Package {name: {p1}})
//					MERGE (p2:Package {name: {p2}})
//					MERGE (p1)-[:DEPEND {s: {sourceVersion}, t: {targetVersion}}]->(p2)`,
//		map[string]interface{}{"p1": dep.PkgName, "p2": dep.Name, "version": dep.Version})
//	return insertErr
//}

func (r *Neo4j) Insert(policies model.Policies) error {
	return nil
}
