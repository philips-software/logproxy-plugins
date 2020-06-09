# logproxy-filter-drop
Example filter that drops messages based on pattern. Configure the pattern using by
passing it via the `FILTER_REGEXP` environment variable.

# building
```
docker build -t logproxy-filter-drop .
```
# deployment
Below is an example `manifest.yml` for deployment to Cloud foundry:

```
---
applications:
- name: logproxy-filter-drop
  docker:
    image: loafoe/logproxy-filter-drop:latest
  env:
    FILTER_REGEXP: "Trigger on this"
    LOGPROXY_QUEUE: channel
    HSDP_LOGINGESTOR_KEY: YourKeyHere
    HSDP_LOGINGESTOR_PRODUCT_KEY: product-key-45c1-here-deadbeafoobar
    HSDP_LOGINGESTOR_SECRET: YourSecretKeyHere
    HSDP_LOGINGESTOR_URL: https://logingestor2-client-test.us-east.philips-healthsuite.com
    TOKEN: SecretT0ken
  routes:
  - route: logproxy-filter-drop.cloud.pcftest.com
  instances: 1
  memory: 128M
```
