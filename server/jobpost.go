package main

import (
	"encoding/json"
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

func (p *Plugin) addJobpost(jobpost Jobpost) {
	p.API.LogInfo(jobpost.CreatedBy)
	jobpostJSON, err1 := json.Marshal(jobpost)
	if err1 != nil {
		p.API.LogError("failed to marshal Jobpost %s", jobpost.ID)
		return
	}
	p.API.LogInfo(string(jobpostJSON))
	err5 := p.API.KVSet(jobpost.ID, jobpostJSON)
	if err5 != nil {
		p.API.LogError("failed KVSet %s", err5)
		return
	}

	bytes, err2 := p.API.KVGet(jobpost.CreatedBy)
	p.API.LogInfo(string(bytes))
	if err2 != nil {
		p.API.LogError("failed KVGet %s", err2)
		return
	}
	var jobposts []string
	if bytes != nil {
		if err3 := json.Unmarshal(bytes, &jobposts); err3 != nil {
			return
		}
		jobposts = append(jobposts, jobpost.ID)
	} else {
		jobposts = []string{jobpost.ID}
	}

	p.API.LogInfo(jobposts[0])
	jobpostsJSON, err4 := json.Marshal(jobposts)
	if err4 != nil {
		p.API.LogError("failed to marshal Jobposts  %s", jobposts)
		return
	}
	p.API.KVSet(jobpost.CreatedBy, jobpostsJSON)
	return
}

func (p *Plugin) addJobpostResponse(postID string, jobpostResponse JobpostResponse) {
	bytes, err1 := p.API.KVGet(postID)
	if err1 != nil {
		p.API.LogError("failed KVGet %s", err1)
		return
	}
	var jobpost Jobpost
	if err2 := json.Unmarshal(bytes, &jobpost); err2 != nil {
		p.API.LogError("failed to unmarshal", err2)
		return
	}
	if jobpost.ExperienceReq && jobpost.MinExperience > jobpostResponse.Experience && jobpost.MaxExperience < jobpostResponse.Experience {
		p.API.LogError("Experience is not matching. Please apply to other jobs.")
		return
	}
	jobpost.JobpostResponses = append(jobpost.JobpostResponses, jobpostResponse)
	jobpostJSON, err3 := json.Marshal(jobpost)
	if err3 != nil {
		p.API.LogError("failed to marshal Jobpost %s", jobpost.ID)
		return
	}
	p.API.LogInfo(string(jobpostJSON))
	err5 := p.API.KVSet(jobpost.ID, jobpostJSON)
	if err5 != nil {
		p.API.LogError("failed KVSet %s", err5)
		return
	}
	return
}
