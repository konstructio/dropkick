package cloudflare

import (
	"context"
	"fmt"
	"strings"

	cloudflarego "github.com/cloudflare/cloudflare-go"
)

func (c *Cloudflare) NukeDNSRecords(ctx context.Context) error {
	c.logger.Infof("removing dns records for domain %q", c.zoneName)

	records, _, err := c.client.ListDNSRecords(ctx, &cloudflarego.ResourceContainer{
		Identifier: c.zoneID,
	}, cloudflarego.ListDNSRecordsParams{})
	if err != nil {
		return fmt.Errorf("unable to list records for domain %q: %w", c.zoneName, err)
	}

	filteredRecords := filterRecords(records, c.subdomain)

	if c.subdomain != "" {
		c.logger.Infof("found %d records for %q", len(filteredRecords), fmt.Sprintf("%s.%s", c.subdomain, c.zoneName))
	} else {
		c.logger.Infof("found %d records for %q", len(filteredRecords), c.zoneName)
	}

	for _, r := range filteredRecords {
		if c.nuke {
			c.logger.Infof("nuke enabled, deleting record %q - %q", r.Type, r.Name)
			if err := c.client.DeleteDNSRecord(ctx, &cloudflarego.ResourceContainer{Identifier: c.zoneID}, r.ID); err != nil {
				return fmt.Errorf("unable to delete record %q: %w", r.Name, err)
			}
		} else {
			c.logger.Warnf("nuke disabled, found record %q - %q", r.Type, r.Name)
		}
	}

	return nil
}

// fetch all txt records with values like "heritage=external-dns,external-dns/owner=default,external-dns/resource=ingress/argo/argo-server"
// fetch all txt records with values like "heritage=external-dns,external-dns/owner=default,external-dns/resource"

func filterRecords(records []cloudflarego.DNSRecord, subdomain string) []cloudflarego.DNSRecord {
	filteredRecords := make([]cloudflarego.DNSRecord, 0, len(records))
	aRecord := "A"
	txtRecord := "TXT"

	for _, r := range records {
		if subdomain != "" {
			if strings.Contains(r.Name, subdomain) && (r.Type == txtRecord || r.Type == aRecord) {
				filteredRecords = append(filteredRecords, r)
			}
		} else {
			if r.Type == txtRecord || r.Type == aRecord {
				filteredRecords = append(filteredRecords, r)
			}
		}
	}
	return filteredRecords
}
