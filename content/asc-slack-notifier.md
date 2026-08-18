---
title: asc-slack-notifier
import_path: go.ngs.io/asc-slack-notifier
repo_url: https://github.com/ngs/asc-slack-notifier
description: Receive App Store Connect webhook notifications and post them to Slack. Deployable to Cloud Run or AWS API Gateway + Lambda.
version: ""
documentation_url: https://pkg.go.dev/go.ngs.io/asc-slack-notifier
license: MIT
author: ngs
created_at: 2026-08-06T23:13:14Z
updated_at: 2026-08-13T04:54:58Z
---

# asc-slack-notifier

Receive [App Store Connect webhook](https://developer.apple.com/documentation/appstoreconnectapi/webhooks)
notifications and post them to Slack.

One binary, two deployment targets: it runs either as a plain HTTP server
(Cloud Run, Kubernetes, anywhere) or behind AWS API Gateway on Lambda, selected
by an environment variable.

## Features

- HMAC-SHA256 verification of the `x-apple-signature` header against the raw
  request body, compared in constant time.
- Block Kit formatted Slack messages, with an emoji per App Store Connect state
  (`READY_FOR_REVIEW` 📝, `PENDING_APPLE_RELEASE` ⏳, `READY_FOR_DISTRIBUTION` ✅,
  `REJECTED` ❌, …).
- Event types Apple has not documented yet are still delivered, using a generic
  key/value rendering — notifications are never silently dropped.
- Optional App Store Connect API enrichment: with an API key configured, version
  state and build upload notifications name the app, version and build number
  and carry an "Open in App Store Connect" button instead of a bare UUID.
- Slack destination is either an Incoming Webhook URL or a bot token plus
  channel (`chat.postMessage`).
- Returns `502` when Slack cannot be reached, so App Store Connect redelivers.
- Structured JSON logging via `log/slog`.
- Minimal dependencies: the standard library plus `aws-lambda-go` and
  `aws-lambda-go-api-proxy`.

## Quick start

```sh
export ASC_WEBHOOK_SECRET='your-webhook-secret'
export SLACK_WEBHOOK_URL='https://hooks.slack.com/services/T000/B000/XXXX'

go run ./cmd/asc-slack-notifier
# listening on :8080, webhook at POST /webhook, health check at GET /health
```

Send a signed request by hand:

```sh
BODY='{"data":{"type":"appStoreVersionAppVersionStateUpdated","id":"1","attributes":{"oldValue":"PREPARE_FOR_SUBMISSION","newValue":"READY_FOR_REVIEW","timestamp":"2025-04-16T05:00:52.745Z"}}}'
SIG=$(printf '%s' "$BODY" | openssl dgst -sha256 -hmac "$ASC_WEBHOOK_SECRET" -hex | awk '{print $NF}')

curl -sS -X POST http://localhost:8080/webhook \
  -H 'Content-Type: application/json' \
  -H "x-apple-signature: hmacsha256=$SIG" \
  -d "$BODY"
```

## Configuration

| Variable | Required | Default | Description |
|---|---|---|---|
| `ASC_WEBHOOK_SECRET` | yes | – | Shared secret registered with the App Store Connect webhook. Startup fails when unset |
| `SLACK_WEBHOOK_URL` | one of | – | Slack Incoming Webhook URL. Wins when both destinations are set |
| `SLACK_BOT_TOKEN` | one of | – | Slack bot token, used with `chat.postMessage`. Requires `SLACK_CHANNEL` |
| `SLACK_CHANNEL` | one of | – | Target channel for `chat.postMessage`, e.g. `#releases` |
| `RUN_MODE` | no | auto | `http` or `lambda`. Auto-detected as `lambda` when `AWS_LAMBDA_FUNCTION_NAME` is set, `http` otherwise |
| `PORT` | no | `8080` | Listen port in `http` mode |
| `WEBHOOK_PATH` | no | `/webhook` | Path notifications are posted to |
| `HEALTH_PATH` | no | `/health` | Path answering health checks |
| `NOTIFY_PING` | no | `true` | Set to `false` to acknowledge webhook pings without posting to Slack |
| `LOG_LEVEL` | no | `info` | `debug`, `info`, `warn` or `error` |
| `ASC_API_KEY_ID` | no | – | App Store Connect API key ID. Enables message enrichment, see below |
| `ASC_API_ISSUER_ID` | no | – | App Store Connect API issuer ID |
| `ASC_API_PRIVATE_KEY` | no | – | PEM contents of the API key, or the whole PEM base64-encoded. Literal `\n` sequences in a plain PEM value are expanded, so a single-line value works |
| `ASC_API_PRIVATE_KEY_PATH` | no | – | Path to the `.p8` key file, read at startup. `ASC_API_PRIVATE_KEY` wins when both are set |

Startup fails when no Slack destination is configured.

`HEALTH_PATH` defaults to `/health` rather than `/healthz` because the Google
frontend in front of Cloud Run can intercept `/healthz` before it reaches the
container. Set `HEALTH_PATH` if your platform expects a different path. It must
differ from `WEBHOOK_PATH`; startup fails when the two collide.

### About `ASC_WEBHOOK_SECRET`

Apple does not issue this secret — you generate it yourself and use the same
value in two places: this service's environment, and the `attributes.secret`
field when you create the webhook with `POST /v1/webhooks`. App Store Connect
signs every delivery with it, and the service recomputes the signature to verify
the request really came from Apple.

```sh
openssl rand -hex 32
```

Any sufficiently long random string works; treat it like a password. If the two
copies drift apart, every delivery is rejected with `401`.

### App Store Connect API enrichment (optional)

A webhook payload identifies the resource it is about by a bare UUID and nothing
else, so an `appStoreVersionAppVersionStateUpdated` notification cannot say
which app or version moved. Configure an App Store Connect API key and the
service looks the resource up before posting, adding **App**, **Version** and
**Build** fields plus an **Open in App Store Connect** button that jumps
straight to the app's distribution page.

Build upload and build state notifications (`buildUploadStateUpdated`,
`buildBetaDetailExternalBuildStateUpdated`) are enriched the same way, from the
build itself: the app, the pre-release version, the build number, and a button
opening the app's TestFlight page.

Create the key in App Store Connect → **Users and Access** → **Integrations** →
**App Store Connect API**. The `Developer` or `App Manager` role is enough — the
service only reads. Download the `.p8` file (it is offered once), and note the
key ID and the issuer ID shown on the same page.

```sh
export ASC_API_KEY_ID='ABCD123456'
export ASC_API_ISSUER_ID='69a6de70-0000-0000-0000-000000000000'
export ASC_API_PRIVATE_KEY_PATH='/secrets/AuthKey_ABCD123456.p8'
# or, for platforms that only carry string secrets:
# export ASC_API_PRIVATE_KEY="$(cat AuthKey_ABCD123456.p8)"
# or base64-encoded, as fastlane's key_content accepts it:
# export ASC_API_PRIVATE_KEY="$(base64 -i AuthKey_ABCD123456.p8)"
```

`ASC_API_PRIVATE_KEY` accepts the PEM as is or base64 encoded; the two are
detected automatically, so no extra flag is needed.

The feature is entirely optional: with none of the `ASC_API_*` variables set,
messages are rendered from the webhook payload alone, exactly as before. Setting
some but not all of them is a configuration error and fails at startup, so a
typo cannot silently disable enrichment. If the API call fails at delivery time
the failure is logged and the Slack message is posted un-enriched — an App Store
Connect outage never costs you a notification.

## Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `${WEBHOOK_PATH}` | Webhook receiver. `401` on signature failure, `400` on malformed payloads, `502` when Slack rejects the message, `200` otherwise |
| `GET` | `${HEALTH_PATH}` | Health check, returns `200 ok` |

The request body is capped at 1 MiB.

## Deploy

The quickest route is to fork this repository and let GitHub Actions deploy it:
set the `DEPLOY_TARGET` repository variable to `cloudrun` or `lambda`, add the
secrets, and every push to `master`/`main` deploys your instance.
[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) walks through the variables, secrets
and one-time platform provisioning. A fork with `DEPLOY_TARGET` unset skips both
jobs and stays green.

The sections below cover deploying by hand instead.

### Deploy to Cloud Run

```sh
PROJECT_ID=your-project
REGION=asia-northeast1
IMAGE="$REGION-docker.pkg.dev/$PROJECT_ID/apps/asc-slack-notifier:latest"

gcloud builds submit --tag "$IMAGE"

# Store the secrets in Secret Manager first:
#   printf '%s' 'your-webhook-secret' | gcloud secrets create asc-webhook-secret --data-file=-
#   printf '%s' 'https://hooks.slack.com/...' | gcloud secrets create slack-webhook-url --data-file=-

gcloud run deploy asc-slack-notifier \
  --image "$IMAGE" \
  --region "$REGION" \
  --allow-unauthenticated \
  --set-secrets "ASC_WEBHOOK_SECRET=asc-webhook-secret:latest,SLACK_WEBHOOK_URL=slack-webhook-url:latest"
```

The service must be publicly reachable so App Store Connect can post to it; the
`x-apple-signature` HMAC is what authenticates the request. Cloud Run injects
`PORT` automatically.

### Deploy to AWS Lambda + API Gateway

Build the deployment package for the `provided.al2023` runtime:

```sh
make lambda-zip                    # dist/asc-slack-notifier-lambda-arm64.zip
make lambda-zip LAMBDA_ARCH=amd64  # x86_64 variant
```

The zip contains a single `bootstrap` binary, as `provided.al2023` requires.

```sh
aws lambda create-function \
  --function-name asc-slack-notifier \
  --runtime provided.al2023 \
  --architectures arm64 \
  --handler bootstrap \
  --role arn:aws:iam::123456789012:role/asc-slack-notifier-role \
  --zip-file fileb://dist/asc-slack-notifier-lambda-arm64.zip \
  --environment 'Variables={ASC_WEBHOOK_SECRET=your-webhook-secret,SLACK_WEBHOOK_URL=https://hooks.slack.com/services/T000/B000/XXXX}'

# HTTP API (payload format 2.0) with a $default route to the function
aws apigatewayv2 create-api \
  --name asc-slack-notifier \
  --protocol-type HTTP \
  --target arn:aws:lambda:ap-northeast-1:123456789012:function:asc-slack-notifier

aws lambda add-permission \
  --function-name asc-slack-notifier \
  --statement-id apigateway-invoke \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn 'arn:aws:execute-api:ap-northeast-1:123456789012:<api-id>/*/*'
```

`RUN_MODE` does not need to be set: `AWS_LAMBDA_FUNCTION_NAME` is present in the
Lambda environment, so lambda mode is detected automatically. The API Gateway
HTTP API payload format 2.0 is expected. Base64-encoded bodies are handled, and
the signature is always verified against the decoded raw bytes.

Equivalent AWS SAM template:

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Resources:
  AscSlackNotifier:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: dist/
      Handler: bootstrap
      Runtime: provided.al2023
      Architectures: [arm64]
      MemorySize: 128
      Timeout: 15
      Environment:
        Variables:
          ASC_WEBHOOK_SECRET: '{{resolve:secretsmanager:asc/webhook-secret}}'
          SLACK_WEBHOOK_URL: '{{resolve:secretsmanager:slack/webhook-url}}'
      Events:
        Webhook:
          Type: HttpApi
          Properties:
            Path: /webhook
            Method: POST
```

## Register the webhook with App Store Connect

Create the webhook with the App Store Connect API. `$TOKEN` is a JWT for your
App Store Connect API key and `$APP_ID` is the app's resource id.

The `attributes.secret` below is **your own** value, not something Apple hands
out — generate it once and give the identical string to both sides:

```sh
openssl rand -hex 32   # use the output as ASC_WEBHOOK_SECRET and as attributes.secret
```

```sh
curl -sS -X POST 'https://api.appstoreconnect.apple.com/v1/webhooks' \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "data": {
      "type": "webhooks",
      "attributes": {
        "name": "Slack notifier",
        "url": "https://your-service.example.com/webhook",
        "secret": "your-webhook-secret",
        "enabled": true,
        "eventTypes": [
          "APP_STORE_VERSION_APP_VERSION_STATE_UPDATED",
          "BUILD_UPLOAD_STATE_UPDATED",
          "BUILD_BETA_DETAIL_EXTERNAL_BUILD_STATE_UPDATED",
          "BETA_FEEDBACK_CRASH_SUBMISSION_CREATED",
          "BETA_FEEDBACK_SCREENSHOT_SUBMISSION_CREATED"
        ]
      },
      "relationships": {
        "app": { "data": { "type": "apps", "id": "'"$APP_ID"'" } }
      }
    }
  }'
