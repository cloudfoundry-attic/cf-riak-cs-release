# riak-release - BOSH Release

This project is a BOSH release for `riak-release`.

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