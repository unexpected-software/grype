# Developing

There are a few useful things to know before diving into the codebase. This project depends on a few things being available like a vulnerability database, which you might want to create manually instead of retrieving a released version.

## Getting started

### Native Development

After cloning do the following:

1. run `make bootstrap` to download go mod dependencies, create the `/.tmp` dir, and download helper utilities.
2. run `make` to run linting, tests, and other verifications to make certain everything is working alright.

Checkout `make help` to see what other actions you can take.

### Docker Development

This depends on Docker and Docker Compose

1. run `docker-compose build grype` to build the local development container
2. run `docker-compose run --rm grype bash` to enter into the container with all the bootstrapped dependencies installed.
3. run `make` to verify everything is installed and working properly

## Relationship to Syft

Grype uses Syft as a library for all-things related to obtaining and parsing the given scan target (pulling container
images, parsing container images, indexing directories, cataloging packages, etc). Releases of Grype should
always use released versions of Syft (commits that are tagged and show up in the GitHub releases page). However,
continually integrating unreleased Syft changes into Grype incrementally is encouraged
(e.g. `go get github.com/anchore/syft@main`) as long as by the time a release is cut the Syft version is updated
to a released version (e.g. `go get github.com/anchore/syft@v<semantic-version>`).

## Inspecting the database

The currently supported database format is Sqlite3. Install `sqlite3` in your system and ensure that the `sqlite3` executable is available in your path. Ask `grype` about the location of the database, which will be different depending on the operating system:

```
$ go run main.go db status
Location:  /Users/alfredo/Library/Caches/grype/db
Built:  2020-07-31 08:18:29 +0000 UTC
Current DB Version:  1
Require DB Version:  1
Status: Valid
```

The database is located within the XDG_CACHE_HOME path. To verify the database filename, list that path:

```
# OSX-specific path
$ ls -alh  /Users/alfredo/Library/Caches/grype/db
total 445392
drwxr-xr-x  4 alfredo  staff   128B Jul 31 09:27 .
drwxr-xr-x  3 alfredo  staff    96B Jul 31 09:27 ..
-rw-------  1 alfredo  staff   139B Jul 31 09:27 metadata.json
-rw-r--r--  1 alfredo  staff   217M Jul 31 09:27 vulnerability.db
```

Next, open the `vulnerability.db` with `sqlite3`:

```
$ sqlite3 /Users/alfredo/Library/Caches/grype/db/vulnerability.db
```

To make the reporting from Sqlite3 easier to read, enable the following:

```
sqlite> .mode column
sqlite> .headers on
```

List the tables:

```
sqlite> .tables
id                      vulnerability           vulnerability_metadata
```

In this example you retrieve a specific vulnerability from the `nvd` namespace:

```
sqlite> select * from vulnerability where (namespace="nvd" and package_name="libvncserver") limit 1;
id             record_source  package_name  namespace   version_constraint  version_format  cpes                                                         proxy_vulnerabilities
-------------  -------------  ------------  ----------  ------------------  --------------  -----------------------------------------------------------  ---------------------
CVE-2006-2450                 libvncserver  nvd         = 0.7.1             unknown         ["cpe:2.3:a:libvncserver:libvncserver:0.7.1:*:*:*:*:*:*:*"]  []
```
