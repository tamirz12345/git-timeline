package main

import (
	"fmt"
	"gitTimeline/pkg/models"
	"gitTimeline/pkg/server"
	"gitTimeline/pkg/storage/metadata"
	"gitTimeline/pkg/storage/timeline"
	"os"

	"github.com/rs/zerolog"
)

var (
	cfg    *models.Config
	logger zerolog.Logger
)

func init() {
	// Load configuration from the environment variables.
	var err error
	cfg, err = models.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration, error: %v", err)
		os.Exit(1)
	}
	// Init logger
	logger = zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Timestamp().Logger()

}

func main() {
	// init zero logger and log starting the server
	logger.Info().Msg("initiating git local repository communicator...")
	localRepositoryCommunicator, err := timeline.NewLocalGitRepositoryCommunicator(fmt.Sprintf("%s/%s", cfg.RepositoriesRootPath, "gitTimeline"), &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initiate git local repository communicator")
	}
	logger.Info().Msg("initiating git storage...")
	gitStorage := timeline.NewGitStorage(localRepositoryCommunicator, &logger)

	logger.Info().Msg("initiating mysql metadata storage...")
	metadataStorage, err := metadata.NewMysqlStorage(cfg.DbHost, cfg.DbUser, cfg.DbPass, "git_timeline", &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initiate mysql metadata storage")
	}
	logger.Info().Msg("initiating timeline server...")
	httpServer := server.NewServer(&logger, gitStorage, metadataStorage, cfg.ServerPort)
	httpServer.Start()
}
