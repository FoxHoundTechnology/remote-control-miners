#!/bin/bash

# Remove unused containers
echo "Removing unused containers..."
sudo docker container prune -f

# Remove unused Docker images
echo "Removing unused images..."
sudo docker image prune -af

# Remove unused volumes
echo "Removing unused volumes..."
sudo docker volume prune -f

# Remove unused networks
echo "Removing unused networks..."
sudo docker network prune -f

echo "Cleanup of unused Docker resources completed."