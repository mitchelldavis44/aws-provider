// awsprovider/aws_provider.go
package awsprovider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchelldavis44/Harmony/pkg/infrastructure"
)

type AWSProvider struct {
	svc *ec2.EC2
}

func NewAWSProvider() infrastructure.Infrastructure {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	svc := ec2.New(sess)

	return &AWSProvider{
		svc: svc,
	}
}

func (a *AWSProvider) CreateResource(name string) error {
	// Use a.svc to create an EC2 instance, S3 bucket, etc.
	// This is just a simplified example and won't actually work.
	_, err := a.svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String("ami-0d52ddcdf3a885741"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *AWSProvider) DeleteResource(name string) error {
	// Use a.svc to delete an EC2 instance, S3 bucket, etc.
	// This is just a simplified example and won't actually work.
	_, err := a.svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(name)},
	})
	if err != nil {
		return err
	}

	return nil
}
