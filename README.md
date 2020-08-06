# Mattermost jobpost-plugin

_**A Plugin that makes task of creating Jobposts in Scaler chat based on Mattermost easier and cleaner**_

### Usage

* `/jobpost` - opens up an [interactive dialog] to create a Jobpost and dm you a google spreadsheet which will track all the response
* `/jobpost list` - displays a list of jobposts created by you
* `/jobpost subscribe x years` - subscribes to jobposts which requires x years of experience where x is an integer
* `/jobpost unsubscribe x years` - unsubscribes to jobposts which requires x years of experience where x is an integer
* `/jobpost resume save https://drive.google.com/file/d/xyz` - saves the resume link which can be prefetched the when you apply to a jobpost
* `/jobpost resume show` - fetch the resume you have saved

### Build
1) Clone the repository

2) Enable the Google APIs:- 
 + * Visit https://developers.google.com/drive/api/v3/quickstart/go and click Enable Drive Api button and get credentials.json
 + * Clone https://github.com/rr250/google-apis
 + * Open google-apis repo
 + * Paste your credentials.json here
 + * Run ```go get```
 + * Run ```go run quickstart.go```
 + * It will prompt you to authorize access through a link which appear in command-line promt
 + * Browse to the provided URL in your web browser
 + * If you are not already logged into your Google account, you will be prompted to log in
 + * Accept all
 + * Copy the code you're given, paste it into the command-line prompt, and press Enter
 + * A token.json file will be created in the folder
 + * Now copy both token.json and credentials.json
  
3) Exit this repository and open the jobpost-plugin repository

4) Paste token.json and credentials.json in assets/credentials/ folder

5) If you want to pre-subscribe users to jobposts according to years of experience, create a usersToSubscribe.csv with`email,year of graduation` format. and paste it in assets/ folder

6) Run

```
make
```

7) This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.github.rr250.jobpost-plugin-0.0.1.tar.gz
```
