// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	AccountId           string `envconfig:"PLUGIN_ACCOUNT_ID"`
	ApiKey              string `envconfig:"PLUGIN_APIKEY"`
	ApplicationId       string `envconfig:"PLUGIN_APPLICATION_ID"`
	ArtifactVersion     string `envconfig:"PLUGIN_ARTIFACT_VERSION"`
	HarnessWebhookId    string `envconfig:"PLUGIN_HARNESS_WEBHOOK_ID"`
	HarnessServiceName  string `envconfig:"PLUGIN_HARNESS_SERVICE_NAME"`
	HarnessArtifactName string `envconfig:"PLUGIN_HARNESS_ARTIFACT_NAME"`
	HarnessUri          string `envconfig:"PLUGIN_HARNESS_URI"`
}

func validateAndSetArgs(args *Args) error {
	if args.AccountId == "" {
		return fmt.Errorf("AccountId must be set in settings")
	}
	if args.ApiKey == "" {
		return fmt.Errorf("ApiKey must be set in settings")
	}
	if args.ApplicationId == "" {
		return fmt.Errorf("ApplictionId must be set in settings")
	}
	if args.ArtifactVersion == "" {
		return fmt.Errorf("ArtifactVersion must be set in settings")
	}
	if args.HarnessWebhookId == "" {
		return fmt.Errorf("HarnessWebhookId must be set in settings")
	}
	if args.HarnessServiceName == "" {
		return fmt.Errorf("HarnessServiceName must be set in settings")
	}
	if args.HarnessArtifactName == "" {
		return fmt.Errorf("HarnessArtifactName must be set in settings")
	}
	if args.HarnessUri == "" {
		return fmt.Errorf("HarnessUri must be set in settings")
	}
	return nil
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	logrus.Debug("Starting Harness Trigger Plug-in")
	err := validateAndSetArgs(&args)
	if err != nil {
		return fmt.Errorf("issues with the parameters passed: %w", err)
	}
	base, err := url.Parse(args.HarnessUri)
	if err != nil {
		return fmt.Errorf("issues parsing url: %w", err)
	}
	base.Path = "/" + args.HarnessWebhookId
	params := url.Values{}
	params.Add("accountId", args.AccountId)
	base.RawQuery = params.Encode()
	postBody, _ := json.Marshal(map[string]string{
		"artifactSourceName": args.HarnessArtifactName,
		"buildNumber":        strconv.Itoa(args.Build.Number),
		"service":            args.HarnessServiceName,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(base.String(), "application/json", responseBody)
	if err != nil {
		return fmt.Errorf("issues making the request.: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("issues reading response: %w", err)
	}
	logrus.Infof("Response: " + string(body))
	logrus.Debug("Complete")
	return nil
}
