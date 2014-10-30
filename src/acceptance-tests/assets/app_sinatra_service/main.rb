require 'json'
require 'uri'
require 'bundler'
require 'fog'

Bundler.require
Excon.defaults[:ssl_verify_peer] = false
class ServiceUnavailableError < StandardError;
end

after do
  headers['Services-Nyet-App'] = 'true'
  headers['App-Signature'] = ENV.fetch('APP_SIGNATURE', "")
end

get '/env' do
  ENV['VCAP_SERVICES']
end

get '/rack/env' do
  ENV['RACK_ENV']
end

get '/timeout/:time_in_sec' do
  t = params['time_in_sec'].to_f
  sleep t
  "waited #{t} sec, should have timed out but maybe your environment has a longer timeout"
end

error ServiceUnavailableError do
  status 503
  headers['Retry-After'] = '5'
  body env['sinatra.error'].message
end

error do
  <<-ERROR
Error: #{env['sinatra.error']}

Backtrace: #{env['sinatra.error'].backtrace.join("\n")}
  ERROR
end

post '/service/blobstore/:service_name/:key' do
  value = request.env["rack.input"].read
  client = blobstore_client(params[:service_name])
  bucket_name = get_bucket_name(params[:service_name])
  bucket = client.directories.get(bucket_name)

  bucket.files.create(
      key: params[:key],
      body: value,
      public: true
  )

  value
end

get '/service/blobstore/:service_name/:key' do
  client = blobstore_client(params[:service_name])
  bucket_name = get_bucket_name(params[:service_name])
  bucket = client.directories.get(bucket_name)
  value = bucket.files.get(params[:key]).body

  value
end

delete '/service/blobstore/:service_name' do
  client = blobstore_client(params[:service_name])
  bucket_name = get_bucket_name(params[:service_name])

  client.directories.get(bucket_name).files.all.each do |file|
    file.destroy
  end

  'successfully_deleted'
end

def blobstore_client(service_name)
  blobstore_service = load_service_by_name(service_name)
  uri = URI(blobstore_service.fetch('uri'))
  key = blobstore_service.fetch('access_key_id')
  secret = blobstore_service.fetch('secret_access_key')

  fog_options = {
      provider: 'AWS',
      path_style: true,
      host: uri.host,
      port: uri.port,
      scheme: uri.scheme,
      aws_access_key_id: key,
      aws_secret_access_key: secret
  }
  Fog::Storage.new(fog_options)
end

def get_bucket_name(service_name)
  blobstore_service = load_service_by_name(service_name)
  uri = URI(blobstore_service.fetch('uri'))
  uri.path.chomp("/").reverse.chomp("/").reverse
end

def load_service_by_name(service_name)
  services = JSON.parse(ENV['VCAP_SERVICES'])
  services.values.each do |v|
    v.each do |s|
      if s["name"] == service_name
        return s["credentials"]
      end
    end
  end
  raise "service with name #{service_name} not found in bound services"
end
