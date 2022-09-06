package shell

import (
	"fmt"
	"go.uber.org/zap"
	"goa.design/goa/v3/http"
	stdhttp "net/http"
)

type DecoderWithLogs struct {
	decoder http.Decoder
	logger  *zap.Logger
	kind    string
}

func RequestDecoderWithLogs(r *stdhttp.Request) http.Decoder {
	return &DecoderWithLogs{
		decoder: http.RequestDecoder(r),
		logger:  Logger(r.Context()),
		kind:    "request",
	}
}

func (d DecoderWithLogs) Decode(v interface{}) error {
	err := d.decoder.Decode(v)
	if err != nil {
		d.logger.Error(fmt.Sprintf("failed to decode %s", d.kind), zap.Error(err))
	}

	return err
}
