# jerkins
An Alexa skill to manage office365 calendars. This project is currently in the **conception** phase right now, so there's really nothing to see here yet.

### feature brainstorm 

#### MVP: What is my next meeting?
- wire up local development environment for Alexa skill
- figure out how to query next meeting with api
- make Alexa answer next meeting with hardcoded access token
- link Alexa account to office365 account with oauth, use that to make API calls to get meetings
- hosting: host app build pipeline on ecs how keith is doing it

#### Other Skill Features
- _When is my next meeting?_
- _When is my next meeting with (x)_
- _Set up (meeting type) meeting with (x) (asap/this week/next week), book a room/ x's office._
- _Reserve (some time, an hour, etc) (today, this week, next week) for (x), invite X, book room Y_
- link multiple accounts, allow different office occupants to identify themselves and make appointments as themselves.

### local development

set up your $GOPATH and clone the repo:
```
export GOPATH="$HOME/dev/funtastic/branches/gopath"
git clone git@github.com:froderick/jerkins.git $GOPATH/src/github.com/froderick/jerkins
```

run the skill server locally:
```
cd $GOPATH/src/github.com/froderick/jerkins
go run main.go
```

set up [localtunnel](https://localtunnel.github.io/www/) so alexa can reach it:
```
npm install -g localtunnel
lt --port 8000
```

`lt` will print out the public https url your skill server is exposed on:
```
your url is: https://${something}.localtunnel.me
```
