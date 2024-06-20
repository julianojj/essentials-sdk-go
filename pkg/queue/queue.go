package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type (
	Queue interface {
		Consume(queueName string, callback func(args []byte) error) error
		Publish(queueURL string, message string) error
	}
	SQS struct {
		client *sqs.Client
	}
)

func NewSQS(
	endpoint string,
	region string,

) Queue {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic(err)
	}
	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
	return &SQS{
		client,
	}
}

func (s *SQS) Consume(queueName string, callback func(args []byte) error) error {
	for {
		ouptut, err := s.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(queueName),
		})
		if err != nil {
			return err
		}
		for _, msg := range ouptut.Messages {
			body := []byte(*msg.Body)
			err := callback(body)
			if err == nil {
				s.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueName),
					ReceiptHandle: msg.ReceiptHandle,
				})
			}
		}
	}
}

func (s *SQS) Publish(queueURL, message string) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(message),
	}
	_, err := s.client.SendMessage(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}
