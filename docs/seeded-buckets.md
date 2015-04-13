## Seeded Buckets

In the context of using cf-riak-cs-release with [CloudFoundry](https://github.com/cloudfoundry/cf-release), Riak-CS
buckets are normally created and managed by the [Service Broker](https://github.com/cloudfoundry/cf-riak-cs-broker/).

It is also possible to deploy Riak-CS without the Service Broker. In this case, buckets must be manually managed using
the [Riak-CS API](http://docs.basho.com/riakcs/latest/references/apis/storage/). However, it may be desirable to have
Riak-CS initialized (seeded) with specific buckets already created. This can be accomplished by specifying a list of
`seeded_buckets` in the riak-cs job properties within the bosh deployment manifest.

```yml
---
jobs:
- name: riak-cs
    properties:
      riak_cs:
        seeded_buckets:
        - seeded_bucket_1
        - seeded_bucket_2
```
