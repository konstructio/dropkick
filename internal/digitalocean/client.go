package digitalocean

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/digitalocean/godo"
	"github.com/konstructio/dropkick/internal/logger"
)

// DigitalOcean is a client for the DigitalOcean API.
type DigitalOcean struct {
	client          *godo.Client   // The underlying DigitalOcean API client.
	s3svc           *s3.S3         // The underlying DigitalOcean Spaces API client.
	nuke            bool           // Whether to nuke resources.
	token           string         // The API token.
	logger          *logger.Logger // The logger instance.
	spacesAccessKey string         // The access key for Spaces.
	spacesSecretKey string         // The secret key for Spaces.
	spacesRegion    string         // The region for Spaces.
}

// Option is a function that configures a DigitalOcean.
type Option func(*DigitalOcean) error

// WithLogger sets the logger for a DigitalOcean.
func WithLogger(logger *logger.Logger) Option {
	return func(c *DigitalOcean) error {
		c.logger = logger
		return nil
	}
}

// WithToken sets the API token for a DigitalOcean.
func WithToken(token string) Option {
	return func(c *DigitalOcean) error {
		c.token = token
		return nil
	}
}

// WithNuke sets whether to nuke resources for a DigitalOcean.
func WithNuke(nuke bool) Option {
	return func(c *DigitalOcean) error {
		c.nuke = nuke
		return nil
	}
}

func WithS3Storage(accessKey, secretKey, region string) Option {
	return func(c *DigitalOcean) error {
		c.spacesAccessKey = accessKey
		c.spacesSecretKey = secretKey
		c.spacesRegion = region
		return nil
	}
}

// New creates a new DigitalOcean with the given options.
// It returns an error if the token or region is not set, or if it fails to
// create the underlying DigitalOcean API client.
func New(ctx context.Context, opts ...Option) (*DigitalOcean, error) {
	c := &DigitalOcean{}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("unable to apply option: %w", err)
		}
	}

	if c.token == "" {
		return nil, errors.New("required token not found")
	}

	// The DigitalOcean API client does not authenticate until a request is made,
	// so we're requesting the account but not using it to verify the token.
	godoClient := godo.NewFromToken(c.token)
	_, _, err := godoClient.Account.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate DigitalOcean client: %w", err)
	}

	c.client = godoClient

	// Set up S3 storage
	if c.spacesAccessKey == "" || c.spacesSecretKey == "" || c.spacesRegion == "" {
		return nil, errors.New("DigitalOcean spaces credentials are not set")
	}

	endpoint, err := generateSpacesEndpoint(c.spacesRegion)
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(c.spacesRegion),
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(c.spacesAccessKey, c.spacesSecretKey, ""),
		S3ForcePathStyle: aws.Bool(false),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create Spaces session against DigitalOcean API: %w", err)
	}

	s3svc := s3.New(sess)
	c.s3svc = s3svc

	// Validate s3 credentials work
	_, err = s3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to validate Spaces credentials: unable to list buckets in the account: %w", err)
	}

	if c.logger == nil {
		c.logger = logger.None
	}

	return c, nil
}

func generateSpacesEndpoint(region string) (string, error) {
	supportedSpacesRegions := [...]string{
		"nyc3",
		"sfo2",
		"sfo3",
		"ams3",
		"fra1",
		"sgp1",
		"blr1",
		"syd1",
	}

	for _, r := range supportedSpacesRegions {
		if r == region {
			return fmt.Sprintf("https://%s.digitaloceanspaces.com", region), nil
		}
	}

	return "", fmt.Errorf("unsupported region %q for DigitalOcean Spaces", region)
}
