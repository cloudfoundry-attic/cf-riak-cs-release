# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "lucid64"
  config.vm.box_url = "http://files.vagrantup.com/lucid64.box"

  config.vm.network :private_network, ip: "192.168.10.10"

  
  config.vm.synced_folder "/Users/brian/Documents/checkouts/riak-release", "/home/vagrant/release"
  config.vm.synced_folder "/Users/brian/Documents/checkouts/nise-bosh-vagrant/lib/nise-bosh-vagrant/../../scripts", "/home/vagrant/scripts"
end
