package internal

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type (
	Client struct {
		ec2 *ec2.EC2
	}

	Instances []Instance

	Instance struct {
		ID    string
		Name  string
		State string
	}
)

func NewClient() *Client {
	// TODO: プロファイル指定できるようにする
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	))
	return &Client{ec2: ec2.New(sess)}
}

func (c Client) FetchAllInstances(ctx context.Context) (Instances, error) {
	result, err := c.ec2.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	var is Instances
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			n := ""
			ts := i.Tags
			for _, t := range ts {
				if *t.Key == "Name" {
					n = *t.Value
					break
				}
			}

			is = append(is, Instance{
				ID:    *i.InstanceId,
				Name:  n,
				State: *i.State.Name,
			})
		}
	}
	return is, nil
}
