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
func (a *AWSProvider) CreateResource(name string, instanceType string, imageID string, securityGroupId string, keyPairName string, subnetId string, iamInstanceProfile string, vpcId string) (string, error) {
    runResult, err := a.svc.RunInstances(&ec2.RunInstancesInput{
    _, err := a.svc.RunInstances(&ec2.RunInstancesInput{
        ImageId:      aws.String(imageID),
        InstanceType: aws.String(instanceType),
        MinCount:     aws.Int64(1),
        MaxCount:     aws.Int64(1),
        KeyName:      aws.String(keyPairName),
        IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
            Name: aws.String(iamInstanceProfile),
        },
        NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
            {
                DeviceIndex:              aws.Int64(0),
                SubnetId:                 aws.String(subnetId),
                Groups:                   []*string{aws.String(securityGroupId)},
                AssociatePublicIpAddress: aws.Bool(true),
            },
        },
    })
    if err != nil {
        return err
    }

    return *runResult.Instances[0].InstanceId, nil
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
