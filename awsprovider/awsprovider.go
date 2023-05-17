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

// Modify the CreateResource function to accept InstanceType and ImageID
func (a *AWSProvider) CreateResource(name string, instanceType string, imageID string) error {
    sgInput := &ec2.DescribeSecurityGroupsInput{
        GroupNames: []*string{
            aws.String("your-security-group-name"),
        },
    }
    sgResult, err := a.svc.DescribeSecurityGroups(sgInput)
    if err != nil {
        return err
    }
    sgId := sgResult.SecurityGroups[0].GroupId

    _, err = a.svc.RunInstances(&ec2.RunInstancesInput{
        ImageId:      aws.String(imageID),
        InstanceType: aws.String(instanceType),
        MinCount:     aws.Int64(1),
        MaxCount:     aws.Int64(1),
        KeyName:      aws.String("your-key-pair-name"),
        SecurityGroupIds: []*string{
            sgId, // Use the security group ID retrieved above
        },
        SubnetId: aws.String("your-subnet-id"),
        IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
            Name: aws.String("your-iam-instance-profile-name"),
        },
    })
	if err != nil {
		return err
	}

	return nil
}

func (a *AWSProvider) DeleteResource(name string) error {
	_, err := a.svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(name)},
	})
	if err != nil {
		return err
	}

	return nil
}
