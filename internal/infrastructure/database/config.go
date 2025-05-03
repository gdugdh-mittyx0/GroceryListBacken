package database

type config struct {
	host     string
	database string
	port     uint
	user     string
	password string
}

func newConfigGorm(host string, database string, password string, user string, port uint) *config {
	return &config{
		host:     host,
		database: database,
		password: password,
		user:     user,
		port:     port,
	}
}
