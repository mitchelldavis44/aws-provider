// awsprovider/aws_provider.go
package awsprovider

import (
    "fmt"
    "time"

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

// Modify the CreateResource function to accept InstanceType, ImageID and Tags
func (a *AWSProvider) CreateResource(name string, instanceType string, imageID string, securityGroupId string, keyPairName string, subnetId string, iamInstanceProfile string, vpcId string, tags map[string]string) (string, error) {
    // Convert the tags to the correct format
    awsTags := []*ec2.Tag{}
    for key, value := range tags {
        awsTags = append(awsTags, &ec2.Tag{
            Key:   aws.String(key),
            Value: aws.String(value),
        })
    }

    // Add tags to the RunInstancesInput
    runResult, err := a.svc.RunInstances(&ec2.RunInstancesInput{
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
        TagSpecifications: []*ec2.TagSpecification{
            {
                ResourceType: aws.String("instance"),
                Tags:         awsTags,
            },
        },
    })
    if err != nil {
        return "", err
    }

    instanceID := runResult.Instances[0].InstanceId
    fmt.Printf("Instance is launching, ID: %s\n", *instanceID)

    // Now we keep checking the instance status until it's running or maximum wait time has been reached
    maxRetries := 40 // you can change this number based on your requirement
    for i := 0; i < maxRetries; i++ {
        // Wait for 15 seconds before checking status again
        time.Sleep(time.Duration(15) * time.Second)

        input := &ec2.DescribeInstancesInput{
            InstanceIds: []*string{instanceID},
        }
        result, err := a.svc.DescribeInstances(input)
        if err != nil {
            return "", err
        }

        if len(result.Reservations) > 0 && len(result.Reservations[0].Instances) > 0 {
            instanceState := result.Reservations[0].Instances[0].State.Name
            if *instanceState == "running" {
                fmt.Printf("Instance is up and running.\n")
                return *instanceID, nil
            } else {
                fmt.Printf("Current state of instance is: %s\n", *instanceState)
            }
        }
    }

    return "", fmt.Errorf("Failed to create instance: Exceeded max wait time of %d seconds", maxRetries*15)
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