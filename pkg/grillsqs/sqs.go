package grillsqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/lovlin-thakkar/swiggy-grill/canned"
)

type SQS struct {
	sqs *canned.SQS
}

func (gs *SQS) Start(ctx context.Context) error {
	sqs, err := canned.NewSQS(ctx)
	if err != nil {
		return err
	}
	gs.sqs = sqs
	return nil
}

func (gs *SQS) Client() sqsiface.SQSAPI {
	return gs.sqs.Client
}

func (gs *SQS) Host() string {
	return gs.sqs.Host
}

func (gs *SQS) Port() string {
	return gs.sqs.Port
}

func (gs *SQS) Region() string {
	return gs.sqs.Region
}

func (gs *SQS) AccessKey() string {
	return gs.sqs.AccessKey
}

func (gs *SQS) SecretKey() string {
	return gs.sqs.SecretKey
}

func (gs *SQS) Stop(ctx context.Context) error {
	return gs.sqs.Container.Terminate(ctx)
}

func (gs *SQS) GetQueueUrl(queueName string) (string, error) {
	queueUrl, err := gs.sqs.Client.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})
	if err != nil {
		return "", fmt.Errorf("error getting queue url for:%s, err=%w", queueName, err)
	}
	return *queueUrl.QueueUrl, nil

}
