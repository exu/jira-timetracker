// Some fun with JIRA API :)
// curl -D -u user:pass -X POST
// --data '{ "comment": "Hard work.", "started":"2016-02-18T01:20:19.843+0000", "timeSpentSeconds": 12000 }'
// -H "Content-Type: application/json" http://jira.pearson.com/rest/api/2/issue/ELTCD-9916/worklog
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

var username = flag.String("u", "", "jira user name")
var password = flag.String("p", "", "jira password")
var id = flag.String("id", "", "jira task/story id")
var message = flag.String("m", "", "time log message")
var duration = flag.String("d", "7h", "time spent on task in duration format e.g. 1h10m")

// Custom time because jira need one true format (server: 500 if not fit)
const JIRA_TIME_FORMAT = "2006-01-02T15:04:05.000Z0700"
const JIRA_URL = "http://jira.pearson.com/rest/api/2/issue/"
const JSON_CONFIG_FILE = ".auth.json"

type Config struct {
	Jira struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	} `json:"jira"`
}

// Custom serializable to json time object
type Time struct {
	time.Time
}

func (time Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Format(JIRA_TIME_FORMAT) + `"`), nil
}

// Payload object
type Payload struct {
	Comment          string `json:"comment"`
	Started          Time   `json:"started"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
}

func req(username, password, id, durationString, message string) error {

	client := &http.Client{}
	url := JIRA_URL + id + "/worklog"

	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return err
	}

	seconds := int(duration.Seconds())

	data, err := json.Marshal(Payload{
		message,
		Time{time.Now()},
		seconds,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)

	return nil
}

func loadFromJSON() (string, string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configFile := usr.HomeDir + "/" + JSON_CONFIG_FILE

	file, _ := os.Open(configFile)
	decoder := json.NewDecoder(file)
	config := &Config{}

	err = decoder.Decode(&config)
	if err != nil {
		log.Println("Decode ./config.json error:", err)
	}

	return config.Jira.User, config.Jira.Pass
}

func init() {
	flag.Parse()
}

func main() {

	if *id == "" {
		out, err := exec.Command("sh", "-c", "git rev-parse --symbolic-full-name --abbrev-ref HEAD").Output()
		if err != nil {
			flag.Usage()
			log.Fatalln("Can't find issue id, and you're probably not in git repo to guess it by feature branch")
		}
		*id = strings.TrimSpace(string(out))
		log.Printf("Time tracking for current branch '%s'", *id)
	}

	if *username == "" && *password == "" {
		*username, *password = loadFromJSON()
	}

	if *username == "" || *password == "" || *id == "" || *duration == "" {
		flag.Usage()
		return
	}

	err := req(*username, *password, *id, *duration, *message)

	if err != nil {
		log.Fatalln(err)
	}
}
