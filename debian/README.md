# autodeb packaging

This repository contains the Debian packaging for autodeb.

It is used to produce two types of packages:
 - Official Debian packages (uploaded to the Debian archive)
 - Unofficial Debian packages (published on gitlab pages)

## Unofficial Debian packages

Unofficial Debian packages are built on gitlab and published on gitlab pages:
 - https://autodeb-team.pages.debian.net/autodeb-packaging/autodeb-server_latest.deb
 - https://autodeb-team.pages.debian.net/autodeb-packaging/autodeb-worker_latest.deb

To trigger a build for the unofficial packages, push a git tag with the name `pages/<version>`.
More information on the build process can be found in `debian/.gitlab-ci.yml`.

## Official Debian packages

Autodeb is not in the Debian archive yet. However, when it is, uploads will be tagged with `debian/<version>`.
