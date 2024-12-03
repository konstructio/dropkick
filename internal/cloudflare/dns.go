package cloudflare

import (
	"context"
	"fmt"
	"strings"

	cloudflarego "github.com/cloudflare/cloudflare-go"
)

func (c *Cloudflare) NukeDNSRecords(ctx context.Context) error {
	c.logger.Infof("removing dns records for domain %q", c.zoneName)

	zones, err := c.client.ListZones(ctx)
	if err != nil {
		return fmt.Errorf("unable to list zones: %w", err)
	}

	for _, zone := range zones {
		if zone.Name == c.zoneName {
			c.zoneID = zone.ID
		}
	}
	if c.zoneID == "" {
		return fmt.Errorf("unable to find zone ID for domain: %q", c.zoneName)
	}

	records, _, err := c.client.ListDNSRecords(ctx, &cloudflarego.ResourceContainer{
		Identifier: c.zoneID,
	}, cloudflarego.ListDNSRecordsParams{})
	if err != nil {
		return fmt.Errorf("unable to list records for domain %q: %w", c.zoneName, err)
	}

	// if subdomain is set filter the records
	var subdomainRecords []cloudflarego.DNSRecord
	if c.subdomain != "" {
		for _, r := range records {
			if strings.Contains(r.Name, c.subdomain) {
				subdomainRecords = append(subdomainRecords, r)
			}
		}
		records = subdomainRecords
		c.logger.Infof("found %d records for %q", len(records), fmt.Sprintf("%s.%s", c.subdomain, c.zoneName))
	} else {
		c.logger.Infof("found %d records for %q", len(records), c.zoneName)
	}

	for _, r := range records {
		if c.nuke {
			c.logger.Infof("nuke enabled, deleting record %q", r.Name)
			if err := c.client.DeleteDNSRecord(ctx, &cloudflarego.ResourceContainer{Identifier: c.zoneID}, r.ID); err != nil {
				return fmt.Errorf("unable to delete record %q: %w", r.Name, err)
			}
		} else {
			c.logger.Warnf("nuke disabled, found record %q", r.Name)
		}
	}

	return nil
}
