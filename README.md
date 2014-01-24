# riak-release - BOSH Release

This project is a BOSH release for `riak-release`.

## Blobs

(note for maintainers of this Bosh Release)

The .zip and .tar.gz files (dependencies of packages) that are stored in the blobstore were obtained as follows:

- download source code (git clone)
- `make package.src`
- grab the resulting tar.gz file from package directory

This should work for stanchion and riak_cs. 
[TODO] verify this for riak.

