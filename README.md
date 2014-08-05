# cf-riak-cs-release

A BOSH release for Riak and Riak CS.

This project is based on [BrianMMcClain/riak-release](https://github.com/BrianMMcClain/riak-release).

## Deployment

### Prerequisites

- A deployment of [BOSH](https://github.com/cloudfoundry/bosh)
- A deployment of [Cloud Foundry](https://github.com/cloudfoundry/cf-release)
- Instructions for installing BOSH and Cloud Foundry can be found at http://docs.cloudfoundry.org/.

### Overview

1. Upload a supported stemcell
1. [Upload a release to the BOSH director](#upload_release)
1. [Create a Deployment Manifest and Deploy](#create_manifest)
  - [BOSH-lite](#bosh-lite)
  - [vSphere](#vsphere)
  - [AWS](#aws)
  - [Deployment Manifest Stub Properties](#stub-properties)
1. [Register the Service Broker with Cloud Foundry](#register_broker)

### Upload a Release<a name="upload_release"></a>

You can use a pre-built final release or build a release from HEAD. Final releases contain pre-compiled packages, making deployment much faster. However, these are created manually and infrequently. To be sure you're deploying the latest code, build a release yourself.

#### Upload a pre-built final BOSH release

1. Check out the tag for the desired version. This is necessary for generating a manifest that matches the code you're deploying.

  ```
  $ cd ~/workspace/cf-riak-cs-release
  $ ./update
  $ git checkout v4
  $ git submodule update --recursive
  ```

1. Run the upload command, referencing one of the config files in the `releases` directory.

  ```
  $ bosh upload release releases/cf-riak-cs-4.yml
  ```

#### Create a BOSH Release from HEAD and Upload:

1. Build a BOSH development release from HEAD

  ```
  $ cd ~/workspace/cf-riak-cs-release
  $ ./update
  $ bosh create release
  ```

  When prompted to name the release, call it `cf-riak-cs`.

1. Upload the release to your bosh environment:

  ```
  $ bosh upload release
  ```

### Create a Manifest and Deploy<a name="create_manifest"></a>

#### BOSH-lite<a name="bosh-lite"></a>

1. Run the script [`bosh-lite/make_manifest`](bosh-lite/make_manifest) to generate your manifest for bosh-lite. This script uses a stub provided for you, `bosh-lite/stub.yml`. For a description of the parameters in this stub, see <a href="#manifest-stub-parameters">Manifest Stub Parameters</a> below.

    ```
    $ ./bosh-lite/make_manifest
    ```
    The manifest will be written to `bosh-lite/manifests/cf-riak-cs-manifest.yml`, which can be modified to change deployment settings. 

1. The `make_manifest` script will set the deployment to `bosh-lite/manifests/cf-riak-cs-manifest.yml` for you, so to deploy you only need to run `bosh deploy`.

#### vSphere<a name="vsphere"></a>

1. Create a stub file called `cf-riak-cs-vsphere-stub.yml` that contains the properties in the example below. For a description of these and other manifest properties, see <a href="#manifest-stub-parameters">Manifest Stub Parameters</a> below. 
  
    This stub differs from the bosh-lite stub in that it requires:

    * Username and password for admin user to support errands
    * Network settings, with 6 static IPs and 6+ dynamic IPs

  ```
  director_uuid: YOUR-DIRECTOR-GUID
  networks:
  - name: riak-cs-network
    subnets:
    - cloud_properties:
        name: YOUR-VSPHERE-NETWORK-NAME
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
    domain: YOUR-CF-SYSTEM-DOMAIN
    nats:
      machines:
      - 10.0.0.15   # IP of nats server
      user: NATS-USERNAME
      password: NATS-PASSWORD
      port: 4222
    cf:
      api_url: https://api.YOUR-CF-SYSTEM-DOMAIN
      apps_domain: YOUR-CF-APP-DOMAIN
      admin_username: CF-ADMIN-USERNAME
      admin_password: CF-ADMIN-PASSWORD
  ```

2. Generate the manifest: `./generate_deployment_manifest vsphere cf-riak-cs-vsphere-stub.yml > cf-riak-cs-vsphere.yml`
To tweak the deployment settings, you can modify the resulting file `cf-riak-cs-vsphere.yml`.

3. To deploy: `bosh deployment cf-riak-cs-vsphere.yml && bosh deploy`

#### AWS<a name="aws"></a>

1. Create a stub file called `cf-riak-cs-aws-stub.yml` that contains the parameters in the example below. For a description of these and other manifest properties, see <a href="#manifest-stub-parameters">Manifest Stub Parameters</a> below.

    This stub differs from the bosh-lite stub in that it requires:

    * Username and password for admin user to support errands
    * Network and resource pool settings

  ```  
  director_uuid: YOUR-DIRECTOR-GUID
  networks:
  - name: riak-cs-network
    subnets:
    - name: riak-cs-subnet
      cloud_properties:
        subnet: YOUR-AWS-SERVICES-SUBNET-ID
  resource_pools:
  - name: riak-pool
    cloud_properties:
      availability_zone: YOUR-PRIMARY-AZ-NAME
  - name: broker-pool
    cloud_properties:
      availability_zone: YOUR-PRIMARY-AZ-NAME
  properties:
    domain: YOUR-CF-SYSTEM-DOMAIN
    nats:
      machines:
      - IP-OF-NATS-SERVER
      user: NATS-USERNAME
      password: NATS-PASSWORD
      port: 4222
    cf:
      api_url: https://api.YOUR-CF-SYSTEM-DOMAIN
      apps_domain: YOUR-CF-APP-DOMAIN
      admin_username: CF-ADMIN-USERNAME
      admin_password: CF-ADMIN-PASSWORD
  ```

1. Generate the manifest: `./generate_deployment_manifest aws cf-riak-cs-aws-stub.yml > cf-riak-cs-aws.yml`
To tweak the deployment settings, you can modify the resulting file `cf-riak-cs-aws.yml`.

1. To deploy: `bosh deployment cf-riak-cs-aws.yml && bosh deploy`

#### Deployment Manifest Stub Parameters<a name="stub-properties"></a>

This section describes the parameters that must be added to manifest stub for the supported environments listed above.

* `director_uuid`: can be found from running `bosh status`

* `properties`
  * `domain`: refers to the system domain that you installed CF against (it should match the domain property from the CF bosh manifest). The value is used to determine both the route advertised by each node in the cluster (see `register_route` below), as well as the route for the broker.

  * `broker`: These properties configure aspects of the broker
    * `name`: the name of the broker to use when registering with a BOSH errand.
    * `username`: username for the broker used for basic auth.
    * `password`: password for the broker used for basic auth.

  * `riak_cs`: These properties control behavior of the Riak CS cluster nodes. As these properties have defaults, it is not necessary to include them in your stub unless you need to change them.
    * `admin_key`: The admin user key for the Riak CS cluster.
    * `admin_secret`: The admin user secret for the Riak CS cluster.
    * `ssl_enabled`: Determines the scheme used by the broker to communicate with riak-cs and the scheme returned in the binding credentials. Defaults to true (`https`).
    * `skip_ssl_validation`: Determines whether or not the service broker should accept self-signed SSL certs from the Riak cluster. Defaults to false.
    * `register_route`: defaults to true. Determines whether each node in the cluster advertises a route. When set to true, all heathly nodes in the cluster can be reached at `riakcs.DOMAIN` (where DOMAIN is the value of the `domain` property above). Having a single route to all healthy nodes allows traffic to be load balanced across the Riak CS nodes. A healthcheck process on each node monitors whether riak and riak-cs are running and the node is a valid member of the cluster. If the healthcheck process determines that a node is not healthy, it will unregister the route for the unhealthy node. 

      When this property is set to false, nodes will not register a route. This is useful when deploying `cf-riak-cs-release` without Cloud Foundry. NOTE: the Riak CS service broker does not yet support `register_route: false`. __When setting `register_route` to false, you must set the instance count of the `cf-riak-cs-broker`, `acceptance-tests`, `broker-registrar`, and `broker-deregistrar` jobs to 0. Also you should omit the `domain` property and all the `cf` properties below.__

  * `cf`: These properties provide information the Riak CS service needs to know about your Cloud Foundry deployment.
    * `api_url`: the CloudController API URL (same thing you use to target using the `cf` CLI).  It's used by a BOSH errand to register the newly deployed broker with CloudController (see below for invocation).
    * `admin_username`: a CloudFoundry admin username. It's used by a BOSH errand to register the newly deployed broker with CloudController (see below for invocation).
    * `admin_password`: a CloudFoundry admin password. It's used by a BOSH errand to register the newly deployed broker with CloudController (see below for invocation).
    * `apps_domain`: the CloudFoundry App Domain. It's used by a BOSH errand to run acceptance tests for this release (see below for invocation).
    * `skip_ssl_validation`: Determines whether or not the service broker should accept self-signed SSL certs from Cloud Foundry when running BOSH errands. Defaults to false.

  * `syslog_aggregator`:
    * `address`: IP address for syslog aggregator
    * `port`: TCP port of syslog aggregator
    * `transport`: Transport to be used when forwarding logs. Valid values are tcp, udp, or relp. Defaults to tcp.
    * `all`: Determines whether syslog data for all processes should be forwarded or only configured jobs. Defaults to false.

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

## Using s3curl to read and write to your Riak CS bucket

Clone s3curl from github:

`git clone https://github.com/rtdp/s3curl`

Add credentials to `~/.s3curl`:
```
%awsSecretAccessKeys = (
    myuser => {
        id => 'my-access-key-id',
        key => 'my-secret-access-key'
    }
);
```

Edit `s3curl.pl` to add `p-riakcs.mydomain` to the known endpoints:

```
...
my @endpoints = ( 's3.amazonaws.com',
                  's3-us-west-1.amazonaws.com',
                  's3-us-west-2.amazonaws.com',
                  's3-us-gov-west-1.amazonaws.com',
                  's3-eu-west-1.amazonaws.com',
                  's3-ap-southeast-1.amazonaws.com',
                  's3-ap-northeast-1.amazonaws.com',
                  's3-sa-east-1.amazonaws.com',
                  'p-riakcs.mydomain');
...
```
*Note: If you never intend on communicating with any of the amazon services, then you can delete the existing entries (the ones beginning with 's3').*

To list bucket contents at service-instance-location:

`./s3curl.pl --id myuser -- http://p-riakcs.mydomain/service-instance-id`

To put contents file to bucket with key `mykey`:

`./s3curl.pl --id myuser --put filename -- http://p-riakcs.mydomain/service-instance-id/mykey`

*Note: curl requires you to escape any special characters in filenames - e.g. filename\\.txt*

To get file with key `mykey` from bucket:

`./s3curl.pl --id myuser -- http://p-riakcs.mydomain/service-instance-id/mykey`

## Limitations

We have not tested changing the structure of a live cluster, e.g. changing the seed node.

This release does not support running Riak without Riak CS.

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

## other

TODO - verify where the `git`, and `erlang` tarfiles came from.

[BOSH lite]: https://github.com/cloudfoundry/bosh-lite

