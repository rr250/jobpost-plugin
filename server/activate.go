package main

import "github.com/mattermost/mattermost-server/v5/model"

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
	p.router = p.InitAPI()
	return nil
}
