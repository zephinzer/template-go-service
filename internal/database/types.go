package database

type ConnectionOpts struct {
	// Scheme is the connection schema to use, when not specified,
	// a default is applied depending on the driver
	Scheme string

	// Username is the user to use for the connection
	Username string

	// Password is the password for the user specified in .Username
	Password string

	// Hosts is the URLs (port-included) where the database server is reachable at
	Hosts []string

	// Database name of the database to use once connected
	Database string

	// Options is an optional hashmap of configuration values. Usage
	// depends on the driver
	Options map[string]string
}
