package joiner

import (
	"context"
	"encoding/json"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers"
	"github.com/lukecold/event-driver/storage"
)

// joiner implements handlers.Handler that joins the events with the same key
// when the sources match the criteria given by Condition.
// The output is of JSON format `{"source1":"content1","source2":"content2",...}`.
type joiner struct {
	condition Condition
	storage   storage.EventStore
}

func New(condition Condition, storage storage.EventStore) *joiner {
	return &joiner{
		condition: condition,
		storage:   storage,
	}
}

func (e *joiner) Process(ctx context.Context, in event.Message, next handlers.CallNext) error {
	// persist input message by key & source
	err := e.storage.Persist(ctx, in.GetKey(), in.GetSource(), in.GetContent())
	if err != nil {
		return err
	}

	// validate sources
	messages, err := e.storage.LookUpByKey(ctx, in.GetKey())
	if err != nil {
		return err
	}
	persistedSources := make([]string, 0, len(messages))
	for _, message := range messages {
		persistedSources = append(persistedSources, message.GetSource())
	}
	if !e.condition.Evaluate(persistedSources) {
		return nil
	}

	// join sources
	contentBySource := make(map[string]string)
	for _, message := range messages {
		contentBySource[message.GetSource()] = message.GetContent()
	}
	jointContent, err := json.Marshal(contentBySource)
	if err != nil {
		return err
	}

	jointEvent := event.NewMessage(in.GetKey(), "composed-event", string(jointContent))

	return next.Call(ctx, jointEvent)
}
