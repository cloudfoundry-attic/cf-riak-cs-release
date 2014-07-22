## riak_backup

This script can be used by a CloudFoundry administrator to back up the data in all service instances (buckets) created by the Riak CS Service Broker.

The script uses admin credentials and fetches the list of organizations, spaces, and service instances. For all service instances of type p-riakcs, it uses `s3cmd` to download all data files from the bucket corresponding to that instance.

The data files are saved in a backup directory whose structure is illustrated by the following example

```
└── orgs
    └── organization-name-0
        └── spaces
            ├── space-name-0
            │   └── service_instances
            │       ├── service-instance-name-0
            │       │   ├── data
            │       │   │   └── files-fetched-from-the-bucket...
            │       │   └── metadata.yml
            │       └── service-instance-name-1
            │           ├── data
            │           │   └── files-fetched-from-the-bucket...
            │           └── metadata.yml
            └── space-name-2
                └── service_instances
                    ├── service-instance-name-2
                    │   ├── data
                    │   │   └── files-fetched-from-the-bucket...
                    │   └── metadata.yml
                    └── service-instance-name-3
                        ├── data
                        │   └── files-fetched-from-the-bucket...
                        └── metadata.yml
```

The `metadata.yml` file in each instance directory lists the `guid` of the service instance, as well as the `name` and `guid` of all apps to which the instance is bound (if any).

```
service_instance_guid: service-instance-guid-0
bound_apps:
- name: app-name-0
  guid: app-guid-0
- name: app-name-1
  guid: app-guid-1
```

### Usage

`riak_backup` interacts with the file system assuming UNIX paths, so it only works on UNIX or OS X based systems. It will not run on Windows.

[s3cmd](http://s3tools.org/s3cmd) is a dependency; it must be installed and must be in your `$PATH` before running `riak_backup`.

A `.s3cfg` file is needed to configure s3cmd to talk to your Riak CS cluster.
An example is provided below. Please adjust the `access_key`, `secret_key`, `host_base`, `host_bucket`, and `proxy_host` parameters to match your Riak CS cluster config.

```
[default]
access_key = admin-key
bucket_location = US
cloudfront_host = cloudfront.amazonaws.com
cloudfront_resource = /2010-07-15/distribution
default_mime_type = binary/octet-stream
delete_removed = False
dry_run = False
encoding = UTF-8
encrypt = False
follow_symlinks = False
force = False
get_continue = False
gpg_command = None
gpg_decrypt = %(gpg_command)s -d --verbose --no-use-agent --batch --yes --passphrase-fd %(passphrase_fd)s -o %(output_file)s %(input_file)s
gpg_encrypt = %(gpg_command)s -c --verbose --no-use-agent --batch --yes --passphrase-fd %(passphrase_fd)s -o %(output_file)s %(input_file)s
gpg_passphrase =
guess_mime_type = True
host_base = p-riakcs.10.244.0.34.xip.io
host_bucket = p-riakcs.10.244.0.34.xip.io/%(bucket)s
human_readable_sizes = False
list_md5 = False
log_target_prefix =
preserve_attrs = True
progress_meter = True
proxy_host = p-riakcs.10.244.0.34.xip.io
proxy_port = 80
recursive = False
recv_chunk = 4096
reduced_redundancy = False
secret_key = admin-secret
send_chunk = 4096
simpledb_host = sdb.amazonaws.com
skip_existing = False
socket_timeout = 300
urlencoding_mode = normal
use_https = False
verbosity = WARNING
```

Run `riak_backup` without any arguments to see the usage.

`riak_backup` uses s3cmd's (sync)[http://s3tools.org/s3cmd-sync] function to fetch files; if the script is interrupted before completing, it will re-generate all metadata files, and will sync data from all buckets (won't download it from scratch).
