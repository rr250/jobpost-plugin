package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	command := strings.Trim(args.Command, " ")

	if strings.Trim(command, " ") == "/"+trigger {

		dialogRequest := model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("/plugins/%s/dialog", manifest.ID),
			Dialog: model.Dialog{
				Title:       "New Job Posting",
				CallbackId:  model.NewId(),
				SubmitLabel: "Create Form",
				Elements: []model.DialogElement{
					{
						DisplayName: "Company",
						Name:        "company",
						Type:        "text",
						SubType:     "text",
					},
					{
						DisplayName: "Job Position",
						Name:        "position",
						Type:        "text",
						SubType:     "text",
					},
					{
						DisplayName: "Job Description",
						Name:        "description",
						Type:        "textarea",
						SubType:     "text",
					},
					{
						DisplayName: "Skills Required",
						Name:        "skills",
						Type:        "text",
						SubType:     "text",
					},
					{
						DisplayName: "Minimum Experience",
						Name:        "minExperience",
						Placeholder: "Write only the year",
						Type:        "text",
						SubType:     "number",
					},
					{
						DisplayName: "Maximum Experience",
						Name:        "maxExperience",
						Placeholder: "Write only the year",
						Type:        "text",
						SubType:     "number",
					},
					{
						DisplayName: "Location",
						Name:        "location",
						Type:        "text",
						SubType:     "text",
					},
					{
						DisplayName: "Resume",
						Placeholder: "Include this field",
						Name:        "resume",
						Type:        "bool",
						Optional:    true,
					},
					{
						DisplayName: "Experience",
						Placeholder: "Include this field",
						Name:        "experience",
						Type:        "bool",
						Optional:    true,
					},
					{
						DisplayName: "Reason on why are you interested",
						Placeholder: "Include this field",
						Name:        "reason",
						Type:        "bool",
						Optional:    true,
					},
					{
						DisplayName: "Post as Anonymous?",
						Placeholder: "Yes, Post as anonymous",
						Name:        "anonymous",
						Type:        "bool",
						Optional:    true,
					},
				},
			},
		}
		if pErr := p.API.OpenInteractiveDialog(dialogRequest); pErr != nil {
			p.API.LogError("Failed opening interactive dialog " + pErr.Error())
		}
	}
	return &model.CommandResponse{}, nil
}
