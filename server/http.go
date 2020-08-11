package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/dialog", p.handleDialog).Methods("POST")
	r.HandleFunc("/applytojob", p.applyToJob).Methods("POST")
	r.HandleFunc("/submit", p.submit).Methods("POST")
	r.HandleFunc("/getjobpostbyid", p.getJobPostByID).Methods("POST")
	r.HandleFunc("/downloadjobpostbyid", p.downloadJobPostByID).Methods("POST")
	return r
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) handleDialog(w http.ResponseWriter, req *http.Request) {

	request := model.SubmitDialogRequestFromJson(req.Body)

	jobpostIDUUID := uuid.New()
	jobpostID := jobpostIDUUID.String()

	company := request.Submission["company"]
	position := request.Submission["position"]
	description := request.Submission["description"]
	skills := request.Submission["skills"]
	minExperience := request.Submission["minExperience"]
	maxExperience := request.Submission["maxExperience"]
	location := request.Submission["location"]
	anonymous := request.Submission["anonymous"]
	companyStr := company.(string)
	positionStr := position.(string)
	descriptionStr := description.(string)
	skillsStr := skills.(string)
	minExperienceStr := minExperience.(string)
	maxExperienceStr := maxExperience.(string)
	locationStr := location.(string)
	minExperienceInt, err11 := strconv.Atoi(minExperienceStr)
	if err11 != nil {
		p.API.LogError("Not an integer", err11)
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("MinExperience is not an integer %s", err11),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
		return
	}
	maxExperienceInt, err12 := strconv.Atoi(maxExperienceStr)
	if err12 != nil {
		p.API.LogError("Not an integer", err12)
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("MaxExperience is not an integer %s", err12),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
		return
	}
	if minExperienceInt < 0 || maxExperienceInt < 0 || maxExperienceInt < minExperienceInt {
		p.API.LogError("Less than zero")
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("Experience less than zero or maxExperience less than minExperience"),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
		return
	}
	var userID string
	if anonymous.(bool) {
		userID = p.botUserID
	} else {
		userID = request.UserId
	}
	postModel1 := &model.Post{
		UserId:    userID,
		ChannelId: request.ChannelId,
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{
				{
					Text: "Company: " + companyStr + "\nPositon: " + positionStr + "\nDescription: " + descriptionStr + "\nSkills: " + skillsStr + "\nExperience Required: " + minExperienceStr + "-" + maxExperienceStr + " years" + "\nLocation: " + locationStr,
					Actions: []*model.PostAction{
						{
							Integration: &model.PostActionIntegration{
								URL: fmt.Sprintf("/plugins/%s/applytojob", manifest.ID),
								Context: model.StringInterface{
									"action":    "applyToJob",
									"submision": request.Submission,
									"jobpostid": jobpostID,
								},
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: "Apply",
						},
					},
				},
			},
		},
	}

	sheetCreate := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: companyStr + " - " + positionStr + " - " + jobpostID,
		},
	}
	sheet, err3 := p.sheetsService.Spreadsheets.Create(sheetCreate).Do()
	if err3 != nil {
		log.Fatalf("Unable to create sheet: %v", err3)
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("Unable to create sheet: %v", err3),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
		return
	}
	log.Println(sheet)
	readRange := "Sheet1!A:Z"
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{
			{
				"Name",
				"Email",
				"Resume",
				"Experience",
				"Reason",
				"FilledAt",
			},
		},
	}
	_, err7 := p.sheetsService.Spreadsheets.Values.Append(sheet.SpreadsheetId, readRange, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err7 != nil {
		log.Fatalf("Unable to append data from sheet: %v", err7)
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("Unable to append data from sheet: %v", err7),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
		return
	}
	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	_, err8 := p.driveService.Permissions.Create(sheet.SpreadsheetId, permission).Do()
	if err8 != nil {
		log.Fatalf("Unable to change drive permission %v", err8)
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("Unable to change drive permission %v", err8),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
	}

	jobpost := Jobpost{
		ID:            jobpostID,
		CreatedBy:     request.UserId,
		CreatedAt:     time.Now(),
		Company:       companyStr,
		Position:      positionStr,
		Description:   descriptionStr,
		Skills:        skillsStr,
		MinExperience: minExperienceInt,
		MaxExperience: maxExperienceInt,
		Location:      locationStr,
		ExperienceReq: request.Submission["experience"].(bool),
		SheetID:       sheet.SpreadsheetId,
		SheetURL:      sheet.SpreadsheetUrl,
	}
	err6 := p.addJobpost(jobpost)
	if err6 != nil {
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   err6.(string),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
		return
	}

	_, err5 := p.API.CreatePost(postModel1)
	if err5 != nil {
		p.API.LogError("failed to create post", err5)
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("failed to create post %s", err5),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
		return
	} else {
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   "Jobpost created. See all your jobposts by command `/jobpost list`",
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
	}

	channel, err9 := p.API.GetDirectChannel(request.UserId, p.botUserID)
	if err9 != nil {
		p.API.LogError("failed to get channel", err9)
	}
	postModel := &model.Post{
		UserId:    p.botUserID,
		ChannelId: channel.Id,
		Message:   fmt.Sprintf("Jobpost created: " + jobpost.Company + " - " + jobpost.Position + "\nCreatedAt: " + jobpost.CreatedAt.Format(time.ANSIC) + "\nTrack the responses through this sheet:-\n" + sheet.SpreadsheetUrl),
	}
	_, err10 := p.API.CreatePost(postModel)
	if err10 != nil {
		p.API.LogError("failed to create post", err10)
	}

	for i := jobpost.MinExperience; i <= jobpost.MaxExperience; i++ {
		p.sendToSubscribers(postModel1, i)
	}
}

