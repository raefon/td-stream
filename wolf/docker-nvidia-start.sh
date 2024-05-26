#/usr/bin/env bash

# Prepare instance to run Wolf via Docker Compose
# - Check Nvidia installed and ready
# - Build Nvidia image
# - Populate volume

set -e

WOLF_COMPOSE_FILE=$1

NVIDIA_VERSION_FILE="/sys/module/nvidia/version"

# Exit gracefully as on initial install as it's expected
# that nvidia version file won't appear until next reboot
# So that nixos-rebuild won't complain
if [ ! -f "$NVIDIA_VERSION_FILE" ]; then
    echo "Nvidia version file does not exists. Does your system needs rebooting? Exiting."
    exit 0
fi

# Check Nvidia driver version
REQUIRED_NVIDIA_VERSION="530.30.02"
CURRENT_NVIDIA_VERSION=$(cat $NVIDIA_VERSION_FILE)
if [[ "$(printf '%s\n' "$REQUIRED_NVIDIA_VERSION" "$CURRENT_NVIDIA_VERSION" | sort -V | head -n1)" != "$REQUIRED_NVIDIA_VERSION" ]]; then
    echo "Your Nvidia driver version ($CURRENT_NVIDIA_VERSION) is less than the required version ($REQUIRED_NVIDIA_VERSION)."
    exit 1
fi

# Build image
curl https://raw.githubusercontent.com/games-on-whales/gow/master/images/nvidia-driver/Dockerfile | \
    docker build -t gow/nvidia-driver:local -f - --build-arg NV_VERSION=$CURRENT_NVIDIA_VERSION /tmp

# Populate nvidia volume
# Create a dummy container to create and populate volume if not already exists
# If volume already exists this won't have any effect
NVIDIA_VOLUME_NAME=nvidia-driver-vol-$CURRENT_NVIDIA_VERSION
docker create --rm --mount source=$NVIDIA_VOLUME_NAME,destination=/usr/nvidia gow/nvidia-driver:local sh

# Check if the nvidia-drm module's modeset parameter is set to Y
if [[ $(sudo cat /sys/module/nvidia_drm/parameters/modeset) != "Y" ]]; then
    echo "The nvidia-drm module's modeset parameter is not set to Y."
    echo "See https://games-on-whales.github.io/wolf/stable/user/quickstart.html for expected setup."
    exit 2
fi

# Sometime /dev/nvidia-caps/* is not present after boot
# Running this seems to make it appear
# TODO enqure why
sudo nvidia-smi > /dev/null # still show errors
sudo nvidia-container-cli --load-kmods info

# Start wolf with matching nvidia volume
SCRIPT_DIR=$(dirname $0)
NVIDIA_VOLUME_NAME=$NVIDIA_VOLUME_NAME docker compose -p wolf -f $WOLF_COMPOSE_FILE up -d