#Clients for Riak CS

As Riak CS is API-compliant with Amazon S3, any Amazon s3 client should allow you to communicate with your Riak CS instance. Here are a few we have validated to work with Riak CS for Pivotal CF.

- [s3curl](#s3curl) - Perl script that wraps curl (Apache2 license)
- [s3cmd](#s3cmd) - Command line S3 client written in Python (GPLv2 license)
- [fog](#fog) - this Ruby library is the swiss-army knife for the cloud (MIT license)
- [java client](#java) - a fork of Basho's java client (Apache2 license forthcoming)

##Prerequisites

You have created a service instance, bound it to an application, and have binding credentials from the VCAP\_SERVICES environment variable.

```
"VCAP_SERVICES":
{
  "p-riakcs": [
    {
      "name": "mybucket",
      "label": "p-riakcs",
      "tags": [
        "riak-cs",
        "s3"
      ],
      "plan": "developer",
      "credentials": {
        "uri": "https://my-access-key-id:url-encoded-my-secret-access-key@p-riakcs.mydomain/service-instance-id",
        "access_key_id": "my-access-key-id",
        "secret_access_key": "my-secret-access-key"
      }
    }
  ]
}
```

##<a id='s3curl'></a>s3curl

[s3curl](https://github.com/rtdp/s3curl) is a Perl script. Clone it from github:

<pre class="terminal">
$ git clone https://github.com/rtdp/s3curl
</pre>

Add credentials to `~/.s3curl`:

```
%awsSecretAccessKeys = (
    myuser => {
        id => 'my-access-key-id',
        key => 'my-secret-access-key'
    }
);
```

Edit `s3curl.pl` to add `p-riakcs.myhostname` to the known endpoints:

```
...
my @endpoints = ( 's3.amazonaws.com',
                  's3-us-west-1.amazonaws.com',
                  's3-us-west-2.amazonaws.com',
                  's3-us-gov-west-1.amazonaws.com',
                  's3-eu-west-1.amazonaws.com',
                  's3-ap-southeast-1.amazonaws.com',
                  's3-ap-northeast-1.amazonaws.com',
                  's3-sa-east-1.amazonaws.com',
                  'p-riakcs.mydomain');
...
```
*Note: If you never intend on communicating with any of the amazon services, then you can delete the existing entries (the ones beginning with 's3').*

To list bucket contents at service-instance-location:

<pre class="terminal">
$ ./s3curl.pl --id myuser -- http://p-riakcs.mydomain/service-instance-id
</pre>

To put contents file to bucket with key `mykey`:

<pre class="terminal">
$ ./s3curl.pl --id myuser --put filename -- http://p-riakcs.mydomain/service-instance-id/mykey
</pre>

*Note: curl requires you to escape any special characters in filenames - e.g. filename\\.txt*

To get file with key `mykey` from bucket:

<pre class="terminal">
$ ./s3curl.pl --id myuser -- http://p-riakcs.mydomain/service-instance-id/mykey
</pre>

##<a id='s3cmd'></a>s3cmd

[s3cmd](http://s3tools.org/s3cmd) is a command line client for connecting to S3 compatible blobstores.

A `.s3cfg` file is needed to configure s3cmd to talk to your Riak CS cluster. An example is provided below. Update `access_key`, `secret_key`, `host_base`, `host_bucket`, and `proxy_host` with your own values.

```
[default]
access_key = my-access-key-id
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
host_base = p-riakcs.my-system-domain.com
host_bucket = p-riakcs.my-system-domain.com/%(bucket)s
human_readable_sizes = False
list_md5 = False
log_target_prefix =
preserve_attrs = True
progress_meter = True
proxy_host = p-riakcs.my-system-domain.com
proxy_port = 80
recursive = False
recv_chunk = 4096
reduced_redundancy = False
secret_key = my-secret-access-key
send_chunk = 4096
simpledb_host = sdb.amazonaws.com
skip_existing = False
socket_timeout = 300
urlencoding_mode = normal
use_https = True
verbosity = WARNING
```

If the `.s3cfg` file is not in your home directory, commands must be run with a flag which specifies the location of it; `-c path/to/.s3cfg`.

List bucket contents:

<pre class="terminal">
$ s3cmd ls s3://service-instance-id
</pre>

s3cmd supports a sync feature which is useful for backing up. Bucket contents can be downloaded with the following command:

<pre class="terminal">
$ s3cmd sync s3://service-instance-id /destination/directory
</pre>

Uploading data to a bucket can be done like this:

<pre class="terminal">
$ s3cmd sync /source/directory/* s3://service-instance-id
</pre>

For more information see: 
- [Simple S3cmd How-To](http://s3tools.org/s3cmd-howto)
- [S3cmd Usage](http://s3tools.org/usage)
- [S3cmd S3 Sync How-To](http://s3tools.org/s3cmd-sync)


##<a id='fog'></a>fog

[Fog](http://fog.io) requires Ruby to be installed.

Install the fog gem:

<pre class="terminal">
$ gem install fog
</pre>

Create a ruby client object (requires fog):

```
require 'fog'

basic_client = Fog::Storage.new(
  provider: 'AWS',
  path_style: true,
  host: 'p-riakcs.mydomain',
  port: 80,
  scheme: 'http',
  aws_access_key_id: 'my-access-key-id',
  aws_secret_access_key: 'my-secret-access-key')
```
*Note: try this in irb with `irb -r 'fog'` and then create the client object.*

To list bucket contents at service-instance-location:

`basic_client.get_bucket('service-instance-id')`

To put text to bucket with key `mykey`:

`basic_client.put_object('service-instance-id','mykey','my text here')`

To get file with key `mykey` from bucket:

`basic_client.get_object('service-instance-id', 'mykey')`

##<a id='java'></a>Java Client

See the [repo README](https://github.com/cloudfoundry-incubator/riakcs-java-client/) for documentation