func (p *Plugin) applyToJob(w http.ResponseWriter, req *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(req.Body)
	submision := request.Context["submision"].(map[string]interface{})
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	user, err := p.API.GetUser(request.UserId)
	userFullName := " "
	userEmail := "@gmail.com"
	if err != nil {
		p.API.LogError("failed to get user", err)
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("failed to get user %s", err),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
	} else {
		userFullName = user.FirstName + " " + user.LastName
		userEmail = user.Email
	}
	userDetails, err1 := p.getResume(request.UserId)
	resumeStr := " "
	if err1 == nil {
		resumeStr = userDetails.Resume
	}
	dialogRequest := model.OpenDialogRequest{
		TriggerId: request.TriggerId,
		URL:       fmt.Sprintf("/plugins/%s/submit", manifest.ID),
		Dialog: model.Dialog{
			Title:       submision["company"].(string) + " - " + submision["position"].(string),
			CallbackId:  model.NewId(),
			SubmitLabel: "Submit",
			State:       request.Context["jobpostid"].(string),
			Elements: []model.DialogElement{
				{
					DisplayName: "Name",
					Name:        "name",
					Type:        "text",
					SubType:     "text",
					Default:     userFullName,
				},
				{
					DisplayName: "Email",
					Name:        "email",
					Type:        "text",
					SubType:     "email",
					Default:     userEmail,
				},
				{
					DisplayName: "Resume",
					Name:        "resume",
					HelpText:    "Put the URL of your resume. Please make sure that it is accessible. To make sure that it is accessible try opening it in incognito mode.",
					Type:        "text",
					SubType:     "text",
					Default:     resumeStr,
					Optional:    !submision["resume"].(bool),
				},
				{
					DisplayName: "Reason on why are you interested",
					Name:        "reason",
					Type:        "textarea",
					SubType:     "text",
					Default:     " ",
					Optional:    !submision["reason"].(bool),
				},
				{
					DisplayName: "Experience",
					Name:        "experience",
					Placeholder: "years of experience",
					Type:        "text",
					SubType:     "text",
					Optional:    !submision["experience"].(bool),
				},
			},
		},
	}
	if pErr := p.API.OpenInteractiveDialog(dialogRequest); pErr != nil {
		p.API.LogError("Failed opening interactive dialog " + pErr.Error())
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("Failed opening interactive dialog " + pErr.Error()),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
	}
}

func (p *Plugin) submit(w http.ResponseWriter, req *http.Request) {

	request := model.SubmitDialogRequestFromJson(req.Body)
	experience := request.Submission["experience"].(string)
	experienceFlt, err1 := strconv.ParseFloat(experience, 64)
	if err1 != nil {
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("Experience format is not correct %s", err1),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
	}
	jobpostResponse := JobpostResponse{
		UserID:     request.UserId,
		Name:       request.Submission["name"].(string),
		Email:      request.Submission["email"].(string),
		Resume:     request.Submission["resume"].(string),
		Reason:     request.Submission["reason"].(string),
		Experience: experienceFlt,
		FilledAt:   time.Now(),
	}
	jobpost, err := p.addJobpostResponse(request.State, jobpostResponse)
	if err != nil {
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   err.(string),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
	} else {
		channel, err9 := p.API.GetDirectChannel(request.UserId, p.botUserID)
		if err9 != nil {
			p.API.LogError("failed to get channel", err9)
		}
		postModel := &model.Post{
			UserId:    p.botUserID,
			ChannelId: channel.Id,
			Message:   "Successfully Applied to " + jobpost.Company + " - " + jobpost.Position + " - " + jobpost.ID,
		}

		p.API.CreatePost(postModel)
		postModel.ChannelId = request.ChannelId
		p.API.SendEphemeralPost(request.UserId, postModel)
	}

}

