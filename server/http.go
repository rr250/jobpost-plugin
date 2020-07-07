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
	// positon := request.Submission["positon"]
	// description := request.Submission["description"]
	// skills := request.Submission["skills"]
	// experience := request.Submission["experience"]
	// location := request.Submission["location"]
	// name := request.Submission["name"]
	// email := request.Submission["email"]
	// resume := request.Submission["resume"]
	// reason := request.Submission["reason"]

	// strs := []string{company.(string), positon.(string), description.(string), skills.(string), experience.(string), location.(string)}
	// postText := strings.Join(strs, "\n")
	// log.Printf(postText)
	postModel := &model.Post{
		UserId:    request.UserId,
		ChannelId: request.ChannelId,
		Message:   company.(string),
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{
				{
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
