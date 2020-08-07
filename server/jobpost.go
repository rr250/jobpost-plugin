package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"google.golang.org/api/sheets/v4"
)

type Jobpost struct {
	ID               string
	CreatedBy        string
	CreatedAt        time.Time
	Company          string
	Position         string
	Description      string
	Skills           string
	MinExperience    int
	MaxExperience    int
	Location         string
	ExperienceReq    bool
	SheetID          string
	SheetURL         string
	JobpostResponses []JobpostResponse
}

type JobpostResponse struct {
	UserID     string
	Name       string
	Email      string
	Resume     string
	Reason     string
	Experience float64
	FilledAt   time.Time
}

type JobPerUser struct {
	JobpostID string
	Details   string
	CreatedAt time.Time
	SheetURL  string
}

type Subscriber struct {
	UserID string
}

type UserDetails struct {
	Resume string
}

func (p *Plugin) addJobpost(jobpost Jobpost) interface{} {
	p.API.LogInfo(jobpost.CreatedBy)
	jobpostJSON, err1 := json.Marshal(jobpost)
	if err1 != nil {
		p.API.LogError("failed to marshal Jobpost %s", jobpost.ID)
		return fmt.Sprintf("failed to marshal Jobpost %s", jobpost.ID)
	}
	p.API.LogInfo(string(jobpostJSON))
	err5 := p.API.KVSet(jobpost.ID, jobpostJSON)
	if err5 != nil {
		p.API.LogError("failed KVSet %s", err5)
		return fmt.Sprintf("failed KVSet %s", err5)
	}

	bytes, err2 := p.API.KVGet(jobpost.CreatedBy)
	p.API.LogInfo(string(bytes))
	if err2 != nil {
		p.API.LogError("failed KVGet %s", err2)
		return fmt.Sprintf("failed KVGet %s", err2)
	}
	jobPerUser := JobPerUser{
		JobpostID: jobpost.ID,
		Details:   jobpost.Company + " - " + jobpost.Position,
		CreatedAt: jobpost.CreatedAt,
		SheetURL:  jobpost.SheetURL,
	}
	var jobposts []JobPerUser
	if bytes != nil {
		if err3 := json.Unmarshal(bytes, &jobposts); err3 != nil {
			return fmt.Sprintf("failed to unmarshal  %s", err3)
		}
		jobposts = append(jobposts, jobPerUser)
	} else {
		jobposts = []JobPerUser{jobPerUser}
	}
	jobpostsJSON, err4 := json.Marshal(jobposts)
	if err4 != nil {
		p.API.LogError("failed to marshal Jobposts  %s", jobposts)
		return fmt.Sprintf("failed to marshal Jobposts  %s", jobposts)
	}
	p.API.KVSet(jobpost.CreatedBy, jobpostsJSON)
	return nil
}

func (p *Plugin) addJobpostResponse(postID string, jobpostResponse JobpostResponse) (Jobpost, interface{}) {
	var jobpost Jobpost
	bytes, err1 := p.API.KVGet(postID)
	if err1 != nil {
		p.API.LogError("failed KVGet %s", err1)
		return jobpost, fmt.Sprintf("failed KVGet %s", err1)
	}
	if err2 := json.Unmarshal(bytes, &jobpost); err2 != nil {
		p.API.LogError("failed to unmarshal", err2)
		return jobpost, fmt.Sprintf("failed to unmarshal  %s", err2)
	}
	if jobpost.ExperienceReq && ((jobpost.MinExperience-1) > int(math.Ceil(jobpostResponse.Experience)) || (jobpost.MaxExperience+1) < int(math.Ceil(jobpostResponse.Experience))) {
		p.API.LogError("Experience is not matching. Please apply to other jobs.")
		return jobpost, fmt.Sprintf("Experience is not matching. Please apply to other jobs.")
	}
	jobpost.JobpostResponses = append(jobpost.JobpostResponses, jobpostResponse)
	readRange := "Sheet1!A:Z"
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{
			{
				jobpostResponse.Name,
				jobpostResponse.Email,
				jobpostResponse.Resume,
				fmt.Sprintf("%.1f years", jobpostResponse.Experience),
				jobpostResponse.Reason,
				jobpostResponse.FilledAt,
			},
		},
	}
	_, err7 := p.sheetsService.Spreadsheets.Values.Append(jobpost.SheetID, readRange, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err7 != nil {
		return jobpost, fmt.Sprintf("Unable to append data from sheet: %v", err7)
	}
	jobpostJSON, err3 := json.Marshal(jobpost)
	if err3 != nil {
		p.API.LogError("failed to marshal Jobpost %s", jobpost.ID)
		return jobpost, fmt.Sprintf("failed to marshal Jobpost %s", jobpost.ID)
	}
	p.API.LogInfo(string(jobpostJSON))
	err5 := p.API.KVSet(jobpost.ID, jobpostJSON)
	if err5 != nil {
		p.API.LogError("failed KVSet %s", err5)
		return jobpost, fmt.Sprintf("failed KVSet %s", err5)
	}
	return jobpost, nil
}

func (p *Plugin) getJobsPerUser(userID string) ([]JobPerUser, interface{}) {
	var jobposts []JobPerUser
	bytes, err2 := p.API.KVGet(userID)
	p.API.LogInfo(string(bytes))
	if err2 != nil {
		p.API.LogError("failed KVGet %s", err2)
		return jobposts, fmt.Sprintf("failed to unmarshal %s", err2)
	}
	if bytes != nil {
		if err3 := json.Unmarshal(bytes, &jobposts); err3 != nil {
			return jobposts, fmt.Sprintf("failed to unmarshal %s", err3)
		}
	} else {
		return jobposts, "No Jobposts found"
	}
	return jobposts, nil
}

