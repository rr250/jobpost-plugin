package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const (
	trigger = "jobpost"
)

const (
	botUserName    = "jobpostbot"
	botDisplayName = "Job Post Bot"
)

// OnActivate register the plugin command
func (p *Plugin) OnActivate() error {
	p.API.RegisterCommand(&model.Command{
		Trigger:          trigger,
		Description:      "Command for JobPost Plugin",
		DisplayName:      "Command for JobPost Plugin",
		AutoComplete:     true,
		AutoCompleteDesc: "Type /jobpost and press enter. For more commands type /jobpost help",
		AutoCompleteHint: "Command for JobPost Plugin",
	})
	botUserID, err := p.ensureBotExists()
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot user")
	}
	p.botUserID = botUserID
	p.router = p.InitAPI()
	p.driveService, p.sheetsService = p.InitGoogleServices()
	return nil
}

func (p *Plugin) ensureBotExists() (string, error) {
	bot := &model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
	}

	return p.Helpers.EnsureBot(bot)
}
