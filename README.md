[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://en.wikipedia.org/wiki/MIT_License)
[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-green.svg)](https://godoc.org/github.com/gford1000-go/awsregions)

# Regions | Retrieve which AWS regions are accessible

The AWS SDK v2 does not provide a simple mechanism to retrieve accessible AWS regions, defined as either having been
"opted-in" by the AWS account owner, or always accessible by default.

The credentials must provide IAM allow access to the action `ec2.DescribeRegions` for the call to be successful.

```go
func main() {
    // These are provided at runtime via secure mechanism - do not hardcode values
    var accessKeyID = "A"
    var secretAccessKey = "B"

    ctx, _ := awsregions.ContextWithRegionsCredentials(context.Background(), awscredentials.NewAWSCredentials("SomeID", accessKeyID, secretAccessKey))

    rgns, _ := awsregions.GetRegions(ctx)

    for _, r := range rgns {
        fmt.Println(r)
    }
}
```
