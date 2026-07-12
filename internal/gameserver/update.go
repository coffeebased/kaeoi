package gameserver

import (
	"errors"
	"fmt"
	"strings"
)

type Update struct {
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

func (r *Update) Normalize() {
	if r.Host != nil {
		host := strings.ToUpper(strings.TrimSpace(*r.Host))
		r.Host = &host
	}

	if r.QueryHost != nil {
		queryHost := strings.ToUpper(strings.TrimSpace(*r.QueryHost))
		r.QueryHost = &queryHost
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

func (r Update) Validate() error {
	if r.Host != nil && *r.Host == "" {
		return errors.New("host is required")
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
