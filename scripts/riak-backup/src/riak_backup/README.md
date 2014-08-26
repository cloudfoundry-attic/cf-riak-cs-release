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

Binaries for linux_amd64, linux_386, and darwin_amd64 are provided in the [bin](../../bin) directory. These were created using (https://github.com/davecheney/golang-crosscompile). 

`riak_backup` interacts with the file system assuming UNIX paths, so it only works on UNIX or OS X based systems. It will not run on Windows.

[s3cmd](http://s3tools.org/s3cmd) is a dependency; it must be installed and must be in your `$PATH` before running `riak_backup`. A `.s3cfg` file is needed to configure `s3cmd` to communicate with your Riak CS cluster. To configure `.s3cfg` you will need values for `access_key` and `secret_key`. These can be found in your deployment manifest for cf-riak-cs-release. 

        riak_cs:
          admin_key: admin-key # configure for access_key in .s3cfg
          admin_secret: admin-secret # configure for secret_key in .s3cfg

For more instructions on using `s3cmd` see [Clients](http://docs.pivotal.io/p-riakcs/clients.html#s3cmd).

Run `riak_backup` without any arguments to see its usage.

`riak_backup` uses s3cmd's [sync](http://s3tools.org/s3cmd-sync) function to fetch files. If the script is interrupted before completing, it can be run again to resume downloading data where it left off (but all metadata files will be regenerated).
 
