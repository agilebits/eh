package example

import (
	"testing"

	"github.com/agilebits/eh/secrets"
	"github.com/hashicorp/hcl"
)

type smtpConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}

type slackConfig struct {
	Channel1 channelConfig
	Channel2 channelConfig
}

type channelConfig struct {
	Hook string
}

type sharedServiceConfig struct {
	Username string
	Password string
}

type appConfig struct {
	SMTP          smtpConfig
	Slack         slackConfig
	SharedService sharedServiceConfig `hcl:"shared_service"`
}

func TestCanReadAndDecodeLocalExampleWithInclude(t *testing.T) {
	contents, err := secrets.Read("./local-example.hcl")
	if err != nil {
		t.Fatal(err)
	}

	hclobj, err := hcl.ParseBytes(contents)
	if err != nil {
		t.Fatal(err)
	}

	var cfg appConfig
	if err := hcl.DecodeObject(&cfg, hclobj); err != nil {
		t.Fatal(err)
	}

	if cfg.SMTP.Username != "AKIAI4JG42A2LILVBNNZ" {
		t.Errorf("unexpected SMTP username: %q", cfg.SMTP.Username)
	}

	if cfg.Slack.Channel1.Hook != "https://slack.com/1" {
		t.Errorf("unexpected slack.channel1.hook  %+q", cfg.Slack.Channel1.Hook)
	}

	if cfg.Slack.Channel2.Hook != "https://slack.com/2" {
		t.Errorf("unexpected slack.channel2.hook %+q", cfg.Slack.Channel2.Hook)
	}

	if cfg.SharedService.Username != "shared-fragment-user" {
		t.Errorf("unexpected shared_service.username %+q", cfg.SharedService.Username)
	}

	if cfg.SharedService.Password != "shared-fragment-password" {
		t.Errorf("unexpected shared_service.password %+q", cfg.SharedService.Password)
	}
}
