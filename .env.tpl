# REG_HOST is optional when exactly one registry is saved in ~/.docker/config.json.
# Required when multiple docker logins exist to disambiguate.
# Can also be set at runtime with the -reg flag: registry --reg myhost.com <command>
# REG_HOST=localhost

REG_SCHEME=https

# REG_USER and REG_PASSWORD are optional when docker login has been run for REG_HOST.
# When REG_USER is set it takes priority over docker credential lookup.
# If neither REG_USER nor docker credentials are found, an error is returned.
# REG_USER=me
# REG_PASSWORD=
