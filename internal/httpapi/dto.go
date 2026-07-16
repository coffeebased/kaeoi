package httpapi

import (
	"time"

	"github.com/coffeebased/kaeoi/internal/gameserver"
)

type gameServerResponse struct {
	Code      string    `json:"code"`
	Address   string    `json:"address"`
	Password  string    `json:"password,omitempty"`
	State     state     `json:"state"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type state struct {
	Application string    `json:"application,omitempty"`
	Version     string    `json:"version,omitempty"`
	MaxPlayers  *int      `json:"maxPlayers,omitempty"`
	PlayerCount *int      `json:"playerCount,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Since       time.Time `json:"since"`
}

func newGameServerResponse(server gameserver.GameServer) gameServerResponse {
	state := state{
		Application: server.State.Application,
		Version:     server.State.Version,
		Title:       server.State.Title,
		Description: server.State.Description,
		Status:      string(server.State.Status),
		Since:       server.State.Since,
	}

	if server.State.MaxPlayers > 0 {
		state.MaxPlayers = &server.State.MaxPlayers
		state.PlayerCount = &server.State.PlayerCount
	}

	return gameServerResponse{
		Code:      server.Code,
		Address:   server.GetDisplayAddress(),
		Password:  server.Password,
		State:     state,
		UpdatedAt: server.Metadata.LastChange,
	}
}
