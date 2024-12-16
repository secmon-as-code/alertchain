# Deployment

## Deploy as a container image

There are several ways to create a container image, but in this section, we will introduce a method to create a container image based on the AlertChain image and add alert policies and action policies. First, create a Dockerfile like the following.

```dockerfile
FROM ghcr.io/secmon-lab/alertchain:v0.0.2

COPY policy /policy

WORKDIR /
EXPOSE 8080
ENTRYPOINT ["/alertchain", "-d", "/policy", "--log-format", "json", "serve", "--addr", "0.0.0.0:8080"]
```

In this Dockerfile, we use the AlertChain image as a base and add alert policies and action policies to the container image. The alert policies and action policies are added to the container image by placing them in the `policy` directory. Also, when starting the AlertChain server, we specify the `--addr` option to make it accessible from outside the container.

For instructions on how to deploy the created image to various runtime environments, please refer to the documentation for each runtime environment.

## Deploy to AWS Lambda

For deploying to AWS Lambda, using CDK makes it easy to deploy. First, install CDK and create a CDK project. For instructions on how to create a project, please refer to [this guide](https://docs.aws.amazon.com/cdk/latest/guide/getting_started.html).

AWS Lambda can be triggered by various events, but the schema of the data received varies depending on the type of event. Therefore, AlertChain needs to create handlers according to the event type. For example, to process events from SQS, specify `NewSQSHandler`, and to process events from Functional URL, specify `NewFunctionalURLHandler`. AlertChain provides functions to create handlers according to the event type.

There are several ways to place policy files in the AWS Lambda execution environment, but in this case, we will use the Go embed package to embed them in the binary. First, write Go code like the following.

```go
package main

import (
        "embed"

        "github.com/aws/aws-lambda-go/lambda"
        ac "github.com/secmon-lab/alertchain/pkg/controller/lambda"
)

//go:embed policy/*
var policyFS embed.FS

func main() {
        lambda.Start(ac.New(
                // Register a handler to process events from SQS
                // Specify the guardduty schema and enable decoding events from SNS
                ac.NewSQSHandler("guardduty", ac.WithDecodeSNS()),

                // Use embed.FS for file reading
                ac.WithReadFile(policyFS.ReadFile),

                // Specify the policy file directory
                ac.WithAlertPolicyDir("policy"),
                ac.WithActionPolicyDir("policy"),
        ))
}
```

This code works as follows:

- Reads files under the `policy` directory, creates alert policies and action policies.
- Processes GuardDuty findings sent to SQS via SNS using CloudWatch Event.

Build this code for AWS Lambda and create a binary named `build/main`.

```bash
$ env GOARCH=amd64 GOOS=linux go build -o ./build/main
```

Next, prepare the CDK code to deploy the created binary. The following is an example written in TypeScript.

```ts
import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as sqs from "aws-cdk-lib/aws-sqs";
import * as eventSources from "aws-cdk-lib/aws-lambda-event-sources";

export class AlertchainCdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const queue = sqs.Queue.fromQueueArn(
      this,
      "alertchain-queue",
      "arn:aws:sqs:ap-northeast-1:111111111:guardduty-alert-queue"
    );

    const f = new lambda.Function(this, "alertchain", {
      runtime: lambda.Runtime.GO_1_X,
      handler: "main",
      code: lambda.Code.fromAsset("./build"),
      timeout: cdk.Duration.seconds(30),
    });
    f.addEventSource(new eventSources.SqsEventSource(queue));
  }
}
```

By deploying this, an AWS Lambda function that processes GuardDuty findings will be created.