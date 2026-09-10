# REG_HOST is optional when exactly one registry is saved in ~/.docker/config.json.
# Required when multiple docker logins exist to disambiguate.
# Can also be set at runtime with the -reg flag: registry --reg myhost.com <command>
# REG_HOST=localhost

REG_SCHEME=https

# REG_USER and REG_PASSWORD are optional when docker login has been run for REG_HOST.
# When REG_USER is set it takes priority over docker credential lookup.
# If neither REG_USER nor Docker credentials are found, requests are sent
# without an Authorization header so the registry can enforce its own policy.
# REG_USER=me
# REG_PASSWORD=

# For a credential helper, configure ~/.docker/config.json instead of keeping
# an inline base64 auth value. Example for dkr.enercal.nc:
# {
#   "credHelpers": {"dkr.enercal.nc": "pass"}
# }
# The helper binary must be installed as docker-credential-pass. Other helper
# names work the same way (secretservice, osxkeychain, ecr-login, ...).
