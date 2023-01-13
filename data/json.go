package data

import (
	"encoding/json"
	"io"
)

// FromJSON method performs unmarshalling of Event from JSON
func (e *Event) FromJSON(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(e)
}

// ToJSON method perfroms marshalling of AvgResponse to JSON
func (ar *AvgResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(ar)
}
