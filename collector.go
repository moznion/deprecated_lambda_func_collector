package collector

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type DeprecatedFunction struct {
	Conf            types.FunctionConfiguration
	DeprecationDate string
}

func CollectDeprecatedLambdaFunctions(ctx context.Context, regions []Region, includeFuturePlanned bool, assumeRoleARN string, extractor *DeprecatedFunctionExtractor) ([]*DeprecatedFunction, error) {
	var deprecatedFunctions []*DeprecatedFunction

	for _, region := range regions {
		cfg, err := config.LoadDefaultConfig(ctx, func(options *config.LoadOptions) error {
			options.Region = string(region)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to load the default AWS config: %w", err)
		}

		if assumeRoleARN != "" {
			stsClient := sts.NewFromConfig(cfg)
			provider := stscreds.NewAssumeRoleProvider(stsClient, assumeRoleARN)
			cfg.Credentials = aws.NewCredentialsCache(provider)
		}

		lambdaSvc := lambda.NewFromConfig(cfg)
		functions, err := lambdaSvc.ListFunctions(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to call the list-functions Lambda API: %w", err)
		}
		deprecatedFunctions = append(deprecatedFunctions, extractor.Extract(functions.Functions, includeFuturePlanned)...)
	}

	return deprecatedFunctions, nil
}
