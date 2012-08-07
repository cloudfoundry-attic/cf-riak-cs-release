require 'riak'

# Create a client interface
client = Riak::Client.new

# Create a client interface that uses Excon
client = Riak::Client.new(:http_backend => :Excon)

# Automatically balance between multiple nodes
client = Riak::Client.new(:host => '127.0.0.1', :http_port => 8100)

puts client.ping

# Retrieve a bucket
bucket = client.bucket("doc")  # a Riak::Bucket

# Get an object from the bucket
object = bucket.get_or_new("test")   # a Riak::RObject

# Change the object's data and save
object.raw_data = "Hello World!"
object.content_type = "text/plain"
object.store

# Reload an object you already have
object.reload                  # Works if you have the key and vclock, using conditional GET
object.reload :force => true   # Reloads whether you have the vclock or not

# Access more like a hash, client[bucket][key]
puts client['doc']['test'].raw_data   # the Riak::RObject