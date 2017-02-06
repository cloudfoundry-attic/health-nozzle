# Introduction

These provides cf pushable app, which read the ingress, egress, and dropped
metrics from the firehose for metron, doppler, and traffic controller, so they
can be read via JSON api endpoint.

# Usage

```sh
git clone https://github.com/cloudfoundry-incubator/health-nozzle
cf push health-nozzle --no-start
cf set-env health-nozzle DOPPLER_ADDR "wss://doppler.system-domain"
cf set-env health-nozzle CF_ACCESS_TOKEN "$(cf oauth-token)"
cf start health-nozzle
curl http://health-nozzle.system-domain
```
