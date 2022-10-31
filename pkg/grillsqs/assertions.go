package grillsqs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/Swiggy/grill"
	"time"
)

func (gs *SQS) AssertCount(queueName string, expectedCount int) grill.Assertion {
	return grill.AssertionFunc(func() error {

		// time for which are messages are considered processed while looping around all messages in the queue.
		// if the visibility timeout is set to 0, we will receive the same messages on each request.
		// with visibility timeout of 1 second, we are able to loop through 1000 messages.
		//If your usecase publishes more messages than that, please contact the dev for increasing the timeout.
		visibilityTimeoutSecond := int64(1)

		queueUrl, err := gs.GetQueueUrl(queueName)
		if err != nil {
			return err
		}

		count := 0

		for {
			out, err := gs.Client().ReceiveMessage(&sqs.ReceiveMessageInput{
				AttributeNames:          nil,
				MaxNumberOfMessages:     aws.Int64(10),
				QueueUrl:                aws.String(queueUrl),
				ReceiveRequestAttemptId: nil,
				VisibilityTimeout:       aws.Int64(visibilityTimeoutSecond),
			})
			if err != nil {
				return err
			}
			if len(out.Messages) == 0 {
				break
			}

			count += len(out.Messages)
		}

		// ensuring messages are placed back in queue for other assertions.
		time.Sleep(time.Second * time.Duration(visibilityTimeoutSecond))

		if count != expectedCount {
			return fmt.Errorf("invalid number of messages, got=%v, want=%v", count, expectedCount)
		}

		return nil
	})
}

func (gs *SQS) AssertMessageCount(queueName string, expectedMessage string, expectedCount int) grill.Assertion {
	return grill.AssertionFunc(func() error {
		// time for which are messages are considered processed while looping around all messages in the queue.
		// if the visibility timeout is set to 0, we will receive the same messages on each request.
		// with visibility timeout of 1 second, we are able to loop through 1000 messages.
		//If your usecase publishes more messages than that, please contact the dev for increasing the timeout.
		visibilityTimeoutSecond := int64(1)

		queueUrl, err := gs.GetQueueUrl(queueName)
		if err != nil {
			return err
		}

		count := 0

		for {
			out, err := gs.Client().ReceiveMessage(&sqs.ReceiveMessageInput{
				AttributeNames:          nil,
				MaxNumberOfMessages:     aws.Int64(10),
				QueueUrl:                aws.String(queueUrl),
				ReceiveRequestAttemptId: nil,
				VisibilityTimeout:       aws.Int64(visibilityTimeoutSecond),
			})
			if err != nil {
				return err
			}
			if len(out.Messages) == 0 {
				break
			}

			for _, m := range out.Messages {
				if *m.Body == expectedMessage {
					count++
				}
			}
		}

		// ensuring messages are placed back in queue for other assertions.
		time.Sleep(time.Second * time.Duration(visibilityTimeoutSecond))

		if count != expectedCount {
			return fmt.Errorf("invalid number of messages, got=%v, want=%v", count, expectedCount)
		}
		return nil
	})
}
