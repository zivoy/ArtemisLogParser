package logparser

import (
	"encoding/binary"
	"errors"
	"io"
	"os"

	"google.golang.org/protobuf/proto"

	pb "artemisLogParser/protobuf"
)

var magic = []byte{0x89, 0x41, 0x52, 0x54}

func readBytes(reader io.Reader, n int) ([]byte, error) {
	bytes := make([]byte, n)
	n, err := reader.Read(bytes)
	return bytes[:n], err
}

func isLogFile(reader io.Reader) bool {
	bytes, err := readBytes(reader, len(magic))
	if err != nil || len(bytes) != len(magic) {
		return false
	}
	for i, b := range bytes {
		if magic[i] != b {
			return false
		}
	}
	return true
}

func Read(file *os.File) (*Game, error) {
	if !isLogFile(file) {
		return nil, errors.New("not a log file")
	}

	game := newGame()

	data, err := readGame(file)
	if err != nil {
		return nil, err
	}
	game.Data = data

	for {
		event, err := readEvent(file)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		game.appendEvent(event)
	}
	return game, nil
}

func readMessage(reader io.Reader) ([]byte, error) {
	length, err := readBytes(reader, 2)
	if err != nil {
		return nil, err
	}
	if len(length) != 2 {
		return nil, errors.New("error reading length")
	}

	n := int(binary.LittleEndian.Uint16(length))
	message, err := readBytes(reader, n)
	return message, err
}

func readGame(reader io.Reader) (*pb.Game, error) {
	gameMessage, err := readMessage(reader)
	if err != nil {
		return nil, err
	}

	game := &pb.Game{}
	err = proto.Unmarshal(gameMessage, game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func readEvent(reader io.Reader) (*pb.AnalyticsEvent, error) {
	eventMessage, err := readMessage(reader)
	if err != nil {
		return nil, err
	}

	event := &pb.AnalyticsEvent{}
	err = proto.Unmarshal(eventMessage, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}
