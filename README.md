# Cloud Foundry Riak CS Service

A BOSH release of an S3-compatible object store for Cloud Foundry using [Riak CS](http://basho.com/riak-cloud-storage/) and a [v2 Service Broker](http://docs.cloudfoundry.org/services/).

This project is based on [BrianMMcClain/riak-release](https://github.com/BrianMMcClain/riak-release).

## Release Notes

The release notes can be found [here](https://github.com/cloudfoundry/cf-riak-cs-release/wiki/Release-Notes).

## <a id='branches'></a>Getting the code

Final releases are designed for public use, and are tagged with a version number of the form "v<N>".

The [**develop**](https://github.com/cloudfoundry/cf-riak-cs-release/tree/develop) branch is where we do active development. Although we endeavor to keep the [**develop**](https://github.com/cloudfoundry/cf-riak-cs-release/tree/develop) branch stable, we do not guarantee that any given commit will deploy cleanly.

The [**release-candidate**](https://github.com/cloudfoundry/cf-riak-cs-release/tree/release-candidate) branch has passed all of our unit, integration, smoke, & acceptance tests, but has not been used in a final release yet. This branch should be fairly stable.

The [**master**](https://github.com/cloudfoundry/cf-riak-cs-release/tree/master) branch points to the most recent stable final release.

At semi-regular intervals a final release is created from the [**release-candidate**](https://github.com/cloudfoundry/cf-riak-cs-release/tree/release-candidate) branch. This final release is tagged and pushed to the [**master**](https://github.com/cloudfoundry/cf-riak-cs-release/tree/master) branch.

Pushing to any branch other than [**develop**](https://github.com/cloudfoundry/cf-riak-cs-release/tree/develop) will create problems for the CI pipeline, which relies on fast forward merges. To recover from this condition follow the instructions [here](https://github.com/cloudfoundry/cf-release/blob/master/docs/fix_commit_to_master.md).

## Development

This BOSH release doubles as a `$GOPATH`. It will automatically be set up for
you if you have [direnv](http://direnv.net) installed.

    # fetch release repo
    mkdir -p ~/workspace
    cd ~/workspace
    git clone https://github.com/cloudfoundry/cf-riak-cs-release.git
    cd cf-riak-cs-release/

    # switch to develop branch (not master!)
    git checkout develop

    # automate $GOPATH and $PATH setup
    direnv allow

    # initialize and sync submodules
    ./update

If you do not wish to use direnv, you can simply `source` the `.envrc` file in the root
of the release repo.  You may manually need to update your `$GOPATH` and `$PATH` variables
as you switch in and out of the directory.

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

You can use a pre-built final release or build a dev release from any of the branches described in <a href="#branches">Getting the Code</a>. 

Final releases are stable releases created periodically for completed features. They also contain pre-compiled packages, which makes deployment much faster. To deploy the latest final release, simply check out the **master** branch. This will contain the latest final release and accompanying materials to generate a manifest. If you would like to deploy an earlier final release, use `git checkout <tag>` to obtain both the release and corresponding manifest generation materials. It's important that the manifest generation materials are consistent with the release.

If you'd like to deploy the latest code, build a release yourself from the **develop** branch.

#### Upload a pre-built final BOSH release

Run the upload command, referencing the latest config file in the `releases` directory. 

  ```
  $ cd ~/workspace/cf-riak-cs-release
  $ git checkout master
  $ ./update
  $ bosh upload release releases/cf-riak-cs-<N>.yml
  ```

If deploying an **older** final release than the latest, check out the tag for the desired version; this is necessary for generating a manifest that matches the code you're deploying.

  ```
  $ cd ~/workspace/cf-riak-cs-release
  $ git checkout v<N>
  $ ./update
  $ bosh upload release releases/cf-riak-cs-<N>.yml
  ```

#### Create and upload a BOSH Release:

1. Checkout one of the branches described in <a href="#branches">Getting the Code</a>. Build a BOSH development release. 

  ```
  $ cd ~/workspace/cf-riak-cs-release
  $ git checkout release-candidate
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

  `URL` specifies where the Cloud Controller will access the Riak CS broker. Use the value of the manifest property `properties.broker.host`.

  For more information, see [Managing Service Brokers](http://docs.cloudfoundry.org/services/managing-service-brokers.html).

1. Then [make the service plan public](http://docs.cloudfoundry.org/services/managing-service-brokers.html#make-plans-public).


## Running Acceptance Tests

To run the Riak CS acceptance tests you will need:
- a running CF instance
- credentials for a CF Admin user
- a deployed Riak CS Release with the broker registered and the plan made public
- a security group granting access to the service for applications

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

To run the acceptance tests manually you will also need an environment variable `$CONFIG` which points to a `.json` file that contains the application domain.

1. Install `go` by following the directions found [here](http://golang.org/doc/install)
2. `cd` into `cf-riak-cs-release/src/acceptance-tests/`
3. Update `cf-riak-cs-release/src/acceptance-tests/integration_config.json`

   The following commands provide a shortcut to configuring `integration_config.json` with values for a [bosh-lite](https://github.com/cloudfoundry/bosh-lite)
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

    Note: `skip_ssl_validation` requires CLI v6.0.2 or newer.
    
    All timeouts in the test suite can be scaled proportionally by changing the timeout_scale factor.

4. Run  the tests

  ```
  $ ./bin/test
  ```

## Security Groups

Since [cf-release](https://github.com/cloudfoundry/cf-release) v175, applications by default cannot to connect to IP addresses on the private network. This may prevents applications from connecting to the Riak CS service. As applications reach the Riak CS service through the router tier in cf-release, create a new security group for the IP configured for the load balancer balancing traffic across your cf-release routers. By default this will be the HAProxy job in cf-release.

1. Add the rule to a file in the following json format; multiple rules are supported.

  ```
  [
      {
        "destination": "10.244.0.34",
        "protocol": "all"
      }
  ]
  ```
- Create a security group from the rule file.
  <pre class="terminal">
  $ cf create-security-group p-riakcs rule.json
  </pre>
- Enable the rule for all apps
  <pre class="terminal">
  $ cf bind-running-security-group p-riakcs
  </pre>

Changes are only applied to new application containers; in order for an existing app to receive security group changes it must be restarted.

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

## Using the Riak CS service

See [Clients for Riak CS](docs/clients.md) for a list of clients that have been validated to work with the service.

[The included test application](src/acceptance-tests/assets/app_sinatra_service), written in Ruby and using the Fog library, is an example of how to use the service with an application.

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
