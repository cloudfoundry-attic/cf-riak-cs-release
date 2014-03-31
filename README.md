# cf-riak-cs-release

A BOSH release for Riak and Riak CS.

Make sure to run `./update` (to update git submodules) before creating the release.

This project is based on [BrianMMcClain/riak-release](https://github.com/BrianMMcClain/riak-release).

## Deployment

1.  First create the release, naming it `cf-riak-cs`.
1.  Then upload the release.
1.  Finally make sure you have uploaded the appropriate stemcell for your deployment (either vsphere or warden)

### Configurations common to both environments
1. SSL:
	* There are two properties under properties.riak-cs called `ssl_enabled` and `skip_ssl_validation`
	* `ssl_enabled` is defaulted to true and `skip_ssl_validation` defaults to false, which will only work if you have valid certs on your CF deployment
	* If you wish to change either of these, you can create a stub with these properties set to your needs:
```
properties:
  riak_cs:
    ssl_enabled: <your value>
    skip_ssl_validation: <your value>
```


### To a BOSH-lite environment

1. Create a stub file called `riak-cs-lite-stub.yml` that contains your director UUID (which you can get from running `bosh status`):

		director_uuid: YOUR-DIRECTOR-GUID-HERE
		properties:
		  domain: your-cf-system-domain-here   # such as 10.244.0.34.xip.io
		  cf:
        api_url: http://api.YOUR-CF-DOMAIN-HERE    # such as http://api.10.244.0.34.xip.io

2. Generate the manifest: `./generate_deployment_manifest warden riak-cs-lite-stub.yml > riak-cs-lite.yml`
To tweak the deployment settings, you can modify the resulting file `riak-cs-lite.yml`.
3. To deploy: `bosh deployment riak-cs-lite.yml && bosh deploy`

### To a vSphere environment

1. Create a stub file called `riak-cs-vsphere-stub.yml` that contains your director UUID (which you can get from running `bosh status`).
It also needs your network settings, with 6 static IPs and 6+ dynamic IPs, like this:

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

1. Generate the manifest: `./generate_deployment_manifest vsphere riak-cs-vsphere-stub.yml > riak-cs-vsphere.yml`
To tweak the deployment settings, you can modify the resulting file `riak-cs-vsphere.yml`.
1. To deploy: `bosh deployment riak-cs-vsphere.yml && bosh deploy`

### To an AWS environment

1. Create a stub file called `riak-cs-aws-stub.yml` that contains your director UUID (which you can get from running `bosh status`).
It also needs your network settings, like this:

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
            admin_username: CF-ADMIN-USERNAMEå
            admin_password: CF-ADMIN-PASSWORD

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


## Caveats

We have not tested changing the structure of a live cluster, e.g. changing the seed node.

## Tests
Instructions for running the acceptance tests:

1. Install `go` by following the directions found [here](http://golang.org/doc/install)
1. Set environment variables `export CF_COLOR=false` and `export CF_VERBOSE_OUTPUT=true`
1. Update `cf-riak-cs-release/test/acceptance-tests/integration_config.json` with the domain of the Cloud Foundry you wish to test against. NOTE: this should not include `http://` or `https://`. E.g. for bosh-lite this might be `10.244.0.34.xip.io`.
1. `cd` into `cf-riak-cs-release/test/acceptance-tests/`
1. Target and log into your CF instance, e.g. `cf login -a http://api.10.244.0.34.xip.io -u admin -p admin`
1. Run `CONFIG=/Users/pivotal/workspace/cf-riak-cs-release/test/acceptance-tests/integration_config.json bin/test`


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
