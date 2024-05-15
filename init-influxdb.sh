# #!/bin/bash

# echo "Waiting for InfluxDB to be ready..."
# until curl -s http://localhost:8086/health | grep -q '"status":"pass"'; do
#     echo "InfluxDB is not ready. Waiting..."
#     sleep 1
# done
# echo "Response from InfluxDB: $(cat /tmp/influx_health_response)"
# echo "InfluxDB is now ready."

# # Continue with the setup process
# if ! influx bucket list --org "${INFLUX_DB_ORG}" --token "${INFLUX_DB_TOKEN}" 2>/dev/null; then
#   echo "Setting up InfluxDB..."
#   influx setup --force \
#                --username "${INFLUX_DB_ADMIN_USER}" \
#                --password "${INFLUX_DB_ADMIN_USER_PASSWORD}" \
#                --org "${INFLUX_DB_ORG}" \
#                --bucket "${INFLUX_DB_BUCKET}" \
#                --retention "${INFLUX_DB_RETENTION}" \
#                --token "${INFLUX_DB_TOKEN}"
#   echo "InfluxDB setup completed."
# else
#   echo "InfluxDB has already been set up."
# fi
