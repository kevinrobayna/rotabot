package shell

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"goa.design/goa/v3/http"
	stdhttp "net/http"
)

type encoder struct {
	encoder http.Encoder
	logger  *zap.Logger
	kind    string
}

func ResponseEncoderWithLogs(ctx context.Context, r stdhttp.ResponseWriter) http.Encoder {
	return &encoder{
		encoder: http.ResponseEncoder(ctx, r),
		logger:  Logger(ctx),
		kind:    "response",
	}
}

func (e encoder) Encode(v interface{}) error {
	err := e.encoder.Encode(v)
	if err != nil {
		e.logger.Error(fmt.Sprintf("failed to decode %s", e.kind), zap.Error(err))
	}

	return err
}

type decoder struct {
	decoder http.Decoder
	logger  *zap.Logger
	kind    string
}

func RequestDecoderWithLogs(r *stdhttp.Request) http.Decoder {
	return &decoder{
		decoder: http.RequestDecoder(r),
		logger:  Logger(r.Context()),
		kind:    "request",
	}
}

func (d decoder) Decode(v interface{}) error {
	err := d.decoder.Decode(v)
	if err != nil {
		d.logger.Error(fmt.Sprintf("failed to decode %s", d.kind), zap.Error(err))
	}

	return err
}
