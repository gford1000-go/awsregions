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

type regionalOptions struct {
	id      awscredentials.CredentialsID
	regions []string
}

// ContextWithRegionsCredentials returns a new context based on the supplied parent, which stores the credentials
// needed to be able to retrieve the set of accessible AWS regions.
func ContextWithRegionsCredentials(ctx context.Context, c *awscredentials.AWSCredentials) (context.Context, error) {

	ctx, err := awscredentials.ContextWithAWSCredentials(ctx, c)
	if err != nil {
		return ctx, err
	}

	v := ctx.Value(regionalCredentialsKey)
	if v != nil {
		if opt, ok := v.(*regionalOptions); ok {
			opt.id = c.ID()
			return ctx, nil
		}
	}

	return context.WithValue(ctx, regionalCredentialsKey, &regionalOptions{
		id:      c.ID(),
		regions: []string{},
	}), nil
}

// ContextWithFixedRegions returns a new context based on the supplied parent, which stores a
// static set of AWS regions.
func ContextWithFixedRegions(ctx context.Context, regions []string) (context.Context, error) {

	v := ctx.Value(regionalCredentialsKey)
	if v != nil {
		if opt, ok := v.(*regionalOptions); ok {
			opt.regions = regions
			return ctx, nil
		}
	}

	return context.WithValue(ctx, regionalCredentialsKey, &regionalOptions{
		regions: regions,
	}), nil
}

// ErrMissingRegionCredentials raised if the context does not hold regional credential information
var ErrMissingRegionCredentials = errors.New("context does not contain regional credentials")

// ErrInvalidRegionalDetails raised if the regional details are invalid
var ErrInvalidRegionalDetails = errors.New("invalid regional details in context")

// newEC2Client constructs a client from details in the context
func newEC2Client(ctx context.Context, id awscredentials.CredentialsID) (*ec2.Client, error) {

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

// GetRegions returns the set of AWS regions that are available.
// Preferentially it will use a static set of regions created by a call to ContextWithFixedRegions(),
// otherwise it will attempt to call AWS to determine the accessible regions, i.e.
// either have been opted-in, or are always accessible, using credentials to ContextWithRegionsCredentials().
// In this second case, the context must contain connection details that allow the IAM action ec2.DescribeRegions
// for the call to be successful.
func GetRegions(ctx context.Context) ([]string, error) {

	v := ctx.Value(regionalCredentialsKey)
	if v == nil {
		return nil, ErrMissingRegionCredentials
	}
	opt, ok := v.(*regionalOptions)
	if !ok {
		return nil, ErrInvalidRegionalDetails
	}

	// Contact AWS to retrieve accessible regions, based on the credentials,
	// if we don't already have a fixed set of regions
	if len(opt.regions) == 0 {
		client, err := newEC2Client(ctx, opt.id)
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

		// Store so that only need to call once; rare that regions accessibility changes
		opt.regions = regions
	}

	return opt.regions, nil
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
