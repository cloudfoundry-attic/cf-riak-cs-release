## Preseeding buckets

It is possible to seed buckets during bosh deployment by adding the bucket names to the manifest as follows:

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

Alternatively, the same manifest can be obtained by merging a stub file with the above contents; this avoids having to edit the spiff-generated manifest directly.
