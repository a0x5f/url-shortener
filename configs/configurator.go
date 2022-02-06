package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Configuration struct {
	Server   ServerConfig
	Postgres PostgresConfig
}

type ServerConfig struct {
	Port        int  `json:"port"`
	UsePostgres bool `json:"postgres"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Dbname   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func ReadConfig() *Configuration {
	serverConfig := ServerConfig{}
	postgresConfig := PostgresConfig{}

	serverConfigFile, err := os.Open("./configs/server.json")
	if err != nil {
		log.Fatal(err)
	}
	defer serverConfigFile.Close()

	serverConfigJSON, _ := ioutil.ReadAll(serverConfigFile)
	json.Unmarshal(serverConfigJSON, &serverConfig)

	if serverConfig.UsePostgres {
		postgresConfigFile, err := os.Open("./configs/postgres.json")
		if err != nil {
			log.Fatal(err)
		}
		defer postgresConfigFile.Close()

		postgresConfigJSON, _ := ioutil.ReadAll(postgresConfigFile)
		json.Unmarshal(postgresConfigJSON, &postgresConfig)
	}

	return &Configuration{
		Server:   serverConfig,
		Postgres: postgresConfig,
	}
}
