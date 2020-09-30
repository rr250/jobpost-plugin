package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (p *Plugin) UserSubscribe() {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		p.API.LogError("failed to get bundle path: %v", err)
	}
	usersCsv, err1 := ioutil.ReadFile(filepath.Join(bundlePath, "assets/usersToSubscribe.csv"))
	if err1 != nil {
		p.API.LogError("Unable to read csv: %v", err1)
		return
	}
	buf := bytes.NewBuffer(usersCsv)
	reader := csv.NewReader(buf)
	users, err2 := reader.ReadAll()
	if err2 != nil {
		p.API.LogError("Unable to read csv from reader: %v", err2)
		return
	}
	for _, user := range users {
		currentYear := time.Now().Year()
		userYear, err3 := strconv.Atoi(user[1])
		if err3 != nil {
			continue
		}
		userExperience := currentYear - userYear
		if userExperience < 0 {
			userExperience = 0
		}
		userDetails, err4 := p.API.GetUserByEmail(user[0])
		if err4 != nil {
			continue
		}
		p.subscribeToExperience(userDetails.Id, userExperience)
	}
	os.Remove(filepath.Join(bundlePath, "assets/usersToSubscribe.csv"))
	p.API.LogInfo("Users Subscribed")
}
