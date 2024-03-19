package civo

import (
	"github.com/civo/civogo"
)

func NewClient(civoToken string, region string) *civogo.Client {
	civoClient, _ := civogo.NewClient(civoToken, region)

	return civoClient
}
