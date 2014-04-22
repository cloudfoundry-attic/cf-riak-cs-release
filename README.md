# cf-riak-cs-release

A BOSH release for Riak and Riak CS.

Make sure to run `./update` (to update git submodules) before creating the release.

This project is based on [BrianMMcClain/riak-release](https://github.com/BrianMMcClain/riak-release).

## Deployment

1.  First create the release, naming it `cf-riak-cs`.
1.  Then upload the release.
1.  Make sure you have uploaded the appropriate stemcell for your deployment (either vsphere or warden)
1.  Create a deployment manifest and deploy, following environment-specific instructions below.

### BOSH-lite environment

1. Create a stub file called `riak-cs-lite-stub.yml` that contains the following configuration parameters.

	```
	director_uuid: YOUR-DIRECTOR-GUID-HERE
	properties:
  	  riak_cs:
	    ssl_enabled: YOUR-CHOICE-HERE #true or false
	    skip_ssl_validation: YOUR-CHOICE-HERE #true or false
	  domain: YOUR-CF-SYSTEM-DOMAIN # such as 10.244.0.34.xip.io for bosh-lite
	  cf:
	    api_url: http://api.YOUR-CF-DOMAIN-HERE # such as http://api.10.244.0.34.xip.io
	```
	
	* Director uuid can be found from running `bosh status`

	* SSL Properties:
		* There are two properties under properties.riak-cs called `ssl_enabled` and `skip_ssl_validation`
		* `ssl_enabled` defaults to true and `skip_ssl_validation` defaults to false, which assumes you have valid certs in your CF deployment
		* If you wish to change either of these put them in a stub file and configure them as needed:

	* Cloud Foundry Properties:

		This release needs to know a little about your CF installation.  The `domain` property refers to the system domain that you installed CF against (it should match the domain property from the CF bosh manifest), and it's used to publish a route for the cluster (e.g.`riakcs.YOUR-CF-SYSTEM-DOMAIN`) and a route for the broker.  The route for the cluster allows traffic to be load balanced across the riak CS nodes.

		The `cf.api_url` parameter refers to the CloudController API URL (same thing you use to target using the `cf` CLI).  It's used by an BOSH errand to register the newly deployed broker with CloudController (see below for invocation).  `cf.admin_username` and `cf.admin_password` are also needed by the BOSH errand to register the broker, but are not required for bosh-lite since the credentials are admin/admin.
	
2. Generate the manifest: `./generate_deployment_manifest warden riak-cs-lite-stub.yml > riak-cs-lite.yml`
To tweak the deployment settings, you can modify the resulting file `riak-cs-lite.yml`.
3. To deploy: `bosh deployment riak-cs-lite.yml && bosh deploy`

### vSphere environment

1. Create a stub file called `riak-cs-vsphere-stub.yml` that contains the parameters described for bosh-lite above. In addition, you must include:

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
	     admin_username: CF-ADMIN-USERNAME
	     admin_password: CF-ADMIN-PASSWORD
	```

2. Generate the manifest: `./generate_deployment_manifest vsphere riak-cs-vsphere-stub.yml > riak-cs-vsphere.yml`
To tweak the deployment settings, you can modify the resulting file `riak-cs-vsphere.yml`.

3. To deploy: `bosh deployment riak-cs-vsphere.yml && bosh deploy`

### AWS environment

1. Create a stub file called `riak-cs-aws-stub.yml` that contains the parameters described for bosh-lite above. In addition, you must include:

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
	  riak_cs:
	    ssl_enabled: YOUR-CHOICE-HERE #true or false
	    skip_ssl_validation: YOUR-CHOICE-HERE #true or false
          domain: your-cf-system-domain-here
          nats:
            machines:
            - IP-OF-NATS-SERVER-HERE
            user: NATS-USERNAME-HERE
            password: NATS-PASSWORD-HERE
            port: 4222
          cf:
            api_url: https://api.YOUR-CF-SYSTEM-DOMAIN-HERE
            admin_username: CF-ADMIN-USERNAME
            admin_password: CF-ADMIN-PASSWORD
	```

1. Generate the manifest: `./generate_deployment_manifest aws riak-cs-aws-stub.yml > riak-cs-aws.yml`
To tweak the deployment settings, you can modify the resulting file `riak-cs-aws.yml`.

1. To deploy: `bosh deployment riak-cs-aws.yml && bosh deploy`

## Registering the broker

### Using BOSH errands

If you're using a new enough BOSH director, stemcell, and CLI to support errands, run the following errand:

        bosh run errand broker-registrar

### Manually
First register the broker using the `cf` CLI.  You have to be logged in as an admin, and the IP of the broker will likely be different on vsphere (use `bosh vms` to find it if necessary)
```
cf create-service-broker riakcs admin admin http://10.244.3.22:8080
```
Then make the [service plan public](http://docs.cloudfoundry.org/services/services/managing-service-brokers.html#make-plans-public).

## De-registering the broker

### Using BOSH errands

If you're using a new enough BOSH director, stemcell, and CLI to support errands, run the following errand:

        bosh run errand broker-deregistrar

## Caveats

We have not tested changing the structure of a live cluster, e.g. changing the seed node.

## Tests

To run the Riak CS Release Acceptance tests, you will need:
- a running CF instance
- credentials for a CF Admin user
- a deployed Riak CS Release with the broker registered and the plan made public
- an environment variable `$CONFIG` which points to a `.json` file that contains the application domain

Instructions for running the acceptance tests:

1. Install `go` by following the directions found [here](http://golang.org/doc/install)
2. `cd` into `cf-riak-cs-release/test/acceptance-tests/`
3. Update `cf-riak-cs-release/test/acceptance-tests/integration_config.json`

The following script will configure these prerequisites for a [bosh-lite](https://github.com/cloudfoundry/bosh-lite)
installation. Replace credentials and URLs as appropriate for your environment.

```bash
#! /bin/bash

cat > integration_config.json <<EOF
{
  "api": "api.10.244.0.34.xip.io",
  "admin_user": "admin",
  "admin_password": "admin",
  "apps_domain": "10.244.0.34.xip.io",
  "riak_cs_scheme" : "http://"
}
EOF
export CONFIG=$PWD/integration_config.json
```

If you are running the tests with version newer than 6.0.2-0bba99f of the Go CLI against bosh-lite or any other environment
using self-signed certificates, add

```
  "skip_ssl_validation": true
```

4. Run  the tests

```bash
./bin/test
```

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
