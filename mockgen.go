//go:build gomock || generate

package quic

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_send_conn_test.go github.com/nxenon/h3spacex SendConn"
type SendConn = sendConn

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_raw_conn_test.go github.com/nxenon/h3spacex RawConn"
type RawConn = rawConn

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_sender_test.go github.com/nxenon/h3spacex Sender"
type Sender = sender

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_stream_internal_test.go github.com/nxenon/h3spacex StreamI"
type StreamI = streamI

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_crypto_stream_test.go github.com/nxenon/h3spacex CryptoStream"
type CryptoStream = cryptoStream

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_receive_stream_internal_test.go github.com/nxenon/h3spacex ReceiveStreamI"
type ReceiveStreamI = receiveStreamI

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_send_stream_internal_test.go github.com/nxenon/h3spacex SendStreamI"
type SendStreamI = sendStreamI

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_stream_getter_test.go github.com/nxenon/h3spacex StreamGetter"
type StreamGetter = streamGetter

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_stream_sender_test.go github.com/nxenon/h3spacex StreamSender"
type StreamSender = streamSender

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_crypto_data_handler_test.go github.com/nxenon/h3spacex CryptoDataHandler"
type CryptoDataHandler = cryptoDataHandler

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_frame_source_test.go github.com/nxenon/h3spacex FrameSource"
type FrameSource = frameSource

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_ack_frame_source_test.go github.com/nxenon/h3spacex AckFrameSource"
type AckFrameSource = ackFrameSource

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_stream_manager_test.go github.com/nxenon/h3spacex StreamManager"
type StreamManager = streamManager

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_sealing_manager_test.go github.com/nxenon/h3spacex SealingManager"
type SealingManager = sealingManager

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_unpacker_test.go github.com/nxenon/h3spacex Unpacker"
type Unpacker = unpacker

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_packer_test.go github.com/nxenon/h3spacex Packer"
type Packer = packer

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_mtu_discoverer_test.go github.com/nxenon/h3spacex MTUDiscoverer"
type MTUDiscoverer = mtuDiscoverer

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_conn_runner_test.go github.com/nxenon/h3spacex ConnRunner"
type ConnRunner = connRunner

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_quic_conn_test.go github.com/nxenon/h3spacex QUICConn"
type QUICConn = quicConn

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_packet_handler_test.go github.com/nxenon/h3spacex PacketHandler"
type PacketHandler = packetHandler

//go:generate sh -c "go run go.uber.org/mock/mockgen -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_packet_handler_manager_test.go github.com/nxenon/h3spacex PacketHandlerManager"

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/h3spacex -destination mock_packet_handler_manager_test.go github.com/nxenon/h3spacex PacketHandlerManager"
type PacketHandlerManager = packetHandlerManager

// Need to use source mode for the batchConn, since reflect mode follows type aliases.
// See https://github.com/golang/mock/issues/244 for details.
//
//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -package quic -self_package github.com/nxenon/h3spacex -source sys_conn_oob.go -destination mock_batch_conn_test.go -mock_names batchConn=MockBatchConn"

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -package quic -self_package github.com/nxenon/h3spacex -self_package github.com/nxenon/h3spacex -destination mock_token_store_test.go github.com/nxenon/h3spacex TokenStore"
//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -package quic -self_package github.com/nxenon/h3spacex -self_package github.com/nxenon/h3spacex -destination mock_packetconn_test.go net PacketConn"
