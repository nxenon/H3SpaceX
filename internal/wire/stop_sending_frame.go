package wire

import (
	"bytes"

	"github.com/nxenon/h3spacex/internal/protocol"
	"github.com/nxenon/h3spacex/internal/qerr"
	"github.com/nxenon/h3spacex/quicvarint"
)

// A StopSendingFrame is a STOP_SENDING frame
type StopSendingFrame struct {
	StreamID  protocol.StreamID
	ErrorCode qerr.StreamErrorCode
}

// parseStopSendingFrame parses a STOP_SENDING frame
func parseStopSendingFrame(r *bytes.Reader, _ protocol.Version) (*StopSendingFrame, error) {
	streamID, err := quicvarint.Read(r)
	if err != nil {
		return nil, err
	}
	errorCode, err := quicvarint.Read(r)
	if err != nil {
		return nil, err
	}

	return &StopSendingFrame{
		StreamID:  protocol.StreamID(streamID),
		ErrorCode: qerr.StreamErrorCode(errorCode),
	}, nil
}

// Length of a written frame
func (f *StopSendingFrame) Length(_ protocol.Version) protocol.ByteCount {
	return 1 + quicvarint.Len(uint64(f.StreamID)) + quicvarint.Len(uint64(f.ErrorCode))
}

func (f *StopSendingFrame) Append(b []byte, _ protocol.Version) ([]byte, error) {
	b = append(b, stopSendingFrameType)
	b = quicvarint.Append(b, uint64(f.StreamID))
	b = quicvarint.Append(b, uint64(f.ErrorCode))
	return b, nil
}
