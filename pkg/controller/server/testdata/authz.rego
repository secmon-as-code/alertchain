package authz.http

default deny = false

deny {
    input.path = "/admin"
}
