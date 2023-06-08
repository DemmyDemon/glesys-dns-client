# glesys-dns-client
Updating my DNS records to point to my home LAN, and other GleSYS DNS stuff.

# Dynamic DNS entries

This is a reimplementation of an ancient perl script I wrote years and years ago. Since I'm doing everything else in Go these days, I might as well do this in Go.

The main task is to figure out what my public IP is, and then update the GleSYS DNS records as needed. This means I can access my home network by a valid hostname *even when the IP changes*. It almost never does, but it's a pain in the smurf when it does.

Fixed by code, because that's how I think.

# Certbot mode
There is also a "Certbot mode"

Certbot, from Let's Encrypt, can call a `manual_auth_hook` program/script. This *used to* call a perl script for me, to set the domain's `_acme-challenge` TXT record.  As is now my custom, this is replaced with Go code.

Since I do this by calling the Glesys DNS API, it made all manner of sense to include it here, even if the two "modes" aren't even needed in the same context.
