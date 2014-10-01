#Warning

##Riak crashes when filesystem is filled

Riak-CS has a leeway time which limits the rate at which garbage is collected after files have been deleted. If users write and delete too many files within this time period, there is a chance that the underlying Riak system will fill up. This is a dangerous scenario, as when this happens, Riak crashes rather than going into read-only mode, and has no easy way to recover. This will result in downtime for the Riak-CS service while Riak is not working.

We have customized the default configuration so that it is optimized for use with a 10GB filesystem, and to approximate a more synchronous deletion and garbage collection. See [configuration.md](configuration.md) for details on the parameters that have been set to achieve this.

###How to recover:

1. Edit the cf-riak-cs-release manifest to give the Riak-CS nodes a larger persistent disk (say 2x the original size).
1. Redeploy with the new manifest.
1. Wait __leeway_interval__ + 3 minutes so that Riak-CS has time to garbage collect, and Riak has time to initiate block merges to reclaim disk space. 
1. Check to see how much disk space has been reclaimed.
1. If desired, redeploy with the original persistent disk size.
