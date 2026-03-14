#!/usr/bin/env bash
set -e

echo "=== Step 0: Ensure running as root ==="
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

echo "=== Step 1: Update packages ==="
apt update
apt install -y curl wget git unzip ca-certificates lsb-release gnupg software-properties-common sudo

echo "=== Step 2: Install Fish shell ==="
apt install -y fish
echo "Changing default shell for root to Fish..."
chsh -s /usr/bin/fish root
# Verify default shell
echo "Default shell for root is now: $(getent passwd root | cut -d: -f7)"

echo "=== Step 3: Install Docker ==="
apt install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
  >/etc/apt/sources.list.d/docker.list

apt update
apt install -y docker-ce docker-ce-cli containerd.io

# Install Docker Compose v2
DOCKER_COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep tag_name | cut -d '"' -f 4)
curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
docker-compose version --short

echo "=== Step 4: Install Go 1.25.1 ==="
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz
rm go1.25.1.linux-amd64.tar.gz

# Export Go environment for this script
export PATH=/usr/local/go/bin:/usr/bin:/bin:$PATH
export GOPATH=/root/go
export PATH=$GOPATH/bin:$PATH

echo "Go version: $(go version)"

echo "=== Step 5: Install Task CLI v3 ==="
curl -L https://taskfile.dev/install.sh | sh -s -- -d
mv ./bin/task /usr/local/bin/task
chmod +x /usr/local/bin/task
task --version

echo "=== Step 6: Install unzip ==="
apt install -y unzip

echo "=== Step 7: Configure Fish environment variables ==="
# Create Fish config if it doesn't exist
mkdir -p /root/.config/fish
cat << 'EOF' > /root/.config/fish/config.fish
# Go environment
set -gx PATH /usr/local/go/bin /usr/bin /bin $PATH
set -gx GOPATH /root/go
set -gx PATH $GOPATH/bin $PATH
EOF

echo "=== Setup Complete! ==="
echo "Go: $(go version)"
echo "Task: $(task --version)"
echo "Docker: $(docker --version)"
echo "Docker Compose: $(docker-compose version --short)"
echo "Default shell for root: $(getent passwd root | cut -d: -f7)"
