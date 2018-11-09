# -*- mode: ruby -*-
# vi: set ft=ruby :

$script = <<SCRIPT
sudo apt install openssh-server
sudo systemctl start sshd
SCRIPT

Vagrant.configure("2") do |config|
  config.vm.provider "virtualbox" do |v|
    v.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
  end
  config.vm.box = "ubuntu/xenial64"
  config.vm.network "forwarded_port", guest: 22, host: 4841
  config.vm.hostname = "target"
  config.vm.provision "shell", inline: $script
end
