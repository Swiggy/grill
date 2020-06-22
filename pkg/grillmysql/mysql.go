package grillmysql

import (
	"context"
	"database/sql"

	"bitbucket.org/swigy/grill/internal/canned"
)

type Mysql struct {
	mysql *canned.Mysql
}

func Start() (*Mysql, error) {
	mysql, err := canned.NewMysql(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Mysql{
		mysql: mysql,
	}, nil
}

func (gm *Mysql) Client() *sql.DB {
	return gm.mysql.Client
}

func (gm *Mysql) Host() string {
	return gm.mysql.Host
}

func (gm *Mysql) Port() string {
	return gm.mysql.Port
}

func (gm *Mysql) Database() string {
	return gm.mysql.Database
}

func (gm *Mysql) Username() string {
	return gm.mysql.Username
}

func (gm *Mysql) Password() string {
	return gm.mysql.Password
}

func (gm *Mysql) Stop() error {
	return gm.mysql.Container.Terminate(context.TODO())
}
