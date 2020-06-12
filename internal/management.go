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

func (m Manager) Do(ctx context.Context) error {
	for _, a := range m.awsAccounts {
		c := NewClient()
		is, err := c.FetchAllInstances(context.Background())
		if err != nil {
			return err
		}

		msgs, err := is.ConvertToMsgMaterials()
		if err != nil {
			return err
		}

		if err := m.senders[a.Sender].Send(msgs); err != nil {
			return err
		}
	}

	return nil
}
