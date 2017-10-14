# jerkins
An Alexa skill to manage office365 calendars. This project is currently in the **conception** phase right now, so there's really nothing to see here yet.

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
