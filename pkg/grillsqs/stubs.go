package grillsqs

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/Swiggy/grill"
)

func (gs *SQS) CreateQueue(createQueueInput *sqs.CreateQueueInput) grill.Stub {
	return grill.StubFunc(func() error {
		if _, err := gs.Client().CreateQueue(createQueueInput); err != nil {
			return err
		}
		return nil
	})
}
