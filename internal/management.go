package internal

import (
	"context"

	"github.com/yyh-gl/go-ec2/internal/sender"
	"github.com/yyh-gl/go-ec2/internal/sender/slack"
)

const SenderVendorSlack = "slack"

type (
	Manager struct {
		awsAccounts map[string]AWSAccount
		senders     map[string]sender.Sender
	}
)

func NewManger(configPath string) *Manager {
	cfg := loadConfigFile(configPath)

	senders := make(map[string]sender.Sender)
	for vendor, ss := range cfg.SenderSet {
		switch vendor {
		case SenderVendorSlack:
			for id, s := range ss {
				senders[id] = slack.NewClient(s)
			}
		}
	}

	return &Manager{
		awsAccounts: cfg.AWSAccountSet,
		senders:     senders,
	}
}

func (m Manager) PrintAllState(ctx context.Context) error {
	for _, a := range m.awsAccounts {
		c := NewClient(a.Profile, a.Region)
		is, err := c.FetchAllInstances(ctx)
		if err != nil {
			return err
		}

		msgs, err := is.ConvertToMsgMaterials()
		if err != nil {
			return err
		}

		if err := m.senders[a.Sender].Send(a.Name, msgs); err != nil {
			return err
		}
	}

	return nil
}

func (m Manager) StopAllInstances(ctx context.Context) error {
	for _, a := range m.awsAccounts {
		c := NewClient(a.Profile, a.Region)
		is, err := c.FetchAllInstances(ctx)
		if err != nil {
			return err
		}

		targetInstaces := make(Instances, 0)
		targetInstanceIDs := make([]*string, 0)
		for _, i := range is {
			if !isExcludedInstance(a.Exclusions, i.ID) {
				i := i
				i.State = "stopping"
				targetInstaces = append(targetInstaces, i)
				targetInstanceIDs = append(targetInstanceIDs, &i.ID)
			}
		}

		if len(targetInstanceIDs) > 0 {
			// TODO: 返ってきたインスタンス情報からメッセージを組み立てる
			if _, err := c.StopInstances(ctx, targetInstanceIDs); err != nil {
				return err
			}

			mtrs, err := targetInstaces.ConvertToMsgMaterials()
			if err != nil {
				return err
			}

			if err := m.senders[a.Sender].Send(a.Name, mtrs); err != nil {
				return err
			}
		}
	}

	return nil
}

func isExcludedInstance(exclusions []string, instanceID string) bool {
	for _, e := range exclusions {
		if instanceID == e {
			return true
		}
	}
	return false
}
