package main

import (
	"encoding/json"
	"fmt"
	"time"
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
	JobpostResponses []JobpostResponse
}

type JobpostResponse struct {
	UserID     string
	Name       string
	Email      string
	Resume     string
	Reason     string
	Experience int
}

type JobPerUser struct {
	JobpostID string
	Details   string
	CreatedAt time.Time
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

func (p *Plugin) addJobpostResponse(postID string, jobpostResponse JobpostResponse) interface{} {
	bytes, err1 := p.API.KVGet(postID)
	if err1 != nil {
		p.API.LogError("failed KVGet %s", err1)
		return fmt.Sprintf("failed KVGet %s", err1)
	}
	var jobpost Jobpost
	if err2 := json.Unmarshal(bytes, &jobpost); err2 != nil {
		p.API.LogError("failed to unmarshal", err2)
		return fmt.Sprintf("failed to unmarshal  %s", err2)
	}
	if jobpost.ExperienceReq && (jobpost.MinExperience > jobpostResponse.Experience || jobpost.MaxExperience < jobpostResponse.Experience) {
		p.API.LogError("Experience is not matching. Please apply to other jobs.")
		return fmt.Sprintf("Experience is not matching. Please apply to other jobs.")
	}
	jobpost.JobpostResponses = append(jobpost.JobpostResponses, jobpostResponse)
	jobpostJSON, err3 := json.Marshal(jobpost)
	if err3 != nil {
		p.API.LogError("failed to marshal Jobpost %s", jobpost.ID)
		return fmt.Sprintf("failed to marshal Jobpost %s", jobpost.ID)
	}
	p.API.LogInfo(string(jobpostJSON))
	err5 := p.API.KVSet(jobpost.ID, jobpostJSON)
	if err5 != nil {
		p.API.LogError("failed KVSet %s", err5)
		return fmt.Sprintf("failed KVSet %s", err5)
	}
	return nil
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
