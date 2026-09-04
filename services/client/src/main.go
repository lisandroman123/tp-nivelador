package main

import (
	"errors"
	"os"

	client "github.com/lisandroman123/tp-nivelador/src/client"
	"github.com/lisandroman123/tp-nivelador/src/client_protocol"
	"github.com/lisandroman123/tp-nivelador/src/logger"
	client_connection "github.com/lisandroman123/tp-nivelador/src/safe_socket"
)

func loadConfig() (client_connection.ClientConfig, error) {
	agencyId := os.Getenv("AGENCY_ID")
	if agencyId == "" {
		return client_connection.ClientConfig{}, errors.New("AGENCY_ID environment variable is required")
	}

	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		return client_connection.ClientConfig{}, errors.New("SERVER_HOST environment variable is required")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		return client_connection.ClientConfig{}, errors.New("SERVER_PORT environment variable is required")
	}

	return client_connection.ClientConfig{
		ServerHost: serverHost,
		ServerPort: serverPort,
		AgencyId:   agencyId,
	}, nil
}

func run() int {
	config, err := loadConfig()
	if err != nil {
		logger.Error("load-config", logger.Fail, "err", err)
		return 1
	}

	connection, err := client_connection.NewClient(config)
	if err != nil {
		logger.Error("client-new", logger.Fail, "err", err)
		return 1
	}
	defer connection.Close()

	protocol := client_protocol.New(connection)

	client := client.New(protocol)

	if err := client.Run(); err != nil {
		logger.Error("client-run", logger.Fail, "err", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
