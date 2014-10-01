#Warning

##Riak-CS crashes when filesystem is filled

Riak-CS has a leeway time which limits the rate at which garbage is collected after files have been deleted. If users write and delete too many files within this time period, there is a chance that the underlying Riak system will fill up. This is a dangerous scenario, as when this happens, Riak-CS crashes rather than going into read-only mode, and has no easy way to recover.

We have customized the default configuration so that it is optimized for use with a 10GB filesystem, and to approximate a more synchronous deletion and garbage collection. See [configuration.md](configuration.md) for details on the parameters that have been set to achieve this.
