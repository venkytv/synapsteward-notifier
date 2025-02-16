# SynapSteward Notifier Documentation

## Overview

The `synapsteward-notifier` program listens to a NATS stream and sends notifications to a [Pushover](https://pushover.net/api) user based on the received messages.

## Usage

### Basic Command

```bash
export PUSHOVER_API_TOKEN=<API_TOKEN>
export PUSHOVER_USER_KEY=<USER_KEY>
synapsteward-notifier
```

### Command Line Options

- `--nats-url`: Specifies the NATS server URL.
  - **Default**: `nats://localhost:4222`
  - **Environment Variable**: `NATS_URL`

- `--nats-stream`: The name of the NATS stream to listen for messages.
  - **Default**: `notifications`
  - **Environment Variable**: `NATS_STREAM`

- `--nats-consumer`: The name of the NATS consumer to create.
  - **Default**: `synapsteward-notifier`
  - **Environment Variable**: `NATS_CONSUMER`

- `--pushover-api-token`: (Required) The API token for the Pushover service.
  - **Environment Variable**: `PUSHOVER_API_TOKEN`

- `--pushover-user-key`: (Required) The user key for the Pushover user.
  - **Environment Variable**: `PUSHOVER_USER_KEY`

- `--debug`: Enable debug logging for troubleshooting purposes.
  - **Environment Variable**: `DEBUG`

### Environment Variables

You can also configure the application using environment variables:

- `NATS_URL`
- `NATS_STREAM`
- `NATS_CONSUMER`
- `PUSHOVER_API_TOKEN`
- `PUSHOVER_USER_KEY`
- `DEBUG`

## Example

Running the notifier with a Pushover API token and user key:

```bash
synapsteward-notifier --pushover-api-token myPushoverApiToken --pushover-user-key myPushoverUserKey
```

Enable debug mode to get more logging information:

```bash
synapsteward-notifier --pushover-api-token myPushoverApiToken --pushover-user-key myPushoverUserKey --debug
```
