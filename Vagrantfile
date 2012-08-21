Vagrant::Config.run do |config|
  config.vm.box = "bosh-solo-0.6.4"
  downloads = "https://github.com/downloads/drnic/bosh-solo/"
  config.vm.box_url = "#{downloads}bosh-solo-0.6.4.box"

  config.vm.network :hostonly, "192.168.10.10"

  config.vm.provision :shell, :path => "scripts/vagrant-setup.sh"

  if local_bosh_src = ENV['BOSH_SRC']
    config.vm.share_folder "bosh-src", "/bosh", local_bosh_src
  end
end