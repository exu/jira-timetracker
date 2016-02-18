# JIRA timetracker for BrownOrcas team!

## Installation

get one of precompiled binaries and put somewhere in your `$PATH`

- linux https://github.com/exu/jira-timetracker/raw/master/jira
- mac https://github.com/exu/jira-timetracker/raw/master/jira_mac

or if you don't trust me check code in jira.go and build it yourself

    go build jira.go -o jira
    cp jira /bin/

## Usage
run jira to get available options

    jira


to run jira tt tool run:

    jira -u yourJiraUsername -p yourJiraPassword -id ELTCD-1111

or cd to your project directory and omit id parameter it'll get it from git feature branch name

    jira -u yourJiraUsername -p yourJiraPassword

default duration of your work is set to 7h (hmmm why? :) ) but you can override it

    jira -u yourJiraUsername -p yourJiraPassword -d 2h30m


If you doesn't like agile time tracking put call to this program into your crontab :P

Voila!
