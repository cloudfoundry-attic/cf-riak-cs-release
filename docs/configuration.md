#Riak-CS Config Options

It is possible to configure RiakCS, Riak, and Stanchion via an "app.config" file corresponding to each job. To change this configuration before building the RiakCS release, you need to change the "<job_name>.app.config.erb" file under each job's template directory.

##Configuring Garbage Collection


Currently, we configure these fields in `riak_cs.app.config.erb` under 'Garbage Collection':

- **leeway_seconds** (Default: 60 seconds) How long a file needs to have been deleted for RiakCS to attempt to garbage collect it
- **gc_interval** (Default: 900 seconds) How often RiakCS runs garbage collection, which notifies Riak of deletions

Additional information on configuring RiakCS can be found [here](http://docs.basho.com/riakcs/latest/cookbooks/configuration/Configuring-Riak-CS/#Garbage-Collection-Settings)

Currently, we configure these fields in `riak.app.config.erb` under 'Bitcask Config':

- **data_root** (Default: "/var/vcap/store/riak/rel/data/bitcask") location of Riak's data directory
- **max_file_size** (Default: 100MB) Maximum size of an open Bitcask file in Riak. These will never be garbage collected while open, so making this smaller makes files eligible for garbage collection more often.
- **dead_bytes_merge_trigger** (Default: 25MB) Number of dead bytes in a Bitcask file to trigger a merge. This merge is the final step of garbage collection, where blocks are defragmented and merged to eliminate dead bytes.
- **dead_bytes_threshold** (Default: 6.25MB) Once a merge is triggered, this threshold determines which blocks will actually be merged. This may be smaller than dead_bytes_merge_trigger.

Additional information on configuring Riak and Bitcask can be found [here](http://docs.basho.com/riak/latest/ops/advanced/backends/bitcask/#Configuring-Bitcask)

By default, we configure the Riak-CS nodes with 10 GB of persistent disk, with many Bitcask data subdirectories, each of which can have an open Bitcask file. To ensure that garbage collection remains possible, we set max_file_size to 100 MB so that we will always have at least some space available to perform garbage collection. We decrease the dead_bytes_merge_trigger and dead_bytes_threshold correspondingly.

When leeway_seconds is too long a wait time, we've seen a Riak node fill up its disk before it is told to garbage collect. So, we set leeway_seconds to one minute to provide adequate time for the deletion to propagate to the other Riak nodes, but to allow deletes to propagate to Riak somewhat synchronously.
