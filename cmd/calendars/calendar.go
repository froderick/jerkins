package main

import (
	"log"
	"io/ioutil"
	"fmt"
	"net/http"
	"os"
)

func main() {

	tok := os.Getenv("TOKEN")

    client := &http.Client{}

	auth := "Bearer " + tok
	fmt.Printf("%s\n", auth)

	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me/events", nil)
	req.Header.Add("Authorization", "Bearer " + tok)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Printf("%s", body)
}
