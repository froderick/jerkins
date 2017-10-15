package main

import (
	"fmt"
	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

func TestIntentHandler(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
	fmt.Println("hi there")
	echoResp.OutputSpeech("Hello world from my new Echo test app!").Card("Hello World", "This is a test card.")
}
