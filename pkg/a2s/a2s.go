// Package a2s provides tools for parsing A2S queries
package a2s

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	ErrA2SSplitPacket        = errors.New("a2s split packets are not supported")
	ErrA2SInvalidResponse    = errors.New("invalid a2s response")
	ErrA2SUnsupportedPacket  = errors.New("unsupported a2s packet")
	ErrA2SChallengeMalformed = errors.New("malformed a2s challenge")
)

const (
	a2sHeaderSinglePacket = uint32(0xffffffff)
	a2sHeaderSplitPacket  = uint32(0xfffffffe)

	a2sInfoRequestType  = byte(0x54)
	a2sInfoResponseType = byte(0x49)
	a2sChallengeType    = byte(0x41)
)

var a2sInfoPayload = []byte("Source Engine Query\x00")

type A2SInfo struct {
	Name       string
	Map        string
	Folder     string
	Game       string
	AppID      uint16
	Players    int
	MaxPlayers int
	Bots       int
	ServerType byte
	Platform   byte
	Visibility byte
	VAC        byte
	Version    string
}

func QueryInfo(conn net.Conn) (A2SInfo, error) {
	body, packetType, err := sendInfo(conn, nil)
	if err != nil {
		return A2SInfo{}, err
	}

	if packetType == a2sChallengeType {
		if len(body) < 4 {
			return A2SInfo{}, ErrA2SChallengeMalformed
		}

		challenge := body[:4]

		body, packetType, err = sendInfo(conn, challenge)
		if err != nil {
			return A2SInfo{}, err
		}
	}

	if packetType != a2sInfoResponseType {
		return A2SInfo{}, fmt.Errorf("%w: 0x%x", ErrA2SUnsupportedPacket, packetType)
	}

	return parseInfo(body)
}

func sendInfo(conn net.Conn, challenge []byte) ([]byte, byte, error) {
	request := make([]byte, 0, 5+len(a2sInfoPayload)+len(challenge))
	request = append(request, 0xff, 0xff, 0xff, 0xff)
	request = append(request, a2sInfoRequestType)
	request = append(request, a2sInfoPayload...)
	request = append(request, challenge...)

	if _, err := conn.Write(request); err != nil {
		return nil, 0, err
	}

	var buf [4096]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		return nil, 0, err
	}

	if n < 5 {
		return nil, 0, ErrA2SInvalidResponse
	}

	header := binary.LittleEndian.Uint32(buf[:])
	if header == a2sHeaderSplitPacket {
		return nil, 0, ErrA2SSplitPacket
	}

	if header != a2sHeaderSinglePacket {
		return nil, 0, ErrA2SInvalidResponse
	}

	packetType := buf[4]
	body := make([]byte, n-5)
	copy(body, buf[5:n])

	return body, packetType, nil
}

func parseInfo(data []byte) (A2SInfo, error) {
	reader := a2sReader{data: data}

	if _, err := reader.readByte(); err != nil {
		return A2SInfo{}, err
	}

	name, err := reader.readString()
	if err != nil {
		return A2SInfo{}, err
	}

	mapName, err := reader.readString()
	if err != nil {
		return A2SInfo{}, err
	}

	folder, err := reader.readString()
	if err != nil {
		return A2SInfo{}, err
	}

	game, err := reader.readString()
	if err != nil {
		return A2SInfo{}, err
	}

	appID, err := reader.readUint16()
	if err != nil {
		return A2SInfo{}, err
	}

	players, err := reader.readByte()
	if err != nil {
		return A2SInfo{}, err
	}

	maxPlayers, err := reader.readByte()
	if err != nil {
		return A2SInfo{}, err
	}

	bots, err := reader.readByte()
	if err != nil {
		return A2SInfo{}, err
	}

	serverType, err := reader.readByte()
	if err != nil {
		return A2SInfo{}, err
	}

	platform, err := reader.readByte()
	if err != nil {
		return A2SInfo{}, err
	}

	visibility, err := reader.readByte()
	if err != nil {
		return A2SInfo{}, err
	}

	vac, err := reader.readByte()
	if err != nil {
		return A2SInfo{}, err
	}

	version, err := reader.readString()
	if err != nil {
		return A2SInfo{}, err
	}

	return A2SInfo{
		Name:       name,
		Map:        mapName,
		Folder:     folder,
		Game:       game,
		AppID:      appID,
		Players:    int(players),
		MaxPlayers: int(maxPlayers),
		Bots:       int(bots),
		ServerType: serverType,
		Platform:   platform,
		Visibility: visibility,
		VAC:        vac,
		Version:    version,
	}, nil
}

type a2sReader struct {
	data   []byte
	offset int
}

func (r *a2sReader) readByte() (byte, error) {
	if r.offset >= len(r.data) {
		return 0, io.ErrUnexpectedEOF
	}

	value := r.data[r.offset]
	r.offset++

	return value, nil
}

func (r *a2sReader) readUint16() (uint16, error) {
	if r.offset+2 > len(r.data) {
		return 0, io.ErrUnexpectedEOF
	}

	value := binary.LittleEndian.Uint16(r.data[r.offset : r.offset+2])
	r.offset += 2

	return value, nil
}

func (r *a2sReader) readString() (string, error) {
	start := r.offset

	for r.offset < len(r.data) {
		if r.data[r.offset] == 0 {
			value := string(r.data[start:r.offset])
			r.offset++
			return value, nil
		}

		r.offset++
	}

	return "", io.ErrUnexpectedEOF
}
