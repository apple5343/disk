#!/bin/bash

echo "Creating Kafka topics..."


kafka-topics --create --if-not-exists --bootstrap-server ${KAFKA_BROKER} \
  --topic uploading_files \
  --partitions 3 \
  --replication-factor 1 \

echo "Available topics:"
kafka-topics --list --bootstrap-server ${KAFKA_BROKER}