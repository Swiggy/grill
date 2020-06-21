package grillconsul

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"bitbucket.org/swigy/grill"
	"github.com/hashicorp/consul/api"
)

// Reads a csv file and puts them in Consul.
// CSV Format - key,value
// Headers are not ignored
func (gc *GrillConsul) SeedFromCSVFile(filepath string) grill.Stub {
	return grill.StubFunc(func() error {
		csvfile, err := os.Open(filepath)
		if err != nil {
			log.Fatalln("Couldn't open the csv file", err)
		}

		r := csv.NewReader(csvfile)

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			_, err = gc.consul.Client.KV().Put(&api.KVPair{Key: record[0], Value: []byte(record[1])}, nil)
			if err != nil {
				return fmt.Errorf("error seeding data in consul, error=%v", err)
			}
		}
		return nil
	})
}

func (gc *GrillConsul) Set(key, value string) grill.Stub {
	return grill.StubFunc(func() error {
		_, err := gc.consul.Client.KV().Put(&api.KVPair{Key: key, Value: []byte(value)}, nil)
		return err
	})
}
