package digitalocean

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var supportedSpacesRegions = [...]string{
	"nyc3",
	"sfo2",
	"sfo3",
	"ams3",
	"fra1",
	"sgp1",
	"blr1",
	"syd1",
}

func generateSpacesEndpoint(region string) (string, error) {
	for _, r := range supportedSpacesRegions {
		if r == region {
			return fmt.Sprintf("https://%s.digitaloceanspaces.com", region), nil
		}
	}

	return "", fmt.Errorf("unsupported region %q for DigitalOcean Spaces", region)
}

func (d *DigitalOcean) NukeS3Storage() error {
	if d.spacesAccessKey == "" || d.spacesSecretKey == "" || d.spacesRegion == "" {
		return fmt.Errorf("DigitalOcean spaces credentials are not set")
	}

	endpoint, err := generateSpacesEndpoint(d.spacesRegion)
	if err != nil {
		return err // nolint: wrapcheck // no need to wrap this error
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(d.spacesRegion),
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(d.spacesAccessKey, d.spacesSecretKey, ""),
		S3ForcePathStyle: aws.Bool(false),
	})
	if err != nil {
		return fmt.Errorf("unable to create Spaces session against DigitalOcean API: %w", err)
	}

	s3svc := s3.New(sess)

	resp, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("unable to list Space buckets for region %q: %w", d.spacesRegion, err)
	}

	for _, bucket := range resp.Buckets {
		objs, err := s3svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: bucket.Name,
		})
		if err != nil {
			return fmt.Errorf("unable to list objects in Space bucket %q, region %q: %w", *bucket.Name, d.spacesRegion, err)
		}

		for _, obj := range objs.Contents {
			if d.nuke {
				_, err := s3svc.DeleteObject(&s3.DeleteObjectInput{
					Bucket: bucket.Name,
					Key:    obj.Key,
				})
				if err != nil {
					return fmt.Errorf("unable to delete object %q from Space bucket %q, region %q: %w", *obj.Key, *bucket.Name, d.spacesRegion, err)
				}
			} else {
				d.logger.Printf("refusing to delete object %q from Space bucket %q, region %q: nuke is not enabled", *obj.Key, *bucket.Name, d.spacesRegion)
			}
		}

		if d.nuke {
			_, err := s3svc.DeleteBucket(&s3.DeleteBucketInput{
				Bucket: bucket.Name,
			})
			if err != nil {
				return fmt.Errorf("unable to delete Space bucket %q, region %q: %w", *bucket.Name, d.spacesRegion, err)
			}
		} else {
			d.logger.Printf("refusing to delete Space bucket %q, region %q: nuke is not enabled\n", *bucket.Name, d.spacesRegion)
		}
	}

	return nil
}
