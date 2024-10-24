package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/2zqa/ssot-specs-collector/metadata"
	"github.com/2zqa/ssot-specs-collector/specs"
	"github.com/charmbracelet/log"
)

type dummySender struct {
}

type dummyJsonSender struct {
}

func (p dummySender) Send(specs specs.PCSpecs, metadata metadata.Metadata, _ context.Context) error {
	log.Info(fmt.Sprintf("Sending \"%v\" to central server\n", specs))
	return nil
}

func (p dummyJsonSender) Send(specs specs.PCSpecs, metadata metadata.Metadata, _ context.Context) error {
	json, err := json.MarshalIndent(specs, "", "  ")
	if err != nil {
		log.Error("Error while marshalling specs", "error", err)
	}
	log.Info(string(json))
	return nil
}
