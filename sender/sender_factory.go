package sender

import (
	"context"
	"fmt"

	"github.com/2zqa/ssot-specs-collector/metadata"
	"github.com/2zqa/ssot-specs-collector/specs"
)

type Sender interface {
	Send(specs.PCSpecs, metadata.Metadata, context.Context) error
}

func NewSender(senderType string) (Sender, error) {
	switch senderType {
	case "dummy":
		return dummySender{}, nil
	case "dummyJson":
		return dummyJsonSender{}, nil
	case "apiClient":
		return APIClientSender{}, nil
	default:
		return nil, fmt.Errorf("invalid sender type: %s", senderType)
	}
}
