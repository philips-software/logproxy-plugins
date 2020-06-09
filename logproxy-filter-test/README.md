# logproxy-filter-test
Simple plugin that prints the ID field and then changes it to "42".

# building
```
docker build -t logproxy-filter-test .
```
# deployment
Below is an example `manifest.yml` for deployment to Cloud foundry:

```
---
applications:
- name: logproxy-filter-test
  docker:
    image: loafoe/logproxy-filter-test:latest
  env:
    LOGPROXY_QUEUE: channel
    HSDP_LOGINGESTOR_KEY: YourKeyHere
    HSDP_LOGINGESTOR_PRODUCT_KEY: product-key-45c1-here-deadbeafoobar
    HSDP_LOGINGESTOR_SECRET: YourSecretKeyHere
    HSDP_LOGINGESTOR_URL: https://logingestor2-client-test.us-east.philips-healthsuite.com
    TOKEN: SecretT0ken
  routes:
  - route: logproxy-filter-test.cloud.pcftest.com
  instances: 1
  memory: 128M
```
