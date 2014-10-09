#RiakCS Behavior

As we investigate various use-case scenarios for RiakCS, we will document our findings here.

##Object Storage

Any named object has multiple versions internal to RiakCS, which are stored in `object manifests`. There is only one 'active' manifest for an object at any given time. When an object is deleted, necessary versions are deleted by removing object manifests.

##Garbage Collection

Garbage collection is initiated via the Garbage Collection daemon, which wakes up at an appointed interval (`gc_interval`) or scheduled time. The GC daemon checks the `riak-cs-gc` bucket for keys, and attempts to delete keys which have been deleted for longer than `leeway_seconds`.

The daemon starts up a worker process to delete the eligible keys. If `gc_paginated_indexes` is set to true, the daemon will send batches of `gc_batch_size` to be deleted by the worker process. The number of concurrent worker processes is bounded by `gc_max_workers`.

Each worker focuses on deleting one block of an object at a time. Once all the blocks in an object are deleted, a separate type of worker will delete the object's manifest version.

### Experiment results

* There appears to be no significant difference in GC between RiakCS 1.5.0 and 1.5.1. Github issues (e.g. [here](https://github.com/basho/riak_cs/pull/949) and [here](https://github.com/basho/riak_cs/issues/946)) and [1.5.1 release notes](http://basho.com/basho-is-pleased-to-announce-the-release-of-riak-cs-1-5/) appear to indicate 1.5.1 performs better - we have not found that to be the case.
* In the following results, NxS files refers to uploading N files of size S all at once, and subsequently deleting them at once, in a single batch.
* Irrespective of number of workers, we can consistently garbage collection 50-60MB files. We successfully garbage collected 1x50MB files and 100x600K files.
* We can garbage collect up to 250MB total at once if we keep the batch size to 1, irrespective of the number of workers. We successfully garbage collected 1x100MB and 400x600K files.
* The Github issues mentioned above indicate that the deletion of a multi-part upload is more error-prone. It is up to the client to disable multi-part upload. For details of how to do this with s3cmd, see [here](http://s3tools.org/kb/item13.htm).
 * Disabling multi-part upload in the client improved the reliability of GC (it succeeded more often), but did not fix it altogether (it still timed out on multiple 200MB files)
