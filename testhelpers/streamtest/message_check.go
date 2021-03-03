package streamtest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageCheck struct {
	Validators []MessagePredicate
}

func (r MessageCheck) Check(msg *message.Message) []error {
	var results []error

	for _, predicate := range r.Validators {
		if !predicate.Matches(context.Background(), msg.Metadata) {
			results = append(results, MessageCheckError{
				Validator: predicate,
			})
		}
	}

	return results
}


type MessageCheckError struct {
	Validator MessagePredicate
}

func (c MessageCheckError) Error() string {
	return fmt.Sprintf("Failed Message validator: %s", c.Validator.Description)
}

type MessagePredicate struct {
	Description string
	Matches stream.MessageFilter
}

func MessageHasMetadata(key, value string) MessagePredicate {
	return MessagePredicate{
		Description: fmt.Sprintf("request.Metadata[%q] == %q", key, value),
		Matches: stream.FilterByMetaData(key, value),
	}
}
