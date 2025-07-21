package internal

import (
	"encoding/json"
	"github.com/Laky-64/gologging"
	"os"
)

func LoadSQLCConfig(path string) SQLCConfig {
	data, err := os.ReadFile(path)
	checkErr(err)

	var config SQLCConfig
	checkErr(json.Unmarshal(data, &config))
	return config
}

func GetPublicSchema(config SQLCConfig) Schema {
	for _, schema := range config.Catalog.Schemas {
		if schema.Name == "public" {
			return schema
		}
	}
	gologging.Fatal("public schema not found")
	return Schema{}
}
