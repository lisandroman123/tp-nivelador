package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
	inputFile := os.Getenv("INPUT_FILE")
	if inputFile == "" {
		return client_connection.ClientConfig{}, errors.New("INPUT_FILE environment variable is required")
	}
	outputFile := os.Getenv("OUTPUT_FILE")
	if outputFile == "" {
		return client_connection.ClientConfig{}, errors.New("OUTPUT_FILE environment variable is required")
	}
	batchSize := os.Getenv("BATCH_SIZE")
	if batchSize == "" {
		return client_connection.ClientConfig{}, errors.New("BATCH_SIZE environment variable is required")
	}
	size, err := strconv.Atoi(batchSize)
	if err != nil {
		return client_connection.ClientConfig{}, err
	}
	return client_connection.ClientConfig{
		ServerHost: serverHost,
		ServerPort: serverPort,
		AgencyId:   agencyId,
		InputFile:  inputFile,
		OutputFile: outputFile,
		BatchSize:  size,
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

	client := client.New(config.InputFile, config.OutputFile, int(config.BatchSize), protocol)
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer stop()

	go func() {
		<-ctx.Done()
		client.Close()
	}()

	if err := client.Run(); err != nil {
		logger.Error("client-run", logger.Fail, "err", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
