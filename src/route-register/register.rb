#!/usr/bin/env ruby

require 'yaml'
require 'nats/client'

ENV['SETTINGS_PATH'] ||= File.expand_path('../register_settings.yml', __FILE__)

settings = YAML.load_file(ENV['SETTINGS_PATH'])

Kernel.at_exit do
  NATS.start(:uri => settings['message_bus_servers']) do
    registrar.shutdown { EM.stop }
  end
end

NATS.start(:uri => settings['message_bus_servers']) do
  registrar.register_with_router
end

def registrar
  Cf::Registrar.new(
    :message_bus_servers => settings['message_bus_servers'],
    :host                => settings['external_ip'],
    :port                => settings['port'],
    :uri                 => settings['external_host'],
    :tags                => { "component" => "Cf-Riak-CS-Node" })
end
