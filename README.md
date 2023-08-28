# Deprecated AWS Lambda Function Collector

A toolchain to scan the deprecated AWS Lambda functions.

## Usage: CLI

```
Usage of deprecated_lambda_func_collector:
  -all-regions
        scan all lambda supported regions; this parameter take priority over the -regions
  -assume-role-arn string
        a role ARN for assume-role
  -output string
        the output format: "json" or "csv" (default "json")
  -regions string
        target regions to scan; this value can have multiple regions as comma-separated string
  -without-future-planned-deprecation
        exclude the future planned deprecations from the result
```

Example:

```
$ AWS_PROFILE="YOUR-AWS-PROFILE" deprecated_lambda_func_collector --regions="ap-northeast-1,us-east-1"
```

NOTE: as the above, `AWS_PROFILE` is required to specify the target AWS environment.

CLI output Example:

```
[
  {
    "arn": "arn:aws:lambda:ap-northeast-1:123456789012:function:go-function",
    "deprecated_runtime": "go1.x",
    "deprecation_date": "Dec 31, 2023"
  },
  ...,
  {
    "arn": "arn:aws:lambda:us-east-1:123456789012:function:node14-function",
    "deprecated_runtime": "nodejs14.x",
    "deprecation_date": "Nov 27, 2023"
  }
]
```

