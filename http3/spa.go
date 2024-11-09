package http3

import (
	"bytes"
	"context"
	"fmt"
	quic "github.com/nxenon/xquic-go"
	"github.com/quic-go/qpack"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CalculateIntegerEncodingLengthValue(l int64) []byte {
	if l < 0 {
		panic("l is lower that 0!")
	}

	var encodedValue []byte

	switch {
	case l <= 63:
		encodedValue = []byte{byte(l)}
	case l <= 16383:
		encodedValue = []byte{0x40 | byte(l>>8), byte(l)}
	case l <= 1073741823:
		encodedValue = []byte{0x80 | byte(l>>24), byte(l >> 16), byte(l >> 8), byte(l)}
	default:
		encodedValue = []byte{0xC0 | byte(l>>56), byte(l >> 48), byte(l >> 40), byte(l >> 32),
			byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}
	}

	return encodedValue
}

func ParseResponseFromStream(biStream quic.Stream) (*http.Response, error) {

	frame1, err := NewXParseNextFrame(biStream, nil)
	if err != nil {
		return &http.Response{}, err
	}

	hf, ok := frame1.(*HeadersFrame)
	if !ok {
		fmt.Println("first frame is not headers frame")
		return &http.Response{}, error(nil)
		//panic("first frame is not headers frame")
	}

	headerBlock := make([]byte, hf.Length)

	if _, err2 := io.ReadFull(biStream, headerBlock); err2 != nil {
		return &http.Response{}, err2
	}
	decoder := qpack.Decoder{}

	hfs, err := decoder.DecodeFull(headerBlock)

	res, err := ResponseFromHeaders(hfs)

	var httpStr quic.Stream
	hstr := NewStream(biStream, nil)
	if _, ok := res.Header["Content-Length"]; ok && res.ContentLength >= 0 {
		httpStr = NewLengthLimitedStream(hstr, res.ContentLength)
	} else {
		httpStr = hstr
	}
	respBody := NewResponseBody(httpStr, nil, nil)

	// Rules for when to set Content-Length are defined in https://tools.ietf.org/html/rfc7230#section-3.3.2.
	_, hasTransferEncoding := res.Header["Transfer-Encoding"]
	isInformational := res.StatusCode >= 100 && res.StatusCode < 200
	isNoContent := res.StatusCode == http.StatusNoContent
	//isSuccessfulConnect := req.Method == http.MethodConnect && res.StatusCode >= 200 && res.StatusCode < 300
	//if !hasTransferEncoding && !isInformational && !isNoContent && !isSuccessfulConnect {
	if !hasTransferEncoding && !isInformational && !isNoContent {
		res.ContentLength = -1
		if clens, ok := res.Header["Content-Length"]; ok && len(clens) == 1 {
			if clen64, err := strconv.ParseInt(clens[0], 10, 64); err == nil {
				res.ContentLength = clen64
			}
		}
	}
	requestGzip := true
	if requestGzip && res.Header.Get("Content-Encoding") == "gzip" {
		res.Header.Del("Content-Encoding")
		res.Header.Del("Content-Length")
		res.ContentLength = -1
		res.Body = NewGzipReader(respBody)
		res.Uncompressed = true
	} else {
		res.Body = respBody
	}

	return res, nil
}

func Print_bytes_in_hex(data []byte) {
	for _, n := range data {
		fmt.Printf("%02x", n) // Prints hexadecimal representation
	}
	fmt.Println()
}

func SendRequestBytesInStream(stream quic.Stream, requestBytes []byte) error {

	_, err := stream.Write(requestBytes)
	return err

}

func ReadFromAllStreams(allStreams map[*http.Request]quic.Stream) map[*http.Request]*http.Response {

	streamsResponseMap := make(map[*http.Request]*http.Response)
	for request, biStream := range allStreams {
		res, err := ReadOneStream(biStream)
		if err != nil {
			fmt.Printf("Stream ID: %d has error (no response)!: %s", biStream.StreamID(), err)
			continue
		}
		streamsResponseMap[request] = res
	}

	return streamsResponseMap

}

func ReadOneStream(biStream quic.Stream) (*http.Response, error) {
	res, err := ParseResponseFromStream(biStream)
	return res, err
}

func GetBidirectionalStream(quicConnection quic.Connection) quic.Stream {
	context := context.Background()
	biStream, err := quicConnection.OpenStreamSync(context)
	if err != nil {
		panic(err)
	}

	return biStream

}

func CloseAllStreams(allStreams map[*http.Request]quic.Stream) {
	for key, _ := range allStreams {
		err := allStreams[key].Close()
		if err != nil {
			//fmt.Printf("Error closing Stream ID %d -> %s\n", value.StreamID(), err)
		}
	}
}

func GetDataFrameBytes(req http.Request) []byte {

	length := CalculateIntegerEncodingLengthValue(req.ContentLength)
	buf := []byte{
		0x00,
	}
	buf = append(buf, length...)

	var requestBodyBuffer bytes.Buffer
	_, err := io.Copy(&requestBodyBuffer, req.Body)
	if err != nil {
		panic(err)
	}

	buf = append(buf, requestBodyBuffer.Bytes()...)

	return buf

}

func GetDataFrameBytesWithLengthMinusLastByteNum(req http.Request, lastByteNum int) []byte {

	length := CalculateIntegerEncodingLengthValue(req.ContentLength - int64(lastByteNum))
	buf := []byte{
		0x00,
	}
	buf = append(buf, length...)

	var requestBodyBuffer bytes.Buffer
	_, err := io.Copy(&requestBodyBuffer, req.Body)
	if err != nil {
		panic(err)
	}

	buf = append(buf, requestBodyBuffer.Bytes()...)

	return buf

}

func GetLastByteDataFrame(b []byte) []byte {

	length := CalculateIntegerEncodingLengthValue(int64(len(b)))
	buf := []byte{
		0x00,
	}
	buf = append(buf, length...)

	buf = append(buf, b...)

	return buf

}

func GetRequestHeadersBytes(req http.Request) []byte {
	buf := &bytes.Buffer{}
	requestWriter := NewXRequestWriter(nil)
	isGzipped := true
	err := requestWriter.NewXwriteHeaders(buf, &req, isGzipped)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func GetRequestFinalPayload(req http.Request) []byte {
	var requestDataBytes []byte
	if req.Body != nil {
		//requestDataBytes := newGetDataFrameBytes(req.Body, req.ContentLength)
		requestDataBytes = GetDataFrameBytes(req)
	}

	finalPayload := append(GetRequestHeadersBytes(req), requestDataBytes...)

	return finalPayload
}

func GetRequestObject(urlString string, method string, headersMap map[string]string, bodyData []byte) (http.Request, error) {
	method = strings.ToUpper(method)
	var req *http.Request
	if method == "GET" {
		reqx, err := http.NewRequest(method, urlString, nil)
		if err != nil {
			return http.Request{}, err
		}
		req = reqx
	} else {
		reqx, err2 := http.NewRequest(method, urlString, bytes.NewReader(bodyData))
		if err2 != nil {
			return http.Request{}, err2
		}
		req = reqx
	}
	for key, value := range headersMap {
		req.Header.Set(key, value)
	}
	return *req, nil
}

func SendLastBytesOfStreams(allStreamsWithLastByte map[quic.Stream][]byte) {
	for s, b := range allStreamsWithLastByte {
		err := SendRequestBytesInStream(s, b)
		if err != nil {
			fmt.Printf("Error occurred in sending last byte of Stream ID: %d -> %s\n", s.StreamID(), err)
		}
	}
}

func SendRequestsWithSinglePacketAttackMethod(quicConn quic.Connection, allRequests []*http.Request, lastByteNum int, sleepMillisecondsBeforeSendingLastByte int) map[*http.Request]*http.Response {

	var allStreams map[*http.Request]quic.Stream
	allStreams = make(map[*http.Request]quic.Stream)
	var allStreamsWithLastByte map[quic.Stream][]byte
	allStreamsWithLastByte = make(map[quic.Stream][]byte)

	for _, request := range allRequests {
		var headersAndDataBytesMinusLastByte []byte

		headersFrameByte := GetRequestHeadersBytes(*request)
		dataFrameBytes := GetDataFrameBytesWithLengthMinusLastByteNum(*request, lastByteNum)

		allDataBytesExceptLastByte := dataFrameBytes[:len(dataFrameBytes)-lastByteNum]

		// all bytes except last byte
		headersAndDataBytesMinusLastByte = append(headersFrameByte, allDataBytesExceptLastByte...)

		finalByte := dataFrameBytes[len(dataFrameBytes)-lastByteNum:] // last byte
		finalByteDataFrame := GetLastByteDataFrame(finalByte)         // last byte data frame

		// send headers+data except last byte
		biStream := GetBidirectionalStream(quicConn)

		allStreamsWithLastByte[biStream] = finalByteDataFrame // for sending last byte
		allStreams[request] = biStream                        // for getting responses
		SendRequestBytesInStream(biStream, headersAndDataBytesMinusLastByte)

	}

	time.Sleep(time.Duration(sleepMillisecondsBeforeSendingLastByte) * time.Millisecond)

	// send all last bytes
	SendLastBytesOfStreams(allStreamsWithLastByte)

	CloseAllStreams(allStreams)
	streamsResponseMap := ReadFromAllStreams(allStreams)

	return streamsResponseMap
}
