package awsregions

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gford1000-go/awscredentials"
)

type regionalCredentials struct {
	key string
}

var regionalCredentialsKey *regionalCredentials = &regionalCredentials{key: "regionalCredentialsKey"}

// ContextWithRegionsCredentials returns a new context based on the supplied, which stores the credentials
// needed to be able to retrieve the set of accessible AWS regions.
func ContextWithRegionsCredentials(ctx context.Context, c *awscredentials.AWSCredentials) (context.Context, error) {

	ctx, err := awscredentials.ContextWithAWSCredentials(ctx, c)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, regionalCredentialsKey, c.ID()), nil
}

// ErrMissingRegionCredentials raised if the context does not hold regional credential information
var ErrMissingRegionCredentials = errors.New("context does not contain regional credentials")

// ErrInvalidRegionalCredentialID raised if the regional credentials ID is invalid
var ErrInvalidRegionalCredentialID = errors.New("invalid ID found for regional credentials")

// newEC2Client constructs a client from details in the context
func newEC2Client(ctx context.Context) (*ec2.Client, error) {

	v := ctx.Value(regionalCredentialsKey)
	if v == nil {
		return nil, ErrMissingRegionCredentials
	}
	id, ok := v.(awscredentials.CredentialsID)
	if !ok {
		return nil, ErrInvalidRegionalCredentialID
	}

	provider, err := awscredentials.GetCredentialsProvider(ctx, id)
	if err != nil {
		return nil, err
	}

	var cfg aws.Config
	if provider != nil {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(provider))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}
	if err != nil {
		return nil, err
	}

	return ec2.NewFromConfig(cfg), nil
}

// GetRegions returns the set of AWS regions that are accessible, i.e.
// either have been opted-in, or are always accessible.
// The context must contain connection details that allow the IAM action ec2.DescribeRegions
// for the call to be successful
func GetRegions(ctx context.Context) ([]string, error) {

	client, err := newEC2Client(ctx)
	if err != nil {
		return nil, err
	}

	output, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	regions := []string{}
	for _, r := range output.Regions {
		if r.OptInStatus != nil && (*r.OptInStatus == "opt-in-not-required" || *r.OptInStatus == "opted-in") {
			regions = append(regions, *r.RegionName)
		}
	}

	return regions, nil
}

// IsUsable returns true if the specified region can be used.
// The context must contain connection details that allow the IAM action ec2.DescribeRegions
// for the call to be successful
func IsUsable(ctx context.Context, region string) (bool, error) {

	regions, err := GetRegions(ctx)
	if err != nil {
		return false, err
	}

	for _, r := range regions {
		if r == region {
			return true, nil
		}
	}

	return false, nil
}
