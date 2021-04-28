# logproxy-filter-replace
Example filter that replaces text in messages based on pattern. Configure the pattern by
passing it via the `FILTER_REGEXP` environment variable and the string to replace by passing it via the `REPLACE_STRING` variable.

This plugin can be useful for obfuscating sensitive information in the log messages.

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
    image: jdelucaa/logproxy-filter-replace:latest
  env:
    FILTER_REGEXP: "Find this"
    REPLACE_STRING: "Replace with this"
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
