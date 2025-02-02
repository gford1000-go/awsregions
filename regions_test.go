package awsregions

import (
	"context"
	"testing"

	"github.com/gford1000-go/awscredentials"
)

func TestGetRegions(t *testing.T) {

	regions := []string{"eu-west-1", "ap-southeast-1"}

	ctx, err := ContextWithFixedRegions(context.Background(), regions)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	retrieved, err := GetRegions(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Unexpected failure to return regions")
	}
	if len(retrieved) != len(regions) {
		t.Fatalf("Mismatch in regions: expected %d, got %d: %v, %v", len(regions), len(retrieved), regions, retrieved)
	}
	for i := 0; i < len(regions); i++ {
		if retrieved[i] != regions[i] {
			t.Fatalf("Mismatch in entries: at %d, expected %s, got %s", i, regions[i], retrieved[i])
		}
	}
}

func TestGetRegions_1(t *testing.T) {

	// Show failure when invalid keys are used

	keyBad := "A"
	secBad := "B"

	ctx, err := ContextWithRegionsCredentials(context.Background(), awscredentials.NewAWSCredentials("IDBad", keyBad, secBad))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	_, err = GetRegions(ctx)
	if err == nil {
		t.Fatalf("Unexpected success: %v", err)
	}
}

// func TestGetRegions_2(t *testing.T) {

// 	// Show success with valid keys

// 	// Add real credentials only when manually testing, and do not check into source code
// 	keyValid := "XXX"
// 	secretValid := "YYY"

// 	ctx, err := ContextWithRegionsCredentials(context.Background(), awscredentials.NewAWSCredentials("IDValid", keyValid, secretValid))
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}

// 	retrieved, err := GetRegions(ctx)
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}

// 	if len(retrieved) == 0 {
// 		t.Fatal("Unexpected response - at least some (nonoptional) regions should be returned")
// 	}
// }

// func TestGetRegions_3(t *testing.T) {

// 	// Show key replacement on repeated calls

// 	// Add bad keys, and then replace with valid keys
// 	key := "A"
// 	sec := "B"

// 	ctx, err := ContextWithRegionsCredentials(context.Background(), awscredentials.NewAWSCredentials("IDBad", key, sec))
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}

// 	// Add real credentials only when manually testing, and do not check into source code
// 	keyValid := "XXX"
// 	secretValid := "YYY"

// 	ctx, err = ContextWithRegionsCredentials(ctx, awscredentials.NewAWSCredentials("IDValid", keyValid, secretValid))
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}

// 	retrieved, err := GetRegions(ctx)
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}

// 	if len(retrieved) == 0 {
// 		t.Fatal("Unexpected response - at least some (nonoptional) regions should be returned")
// 	}
// }

// func TestGetRegions_4(t *testing.T) {

// 	// Show key replacement on repeated calls

// 	// Add real credentials only when manually testing, and do not check into source code
// 	keyValid := "XXX"
// 	secretValid := "YYY"

// 	ctx, err := ContextWithRegionsCredentials(context.Background(), awscredentials.NewAWSCredentials("IDValid", keyValid, secretValid))
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}

// 	// Add bad keys, and then replace with valid keys
// 	keyBad := "A"
// 	secBad := "B"

// 	ctx, err = ContextWithRegionsCredentials(ctx, awscredentials.NewAWSCredentials("IDBad", keyBad, secBad))
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}

// 	_, err = GetRegions(ctx)
// 	if err == nil {
// 		t.Fatal("Unexpected success")
// 	}
// }

func TestIsUsable(t *testing.T) {

	regions := []string{"eu-west-1", "ap-southeast-1"}

	ctx, err := ContextWithFixedRegions(context.Background(), regions)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	type testData struct {
		region   string
		expected bool
	}

	tests := []testData{
		{
			region:   "eu-west-1",
			expected: true,
		},
		{
			region:   "invalid-region",
			expected: false,
		},
	}

	for _, test := range tests {
		ok, err := IsUsable(ctx, test.region)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if ok && !test.expected {
			t.Fatalf("Unexpected result - returnd usable, when %s is not usable", test.region)
		}
		if !ok && test.expected {
			t.Fatalf("Unexpected result - returnd unusable, when %s is usable", test.region)
		}
	}
}

func TestIsUsable_1(t *testing.T) {

	regions := []string{"ap-southeast-1"}

	ctx, err := ContextWithFixedRegions(context.Background(), regions)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	testRegion := "eu-west-1"

	ok, err := IsUsable(ctx, testRegion)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if ok {
		t.Fatal("Unexpectedly received usable, when not usable")
	}

	ctx, err = ContextWithFixedRegions(ctx, append(regions, testRegion))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	ok, err = IsUsable(ctx, testRegion)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("Unexpectedly received unusable, when should now be usable")
	}
}
