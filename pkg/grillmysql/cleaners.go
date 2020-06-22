package grillmysql

import (
	"fmt"

	"bitbucket.org/swigy/grill"
)

func (gm *Mysql) DeleteTable(tableName string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := gm.Client().Exec(fmt.Sprintf("DROP TABLE %s", tableName))
		return err
	})
}
