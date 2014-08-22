# Cloud Foundry Riak CS Service

A BOSH release of an S3-compatible object store for Cloud Foundry using [Riak CS](http://basho.com/riak-cloud-storage/) and a [v2 Service Broker](http://docs.cloudfoundry.org/services/).

This project is based on [BrianMMcClain/riak-release](https://github.com/BrianMMcClain/riak-release).

## Release Notes

[Release Notes](https://github.com/cloudfoundry/cf-riak-cs-release/wiki/Release-Notes)

## Deployment

### Prerequisites

- A deployment of [BOSH](https://github.com/cloudfoundry/bosh)
- A deployment of [Cloud Foundry](https://github.com/cloudfoundry/cf-release)
- Instructions for installing BOSH and Cloud Foundry can be found at http://docs.cloudfoundry.org/.

### Overview

1. Upload Stemcell
1. [Upload Release](#upload_release)
1. [Create Manifest and Deploy](#create_manifest)
  - [BOSH-lite](#bosh-lite)
  - [vSphere](#vsphere)
  - [AWS](#aws)
  - [Deployment Manifest Stub Properties](#stub-properties)
1. [Register the Service Broker](#register_broker)

### Upload Release<a name="upload_release"></a>

You can use a pre-built final release or build a release from HEAD. Final releases contain pre-compiled packages, making deployment much faster. However, these are created manually and infrequently. To be sure you're deploying the latest code, build a release yourself.

#### Upload a pre-built final BOSH release

1. Check out the tag for the desired version. This is necessary for generating a manifest that matches the code you're deploying.

  ```
  $ cd ~/workspace/cf-riak-cs-release
  $ ./update
  $ git checkout v5
  $ git submodule update --recursive
  ```

1. Run the upload command, referencing one of the config files in the `releases` directory.

  ```
  $ bosh upload release releases/cf-riak-cs-5.yml
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

### Create Manifest and Deploy<a name="create_manifest"></a>

#### BOSH-lite<a name="bosh-lite"></a>

1. Run the script [`bosh-lite/make_manifest`](bosh-lite/make_manifest) to generate your manifest for bosh-lite. This script uses a stub provided for you, `bosh-lite/stub.yml`.

```
    $ ./bosh-lite/make_manifest
```

The manifest will be written to `bosh-lite/manifests/cf-riak-cs-manifest.yml`, which can be modified to change deployment settings.

1. The `make_manifest` script will set the deployment to `bosh-lite/manifests/cf-riak-cs-manifest.yml` for you, so to deploy you only need to run:
  ```
  $ bosh deploy
  ```

#### vSphere<a name="vsphere"></a>

1. Create a stub file called cf-riak-cs-vsphere-stub.yml by copying and modifying the [sample_vsphere_stub.yml](https://github.com/cloudfoundry/cf-riak-cs-release/blob/master/templates/sample_stubs/sample_vsphere_stub.yml) in templates/sample_stubs.

2. Generate the manifest:
  ```
  $ ./generate_deployment_manifest vsphere cf-riak-cs-vsphere-stub.yml > cf-riak-cs-vsphere.yml
  ```
To tweak the deployment settings, you can modify the resulting file `cf-riak-cs-vsphere.yml`.

3. To deploy:
  ```
  $ bosh deployment cf-riak-cs-vsphere.yml && bosh deploy
  ```

#### AWS<a name="aws"></a>

1. Create a stub file called cf-riak-cs-aws-stub.yml by copying and modifying the [sample_aws_stub.yml](https://github.com/cloudfoundry/cf-riak-cs-release/blob/master/templates/sample_stubs/sample_aws_stub.yml) in templates/sample_stubs.

1. Generate the manifest:
  ```
  $ ./generate_deployment_manifest aws cf-riak-cs-aws-stub.yml > cf-riak-cs-aws.yml
  ```
To tweak the deployment settings, you can modify the resulting file `cf-riak-cs-aws.yml`.

1. To deploy:
  ```
  $ bosh deployment cf-riak-cs-aws.yml && bosh deploy
  ```

#### Deployment Manifest Properties<a name="stub-properties"></a>

Manifest properties are described in the `spec` file for each job; see [jobs](jobs).

You can find your director_uuid by running `bosh status`.

## Register the Service Broker<a name="register_broker"></a>

### Using BOSH errands

BOSH errands were introduced in version 2366 of the BOSH CLI, BOSH Director, and stemcells.
  ```
  $ bosh run errand broker-registrar
  ```
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
$ bosh run errand acceptance-tests
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

  ```
  $ ./bin/test
  ```

## De-registering the broker

The following commands are destructive and are intended to be run in conjuction with deleting your BOSH deployment.

### Using BOSH errands

BOSH errands were introduced in version 2366 of the BOSH CLI, BOSH Director, and stemcells.

This errand runs the two commands listed in the manual section below from a BOSH-deployed VM. This errand should be run before deleting your BOSH deployment. If you have already deleted your deployment follow the manual instructions below.

```
$ bosh run errand broker-deregistrar
```

#### Manually

Run the following:

```
$ cf purge-service-offering p-riakcs
$ cf delete-service-broker p-riakcs
```

## Using s3curl to read and write to your Riak CS bucket

Clone s3curl from github:
```
$ git clone https://github.com/rtdp/s3curl
```
Add credentials to `~/.s3curl`:
```bash
%awsSecretAccessKeys = (
    myuser => {
        id => 'my-access-key-id',
        key => 'my-secret-access-key'
    }
);
```

Edit `s3curl.pl` to add `p-riakcs.mydomain` to the known endpoints:

```bash
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
```
$ ./s3curl.pl --id myuser -- http://p-riakcs.mydomain/service-instance-id
```
To put contents file to bucket with key `mykey`:
```
$ ./s3curl.pl --id myuser --put filename -- http://p-riakcs.mydomain/service-instance-id/mykey
```
*Note: curl requires you to escape any special characters in filenames - e.g. filename\\.txt*

To get file with key `mykey` from bucket:
```
$ ./s3curl.pl --id myuser -- http://p-riakcs.mydomain/service-instance-id/mykey
```
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
