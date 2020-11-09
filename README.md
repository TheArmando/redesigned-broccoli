# redesigned-broccoli

Quickstart:
1. set environment variable as described in env-sample with export MAXMIND_LICENSE_KEY=
2. go run ./cmd/run/main.go


Assumptions
1. This initial iteration of the API will be internal facing only. So for now API security mechanisms to authenicate the caller will not be in scope
2. Maxmind's Geolite2 country blocks distinguish between registered country and represented country. Which country to use is a business requirement, for now we'll include both. If either country is not in the whitelist, it will be rejected
3. For simplicity we'll assume the upstream is providing the country list in ISO 3166-1 alpha-2. 
    a. If that's not the case we can easily lookup whichever format the upstream provides and map them to ISO 3166-1 alpha-2
    b. For simplicity we'll assume the countries are valid. If we have time we'll add validation on the prototype


Microservice workflow:
1. On startup, download bulk csvs and parse them into memory
    a. Create a map of 3166-1 alpha-2 codes to lists of IPNets


IP Address Lookup
We have two csvs of IP address to country blocks. One in IPv4 the the other in IPv6.

1. Identify whether the IP address provided is in IPv4 or in IPv6
2. Lookup the IP address by referencing the whitelist
    a. Lookup each ISO 3166-1 alpha-2 from the upstream whitelist to get the IPNets
    b. Iterate through each IPNet



Resources:
[0] https://www.juniper.net/documentation/en_US/junos/topics/topic-map/security-interface-ipv4-ipv6-protocol.html
[1] https://stackoverflow.com/questions/19882961/go-golang-check-ip-address-in-range
[3] https://github.com/go-chi/chi
[4] https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
