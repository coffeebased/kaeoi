package gameserver

import (
	"errors"
	"fmt"
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
	Application      string
	Version          string
	MaxPlayers       int
	Title            string
	Description      string
}

func (r *CreateRequest) Normalize() {
	r.Code = strings.ToUpper(strings.TrimSpace(r.Code))
	r.Host = strings.TrimSpace(r.Host)
	r.QueryHost = strings.TrimSpace(r.QueryHost)

	if r.Kind == "" {
		r.Kind = KindNone
	}

	r.Password = strings.TrimSpace(r.Password)
	r.Application = strings.TrimSpace(r.Application)
	r.Version = strings.TrimSpace(r.Version)
	r.Title = strings.TrimSpace(r.Title)
	r.Description = strings.TrimSpace(r.Description)
}

func (r CreateRequest) Validate() error {
	if r.Code == "" {
		return errors.New("code is required")
	}

	if r.Code != strings.ToUpper(r.Code) {
		return errors.New("code must be upper case")
	}

	if r.Host == "" {
		return errors.New("host is required")
	}

	if r.Port <= 0 || r.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}

	if r.QueryPort < 0 || r.QueryPort > 65535 {
		return errors.New("query port must be between 0 and 65535")
	}

	if !r.Kind.valid() {
		return fmt.Errorf("invalid kind: %s", r.Kind)
	}

	if r.QueryHost != "" && r.QueryPort == 0 {
		return errors.New("query port must be provided along query host")
	}

	return nil
}

type UpdateRequest struct {
	Host             *string
	Port             *int
	QueryHost        *string
	QueryPort        *int
	Kind             *Kind
	Password         *string
	DisplaysIP       *bool
	DisplaysPort     *bool
	Ignore           *bool
	OverrideDefaults *bool
	Application      *string
	Version          *string
	MaxPlayers       *int
	Title            *string
	Description      *string
}

func (r *UpdateRequest) Normalize() {
	if r.Host != nil {
		host := strings.ToUpper(strings.TrimSpace(*r.Host))
		r.Host = &host
	}

	if r.QueryHost != nil {
		queryHost := strings.ToUpper(strings.TrimSpace(*r.QueryHost))
		r.QueryHost = &queryHost
	}

	if r.Kind != nil && *r.Kind == "" {
		kindNone := KindNone
		r.Kind = &kindNone
	}

	if r.Password != nil {
		password := strings.ToUpper(strings.TrimSpace(*r.Password))
		r.Password = &password
	}

	if r.Application != nil {
		application := strings.ToUpper(strings.TrimSpace(*r.Application))
		r.Application = &application
	}

	if r.Version != nil {
		version := strings.ToUpper(strings.TrimSpace(*r.Version))
		r.Version = &version
	}

	if r.Title != nil {
		title := strings.ToUpper(strings.TrimSpace(*r.Title))
		r.Title = &title
	}

	if r.Description != nil {
		description := strings.ToUpper(strings.TrimSpace(*r.Description))
		r.Description = &description
	}
}

func (r UpdateRequest) Validate() error {
	if r.Host != nil && *r.Host == "" {
		return errors.New("Host is required")
	}

	if r.Port != nil && (*r.Port <= 0 || *r.Port > 65535) {
		return errors.New("port must be between 1 and 65535")
	}

	if r.QueryPort != nil && (*r.QueryPort < 0 || *r.QueryPort > 65535) {
		return errors.New("query port must be between 0 and 65535")
	}

	if r.Kind != nil && !r.Kind.valid() {
		return fmt.Errorf("invalid kind: %s", *r.Kind)
	}

	if r.QueryHost != nil && *r.QueryHost != "" && (r.QueryPort == nil || *r.QueryPort == 0) {
		return errors.New("query port must be provided along query host")
	}

	if r.QueryPort != nil && *r.QueryPort == 0 && (r.QueryHost == nil || *r.QueryHost != "") {
		return errors.New("query host must be set to empty string when query port is set to 0")
	}

	return nil
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
