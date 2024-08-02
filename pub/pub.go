/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

type Olympians struct {
	Olympians []Olympian `json:"olympians"`
}

type Olympian struct {
	ID     int         `json:"id"`
	Name   string      `json:"name"`
	Sex    string      `json:"sex"`
	Age    interface{} `json:"age"`
	NOC    string      `json:"noc"`
	Year   int         `json:"year"`
	Season string      `json:"season"`
	Medal  string      `json:"medal"`
}

var (
	// set the environment as instructions.
	pubsubName = "olympians"
	topicName  = "athletes"
)

func main() {
	ctx := context.Background()
	publishEventMetadata := map[string]string{
		"cloudevent.type": "important",
	}

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Open and read olympians data
	jsonFile, err := os.Open("../olympics.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened olympics.json")
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("failed to read json file, error: %v", err)
		return
	}
	var olympians Olympians
	if err := json.Unmarshal(byteValue, &olympians); err != nil {
		panic(err)
	}

	// Publish
	for i := 0; i < len(olympians.Olympians); i++ {
		if olympians.Olympians[i].NOC == "CAN" {
			fmt.Println("Olympian is from Canada: " + olympians.Olympians[i].Name)

			// If the Olympian is from Canada, publish as an important event
			if err := client.PublishEvent(ctx, pubsubName, topicName, olympians.Olympians[i], dapr.PublishEventWithMetadata(publishEventMetadata)); err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Olympian is NOT from Canada: " + olympians.Olympians[i].Name)
			if err := client.PublishEvent(ctx, pubsubName, topicName, olympians.Olympians[i]); err != nil {
				panic(err)
			}
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Println("Done (CTRL+C to Exit)")
}