func (p *Plugin) getJobPost(jobpostID string) (Jobpost, interface{}) {
	var jobpost Jobpost
	bytes, err1 := p.API.KVGet(jobpostID)
	if err1 != nil {
		p.API.LogError("failed KVGet %s", err1)
		return jobpost, fmt.Sprintf("failed KVGet %s", err1)
	}
	if err2 := json.Unmarshal(bytes, &jobpost); err2 != nil {
		p.API.LogError("failed to unmarshal", err2)
		return jobpost, fmt.Sprintf("failed to unmarshal %s", err2)
	}
	return jobpost, nil
}

func (p *Plugin) subscribeToExperience(userID string, year int) interface{} {
	bytes, err2 := p.API.KVGet("year-" + strconv.Itoa(year))
	p.API.LogInfo(string(bytes))
	if err2 != nil {
		p.API.LogError("failed KVGet %s", err2)
		return fmt.Sprintf("failed KVGet %s", err2)
	}
	var subscribers []Subscriber
	newSubscriber := Subscriber{
		UserID: userID,
	}
	if bytes != nil {
		if err3 := json.Unmarshal(bytes, &subscribers); err3 != nil {
			return fmt.Sprintf("failed to unmarshal  %s", err3)
		}
		for _, subscriber := range subscribers {
			if subscriber.UserID == userID {
				return "You are already subscribed"
			}
		}

		subscribers = append(subscribers, newSubscriber)
	} else {
		subscribers = []Subscriber{newSubscriber}
	}
	subscribersJSON, err4 := json.Marshal(subscribers)
	if err4 != nil {
		p.API.LogError("failed to marshal Jobposts  %s", subscribers)
		return fmt.Sprintf("failed to marshal Jobposts  %s", subscribers)
	}
	p.API.KVSet("year-"+strconv.Itoa(year), subscribersJSON)
	return nil
}

func (p *Plugin) unSubscribeToExperience(userID string, year int) interface{} {
	bytes, err2 := p.API.KVGet("year-" + strconv.Itoa(year))
	p.API.LogInfo(string(bytes))
	if err2 != nil {
		p.API.LogError("failed KVGet %s", err2)
		return fmt.Sprintf("failed KVGet %s", err2)
	}
	var subscribers []Subscriber
	if bytes != nil {
		if err3 := json.Unmarshal(bytes, &subscribers); err3 != nil {
			return fmt.Sprintf("failed to unmarshal  %s", err3)
		}
		if len(subscribers) == 0 {
			return fmt.Sprintf("You are not a subscriber")
		}
		p.API.LogError("Length of subscribers  %d", len(subscribers))
		for i := range subscribers {
			if subscribers[i].UserID == userID {
				subscribers = append(subscribers[:i], subscribers[i+1:]...)
			} else {
				return fmt.Sprintf("You are not a subscriber")
			}
		}
		subscribersJSON, err4 := json.Marshal(subscribers)
		if err4 != nil {
			p.API.LogError("failed to marshal Jobposts  %s", subscribers)
			return fmt.Sprintf("failed to marshal Jobposts  %s", subscribers)
		}
		p.API.KVSet("year-"+strconv.Itoa(year), subscribersJSON)
	} else {
		return fmt.Sprintf("You are not a subscriber")
	}
	return nil
}

func (p *Plugin) sendToSubscribers(postModel *model.Post, year int) {
	bytes, err2 := p.API.KVGet("year-" + strconv.Itoa(year))
	p.API.LogInfo(string(bytes))
	if err2 != nil {
		p.API.LogError("failed KVGet %s", err2)
		return
	}
	var subscribers []Subscriber
	if bytes != nil {
		if err3 := json.Unmarshal(bytes, &subscribers); err3 != nil {
			return
		}
		for _, subscriber := range subscribers {
			channel, err1 := p.API.GetDirectChannel(subscriber.UserID, p.botUserID)
			if err1 == nil {
				postModel.ChannelId = channel.Id
				p.API.CreatePost(postModel)
			}
		}
	}
	return
}

func (p *Plugin) saveResume(userID string, resume string) interface{} {
	bytes, err2 := p.API.KVGet("resume-" + userID)
	p.API.LogInfo(string(bytes))
	if err2 != nil {
		p.API.LogError("failed KVGet %s", err2)
		return fmt.Sprintf("failed KVGet %s", err2)
	}
	var userDetails UserDetails
	if bytes != nil {
		if err3 := json.Unmarshal(bytes, &userDetails); err3 != nil {
			return fmt.Sprintf("failed to unmarshal  %s", err3)
		}
		userDetails.Resume = resume
	} else {
		userDetails.Resume = resume
	}
	userDetailsJSON, err4 := json.Marshal(userDetails)
	if err4 != nil {
		p.API.LogError("failed to marshal Jobposts  %s", userDetails)
		return fmt.Sprintf("failed to marshal Jobposts  %s", userDetails)
	}
	p.API.KVSet("resume-"+userID, userDetailsJSON)
	return nil
}

func (p *Plugin) getResume(userID string) (UserDetails, interface{}) {
	var userDetails UserDetails
	bytes, err2 := p.API.KVGet("resume-" + userID)
	p.API.LogInfo(string(bytes))
	if err2 != nil {
		p.API.LogError("failed KVGet %s", err2)
		return userDetails, fmt.Sprintf("failed KVGet %s", err2)
	}

	if bytes != nil {
		if err3 := json.Unmarshal(bytes, &userDetails); err3 != nil {
			return userDetails, fmt.Sprintf("failed to unmarshal  %s", err3)
		}
	} else {
		return userDetails, "No resume found"
	}
	return userDetails, nil
}
