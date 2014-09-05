require 'fog'

bucket_name = ARGV[0]
host = ARGV[1]
admin_access_key = ARGV[2]
admin_secret_key = ARGV[3]

client = Fog::Storage.new(
  provider: 'AWS',
  path_style: true,
  host: host,
  port: 8080,
  scheme: 'http',
  aws_access_key_id: admin_access_key,
  aws_secret_access_key: admin_secret_key)
puts "Fog client: #{client.inspect}"

begin
  puts "Bucket name: #{bucket_name}"
  client.directories.create(key: bucket_name)
  puts "Created bucket #{bucket_name}"
rescue => e
  puts "Error creating bucket: #{e.message}"
end



