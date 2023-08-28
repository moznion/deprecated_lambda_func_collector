package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/moznion/deprecated_lambda_func_collector"
)

func main() {
	var commaSeparatedRegionsString string
	var assumeRoleARN string
	var outputStyle string
	var allRegions bool
	var excludeFuturePlannedDeprecation bool
	flag.StringVar(&commaSeparatedRegionsString, "regions", "", "target regions to scan; this value can have multiple regions as comma-separated string")
	flag.StringVar(&assumeRoleARN, "assume-role-arn", "", "a role ARN for assume-role")
	flag.StringVar(&outputStyle, "output", "json", "the output format: \"json\" or \"csv\" (default: \"json\")")
	flag.BoolVar(&allRegions, "all-regions", false, "scan all regions; this parameter take priority over the `-regions`.")
	flag.BoolVar(&excludeFuturePlannedDeprecation, "without-future-planned-deprecation", false, "exclude the future planned deprecations from the result (default: false)")
	flag.Parse()

	if commaSeparatedRegionsString == "" && !allRegions {
		flag.Usage()
		log.Fatalf("[ERROR] missing the mandatory `-regions` or `-all-regions` parameter")
	}

	var presenter collector.Presenter
	if outputStyle == "json" {
		presenter = &collector.JSONPresenter{}
	} else if outputStyle == "csv" {
		presenter = &collector.CSVPresenter{}
	} else {
		flag.Usage()
		log.Fatalf("[ERROR] invalid `-output` parameter; this must be \"json\" or \"csv\"")
	}

	regions := collector.AllLambdaSupportedRegions
	if !allRegions {
		regionStrings := strings.Split(commaSeparatedRegionsString, ",")
		regions = make([]collector.Region, len(regionStrings))
		for i, r := range regionStrings {
			regions[i] = collector.Region(strings.TrimSpace(r))
		}
	}

	ctx := context.Background()

	extractor, err := collector.NewDeprecatedFunctionExtractor()
	if err != nil {
		log.Fatalf("failed to make a deprecated functions exttractor; %s", err)
	}

	deprecatedFunctions, err := collector.CollectDeprecatedLambdaFunctions(ctx, regions, !excludeFuturePlannedDeprecation, assumeRoleARN, extractor)
	if err != nil {
		log.Fatalf("failed to collect the deprecated functions; %s", err)
	}

	out, err := presenter.Render(deprecatedFunctions)
	if err != nil {
		log.Fatalf("failed to render the deprecated functions' result; %s", err)
	}

	fmt.Printf("%s\n", out)
}