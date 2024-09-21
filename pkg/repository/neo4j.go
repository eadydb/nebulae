package repository

type Neo4jDriver struct{}

func NewNeo4jDriver() *Neo4jDriver {
	return &Neo4jDriver{}
}

func (n *Neo4jDriver) Connect() error {
	return nil
}

func (n *Neo4jDriver) Close() error {
	return nil
}
