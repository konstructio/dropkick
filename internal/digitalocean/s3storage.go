package digitalocean

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

func (d *DigitalOcean) NukeS3Storage() error {
	d.logger.Printf("listing Space buckets for region %q", d.spacesRegion)

	resp, err := d.s3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("unable to list Space buckets for region %q: %w", d.spacesRegion, err)
	}

	d.logger.Printf("found %d Space buckets", len(resp.Buckets))

	for _, bucket := range resp.Buckets {
		d.logger.Printf("found Space bucket %q, region %q", *bucket.Name, d.spacesRegion)

		objs, err := d.s3svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: bucket.Name,
		})
		if err != nil {
			return fmt.Errorf("unable to list objects in Space bucket %q, region %q: %w", *bucket.Name, d.spacesRegion, err)
		}

		for _, obj := range objs.Contents {
			if d.nuke {
				d.logger.Printf("deleting object %q from Space bucket %q, region %q", *obj.Key, *bucket.Name, d.spacesRegion)
				_, err := d.s3svc.DeleteObject(&s3.DeleteObjectInput{
					Bucket: bucket.Name,
					Key:    obj.Key,
				})
				if err != nil {
					return fmt.Errorf("unable to delete object %q from Space bucket %q, region %q: %w", *obj.Key, *bucket.Name, d.spacesRegion, err)
				}
				outputwriter.WriteStdout("deleted object %q from Space bucket %q, region %q", *obj.Key, *bucket.Name, d.spacesRegion)
			} else {
				d.logger.Printf("refusing to delete object %q from Space bucket %q, region %q: nuke is not enabled", *obj.Key, *bucket.Name, d.spacesRegion)
			}
		}

		if d.nuke {
			d.logger.Printf("deleting Space bucket %q, region %q", *bucket.Name, d.spacesRegion)
			_, err := d.s3svc.DeleteBucket(&s3.DeleteBucketInput{
				Bucket: bucket.Name,
			})
			if err != nil {
				return fmt.Errorf("unable to delete Space bucket %q, region %q: %w", *bucket.Name, d.spacesRegion, err)
			}
			outputwriter.WriteStdout("deleted Space bucket %q, region %q", *bucket.Name, d.spacesRegion)
		} else {
			d.logger.Printf("refusing to delete Space bucket %q, region %q: nuke is not enabled\n", *bucket.Name, d.spacesRegion)
		}
	}

	return nil
}
