# Riak CS Test App

This app is used by the acceptance tests and is an example of how you can use the service from an app.

The app allows you to write and read data in the form of key:value pairs to a Riak CS bucket using RESTful endpoints.

### Getting Started

Install the app by pushing it to your Cloud Foundry and binding with the Riak CS service

Example:

    $ cf push riaktest --no-start
    $ cf create-service p-riakcs developer mybucket
    $ cf bind-service riaktest mybucket
    $ cf restart riaktest

### Endpoints

#### PUT /:key

Stores the key:value pair in the Riak CS bucket. Example:

    $ curl -X POST riaktest.my-cloud-foundry.com/service/blobstore/mybucket/foo -d 'bar'
    success


#### GET /:key

Returns the value stored in the Riak CS bucket for a specified key. Example:

    $ curl -X GET riaktest.my-cloud-foundry.com/service/blobstore/mybucket/foo
    bar

#### DELETE /:key

Deletes the bucket. Example:

    $ curl -X DELETE riaktest.my-cloud-foundry.com/service/blobstore/mybucket
    success

Once you've deleted your bucket, you should unbind and delete the service instance, as these are references in Cloud Foundry to an instance which no longer exists.

    $ cf unbind-service riaktest mybucket
    $ cf delete-service mybucket
    $ cf restart riaktest
