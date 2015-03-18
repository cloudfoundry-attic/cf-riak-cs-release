## Acceptance Tests

The acceptance tests are for developers to validate changes to the Riak CS Release.

To run the Riak CS Release Acceptance tests, you will need:
- a running CF instance
- credentials for a CF Admin user
- a deployed Riak CS Release [with the broker registered](/README.md#register_broker) and the plan made public
- a [security group](/README.md#security-groups) granting access to the service for applications

### BOSH errand

BOSH errands were introduced in version 2366 of the BOSH CLI, BOSH Director, and stemcells.

The acceptance-tests deployment job properties must be specified in the deployment manifest. See the [acceptance-test job configuration](http://bosh.io/jobs/acceptance-tests?source=github.com/cloudfoundry/cf-riak-cs-release&version=8) for descriptions and defaults. For the latest auto-generated config docs, navigate from the [cf-riak-cs-release page](http://bosh.io/releases/github.com/cloudfoundry/cf-riak-cs-release) > Version X > acceptance-tests (job).

To run the acceptance tests via bosh errand:

```
$ bosh run errand acceptance-tests
```

### Manually

The acceptance tests can also be run manually. An advantage to doing this is that output will be streamed in real-time (bosh errand output is printed all at once when they finish running).

Instructions:

1. Install **Go** by following the directions found [here](http://golang.org/doc/install)
1. Export the environment variable `$GOPATH` set to the `cf-riak-cs-release` directory to manage Golang dependencies. For more information, see [here](https://github.com/cloudfoundry/cf-riak-cs-release/tree/release-candidate#development).
1. Change to the acceptance-tests directory:

    ```
    $ cd cf-riak-cs-release/src/github.com/cloudfoundry-incubator/cf-riak-cs-acceptance-tests/
    ```

1. Install [Ginkgo](http://onsi.github.io/ginkgo/):

    ```
    $ go get github.com/onsi/ginkgo/ginkgo
    ```

1. Configure the tests by creating `integration_config.json` and setting the environment variable `$CONFIG` to point to it. The following commands provide a shortcut to configuring `integration_config.json` with values for a [bosh-lite](https://github.com/cloudfoundry/bosh-lite)
deployment. Copy and paste this into your terminal, then open the resulting `integration_config.json` in an editor to replace values as appropriate for your environment.

  ```bash
  cat > integration_config.json <<EOF
  {
    "api":                 "api.10.244.0.34.xip.io",
    "admin_user":          "admin",
    "admin_password":      "admin",
    "apps_domain":         "10.244.0.34.xip.io",
    "riak_cs_host":        "p-riakcs.10.244.0.34.xip.io",
    "riak_cs_scheme" :     "https://",
    "service_name":        "p-riakcs",
    "plan_name":           "developer",
    "broker_host":         "p-riakcs-broker.10.244.0.34.xip.io",
    "timeout_scale":       1.0,
    "skip_ssl_validation": true
  }
  EOF
  export CONFIG=$PWD/integration_config.json
  ```

  When `skip_ssl_validation: true`, commands run by the tests will accept self-signed certificates from Cloud Foundry. This option requires v6.0.2 or newer of the CLI.

  All timeouts in the test suite can be scaled proportionally by changing the `timeout_scale` factor.

4. Run  the tests

  ```
  $ ./bin/test
  ```
