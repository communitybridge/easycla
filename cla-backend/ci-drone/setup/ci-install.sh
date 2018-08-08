#!/usr/bin/env bash
set -e

[ -z "$DOMAIN" ] && echo "Need to set DOMAIN environment variable" && exit 1;
[ -z "$GITHUB_CLIENT_ID" ] && echo "Need to set GITHUB_CLIENT_ID environment variable" && exit 1;
[ -z "$GITHUB_CLIENT_SECRET" ] && echo "Need to set GITHUB_CLIENT_SECRET environment variable" && exit 1;

cd setup
# Adds Docker's signing key to apt-key, (so we trust Docker's PPA)
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

# Adds Docker's PPA to apt-get. This let's us install docker using apt-get.
sudo add-apt-repository -y "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update

# Actually install docker commuinity edition
sudo apt-get install -y docker-ce

# Install docker compose. Yes, this is how docker recommend we install it: https://docs.docker.com/compose/install/#install-compose
# Note, if you want to upgrade the docker-compose version, you need to bump the 1.21 part of the url.
sudo curl -L https://github.com/docker/compose/releases/download/1.21.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install nginx. We will use this to proxy content to drone, and force HTTPS upgrades.
sudo apt-get install -y nginx

# UFW is the ubuntu firewall. We let NGINX https traffic pass through.
sudo ufw allow 'Nginx Full'

# While we are here, let's install the Let's Encrypt CERT BOT, to automatically create/manage
# our TLS certificate. We won't actually run it just yet.
sudo add-apt-repository -y ppa:certbot/certbot
sudo apt-get update
sudo apt-get install -y python-certbot-nginx

# Copy our docker source files. These files directly configure drone
sudo rm -rf /etc/drone
sudo mkdir /etc/drone
sudo cp docker-compose.yml /etc/drone

# Generate a secret for the drone agent and drone server to share.
# This isn't too important, as both the agent and server are on the same machine, talking behind a firewall.
DRONE_SECRET=$(openssl rand -base64 64)

# We inject DRONE_SECRET,DOMAIN, GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET into these two env files.
# They will be read by drone when setting up the servers/agents.
eval "echo \"$(<agent.env.template)\"" > agent.env.temp
eval "echo \"$(<server.env.template)\"" > server.env.temp
sudo cp agent.env.temp /etc/drone/agent.env
sudo cp server.env.temp /etc/drone/server.env
rm agent.env.temp
rm server.env.temp

# Copy the systemctl service for drone. This will autoboot the service whenever the instance starts up
sudo cp drone.service /etc/systemd/system/drone.service

# Overwrite the nginx configuration with a newly injected DOMAIN
eval "DOMAIN='${DOMAIN}'; echo \"$(<nginx-default.template)\"" > nginx-default.template.temp
sudo cp nginx-default.template.temp /etc/nginx/sites-enabled/default
rm nginx-default.template.temp

# Reboot nginx and drone to update with changes
sudo systemctl enable drone
sudo systemctl restart nginx
sudo systemctl start drone