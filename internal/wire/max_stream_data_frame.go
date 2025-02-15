package wire

import (
	"bytes"

	"github.com/nxenon/h3spacex/internal/protocol"
	"github.com/nxenon/h3spacex/quicvarint"
)

// A MaxStreamDataFrame is a MAX_STREAM_DATA frame
type MaxStreamDataFrame struct {
	StreamID          protocol.StreamID
	MaximumStreamData protocol.ByteCount
}

func parseMaxStreamDataFrame(r *bytes.Reader, _ protocol.Version) (*MaxStreamDataFrame, error) {
	sid, err := quicvarint.Read(r)
	if err != nil {
		return nil, err
	}
	offset, err := quicvarint.Read(r)
	if err != nil {
		return nil, err
	}

	return &MaxStreamDataFrame{
		StreamID:          protocol.StreamID(sid),
		MaximumStreamData: protocol.ByteCount(offset),
	}, nil
}

func (f *MaxStreamDataFrame) Append(b []byte, version protocol.Version) ([]byte, error) {
	b = append(b, maxStreamDataFrameType)
	b = quicvarint.Append(b, uint64(f.StreamID))
	b = quicvarint.Append(b, uint64(f.MaximumStreamData))
	return b, nil
}

// Length of a written frame
func (f *MaxStreamDataFrame) Length(version protocol.Version) protocol.ByteCount {
	return 1 + quicvarint.Len(uint64(f.StreamID)) + quicvarint.Len(uint64(f.MaximumStreamData))
}
