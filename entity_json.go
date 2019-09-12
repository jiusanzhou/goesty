package goesty

// TODO: use build flag to use jsoniter
import "encoding/json"

var (
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)
