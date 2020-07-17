package main

import (
	"fmt"
	"log"
	"strings"
	"time"

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
	} else if strings.Trim(command, " ") == "/"+trigger+" list" {
		log.Println("list")
		jobposts, err := p.getJobsPerUser(args.UserId)
		if err == nil {
			postModel := &model.Post{
				UserId:    args.UserId,
				ChannelId: args.ChannelId,
				Message:   "Jobposts :-",
				Props: model.StringInterface{
					"attachments": []*model.SlackAttachment{},
				},
			}
			for _, jobpost := range jobposts {
				attachment := &model.SlackAttachment{
					Text: "Jobpost:" + jobpost.Details + "\nCreatedAt:" + jobpost.CreatedAt.Format(time.ANSIC),
					Actions: []*model.PostAction{
						{
							Integration: &model.PostActionIntegration{
								URL: fmt.Sprintf("/plugins/%s/getjobpostbyid", manifest.ID),
								Context: model.StringInterface{
									"action":    "getjobpostbyid",
									"jobpostid": jobpost.JobpostID,
								},
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: "Fetch Jobpost",
						},
					},
				}
				postModel.Props["attachments"] = append(postModel.Props["attachments"].([]*model.SlackAttachment), attachment)
			}
			p.API.SendEphemeralPost(args.UserId, postModel)

		} else {
			postModel := &model.Post{
				UserId:    args.UserId,
				ChannelId: args.ChannelId,
				Message:   "No Jobposts",
			}
			p.API.SendEphemeralPost(args.UserId, postModel)
		}
	}
	return &model.CommandResponse{}, nil
}
