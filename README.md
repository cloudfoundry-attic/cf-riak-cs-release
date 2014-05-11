# cf-riak-cs-release

A BOSH release for Riak and Riak CS.

Make sure to run `./update` (to update git submodules) before creating the release.

This project is based on [BrianMMcClain/riak-release](https://github.com/BrianMMcClain/riak-release).

## Deployment

1. [Upload a release to the BOSH director](#upload_release)
1.  Upload the appropriate stemcell for your deployment (warden, vsphere, or aws), if it has not already been uploaded.
1. [Create a deployment manifest and deploy, following environment-specific instructions below.](#create_manifest)
  - [BOSH-lite](#bosh-lite)
  - [vSphere](#vsphere)
  - [AWS](#aws)
1. [Register the service broker with Cloud Foundry](#register_broker)

### Upload a Release<a name="upload_release"></a>

You can use a pre-built final release or build a release from HEAD. Final releases contain pre-compiled packages, making deployment much faster. However, these are created manually and infrequently. To be sure you're deploying the latest code, build a release yourself.

#### Upload a pre-built final BOSH release

1. Check out the tag for the desired version. This is necessary for generating a manifest that matches the code you're deploying.

  ```bash
  cd ~/workspace/cf-riak-cs-release
  ./update
  git checkout v1
  ```

1. Run the upload command, referencing one of the config files in the `releases` directory.

  ```bash
  bosh upload release releases/cf-riak-cs-1.yml
  ```

#### Create a BOSH Release from HEAD and Upload:

1. Build a BOSH development release from HEAD

  ```bash
  cd ~/workspace/cf-riak-cs-release
  ./update
  bosh create release
  ```

  When prompted to name the release, call it `cf-riak-cs`.

1. Upload the release to your bosh environment:

  ```bash
  bosh upload release
  ```

### Create a Manifest and Deploy<a name="create_manifest"></a>

#### BOSH-lite<a name="bosh-lite"></a>

1. Run the script [`bosh-lite/make_manifest`](bosh-lite/make_manifest) to generate your manifest for bosh-lite. This script uses a stub provided for you in `bosh-lite/stub.yml`. For a description of the parameters in the stub, see <a href="#manifest-stub-parameters">Manifest Stub Parameters</a> below.

    ```
    $ ./bosh-lite/make_manifest
    ```
    The manifest will be written to `bosh-lite/manifests/cf-riak-cs-manifest.yml`, which can be modified to change deployment settings. 

1. The `make_manifest` script will set the deployment to `bosh-lite/manifests/cf-riak-cs-manifest.yml` for you, so to deploy you only need to run `bosh deploy`.

#### vSphere<a name="vsphere"></a>

1. Create a stub file called `cf-riak-cs-vsphere-stub.yml` that contains the properties in the example below. For a description of these parameters, see <a href="#manifest-stub-parameters">Manifest Stub Parameters</a> below. 
  
    This stub differs from the bosh-lite stub in that it requires:

    * Username and password for admin user to support errands
    * Network settings, with 6 static IPs and 6+ dynamic IPs

  ```
  director_uuid: YOUR-DIRECTOR-GUID-HERE
  networks:
  - name: riak-cs-network
    subnets:
    - cloud_properties:
        name: VM Network  # name of vsphere network
      dns:
      - 8.8.8.8
      gateway: 10.0.0.1
      range: 10.0.0.0/24
      reserved:           # IPs that bosh should not use inside your subnet range
      - 10.0.0.2-10.0.0.99
      - 10.0.0.115-10.0.0.254
      static:
      - 10.0.0.100
      - 10.0.0.101
      - 10.0.0.102
      - 10.0.0.103
      - 10.0.0.104
      - 10.0.0.105
  properties:
    domain: your-cf-system-domain-here
    nats:
      machines:
      - 10.0.0.15   # IP of nats server
      user: nats-username-here
      password: nats-password-here
      port: 4222
    cf:
      api_url: https://api.YOUR-CF-DOMAIN-HERE
      apps_domain: YOUR-APP-DOMAIN-HERE
      admin_username: CF-ADMIN-USERNAME
      admin_password: CF-ADMIN-PASSWORD
  ```

2. Generate the manifest: `./generate_deployment_manifest vsphere cf-riak-cs-vsphere-stub.yml > cf-riak-cs-vsphere.yml`
To tweak the deployment settings, you can modify the resulting file `cf-riak-cs-vsphere.yml`.

3. To deploy: `bosh deployment cf-riak-cs-vsphere.yml && bosh deploy`

#### AWS<a name="aws"></a>

1. Create a stub file called `cf-riak-cs-aws-stub.yml` that contains the parameters in the example below. For a description of these parameters, see <a href="#manifest-stub-parameters">Manifest Stub Parameters</a> below.

    This stub differs from the bosh-lite stub in that it requires:

    * Username and password for admin user to support errands
    * Network and resource pool settings

  ```  
  director_uuid: YOUR-DIRECTOR-GUID-HERE
  networks:
  - name: riak-cs-network
    subnets:
    - name: riak-cs-subnet
      cloud_properties:
        subnet: YOUR-AWS-SERVICES-SUBNET-ID-HERE
  resource_pools:
  - name: riak-pool
    cloud_properties:
      availability_zone: YOUR-PRIMARY-AZ-NAME-HERE
  - name: broker-pool
    cloud_properties:
      availability_zone: YOUR-PRIMARY-AZ-NAME-AGAIN
  properties:
    domain: your-cf-system-domain-here
    nats:
      machines:
      - IP-OF-NATS-SERVER-HERE
      user: NATS-USERNAME-HERE
      password: NATS-PASSWORD-HERE
      port: 4222
    cf:
      api_url: https://api.YOUR-CF-SYSTEM-DOMAIN-HERE
      apps_domain: YOUR-APP-DOMAIN-HERE
      admin_username: CF-ADMIN-USERNAME
      admin_password: CF-ADMIN-PASSWORD
  ```

1. Generate the manifest: `./generate_deployment_manifest aws cf-riak-cs-aws-stub.yml > cf-riak-cs-aws.yml`
To tweak the deployment settings, you can modify the resulting file `cf-riak-cs-aws.yml`.

1. To deploy: `bosh deployment cf-riak-cs-aws.yml && bosh deploy`

### Manifest Stub Parameters

This section describes the parameters that must be added to manifest stub for the supported environments listed above.

* `director_uuid`: can be found from running `bosh status`

* `properties`
  * `domain`: refers to the system domain that you installed CF against (it should match the domain property from the CF bosh manifest). The value is used to determine both the route advertised by each node in the cluster (see `register_route` below), as well as the route for the broker.

  * `riak_cs`: These properties control behavior of the Riak CS cluster nodes. As these properties have defaults, it is not necessary to include them in your stub unless you need to change them.
    * `ssl_enabled` defaults to true 
    * `skip_ssl_validation` defaults to false, which assumes you have valid certs in your CF deployment
    * `register_route`: defaults to true. Determines whether each node in the cluster advertises a route. When set to true, all heathly nodes in the cluster can be reached at `riakcs.DOMAIN` (where DOMAIN is the value of the `domain` property above). Having a single route to all healthy nodes allows traffic to be load balanced across the Riak CS nodes. A healthcheck process on each node monitors whether riak and riak-cs are running and the node is a valid member of the cluster. If the healthcheck process determines that a node is not healthy, it will unregister the route for the unhealthy node. 

      When this property is set to false, nodes will not register a route. This is useful when deploying `cf-riak-cs-release` without Cloud Foundry. NOTE: the Riak CS service broker does not yet support `register_route: false`. __When setting `register_route` to false, you must set the instance count of the `cf-riak-cs-broker`, `acceptance-tests`, `broker-registrar`, and `broker-deregistrar` jobs to 0. Also you should omit the `domain` property and all the `cf` properties below.__

  * `cf`: These properties provide information the Riak CS service needs to know about your Cloud Foundry deployment.
    * `api_url`: the CloudController API URL (same thing you use to target using the `cf` CLI).  It's used by a BOSH errand to register the newly deployed broker with CloudController (see below for invocation).
    * `admin_username`: a CloudFoundry admin username. It's used by a BOSH errand to register the newly deployed broker with CloudController (see below for invocation).
    * `admin_password`: a CloudFoundry admin password. It's used by a BOSH errand to register the newly deployed broker with CloudController (see below for invocation).
    * `apps_domain`: the CloudFoundry App Domain. It's used by a BOSH errand to run acceptance tests for this release (see below for invocation).

## Register the Service Broker<a name="register_broker"></a>

### Using BOSH errands

BOSH errands were introduced in version 2366 of the BOSH CLI, BOSH Director, and stemcells. 

        bosh run errand broker-registrar
        
Note: the broker-registrar errand will fail if the broker has already been registered, and the broker name does not match the manifest property `broker.name`. Use the `cf rename-service-broker` CLI command to change the broker name to match the manifest property then this errand will succeed. 

### Manually

1. First register the broker using the `cf` CLI.  You must be logged in as an admin.

  ```
  $ cf create-service-broker p-riakcs BROKER_USERNAME BROKER_PASSWORD URL
  ```
    
  `BROKER_USERNAME` and `BROKER_PASSWORD` are the credentials Cloud Foundry will use to authenticate when making API calls to the service broker. Use the values for manifest properties `properties.broker.username` and `properties.broker.password`. 
  
  `URL` specifies where the Cloud Controller will access the MySQL broker. Use the value of the manifest property `properties.broker.host`.

  For more information, see [Managing Service Brokers](http://docs.cloudfoundry.org/services/managing-service-brokers.html).

1. Then [make the service plan public](http://docs.cloudfoundry.org/services/services/managing-service-brokers.html#make-plans-public).


## Running Acceptance Tests

To run the Riak CS Release Acceptance tests, you will need:
- a running CF instance
- credentials for a CF Admin user
- a deployed Riak CS Release with the broker registered and the plan made public
- an environment variable `$CONFIG` which points to a `.json` file that contains the application domain

### Using BOSH errands

BOSH errands were introduced in version 2366 of the BOSH CLI, BOSH Director, and stemcells. 

The following properties must be included in the manifest (most will be there by default):
- cf.api_url:
- cf.admin_username:
- cf.admin_password:
- cf.apps_domain:
- cf.skip_ssl_validation:
- riak_cs.ssl_enabled:
- broker.host:
- external_riakcs_host:

```
bosh run errand acceptance-tests
```

### Manually

1. Install `go` by following the directions found [here](http://golang.org/doc/install)
2. `cd` into `cf-riak-cs-release/test/acceptance-tests/`
3. Update `cf-riak-cs-release/test/acceptance-tests/integration_config.json`

   The following script will configure these prerequisites for a [bosh-lite](https://github.com/cloudfoundry/bosh-lite)
installation. Replace credentials and URLs as appropriate for your environment.

    ```bash
    #! /bin/bash
    
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
      "skip_ssl_validation": true
    }
    EOF
    export CONFIG=$PWD/integration_config.json
    ```

    Note: `skip_ssl_validation` requires CLI v6.0.2 or newer.

4. Run  the tests

    ```bash
    ./bin/test
    ```

## De-registering the broker

The following commands are destructive and are intended to be run in conjuction with deleting your BOSH deployment.

### Using BOSH errands

BOSH errands were introduced in version 2366 of the BOSH CLI, BOSH Director, and stemcells. 

This errand runs the two commands listed in the manual section below from a BOSH-deployed VM. This errand should be run before deleting your BOSH deployment. If you have already deleted your deployment follow the manual instructions below.

    bosh run errand broker-deregistrar

#### Manually

Run the following:

```bash
cf purge-service-offering p-riakcs
cf delete-service-broker p-riakcs
```

## Caveats

We have not tested changing the structure of a live cluster, e.g. changing the seed node.

## Blobs

See [Bosh Blobstore](http://docs.cloudfoundry.com/docs/running/bosh/components/blobstore.html) for blobstore configuration.

To update a blob:

1. Remove its entry from `config/blobs.yml`
2. Remove the cached blob from `.blobs/` (you can find it by checking the symlink in `blobs/<package>/`)
3. Copy the new blob file into `blobs/<package>/`
4. Upload the new blob: `bosh upload blobs`

### riak

Clone the [riak repository](https://github.com/basho/riak), check out the desired tag, and `make dist`.
The resulting `tar.gz` file can be found in the working directory.

### riak-cs

Clone the [riak_cs repository](https://github.com/basho/riak_cs), check out the desired tag, and `make package.src`.
The resulting `tar.gz` file can be found in the `package/` directory.

### stanchion

Clone the [stanchion repository](https://github.com/basho/stanchion), check out the desired tag, and `make package.src`.
The resulting `tar.gz` file can be found in the `package/` directory.


### other

TODO - verify where the `git`, and `erlang` tarfiles came from.

## TODO

- The settings for the Riak job in this release are configured with options suggested by Basho for deploying Riak in a Riak CS cluster.  We could add an option to configure Riak for standalone operation (when a manifest includes only Riak but not Riak CS)

[BOSH lite]: https://github.com/cloudfoundry/bosh-lite

