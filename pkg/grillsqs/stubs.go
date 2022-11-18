package grillsqs

import (
	"github.com/Swiggy/grill"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (gs *SQS) CreateQueue(createQueueInput *sqs.CreateQueueInput) grill.Stub {
	return grill.StubFunc(func() error {
		if _, err := gs.Client().CreateQueue(createQueueInput); err != nil {
			return err
		}
		return nil
	})
}
