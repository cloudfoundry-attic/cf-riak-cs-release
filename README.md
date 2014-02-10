# riak-release

A BOSH release for Riak and Riak CS.

Make sure to run `git submodule update --init` before creating the release.

This project is based on [BrianMMcClain/riak-release](https://github.com/BrianMMcClain/riak-release).


## Deployment

After you create and upload the release, the `generate_deployment_manifest` script can generate the deployment manifest based on a deployment template. The available templates are: `warden`(bosh-lite) and `vsphere`.

1. Put your director UUID into `templates/riak-cs-service.yml`
2. Generate the manifest: `./generate_deployment_manifest <template> > riak-cs-service.yml`
To tweak the deployment settings, you can modify the resulting file `riak-cs-service.yml`.
3. To deploy: `bosh deployment riak-cs-service.yml && bosh deploy`


## Caveats

We have not tested changing the structure of a live cluster, e.g. changing the seed node.

## Tests
Instructions for running the cf-service-acceptance tests under the tests/ directory:

1. Install go by following the directions found [here](http://golang.org/doc/install)
2. Set environment variables `CF_COLOR=false` and `CF_VERBOSE_OUTPUT=true`
3. `cd` into this directory: `riak-release/test/cf-service-acceptance-tests/apps/`
4. Alter `integration_config.json` to use the domain of the Cloud Foundry you wish to test against 
5. Run `CONFIG=/Users/pivotal/workspace/riak-release/test/cf-service-acceptance-tests/integration_config.json go test`

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
