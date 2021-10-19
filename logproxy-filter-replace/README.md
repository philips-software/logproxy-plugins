# logproxy-filter-replace

Example filter that replaces texts in messages based on patterns. Configure the patterns by
passing them via the `FILTER_CONFIG` environment variable. The content of the variable must be a base64 encoded json.

The json should look like the example below and be encoded to base64 format.

```json
[
  {
    "pattern": "(([A-Z])\\w+)",
    "replace": "<FooBar>"
  }
]
```

Note: All invalid regular expressions will be skipped.

The filter will look for all the patterns in this list and replace by the corresponding `replace` string.

This can be useful for obfuscating sensitive information in the log messages.

# building

```
docker build -t logproxy-filter-replace .
```

# deployment

Below is an example `manifest.yml` for deployment to Cloud foundry:

```
---
applications:
- name: logproxy-filter-replace
  docker:
    image: philipssoftware/logproxy-filter-replace:latest
  env:
    FILTER_CONFIG: "Ww0KICB7DQogICAgInBhdHRlcm4iOiAiKChbQS1aXSlcXHcrKSIsDQogICAgInJlcGxhY2UiOiAiPEZvb0Jhcj4iDQogIH0NCl0NCg=="
    LOGPROXY_QUEUE: channel
    HSDP_LOGINGESTOR_KEY: YourKeyHere
    HSDP_LOGINGESTOR_PRODUCT_KEY: product-key-45c1-here-deadbeafoobar
    HSDP_LOGINGESTOR_SECRET: YourSecretKeyHere
    HSDP_LOGINGESTOR_URL: https://logingestor2-client-test.us-east.philips-healthsuite.com
    TOKEN: SecretT0ken
  routes:
  - route: logproxy-filter-replace.cloud.pcftest.com
  instances: 1
  memory: 128M
```
