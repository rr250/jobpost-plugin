package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/dialog", p.handleDialog).Methods("POST")
	return r
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) handleDialog(w http.ResponseWriter, req *http.Request) {

	request := model.SubmitDialogRequestFromJson(req.Body)

	// user, uErr := p.API.GetUser(request.UserId)
	// if uErr != nil {
	// 	p.API.LogError(uErr.Error())
	// 	return
	// }

	company := request.Submission["company"]
	position := request.Submission["position"]
	description := request.Submission["description"]
	skills := request.Submission["skills"]
	experience := request.Submission["experience"]
	location := request.Submission["location"]
	anonymous := request.Submission["anonymous"]
	companyStr := company.(string)
	positionStr := position.(string)
	descriptionStr := description.(string)
	skillsStr := skills.(string)
	experienceStr := experience.(string)
	locationStr := location.(string)
	var userID string
	if anonymous.(bool) {
		userID = p.botUserID
	} else {
		userID = request.UserId
	}
	postModel := &model.Post{
		UserId:    userID,
		ChannelId: request.ChannelId,
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{
				{
					Text: "Company: " + companyStr + "\nPositon: " + positionStr + "\nDescription: " + descriptionStr + "\nSkills: " + skillsStr + "\nExperience Required: " + experienceStr + "\nLocation: " + locationStr,
					Actions: []*model.PostAction{
						{
							Integration: &model.PostActionIntegration{
								// Context: model.StringInterface{
								// 	"reminder_id":   r.Reminder.Id,
								// 	"occurrence_id": r.Reminder.Occurrences[0].Id,
								// 	"action":        "delete/ephemeral",
								// },
								URL: fmt.Sprintf("/plugins/%s/fillform  ", manifest.ID),
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: "Fill the form",
						},
					},
				},
			},
		},
	}

	_, err5 := p.API.CreatePost(postModel)
	if err5 != nil {
		log.Fatalln(err5)
	}

}
