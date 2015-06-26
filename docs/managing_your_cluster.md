#Managing your Cluster

##Stanchion
###What is Stanchion?
Stanchion is an application used by Riak CS to manage the serialization of requests, which enables Riak CS to manage globally unique entities like **users** and **bucket names**. Futher details can be found on [Basho's website](https://docs.basho.com/riakcs/latest/theory/stanchion/).

###Scaling Stanchion

Stanchion cannot be scaled. There should one and only one Stanchion node within a riak-cs cluster.

Note: If the cluster is managed by [BOSH](https://bosh.io/), it (BOSH) must be configured in [resurrect](https://bosh.io/docs/sysadmin-commands.html#health) mode to auto-recover the Stanchion node in the event of node failure.

##Service Brokers
###Scaling Service Brokers
If the number of service brokers available in the cluster is greater than 1, the [GoRouter](https://github.com/cloudfoundry/gorouter) will distribute requests to the Service Brokers in a [Round-robin](https://en.wikipedia.org/wiki/Round-robin_scheduling) fashion.


##Riak Nodes
###Scaling Up
The default replication factor for Riak buckets is 3. This means that for clusters with 3 or fewer nodes, each node contains a complete replica of all the data. Increasing the number of nodes beyond 3 will cause the data to be redistributed evenly among the nodes, but there will still be only 3 copies total of each key-value pair.

###Scaling Down
Reducing the number of Riak nodes in the cluster will also cause the data to be redistributed evenly among the remaining nodes. If the node count is reduced to 1, the single remaining node will contain all of the data.

Further details of Riak's replication and configurable properties can be found in [Basho's documentation](http://docs.basho.com/riak/latest/theory/concepts/Replication).

##Further Notes
####Riak vs Riak CS
Both Riak CS and Riak are, at their core, places to store objects. Both are open source and both are designed to be used in a cluster of servers for availability and scalability.

The fundamental distinction between the two is simple: Riak CS can be used for storing very large objects, into the terabyte size range, while Riak is optimized for fast storage and retrieval of small objects (typically no more than a few megabytes).

Additional details of the distinction between the two can be found on [Basho's website](http://basho.com/posts/technical/riak-cs-vs-riak/)