func (p *Plugin) getJobPostByID(w http.ResponseWriter, req *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(req.Body)
	jobpostID := request.Context["jobpostid"].(string)
	log.Println(jobpostID)
	jobpost, err := p.getJobPost(jobpostID)
	if err == nil {
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   "Company: " + jobpost.Company + "\nPositon: " + jobpost.Position + "\nDescription: " + jobpost.Description + "\nSkills: " + jobpost.Skills + "\nExperience Required: " + strconv.Itoa(jobpost.MinExperience) + "-" + strconv.Itoa(jobpost.MaxExperience) + " years" + "\nLocation: " + jobpost.Location,
			Props: model.StringInterface{
				"attachments": []*model.SlackAttachment{},
			},
		}
		for _, jobpostResponse := range jobpost.JobpostResponses {
			attachment := &model.SlackAttachment{
				Text: "Name: " + jobpostResponse.Name + "\nEmail: " + jobpostResponse.Email + "\nResume: " + jobpostResponse.Resume + "\nExperience: " + fmt.Sprintf("%.1f years", jobpostResponse.Experience) + "\nReason:" + jobpostResponse.Reason,
			}
			postModel.Props["attachments"] = append(postModel.Props["attachments"].([]*model.SlackAttachment), attachment)
		}
		attachment := &model.SlackAttachment{
			Actions: []*model.PostAction{
				{
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("/plugins/%s/downloadjobpostbyid", manifest.ID),
						Context: model.StringInterface{
							"action":    "downloadjobpostbyid",
							"jobpostid": jobpost.ID,
						},
					},
					Type: model.POST_ACTION_TYPE_BUTTON,
					Name: "Download Responses as csv",
				},
			},
		}
		postModel.Props["attachments"] = append(postModel.Props["attachments"].([]*model.SlackAttachment), attachment)
		p.API.SendEphemeralPost(request.UserId, postModel)
	} else {
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   err.(string),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
	}
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func (p *Plugin) downloadJobPostByID(w http.ResponseWriter, req *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(req.Body)
	jobpostID := request.Context["jobpostid"].(string)
	log.Println(jobpostID)
	jobpost, err := p.getJobPost(jobpostID)
	if err == nil {
		buf := bytes.NewBuffer(nil)
		writer := csv.NewWriter(buf)
		err8 := writer.Write([]string{"Name", "Email", "Resume", "Experience", "Reason", "Applied At"})
		if err8 != nil {
			p.API.LogError("Cannot write to file", err)
			postModel := &model.Post{
				UserId:    request.UserId,
				ChannelId: request.ChannelId,
				Message:   fmt.Sprintf("Cannot write to file %s", err),
			}
			p.API.SendEphemeralPost(request.UserId, postModel)
		}
		for _, jobpostResponse := range jobpost.JobpostResponses {
			jobpostResponseCsv := []string{jobpostResponse.Name, jobpostResponse.Email, jobpostResponse.Resume, fmt.Sprintf("%.1f years", jobpostResponse.Experience), jobpostResponse.Reason, jobpostResponse.FilledAt.Local().Format(time.RFC3339Nano)}
			err1 := writer.Write(jobpostResponseCsv)
			if err1 != nil {
				p.API.LogError("Cannot write to file %s", err1)
				postModel := &model.Post{
					UserId:    request.UserId,
					ChannelId: request.ChannelId,
					Message:   fmt.Sprintf("Cannot write to file %s", err1),
				}
				p.API.SendEphemeralPost(request.UserId, postModel)
			}
		}
		writer.Flush()
		data := buf.Bytes()
		channel, err7 := p.API.GetDirectChannel(request.UserId, p.botUserID)
		var channelID string
		if err7 != nil {
			p.API.LogError("failed to get channel", err7)
			postModel := &model.Post{
				UserId:    request.UserId,
				ChannelId: request.ChannelId,
				Message:   fmt.Sprintf("failed to get channel  %s", err7),
			}
			p.API.SendEphemeralPost(request.UserId, postModel)
		}
		channelID = channel.Id
		fileInfo, err3 := p.API.UploadFile(data, channelID, jobpost.Company+"-"+jobpost.Position+"-"+strconv.Itoa(int(time.Now().Unix()))+".csv")
		if err3 != nil {
			p.API.LogError("Cannot upload file  %s", err3)
			postModel := &model.Post{
				UserId:    request.UserId,
				ChannelId: request.ChannelId,
				Message:   fmt.Sprintf("Cannot upload file  %s", err3),
			}
			p.API.SendEphemeralPost(request.UserId, postModel)
		}
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: channelID,
			FileIds:   []string{fileInfo.Id},
		}

		_, err6 := p.API.CreatePost(postModel)
		if err6 != nil {
			p.API.LogError("failed to create post", err6)
			postModel := &model.Post{
				UserId:    request.UserId,
				ChannelId: request.ChannelId,
				Message:   fmt.Sprintf("failed to create post %s", err6),
			}
			p.API.SendEphemeralPost(request.UserId, postModel)
		} else {
			postModel := &model.Post{
				UserId:    request.UserId,
				ChannelId: request.ChannelId,
				Message:   "CSV file sent as a direct message",
			}
			p.API.SendEphemeralPost(request.UserId, postModel)
		}
	} else {
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   err.(string),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
	}
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func writePostActionIntegrationResponseOk(w http.ResponseWriter, response *model.PostActionIntegrationResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response.ToJson())
}
