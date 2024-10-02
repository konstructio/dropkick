package civo

import (
	"context"
	"fmt"
	"testing"

	"github.com/konstructio/dropkick/internal/logger"
)

func TestNew(t *testing.T) {
	cases := []struct {
		Name       string
		Opts       []Option
		Token      string
		Region     string
		APIURL     string
		Context    context.Context
		Logger     *logger.Logger
		Nuke       bool
		NameFilter string
		WantErr    bool
	}{
		{
			Name:       "all good",
			Token:      "token",
			Region:     "region",
			APIURL:     "https://api.example.com",
			Context:    context.Background(),
			Logger:     logger.None,
			Nuke:       true,
			NameFilter: "filter",
			WantErr:    false,
		},
		{
			Name: "errored out option",
			Opts: []Option{
				func(c *Civo) error {
					return fmt.Errorf("error")
				},
			},
			WantErr: true,
		},
		{
			Name:    "missing token",
			Token:   "",
			WantErr: true,
		},
		{
			Name:    "missing region",
			Token:   "token",
			Region:  "",
			WantErr: true,
		},
		{
			Name:    "impossible client using invalid URL",
			Token:   "token",
			Region:  "region",
			APIURL:  "#@$%^&*",
			WantErr: true,
		},
		{
			Name:    "missing context",
			Token:   "token",
			Region:  "region",
			Context: nil,
			Logger:  nil,
		},
		{
			Name:    "default logger",
			Token:   "token",
			Region:  "region",
			Logger:  nil,
			WantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(tt *testing.T) {
			opts := append([]Option{}, tc.Opts...)

			if tc.Token != "" {
				opts = append(opts, WithToken(tc.Token))
			}

			if tc.Region != "" {
				opts = append(opts, WithRegion(tc.Region))
			}

			if tc.APIURL != "" {
				opts = append(opts, WithAPIURL(tc.APIURL))
			}

			if tc.Logger != nil {
				opts = append(opts, WithLogger(tc.Logger))
			}

			if tc.NameFilter != "" {
				opts = append(opts, WithNameFilter(tc.NameFilter))
			}

			opts = append(opts, WithNuke(tc.Nuke))

			client, err := New(opts...)

			if tc.WantErr {
				if err == nil {
					tt.Fatal("expected err to not be nil")
				}

				return
			}

			if err != nil {
				tt.Fatalf("expected err to be nil, got %v", err)
			}

			if client == nil {
				tt.Fatal("expected client to not be nil")
			}

			if client.client == nil {
				tt.Fatal("expected client.client to not be nil")
			}

			if client.token != tc.Token {
				tt.Fatalf("expected client.token to be %q, got %q", tc.Token, client.token)
			}

			if client.region != tc.Region {
				tt.Fatalf("expected client.region to be %q, got %q", tc.Region, client.region)
			}

			if tc.APIURL == "" && client.apiURL != civoAPIURL {
				tt.Fatalf("expected client.apiURL to be %q, got %q", civoAPIURL, client.apiURL)
			}

			if tc.APIURL != "" && client.apiURL != tc.APIURL {
				tt.Fatalf("expected client.apiURL to be %q, got %q", tc.APIURL, client.apiURL)
			}

			if tc.Logger == nil && client.logger != logger.None {
				tt.Fatalf("expected client.logger to be the default logger, got: %#v", client.logger)
			}

			if tc.Logger != nil && client.logger != tc.Logger {
				tt.Fatalf("expected client.logger to be %#v, got %#v", tc.Logger, client.logger)
			}

			if client.nuke != tc.Nuke {
				tt.Fatalf("expected client.nuke to be %v, got %v", tc.Nuke, client.nuke)
			}

			if client.nameFilter != tc.NameFilter {
				tt.Fatalf("expected client.nameFilter to be %q, got %q", tc.NameFilter, client.nameFilter)
			}
		})
	}
}
