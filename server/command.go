package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	trigger = "form"
)

// type Field struct {
// 	Title string `json:"title"`
// 	Ref   string `json:"ref"`
// 	Type  string `json:"type"`
// }

// type Form struct {
// 	Title  string  `json:"title"`
// 	Fields []Field `json:"fields"`
// }

// ExecuteCommand
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
						DisplayName: "Experience",
						Name:        "experience",
						Type:        "text",
						SubType:     "text",
					},
					{
						DisplayName: "Location",
						Name:        "location",
						Type:        "text",
						SubType:     "text",
					},
					{
						DisplayName: "Name",
						Placeholder: "Include this field",
						Name:        "name",
						Type:        "bool",
						Optional:    true,
					},
					{
						DisplayName: "Email",
						Placeholder: "Include this field",
						Name:        "email",
						Type:        "bool",
						Optional:    true,
					},
					{
						DisplayName: "Resume",
						Placeholder: "Include this field",
						Name:        "resume",
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
	// userID := args.UserId
	// channelID := args.ChannelId

	// form := &Form{
	// 	Title: "Application for XYZ company",
	// 	Fields: []Field{
	// 		Field{
	// 			Title: "Name",
	// 			Ref:   "name",
	// 			Type:  "short-text",
	// 		},
	// 		Field{
	// 			Title: "Email",
	// 			Ref:   "email",
	// 			Type:  "short-text",
	// 		},
	// 		Field{
	// 			Title: "Resume drive link",
	// 			Ref:   "resume",
	// 			Type:  "short-text",
	// 		},
	// 	},
	// }

	// requestBody, err1 := json.Marshal(form)

	// if err1 != nil {
	// 	log.Fatalln(err1)
	// }

	// timeout := time.Duration(5 * time.Second)
	// client := http.Client{
	// 	Timeout: timeout,
	// }

	// request, err2 := http.NewRequest("POST", "https://api.typeform.com/forms", bytes.NewBuffer(requestBody))
	// request.Header.Set("Authorization", "Bearer 4aCgqf71W9eQUyDKpvoukbkrzsyW3YSsqvddfNJWHiAV")

	// if err2 != nil {
	// 	log.Fatalln(err2)
	// }

	// resp, err3 := client.Do(request)

	// if err3 != nil {
	// 	log.Fatalln(err3)
	// }

	// body, err4 := ioutil.ReadAll(resp.Body)

	// if err4 != nil {
	// 	log.Fatalln(err4)
	// }

	// postModel := &model.Post{
	// 	UserId:    userID,
	// 	ChannelId: channelID,
	// 	Message:   string(body),
	// }

	// _, err5 := p.API.CreatePost(postModel)
	// if err5 != nil {
	// 	return nil, err5
	// }

	return &model.CommandResponse{}, nil
}
