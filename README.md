# Overview

Simple oppinionated go library to help with common tasks when implementing REST servers and clients.


## Config

Go services run in different environments: in production and on developer's
machine, inside a docker container or packaged as an RPM and deployed to
the server to run under systemd, or even installed using NixOS flake.

For different environments it is usually convenient to have service configured
differently.

### Config format

This module allows defining service configuration in one place declaratively
with struct tags and load it using different methods depending on current
needs.

The same configuration can be loaded:
 * from .toml file
 * from .env file
 * from environment variables

### Validation

Library allows you to specify validation code for different parts of
your configuration consistently.
