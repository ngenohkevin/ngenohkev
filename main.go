package main

import (
	"github.com/ngenohkevin/ngenohkev/internals/server"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	port := 9000

	//create new server instance
	srv, err := server.NewServer(logger, port)
	if err != nil {
		logger.Error("Failed to create server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	//start the server
	if err := srv.Start(); err != nil {
		logger.Error("Server encountered an error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
