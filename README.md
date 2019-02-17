# terraform-provider-verisignmdns

[![Project Status: Concept – Minimal or no implementation has been done yet, or the repository is only intended to be a limited example, demo, or proof-of-concept.](https://www.repostatus.org/badges/latest/concept.svg)](https://www.repostatus.org/#concept)

**Completely unusable WIP/concept code**

This is the very beginning of my first attempt at writing a terraform provider
and first attempt at writing Go: a provider that manages Resource Records (RRs)
in Verisign Managed DNS (MDNS).

I started the very bare skeleton of the provider itself, and a Python/Flask
mocked Verisign MDNS API (since Verisign doesn't appear to have a sandbox).
After some consideration, I decided that my need for this is so infrequent that
it's not worth writing a terraform provider, and I'd just write a quick minimal
script and be on my way. I might come back to it at some point.

For the actual WIP provider skeleton and API Mock skeleton, see [the WIP branch](https://github.com/jantman/terraform-provider-verisignmdns/tree/WIP).

## References

Helpful links if I (or anyone else) picks this up:

* [Verisign MDNS ReST API docs](https://mdns.verisign.com/rest/rest-doc/index.html) - this provider would (at least initially) only work with the ``/api/v1/accounts/{accountId}/zones/{zoneName}/rr`` and ``/api/v1/accounts/{accountId}/zones/{zoneName}/rr/{resourceRecordId}`` endpoints
* [Provider Plugins - Terraform by HashiCorp](https://www.terraform.io/docs/plugins/provider.html)
* Where I left off: "Implement Create" [Writing Custom Providers - Guides - Terraform by HashiCorp](https://www.terraform.io/docs/extend/writing-custom-providers.html#implement-create)
* [Home - Extending Terraform - Terraform by HashiCorp](https://www.terraform.io/docs/extend/index.html)
* [terraform-providers/terraform-provider-template: Terraform template provider](https://github.com/terraform-providers/terraform-provider-template)
* a random, simple example provider: [terraform-provider-arukas](https://github.com/terraform-providers/terraform-provider-arukas)
* a generic ReST API provider: [Mastercard/terraform-provider-restapi: A terraform provider to manage objects in a RESTful API](https://github.com/Mastercard/terraform-provider-restapi)

## Development

See [Provider Plugins - Terraform by HashiCorp](https://www.terraform.io/docs/plugins/provider.html)
