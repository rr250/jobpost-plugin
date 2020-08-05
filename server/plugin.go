package main

import (
	"sync"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

type Plugin struct {
	plugin.MattermostPlugin
	router            *mux.Router
	botUserID         string
	configurationLock sync.RWMutex
	configuration     *configuration
	driveService      *drive.Service
	sheetsService     *sheets.Service
}
