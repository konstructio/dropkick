package cloudflare

import (
	"context"
	"fmt"
	"strings"

	cloudflarego "github.com/cloudflare/cloudflare-go"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

func (c *Cloudflare) NukeDNSRecords(ctx context.Context) error {
	c.logger.Infof("listing DNS records for %q", c.getFullName())

	records, _, err := c.client.ListDNSRecords(ctx, &cloudflarego.ResourceContainer{
		Identifier: c.zoneID,
	}, cloudflarego.ListDNSRecordsParams{})
	if err != nil {
		return fmt.Errorf("unable to list records for domain %q: %w", c.zoneName, err)
	}

	c.logger.Infof("found %d records for %q", len(records), c.getFullName())

	filteredRecords := filterRecords(records, c.getFullName())

	c.logger.Infof("applied filter %q to %d records, got %d records suitable for deletion", c.getFullName(), len(records), len(filteredRecords))

	for _, r := range filteredRecords {
		if c.nuke {
			c.logger.Infof("deleting %s record %q", r.Type, r.Name)
			if err := c.client.DeleteDNSRecord(ctx, &cloudflarego.ResourceContainer{Identifier: c.zoneID}, r.ID); err != nil {
				return fmt.Errorf("unable to delete record %q: %w", r.Name, err)
			}
			outputwriter.WriteStdoutf("deleted %s record %q", r.Type, r.Name)
		} else {
			c.logger.Warnf("refusing to delete %s record %q: nuke is not enabled", r.Type, r.Name)
		}
	}

	return nil
}

// fetch all txt records with values like "heritage=external-dns,external-dns/owner=default,external-dns/resource=ingress/argo/argo-server"
// fetch all txt records with values like "heritage=external-dns,external-dns/owner=default,external-dns/resource"

func filterRecords(records []cloudflarego.DNSRecord, suffix string) []cloudflarego.DNSRecord {
	aRecord := "A"
	txtRecord := "TXT"

	filteredRecords := make([]cloudflarego.DNSRecord, 0, len(records))
	for _, r := range records {
		if strings.HasSuffix(r.Name, suffix) && (r.Type == txtRecord || r.Type == aRecord) {
			filteredRecords = append(filteredRecords, r)
		}
	}

	return filteredRecords
}
