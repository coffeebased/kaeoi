// Package gameserver provides monitoring functionality for game servers
package gameserver

import (
	"net"
	"strconv"
	"strings"
	"time"
)

type Kind string

const (
	KindSteamA2S      Kind = "steam_a2s"
	KindMinecraftJava Kind = "minecraft_java"
)

func (k Kind) valid() bool {
	return k == KindSteamA2S ||
		k == KindMinecraftJava
}

type Status string

const (
	StatusUnknown Status = "unknown"
	StatusOnline  Status = "online"
	StatusOffline Status = "offline"
)

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
	Code             string
	Host             string
	Port             int
	QueryHost        string
	QueryPort        int
	Kind             Kind
	Password         string
	DisplaysIP       bool
	DisplaysPort     bool
	Ignore           bool
	OverrideDefaults bool
	Defaults         Defaults
	State            State
}

func (gs GameServer) GetDisplayAddress() string {
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

func (gs GameServer) GetQueryAddress() string {
	host := gs.Host
	if gs.QueryHost != "" {
		host = gs.QueryHost
	}

	port := gs.Port
	if gs.QueryPort != 0 {
		port = gs.QueryPort
	}

	return net.JoinHostPort(host, strconv.Itoa(port))
}
