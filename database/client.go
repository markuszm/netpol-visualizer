package database

import (
	"github.com/markuszm/netpol-visualizer/model"
)

type Client interface {
	Insert(policies model.Policies) error
}

//func (r *Neo4jClient) exec(query string, args map[string]interface{}) (int64, error) {
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
