package gameserver

import (
	"net"
	"strconv"
	"strings"
	"time"
)

type Kind string

const (
	KindNone          Kind = "none"
	KindGenericTCP    Kind = "generic_tcp"
	KindGenericUDP    Kind = "generic_udp"
	KindSteamA2S      Kind = "steam_a2s"
	KindMinecraftJava Kind = "minecraft_java"
)

func (k Kind) valid() bool {
	return k == KindNone ||
		k == KindGenericTCP ||
		k == KindGenericUDP ||
		k == KindSteamA2S ||
		k == KindMinecraftJava
}

type Status string

const (
	StatusUnknown Status = "unknown"
	StatusOnline  Status = "online"
	StatusOffline Status = "offline"
)

func (s Status) valid() bool {
	return s == StatusUnknown ||
		s == StatusOnline ||
		s == StatusOffline
}

type CreateRequest struct {
	Code         string
	Host         string
	Port         int
	QueryHost    string
	QueryPort    int
	Kind         Kind
	Password     string
	DisplaysIP   bool
	DisplaysPort bool
	Application  string
	Version      string
	MaxPlayers   int
	Title        string
	Description  string
}

type UpdateRequest struct {
	Host         *string
	Port         *int
	QueryHost    *string
	QueryPort    *int
	Kind         *Kind
	Password     *string
	DisplaysIP   *bool
	DisplaysPort *bool
	Application  *string
	Version      *string
	MaxPlayers   *int
	Title        *string
	Description  *string
}

type Defaults struct {
	Application string
	Version     string
	MaxPlayers  int
	Title       string
	Description string
}

type State struct {
	IP          string
	Application string
	Version     string
	MaxPlayers  int
	PlayerCount int
	Title       string
	Description string
	Message     string
	Status      Status
	Since       time.Time
}

type GameServer struct {
	Code         string
	Host         string
	Port         int
	QueryHost    string
	QueryPort    int
	Kind         Kind
	Password     string
	DisplaysIP   bool
	DisplaysPort bool
	Defaults     Defaults
	State        State
}

func (gs *GameServer) GetDisplayAddress() string {
	host := strings.TrimSpace(gs.Host)
	if host == "" {
		return ""
	}

	if gs.DisplaysIP && net.ParseIP(host) == nil {
		statusIP := strings.TrimSpace(gs.State.IP)
		if net.ParseIP(statusIP) == nil {
			return "unreachable"
		}

		host = statusIP
	}

	if gs.DisplaysPort {
		return net.JoinHostPort(host, strconv.Itoa(gs.Port))
	}

	return host
}
