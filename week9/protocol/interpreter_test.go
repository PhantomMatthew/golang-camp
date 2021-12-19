package protocol

import (
	"log"
	"testing"
)

func TestDecoder(t *testing.T) {
	content := "This is a goim packet"
	version := 1
	operation := 1
	sequence := 0
	packet := New(version, operation, sequence, []byte(content))

	request := Encoder(packet)
	response, err := Decoder(request)

	// validate decode result
	if err != nil ||
		packet.PackageLength != getProtocolHeaderSize()+len(content) ||
		packet.HeaderLength != getProtocolHeaderSize() ||
		packet.ProtocolVersion != version ||
		packet.OperationCode != operation ||
		packet.SequenceId != sequence ||
		string(packet.Content) != content {
		t.Error("NOT PASS")
	}
	log.Printf("%#v", response)
}
