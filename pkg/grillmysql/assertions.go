package grillmysql

import (
	"fmt"

	"github.com/swiggy-private/grill"
)

func (gm *Mysql) AssertCount(tableName string, expectedCount int) grill.Assertion {
	return grill.AssertionFunc(func() error {
		result, err := gm.Client().Query(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName))
		if err != nil {
			return err
		}
		if !result.Next() {
			return fmt.Errorf("no result returned")
		}

		count := 0
		result.Scan(&count)

		if count != expectedCount {
			return fmt.Errorf("invalid number of rows in table, got=%v, want=%v", count, expectedCount)
		}
		return nil
	})
}
