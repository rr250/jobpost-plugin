package main

import (
	"fmt"
	"log"
	"strconv"
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
				SubmitLabel: "Create Jobpost",
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
						Placeholder: "years of experience",
						Type:        "text",
						SubType:     "number",
					},
					{
						DisplayName: "Maximum Experience",
						Name:        "maxExperience",
						Placeholder: "years of experience",
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
						DisplayName: "Please check the fields which you want to include in your jobpost form.",
						Placeholder: "Resume",
						Name:        "resume",
						Type:        "bool",
						Optional:    true,
					},
					{
						Placeholder: "Experience",
						Name:        "experience",
						Type:        "bool",
						Optional:    true,
					},
					{
						Placeholder: "Reason on why is applicant interested",
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
			postModel := &model.Post{
				UserId:    args.UserId,
				ChannelId: args.ChannelId,
				Message:   fmt.Sprintf("Failed opening interactive dialog " + pErr.Error()),
			}
			p.API.SendEphemeralPost(args.UserId, postModel)
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
					Text: "Jobpost: " + jobpost.Details + "\nCreatedAt: " + jobpost.CreatedAt.Format(time.ANSIC),
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
							Name: "Fetch Applicants",
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
				Message:   err.(string),
			}
			p.API.SendEphemeralPost(args.UserId, postModel)
		}
	} else if splitText := strings.Split(strings.Trim(command, " "), " "); splitText[1] == "subscribe" {
		log.Println("subscribe")
		year, err1 := strconv.Atoi(splitText[2])
		if err1 != nil {
			postModel := &model.Post{
				UserId:    args.UserId,
				ChannelId: args.ChannelId,
				Message:   "Wrong Command Syntax. Subscribe Command Example: /" + trigger + " subscribe 2 years",
			}
			p.API.SendEphemeralPost(args.UserId, postModel)
		} else {
			err := p.subscribeToExperience(args.UserId, year)
			if err == nil {
				postModel := &model.Post{
					UserId:    args.UserId,
					ChannelId: args.ChannelId,
					Message:   "Subscribed",
				}
				p.API.SendEphemeralPost(args.UserId, postModel)
			} else {
				postModel := &model.Post{
					UserId:    args.UserId,
					ChannelId: args.ChannelId,
					Message:   err.(string),
				}
				p.API.SendEphemeralPost(args.UserId, postModel)
			}
		}
	} else if splitText := strings.Split(strings.Trim(command, " "), " "); splitText[1] == "unsubscribe" {
		log.Println("unsubscribe")
		year, err1 := strconv.Atoi(splitText[2])
		if err1 != nil {
			postModel := &model.Post{
				UserId:    args.UserId,
				ChannelId: args.ChannelId,
				Message:   "Wrong Command Syntax. Unsubscribe Command Example: /" + trigger + " unsubscribe 2 years",
			}
			p.API.SendEphemeralPost(args.UserId, postModel)
		} else {
			err := p.unSubscribeToExperience(args.UserId, year)
			if err == nil {
				postModel := &model.Post{
					UserId:    args.UserId,
					ChannelId: args.ChannelId,
					Message:   "Unsubscribed",
				}
				p.API.SendEphemeralPost(args.UserId, postModel)
			} else {
				postModel := &model.Post{
					UserId:    args.UserId,
					ChannelId: args.ChannelId,
					Message:   err.(string),
				}
				p.API.SendEphemeralPost(args.UserId, postModel)
			}
		}
	} else if strings.Trim(command, " ") == "/"+trigger+" help" {
		postModel := &model.Post{
			UserId:    args.UserId,
			ChannelId: args.ChannelId,
			Message:   "* `/jobpost` - opens up an [interactive dialog] to create a Jobpost \n* `/jobpost list` - displays a list of jobposts created by you \n* `/jobpost subscribe x years` - subscribes to jobposts which requires x years of experience where x is an integer \n* `/jobpost unsubscribe x years` - unsubscribes to jobposts which requires x years of experience where x is an integer",
		}
		p.API.SendEphemeralPost(args.UserId, postModel)
	}

	return &model.CommandResponse{}, nil
}
