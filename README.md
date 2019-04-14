# ABANDONED - Terraform Verisign MDNS Provider

[![Project Status: Abandoned – Initial development has started, but there has not yet been a stable, usable release; the project has been abandoned and the author(s) do not intend on continuing development.](https://www.repostatus.org/badges/latest/abandoned.svg)](https://www.repostatus.org/#abandoned)

This is a terraform provider for the [Verisign MDNS ReST API](https://mdns.verisign.com/rest/rest-doc/index.html).

## PROJECT ABANDONED

This project was not yet in a usable state as of April, 2019. In early April, my employer was notified that all Verisign Managed DNS contracts have been [acquired by Neustar](https://www.home.neustar/about-us/news-room/press-releases/2018/VerisignSecurityServices) and Neustar is requesting that customers migrate off of the Verisign Managed DNS platform by the end of their current contacts. As Verisign MDNS is going away and Neustar UltraDNS has an [official terraform provider](https://www.terraform.io/docs/providers/ultradns/index.html), I'm discontinuing all work on this provider.

## Using The Provider

### Provider Configuration

The provider requires some explicit configuration. In order to simplify the provider code, a separate provider instance must be configured for every Zone that you want to manage resource records in (this is because terraform stores a simple string unique ID, but Verisign's ReST API paths include the Account ID, Zone Name, and record ID).

```
provider "verisignmdns" {
  token      = "YourApiToken"
  account_id = "YourAccountId"
  zone_name  = "example.com"
}
```

The provider configuration options are as follows:

* ``token`` - Your Verisign MDNS API token. Can also be set with the ``VERISIGN_MDNS_API_TOKEN`` environment variable.
* ``account_id`` - Your Verisign MDNS Account ID. Can also be set with the ``VERISIGN_ACCOUNT_ID`` environment variable.
* ``zone_name`` - The Zone that this provider will manage records in. Can also be set with the ``VERISIGN_ZONE_NAME`` environment variable.
* ``timeout`` - _(optional)_ The ReST API call timeout in seconds. Defaults to 900. Can also be set with the ``VERISIGN_MDNS_TIMEOUT`` environment variable.
* ``debug`` - _(optional)_ Whether or not to enable debug-level logging for this provider. "true" or "false", defaults to "false". Can also be set with the ``VERISIGN_MDNS_DEBUG`` environment variable.
* ``api_url`` - _(optional)_ The base URL to the Verisign MDNS API. Really only useful for acceptance tests of the provider itself. Defaults to ``https://mdns.verisign.com/mdns-web/api/``. Can also be set with the ``VERISIGN_MDNS_API_URL`` environment variable.

### Resources

#### verisignmdns_rr

This resource manages a single Resource Record in the Zone the provider is configured for. This is currently the only resource that the provider supports.

This will create an A record at "foo.example.com" with a value of "1.2.3.4":

```
resource "verisignmdns_rr" "foo" {
  record_name = "foo.example.com"
  record_type = "A"
  record_data = "1.2.3.4"
}
```

The resource supports the following parameters:

* ``record_name`` - the name of the resource record (FQDN) without a trailing dot.
* ``record_type`` - the type of record, i.e. "A", "AAAA", "CNAME", etc.
* ``record_data`` - the value of the record.

__PLEASE NOTE__ that there are currently a few limitations:

* The provider does not currently support TTL fields.
* Only ``record_data`` can be changed in place; changes to ``record_name`` and ``record_type`` require a destroy and replacement. This is a limitation of the Verisign MDNS API, not the terraform provider.

### Importing

This provider currently supports importing existing records. To import them you will need to know the Verisign MDNS resourceRecordId, which can be found from the [Verisign MDNS ReST API](https://mdns.verisign.com/rest/rest-doc/index.html) (or by mousing over the record link in their web UI). Resource record API paths are in the form ``/api/v1/accounts/{accountId}/zones/{zoneName}/rr/{resourceRecordId}``.

```bash
terraform import verisignmdns_rr.foo resourceRecordId
```

## Development

See [Provider Plugins - Terraform by HashiCorp](https://www.terraform.io/docs/plugins/provider.html)

### Building

Clone repository to: `$GOPATH/src/github.com/jantman/terraform-provider-verisignmdns`

```sh
$ mkdir -p $GOPATH/src/github.com/jantman; cd $GOPATH/src/github.com/jantman
$ git clone git@github.com:jantman/terraform-provider-verisignmdns
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/jantman/terraform-provider-verisignmdns
$ go build -o terraform-provider-verisignmdns
```

### Developing

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ go install
...
$ $GOPATH/bin/terraform-provider-verisignmdns
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. Alternatively,
you can run against a best-effort API mock using the instructions below.

```sh
$ make testacc
```

### Mock API

This provider also includes a Python/Flask mocked Verisign MDNS API, since Verisign
doesn't appear to offer a sandbox environment. To start up the mocked API server:

```bash
cd mockapi
python3 -mvenv .venv
source .venv/bin/activate
pip install -r requirements.txt
./apimock.py
```

In another window, ``source scripts/setup_for_dev.sh`` and begin your development.
