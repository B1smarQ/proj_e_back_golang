package postgres

type DatabaseDriver interface {
	Create() error
}

type PostgresDriver struct {
	
}

func (pg PostgresDriver) Create() error {
	return nil
}
