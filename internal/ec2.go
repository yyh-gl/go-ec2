package internal

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type (
	Client struct {
		ec2 *ec2.EC2
	}

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

func (c Client) ShowAllInstances(ctx context.Context) error {
	is, err := c.fetchAllInstances(ctx)
	if err != nil {
		return err
	}

	for _, i := range is {
		fmt.Println("========================")
		fmt.Println(i)
		fmt.Println("========================")
	}

	return nil
}

func (c Client) fetchAllInstances(ctx context.Context) ([]Instance, error) {
	result, err := c.ec2.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	is := make([]Instance, 0)
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
