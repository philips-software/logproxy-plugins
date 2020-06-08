# logproxy-filter-email
Example filter that triggers an email

# dependencies
An SMTP service binding is required

# environment
| Variable | Description | Required |
|----------|-------------|----------|
| EMAIL\_TO | Send email to this address | Yes |
| EMAIL\_FROM | Sender email address | Yes |
| EMAIL\_SUBJECT | Subject of the email | Yes |
| FILTER\_REGEXP | Regular expression to trigger on | Yes |

# building
```
docker build -t logproxy-filter-email .
```
# deployment
Below is an example `manifest.yml` for deployment to Cloud foundry:

```
---
applications:
- name: logproxy-filter-email
  docker:
    image: loafoe/logproxy-filter-email:latest
  env:
    EMAIL_TO: "your_inbox@philips.com"
    EMAIL_FROM: "no_reply@philips.com"
    EMAIL_SUBJECT: "Alert email"
    FILTER_REGEXP: "Trigger on this"
    LOGPROXY_QUEUE: channel
    LOGPROXY_DELIVERY: none
    TOKEN: SecretT0ken
  routes:
  - route: logproxy-filter-email.cloud.pcftest.com
  instances: 1
  memory: 128M
  services:
  - smtp
```

