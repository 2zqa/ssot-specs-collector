package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/2zqa/ssot-specs-collector/metadata"
	"github.com/2zqa/ssot-specs-collector/sender"
	"github.com/2zqa/ssot-specs-collector/specs"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v3"
)

const projectName = "ssot-specs-collector"

func main() {
	fs := flag.NewFlagSet(projectName, flag.ExitOnError)
	defaultConfigPath := filepath.Join("/etc", projectName, "config")
	var (
		apiKey = fs.String("api-key", "", "api key for sending specs to the server")
		debug  = fs.Bool("debug", false, "log debug information")
		_      = fs.String("config", defaultConfigPath, "config file")
	)

	var deviceUUID uuid.UUID
	fs.Func("uuid", "uuid to identify device with in database", func(s string) error {
		var err error
		deviceUUID, err = uuid.Parse(s)
		return err
	})

	if err := ff.Parse(fs, os.Args[1:], ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser)); err != nil {
		log.Fatal(err)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if *apiKey == "" || deviceUUID == uuid.Nil {
		log.Fatal("api-key and uuid are required")
	}

	metadata := metadata.Metadata{
		APIKey: *apiKey,
		UUID:   deviceUUID,
	}

	// Fetch and send specs
	specs, err := specs.Fetch()
	if err != nil {
		log.Warn("Something went wrong while trying to fetch specs", "err", err)
	}

	s, err := sender.NewSender("apiClient")
	if err != nil {
		log.Fatal(err)
	}

	err = s.Send(specs, metadata, context.Background())
	if err != nil {
		log.Fatal("Could not send specs", "err", err)
	}
}
