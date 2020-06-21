package grillmysql

import (
	"context"
	"database/sql"

	"bitbucket.org/swigy/grill/internal/canned"
)

type GrillMysql struct {
	mysql *canned.Mysql
}

func Start() (*GrillMysql, error) {
	mysql, err := canned.NewMysql(context.TODO())
	if err != nil {
		return nil, err
	}

	return &GrillMysql{
		mysql: mysql,
	}, nil
}

func (gm *GrillMysql) Client() *sql.DB {
	return gm.mysql.Client
}

func (gm *GrillMysql) Host() string {
	return gm.mysql.Host
}

func (gm *GrillMysql) Port() string {
	return gm.mysql.Port
}

func (gm *GrillMysql) Database() string {
	return gm.mysql.Database
}

func (gm *GrillMysql) Username() string {
	return gm.mysql.Username
}

func (gm *GrillMysql) Password() string {
	return gm.mysql.Password
}

func (gm *GrillMysql) Stop() error {
	return gm.mysql.Container.Terminate(context.TODO())
}
