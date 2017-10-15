/*
	The main entrypoint for the alexa skill server.
 */
package main

import (
	alexa "github.com/mikeflynn/go-alexa/skillserver"
	"github.com/froderick/jerkins"
)

var Applications = map[string]interface{}{
	"/echo/jerkins": alexa.EchoApplication{ // Route
		AppID:    "amzn1.ask.skill.562b57a4-15ee-4f12-9614-0ec8da789955", // Echo App ID from Amazon Dashboard
		OnIntent: jerkins.TestIntentHandler,
		OnLaunch: jerkins.TestIntentHandler,
	},
}

func main() {
	alexa.Run(Applications, "8081")
}
