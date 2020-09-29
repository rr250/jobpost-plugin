package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

func (p *Plugin) getClient(config *oauth2.Config, bundlePath string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "assets/credentials/token.json"
	tok, err := tokenFromFile(tokFile, bundlePath)
	if err != nil {
		p.API.LogError("Unable to read token: %v", err)
		// tok = getTokenFromWeb(config)
		p.saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
// func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
// 	fmt.Printf("Go to the following link in your browser then type the "+
// 		"authorization code: \n%v\n", authURL)

// 	var authCode string
// 	if _, err := fmt.Scan(&authCode); err != nil {
// 		p.API.LogError("Unable to read authorization code: %v", err)
// 	}

// 	tok, err := config.Exchange(context.TODO(), authCode)
// 	if err != nil {
// 		p.API.LogError("Unable to retrieve token from web: %v", err)
// 	}
// 	return tok
// }

// Retrieves a token from a local file.
func tokenFromFile(file string, bundlePath string) (*oauth2.Token, error) {
	f, err := os.Open(filepath.Join(bundlePath, file))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func (p *Plugin) saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		p.API.LogError("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (p *Plugin) InitGoogleServices() (*drive.Service, *sheets.Service) {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		p.API.LogError("failed to get bundle path: %v", err)
	}
	b, err1 := ioutil.ReadFile(filepath.Join(bundlePath, "assets/credentials/credentials.json"))
	if err1 != nil {
		p.API.LogError("Unable to read client secret file: %v", err1)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err2 := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets "+drive.DriveScope)
	if err2 != nil {
		p.API.LogError("Unable to parse client secret file to config: %v", err2)
	}
	client := p.getClient(config, bundlePath)

	driveService, err3 := drive.New(client)
	if err3 != nil {
		p.API.LogError("Unable to retrieve Drive client: %v", err3)
	}

	sheetsService, err4 := sheets.New(client)
	if err4 != nil {
		p.API.LogError("Unable to retrieve Sheets client: %v", err4)
	}
	p.API.LogInfo("Google APIs activated")
	return driveService, sheetsService
}
