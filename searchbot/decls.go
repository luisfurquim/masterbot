package searchbot

import (
   "errors"
   "github.com/luisfurquim/goose"
)

const (
   bufsz = 4096
)

var ErrTmoutFetchingSwagger          = errors.New("Error timeout fetching swagger.json")
var ErrHttpStatusFetchingSwagger     = errors.New("Error HTTP status fetching swagger.json")
var ErrDecodingSwagger               = errors.New("Error decoding swagger.json")
var ErrUndefinedField                = errors.New("Error undefined field")
var ErrMarshalingRequestBody         = errors.New("Error marshaling request body")
var ErrUnmarshalingRequestBody       = errors.New("Error unmarshaling request body")
var ErrQueryingSearchBot             = errors.New("Error querying search bot")
var ErrAssemblyingRequest            = errors.New("Error assemblying http request")
var ErrReadingResponseBody           = errors.New("Error reading response body")

type SearchbotG struct {
   Search goose.Alert    `json:"Search"`
   Taxonomy goose.Alert  `json:"Taxonomy"`
}

var Goose SearchbotG