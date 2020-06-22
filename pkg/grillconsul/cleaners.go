package grillconsul

import "bitbucket.org/swigy/grill"

func (gc *Consul) DeleteAllKeys() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := gc.consul.Client.KV().DeleteTree("", nil)
		return err
	})
}