```

Available event types:

| `WebhookEventType` | Payload `data.type` |
|---|---|
| `ALTERNATIVE_DISTRIBUTION_PACKAGE_AVAILABLE_UPDATED` | `alternativeDistributionPackageAvailableUpdated` |
| `ALTERNATIVE_DISTRIBUTION_PACKAGE_VERSION_CREATED` | `alternativeDistributionPackageVersionCreated` |
| `ALTERNATIVE_DISTRIBUTION_TERRITORY_AVAILABILITY_UPDATED` | `alternativeDistributionTerritoryAvailabilityUpdated` |
| `APP_STORE_VERSION_APP_VERSION_STATE_UPDATED` | `appStoreVersionAppVersionStateUpdated` |
| `BACKGROUND_ASSET_VERSION_APP_STORE_RELEASE_STATE_UPDATED` | `backgroundAssetVersionAppStoreReleaseStateUpdated` |
| `BACKGROUND_ASSET_VERSION_EXTERNAL_BETA_RELEASE_STATE_UPDATED` | `backgroundAssetVersionExternalBetaReleaseStateUpdated` |
| `BACKGROUND_ASSET_VERSION_INTERNAL_BETA_RELEASE_CREATED` | `backgroundAssetVersionInternalBetaReleaseCreated` |
| `BACKGROUND_ASSET_VERSION_STATE_UPDATED` | `backgroundAssetVersionStateUpdated` |
| `BETA_FEEDBACK_CRASH_SUBMISSION_CREATED` | `betaFeedbackCrashSubmissionCreated` |
| `BETA_FEEDBACK_SCREENSHOT_SUBMISSION_CREATED` | `betaFeedbackScreenshotSubmissionCreated` |
| `BUILD_BETA_DETAIL_EXTERNAL_BUILD_STATE_UPDATED` | `buildBetaDetailExternalBuildStateUpdated` |
| `BUILD_UPLOAD_STATE_UPDATED` | `buildUploadStateUpdated` |

## Test with a webhook ping

App Store Connect can send a test delivery to a registered webhook:

```sh
curl -sS -X POST 'https://api.appstoreconnect.apple.com/v1/webhookPings' \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "data": {
      "type": "webhookPings",
      "relationships": {
        "webhook": { "data": { "type": "webhooks", "id": "'"$WEBHOOK_ID"'" } }
      }
    }
  }'
```

The service answers `200` and posts a short "webhook ping received" message to
Slack. Set `NOTIFY_PING=false` to acknowledge pings without notifying Slack.

Recent deliveries and their responses can be inspected through
`GET /v1/webhooks/{id}/deliveries`.

## Development

```sh
make test          # go test ./...
make check         # go vet ./... && go test ./...
make cover         # coverage summary
make build         # bin/asc-slack-notifier
make docker-build  # container image
make lambda-zip    # Lambda deployment package
```

## License

MIT © 2026 Atsushi NAGASE
