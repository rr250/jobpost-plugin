package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const (
	trigger = "form"
)

const (
	botUserName    = "jobpostbot"
	botDisplayName = "Job Post Bot"
)

// OnActivate register the plugin command
func (p *Plugin) OnActivate() error {
	p.API.RegisterCommand(&model.Command{
		Trigger:          trigger,
		Description:      "Make a form",
		DisplayName:      "Make a form",
		AutoComplete:     true,
		AutoCompleteDesc: "Make a form (Use it by clicking reply first then slash command)",
		AutoCompleteHint: "Make a form",
	})
	botUserID, err := p.ensureBotExists()
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot user")
	}
	p.botUserID = botUserID
	p.router = p.InitAPI()
	return nil
}

func (p *Plugin) ensureBotExists() (string, error) {
	bot := &model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
	}

	return p.Helpers.EnsureBot(bot)
}
