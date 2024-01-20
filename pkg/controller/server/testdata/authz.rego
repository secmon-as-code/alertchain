package authz.http

default deny = true

deny := false { allow }

allow {
	input.path == "/health"
}

jwks_request(url) := http.send({
    "url": url,
    "method": "GET",
    "force_cache": true,
    "force_cache_duration_seconds": 3600 # Cache response for an hour
}).raw_body

# Allow request to /api if the request contains a valid JWT token
allow {
    # Check path
	startswith(input.path, "/alert/")

    # Extract token from Authorization header
    authHdr := input.header["Authorization"]
    count(authHdr) == 1
    authHdrValues := split(authHdr[0], " ")
    count(authHdrValues) == 2
    lower(authHdrValues[0]) == "bearer"
    token := authHdrValues[1]

    # Get JWKS of google
    jwks := jwks_request("https://www.googleapis.com/oauth2/v3/certs")

    # Verify token
    io.jwt.verify_rs256(token, jwks)
    claims := io.jwt.decode(token)
    claims[1]["iss"] == "https://accounts.google.com"
    claims[1]["email"] == "__GOOGLE_CLOUD_ACCOUNT_EMAIL__"
    time.now_ns() / (1000 * 1000 * 1000) < claims[1]["exp"]
}

allow {
    input.path == "/alert/raw/cloudstrike_hawk"

    sigHdr := input.header["X-Cs-Primary-Signature"]
    count(sigHdr) == 1
    sig := hex.encode(base64.decode(sigHdr[0]))

    tsHdr := input.header["X-Cs-Delivery-Timestamp"]
    count(tsHdr) == 1
    ts := tsHdr[0]

    data := concat("", [input.body, ts])
    sign := crypto.hmac.sha256(data, input.env.CLOUDSTRIKE_HAWK_KEY)
    crypto.hmac.equal(sign, sig)
}
