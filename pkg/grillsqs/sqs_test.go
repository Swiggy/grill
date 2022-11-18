package grillsqs

import (
	"context"
	"fmt"
	"github.com/Swiggy/grill"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"testing"
	"time"
)

func Test_GrillSQS(t *testing.T) {
	helper := &SQS{}
	if err := helper.Start(context.TODO()); err != nil {
		t.Errorf("error starting sqs grill, error=%v", err)
		return
	}

	tests := []grill.TestCase{
		{
			Name: "Test_GrillSQS_ProduceCount",
			Stubs: []grill.Stub{
				helper.CreateQueue(&sqs.CreateQueueInput{QueueName: aws.String("test-queue")}),
			},
			Action: func() interface{} {
				queueUrl, err := helper.GetQueueUrl("test-queue")
				if err != nil {
					return err
				}
				_, err = helper.Client().SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("test message body"),
					QueueUrl:    aws.String(queueUrl),
				})
				return err
			},
			Assertions: []grill.Assertion{
				grill.AssertOutput(nil),
				grill.Try(time.Second*30, 3, helper.AssertCount("test-queue", 1)),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteQueues("test-queue"),
			},
		},
		{
			Name: "Test_GrillSQS_ProduceCount_1000Messages",
			Stubs: []grill.Stub{
				helper.CreateQueue(&sqs.CreateQueueInput{QueueName: aws.String("test-queue")}),
			},
			Action: func() interface{} {
				queueUrl, err := helper.GetQueueUrl("test-queue")
				if err != nil {
					return err
				}

				for i := 0; i < 10; i++ {
					for i := 0; i < 100; i++ {
						var entries []*sqs.SendMessageBatchRequestEntry
						entries = append(entries, &sqs.SendMessageBatchRequestEntry{
							Id:          aws.String(fmt.Sprintf("%d", i)),
							MessageBody: aws.String(fmt.Sprintf("message:%d", i)),
						})
						if _, err = helper.Client().SendMessageBatch(&sqs.SendMessageBatchInput{
							Entries:  entries,
							QueueUrl: aws.String(queueUrl),
						}); err != nil {
							return err
						}
					}
				}
				return nil
			},
			Assertions: []grill.Assertion{
				grill.AssertOutput(nil),
				grill.Try(time.Second*10, 3, helper.AssertCount("test-queue", 1000)),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteQueues("test-queue"),
			},
		},
		{
			Name: "Test_GrillSQS_AssertMessageCount",
			Stubs: []grill.Stub{
				helper.CreateQueue(&sqs.CreateQueueInput{QueueName: aws.String("test-queue")}),
			},
			Action: func() interface{} {
				queueUrl, err := helper.GetQueueUrl("test-queue")
				if err != nil {
					return err
				}
				var entries []*sqs.SendMessageBatchRequestEntry
				for i := 0; i < 3; i++ {
					entries = append(entries, &sqs.SendMessageBatchRequestEntry{
						Id:          aws.String(fmt.Sprintf("%d", i)),
						MessageBody: aws.String(fmt.Sprintf("message:%d", i)),
					})
				}
				if _, err = helper.Client().SendMessageBatch(&sqs.SendMessageBatchInput{
					Entries:  entries,
					QueueUrl: aws.String(queueUrl),
				}); err != nil {
					return err
				}

				return nil
			},
			Assertions: []grill.Assertion{
				grill.AssertOutput(nil),
				grill.Try(time.Second*10, 3, grill.AssertionFunc(func() error {
					if err := helper.AssertMessageCount("test-queue", "message:0", 1).Assert(); err != nil {
						return err
					}
					return helper.AssertMessageCount("test-queue", "message:1", 1).Assert()
				})),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteQueues("test-queue"),
			},
		},
	}

	grill.Run(t, tests)
}
