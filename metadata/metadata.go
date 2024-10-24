package metadata

import "github.com/google/uuid"

type Metadata struct {
	UUID   uuid.UUID
	APIKey string
}
