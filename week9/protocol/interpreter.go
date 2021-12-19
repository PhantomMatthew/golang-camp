package protocol

import (
	"encoding/binary"
	"errors"
)

const (
	PACKAGE_LENGTH_DATA_SIZE   = 4
	HEADER_LENGTH_DATA_SIZE    = 2
	PROTOCOL_VERSION_DATA_SIZE = 2
	OPERATION_CODE_DATA_SIZE   = 4
	SEQUENCE_ID_DATA_SIZE      = 4
)

var ErrorBrokenPacket = errors.New("error: broken package")

type GoimPacket struct {
	PackageLength   int
	HeaderLength    int
	ProtocolVersion int
	OperationCode   int
	SequenceId      int
	Content         []byte
}

func New(version, operationCode, sequenceId int, content []byte) *GoimPacket {
	return &GoimPacket{
		PackageLength:   len(content) + getProtocolHeaderSize(),
		HeaderLength:    getProtocolHeaderSize(),
		ProtocolVersion: version,
		OperationCode:   operationCode,
		SequenceId:      sequenceId,
		Content:         content,
	}
}

func Encoder(packet *GoimPacket) []byte {
	request := make([]byte, packet.PackageLength)

	binary.BigEndian.PutUint32(
		request[:PACKAGE_LENGTH_DATA_SIZE],
		uint32(packet.PackageLength),
	)

	binary.BigEndian.PutUint16(
		request[PACKAGE_LENGTH_DATA_SIZE:PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE],
		uint16(getProtocolHeaderSize()),
	)

	binary.BigEndian.PutUint16(
		request[PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE:PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE+PROTOCOL_VERSION_DATA_SIZE],
		uint16(packet.ProtocolVersion),
	)

	binary.BigEndian.PutUint32(
		request[PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE+PROTOCOL_VERSION_DATA_SIZE:PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE+PROTOCOL_VERSION_DATA_SIZE+OPERATION_CODE_DATA_SIZE],
		uint32(packet.OperationCode),
	)

	binary.BigEndian.PutUint32(
		request[PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE+PROTOCOL_VERSION_DATA_SIZE+OPERATION_CODE_DATA_SIZE:PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE+PROTOCOL_VERSION_DATA_SIZE+OPERATION_CODE_DATA_SIZE+SEQUENCE_ID_DATA_SIZE],
		uint32(packet.SequenceId),
	)

	copy(request[getProtocolHeaderSize():], packet.Content)

	return request
}

func Decoder(message []byte) (*GoimPacket, error) {
	if len(message) <= getProtocolHeaderSize() {
		return nil, ErrorBrokenPacket
	}

	packageLength := binary.BigEndian.Uint32(message[:PACKAGE_LENGTH_DATA_SIZE])

	headerLength := binary.BigEndian.Uint16(message[PACKAGE_LENGTH_DATA_SIZE : PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE])

	protocolVersion := binary.BigEndian.Uint16(message[PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE : PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE+PROTOCOL_VERSION_DATA_SIZE])

	operationCode := binary.BigEndian.Uint32(message[PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE+PROTOCOL_VERSION_DATA_SIZE : PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE+PROTOCOL_VERSION_DATA_SIZE+OPERATION_CODE_DATA_SIZE])

	sequenceId := binary.BigEndian.Uint32(message[PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE+PROTOCOL_VERSION_DATA_SIZE+OPERATION_CODE_DATA_SIZE : PACKAGE_LENGTH_DATA_SIZE+HEADER_LENGTH_DATA_SIZE+PROTOCOL_VERSION_DATA_SIZE+OPERATION_CODE_DATA_SIZE+SEQUENCE_ID_DATA_SIZE])

	content := message[getProtocolHeaderSize():]
	return &GoimPacket{
		PackageLength:   int(packageLength),
		HeaderLength:    int(headerLength),
		ProtocolVersion: int(protocolVersion),
		OperationCode:   int(operationCode),
		SequenceId:      int(sequenceId),
		Content:         content,
	}, nil
}

func getProtocolHeaderSize() int {
	return PACKAGE_LENGTH_DATA_SIZE +
		HEADER_LENGTH_DATA_SIZE +
		PROTOCOL_VERSION_DATA_SIZE +
		OPERATION_CODE_DATA_SIZE +
		SEQUENCE_ID_DATA_SIZE
}
