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
	"log"
	"net/http"
	"strconv"

	"github.com/dapr/go-sdk/client"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

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

type App struct {
	client dapr.Client
}

// Subscription to tell the dapr what topic to subscribe.
//   - PubsubName: is the name of the component configured in the metadata of pubsub.yaml.
//   - Topic: is the name of the topic to subscribe.
//   - Route: tell dapr where to request the API to publish the message to the subscriber when get a message from topic.
//   - Match: (Optional) The CEL expression to match on the CloudEvent to select this route.
//   - Priority: (Optional) The priority order of the route when Match is specified.
//     If not specified, the matches are evaluated in the order in which they are added.
var defaultSubscription = &common.Subscription{
	PubsubName: "olympians",
	Topic:      "athletes",
	Route:      "/athletes",
}

var importantSubscription = &common.Subscription{
	PubsubName: "olympians",
	Topic:      "athletes",
	Route:      "/important",
	Match:      `event.type == "important"`,
	Priority:   1,
}

var (
	// set the environment as instructions.
	kvStoreName = "kvstore"
)

var app App

func main() {
	//Create new Dapr service
	server := daprd.NewService(":8080")

	//Create new Dapr client
	client, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	app = App{
		client: client,
	}

	if err := server.AddTopicEventHandler(defaultSubscription, eventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	if err := server.AddTopicEventHandler(importantSubscription, importantEventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listening: %v", err)
	}
}

// Receive regular events
func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	log.Printf("Received event - PubsubName: %s, Topic: %s", e.PubsubName, e.Topic)
	return false, nil
}

// Receive and persist important events
func importantEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	log.Printf("Received event (important) - PubsubName: %s, Topic: %s", e.PubsubName, e.Topic)

	var olympian Olympian
	if err := json.Unmarshal(e.RawData, &olympian); err != nil {
		log.Fatalf("Error unmarshaling %v", err)
	}

	fmt.Println("Persisting athlete to state store", olympian)
	if err := app.client.SaveState(ctx, kvStoreName, strconv.Itoa(olympian.ID), e.RawData, nil); err != nil {
		panic(err)
	}

	fmt.Println("Persisted Canadian athlete to state store")
	return false, nil
}
