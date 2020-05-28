#!/usr/bin/env bash

# go modules are downloaded readonly, which makes them tricky to remove.
# If we leave them behind then Jenkins won't be able to delete them.

find .gomodules -type f -exec chmod 644 {} \;
find .gomodules -type d -exec chmod 755 {} \;
rm -rf .gomodules
