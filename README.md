# Event Driver

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lukecold/event-driver)
[![Go Project Version](https://badge.fury.io/go/github.com%2Flukecold%2Fevent-driver.svg)](https://badge.fury.io/go/github.com%2Flukecold%2Fevent-driver)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=lukecold_event-driver&metric=alert_status)](https://sonarcloud.io/dashboard?id=lukecold_event-driver)
![Sonar Coverage](https://img.shields.io/sonar/coverage/lukecold_event-driver?server=https%3A%2F%2Fsonarcloud.io)


Event Driver is a lightweight and flexible event-driven programming framework for managing and handling events in your applications. It provides a simple and intuitive API to facilitate communication between different components or modules in your software.

## Features

- **Event-driven Architecture**: Easily implement an event-driven architecture in your application.
- **Custom Handlers**: Define and dispatch custom handlers tailored to your application's needs.
- **Pipeline Structure**: Simply put your handlers in order and expect it to work like a pipeline.
- **Asynchronous Support**: Handle events asynchronously for improved performance and responsiveness.
- **Lightweight and Easy to Use**: Minimalistic design for quick integration and usage.

## Usage

The following showcases an example of building a message processing pipeline, and convert it into a KNative event handler.
```golang
package main

import (
	"context"
	"log"

	"github.com/lukecold/event-driver/convert"
	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers/cache"
	"github.com/lukecold/event-driver/handlers/joiner"
	"github.com/lukecold/event-driver/handlers/transformer"
	"github.com/lukecold/event-driver/pipeline"
	"github.com/lukecold/event-driver/storage"
)

func main() {
	ctx := context.Background()
	renameSources, err := transformer.RenameSources(map[string][]string{"source1": {"alias1", "alias2"}})
	if err != nil {
		log.Panic("failed to create 'RenameSources' transformer", err)
	}
	myPipeline := pipeline.New().
		WithNextHandler(renameSources).
		WithNextHandler(joiner.New(joiner.MatchAll("source1", "source2"), storage.NewInMemoryStore())).
		WithNextHandler(cache.New(storage.NewInMemoryStore(), cache.SkipOnConflict()))

	// if convert to cloud event handler
	handleKNativeEvent := convert.ToKNativeEventHandler(
		convert.CloudEventToInput,
		myPipeline,
		convert.OutputToCloudResult)
	cloudEventClient.StartReceiver(ctx, handleKNativeEvent)
}
```
