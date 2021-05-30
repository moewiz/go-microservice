package data

import (
	"encoding/json"
	"io"
)

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance then json.Unmarshal as it does not have to
// buffer the output into an in-memory slice of bytes
// This reduces allocations and the overheads of the service
// @see https://golang.org/pkg/encoding/json/#NewEncoder
func ToJSON(i interface{}, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(i)
}

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
func FromJSON(i interface{}, r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(i)
}
