package grillmysql

import (
	"fmt"

	"bitbucket.org/swigy/grill"
)

func (gm *GrillMysql) DeleteTable(tableName string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := gm.Client().Exec(fmt.Sprintf("DROP TABLE %s", tableName))
		return err
	})
}
