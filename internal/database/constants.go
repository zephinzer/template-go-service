package database

var TypesAvailable = []string{
	TypeMongo,
	TypeMySQL,
	TypePostgres,
}

const (
	TypeMongo    = "mongo"
	TypeMySQL    = "mysql"
	TypePostgres = "postgres"
)
