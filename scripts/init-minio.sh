#!/bin/sh

sleep 10

if ! mc ls myminio/user-files >/dev/null 2>&1; then
    mc mb myminio/user-files
    echo "Bucket 'user-files' created."
else
    echo "Bucket 'user-files' already exists."
fi
