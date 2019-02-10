# Slack Bot of Google Cloud Functions for Go

# Deploy

## HelloCommand

* `$VERIFICATION_TOKEN` : Set the Verification Token of Slack.

```bash
$ gcloud functions deploy hello \
--entry-point HelloCommand \
--runtime go111 \
--set-env-vars VERIFICATION_TOKEN=$VERIFICATION_TOKEN \
--trigger-http
```

## WeatherCommand

* `$VERIFICATION_TOKEN` : Set the Verification Token of Slack.
* `$APP_ID` : Set API key of OpenWeatherMap.

```bash
$ gcloud functions deploy weather \
--entry-point WeatherCommand \
--runtime go111 \
--set-env-vars VERIFICATION_TOKEN=$VERIFICATION_TOKEN,APP_ID=$APP_ID \
--trigger-http
```
