# riak-release - BOSH Release

This project is a BOSH release for `riak-release`.

## Deploying

### Configuring admin user

To create an admin user:

1. set `anonymous_user_creation` to `true` in the manifest and deploy the release
2. create an admin user (see instructions in the example manifest) 
3. set `anonymous_user_creation` back to `false` and add the admin user's `admin_key` and `admin_secret` to the manifest
4. re-deploy the release

## Blobs

(note for maintainers of this Bosh Release)

The `riak-cs-*.tar.gz` and `stanchion-*.tar.gz` files (dependencies of packages) that are stored in the blobstore were obtained as follows:

- download source code (git clone)
- `make package.src`
- grab the resulting `tar.gz` file from package directory

This should work for `stanchion` and `riak_cs`. 

TODO - verify where the `riak`, `git`, and `erlang` tarfiles came from.

## TODO

- The settings for the Riak job in this release are configured with options suggested by Basho for deploying Riak in a Riak CS cluster.  We could add an option to configure Riak for standalone operation (when a manifest includes only Riak but not Riak CS) 

- Add automatic creation of admin user upon initial deploy