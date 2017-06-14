# BOSH Compile Release Cache

A go webserver that caches compiled releases to S3.

### Requirements:

- [BOSH V2 CLI](http://bosh.io/docs/cli-v2.html#install) in the path as `bosh-cli`
- The following environment variables:
```
AWS_DEFAULT_REGION
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
BOSH_CLIENT
BOSH_CLIENT_SECRET
BOSH_HOST
BOSH_CA_CERT
S3_BUCKET
```