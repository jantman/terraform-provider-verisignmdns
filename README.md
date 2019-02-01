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
