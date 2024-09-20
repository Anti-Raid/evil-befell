// Package loc provides functions for handling location metadata.
package loc

import (
	"encoding/json"
	"strings"
)

// The location string can be converted to a LocMetadata object
//
// Format of the location string: "route_id?[JSON data]"
type LocMetadata struct {
	ID   string
	Data map[string]string
}

func (loc *LocMetadata) MarshalJSON() ([]byte, error) {
	// Format using FormatLocMetadata
	return json.Marshal(FormatLocMetadata(loc))
}

func (loc *LocMetadata) UnmarshalJSON(data []byte) error {
	var locStr string

	// Unmarshal the JSON data
	if err := json.Unmarshal(data, &locStr); err != nil {
		return err
	}

	// Parse using ParseLocMetadata
	meta, err := ParseLocMetadata(locStr)

	if err != nil {
		return err
	}

	*loc = *meta
	return nil
}

// Convert a location string to a LocMetadata object
func ParseLocMetadata(loc string) (*LocMetadata, error) {
	// Split the location string into the route ID and the JSON data
	parts := strings.Split(loc, "?")

	// Create a new LocMetadata object
	meta := LocMetadata{
		ID: parts[0],
	}

	// If there is JSON data, parse it
	if len(parts) > 1 {
		if err := json.Unmarshal([]byte(parts[1]), &meta.Data); err != nil {
			return nil, err
		}
	}

	return &meta, nil
}

func FormatLocMetadata(loc *LocMetadata) string {
	if loc == nil {
		return ""
	} else if len(loc.Data) == 0 {
		return loc.ID
	}

	// Convert the JSON data to a string
	data, _ := json.Marshal(loc.Data)

	// Return the formatted location string
	return loc.ID + "?" + string(data)
}
