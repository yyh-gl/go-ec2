package internal

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/yyh-gl/go-ec2/internal/sender"
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

	StopInstances []StopInstance

	StopInstance struct {
		InstanceID    string        `json:"InstanceID"`
		CurrentState  CurrentState  `json:"CurrentState"`
		PreviousState PreviousState `json:"PreviousState"`
	}

	CurrentState struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
	}

	PreviousState struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
	}
)

func NewClient(profile, region string) *Client {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config:  aws.Config{Region: &region},
	}))
	return &Client{ec2: ec2.New(sess)}
}

func (c Client) StopInstances(ctx context.Context, instanceIDs []*string) (StopInstances, error) {
	resp, err := c.ec2.StopInstancesWithContext(ctx, &ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return nil, err
	}

	stopInstances := make(StopInstances, len(resp.StoppingInstances))
	for i, si := range resp.StoppingInstances {
		stopInstances[i] = StopInstance{
			InstanceID: *si.InstanceId,
			CurrentState: CurrentState{
				Code: *si.CurrentState.Code,
				Name: *si.CurrentState.Name,
			},
			PreviousState: PreviousState{
				Code: *si.PreviousState.Code,
				Name: *si.PreviousState.Name,
			},
		}
	}
	return stopInstances, nil
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

func (i Instance) ConvertToMsgMaterial() (*sender.Material, error) {
	m := sender.Material(i)
	return &m, nil
}

func (is Instances) ConvertToMsgMaterials() (sender.Materials, error) {
	msgs := make(sender.Materials, len(is))
	for ix, i := range is {
		msg, err := i.ConvertToMsgMaterial()
		if err != nil {
			return nil, err
		}
		msgs[ix] = *msg
	}
	return msgs, nil
}
