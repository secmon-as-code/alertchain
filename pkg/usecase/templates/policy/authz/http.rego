package authz.http

deny := false

## Add allow rules if needed
#
# deny := false if {
# 	allow
# }
# allow if {
# 	input.path == "/health"
# }
# allow if {
# 	input.header.Authorization == "Bearer XXX"
# }
