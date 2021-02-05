# DKAN New Dataset Tweeter

Sends automated tweets containing newly created datasets in the [open data portal of the city of Münster](https://opendata.stadt-muenster.de/).

## Operation

Operation is powered by GitHub Actions & a S3 compatible object storage.

## config.json guide

Good to know

### twitter

Configuration values for communication with the Twitter API. You have to register an app at developer.twitter.com

After you have registered an app at developer.twitter.com, obtain values for `consumerKey`, `consumerSecret`, `accessToken`, `accessSecret` and set them in your `config.json`.

### s3

Configuration values of the S3 client. The S3 bucket will be used for storing json state files.

### endpoint

S3 endpoint URL. For example `https://s3.fr-par.scw.cloud` (Scaleway) or `https://s3.amazonaws.com` (AWS)

### region

S3 region. For example `fr-par` (Scaleway) or `eu-central-1` (AWS)

### bucket

The name of the bucket to be used for storage of json files.

### path

The path to the **directory** used for storage of json files. Should not end with a `/` slash.

### accessKeyId & secretAccessKey

Secret credentials to be used by the S3 client.
