package grillsqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/lovlin-thakkar/swiggy-grill"
)

func (gs *SQS) DeleteQueues(queues ...string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		for _, queue := range queues {
			queueUrl, err := gs.GetQueueUrl(queue)
			if err != nil {
				return err
			}
			if _, err := gs.sqs.Client.DeleteQueue(&sqs.DeleteQueueInput{QueueUrl: aws.String(queueUrl)}); err != nil {
				return err
			}
		}
		return nil
	})
}
