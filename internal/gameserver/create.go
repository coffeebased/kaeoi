package gameserver

import (
	"errors"
	"fmt"
	"strings"
)

type Create struct {
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

func (r *Create) Normalize() {
	r.Code = strings.ToUpper(strings.TrimSpace(r.Code))
	r.Host = strings.TrimSpace(r.Host)
	r.QueryHost = strings.TrimSpace(r.QueryHost)
	r.Password = strings.TrimSpace(r.Password)
	r.Application = strings.TrimSpace(r.Application)
	r.Version = strings.TrimSpace(r.Version)
	r.Title = strings.TrimSpace(r.Title)
	r.Description = strings.TrimSpace(r.Description)
}

func (r Create) Validate() error {
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
