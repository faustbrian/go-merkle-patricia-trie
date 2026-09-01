SHELL := /usr/bin/env bash

GOLIB ?= golib

.PHONY: check ci cohesion inventory repository-check specification-check

check:
	$(GOLIB) check --all

ci:
	$(GOLIB) repository check
	$(GOLIB) specification check
	$(GOLIB) cohesion check
	$(GOLIB) check --all

cohesion:
	$(GOLIB) cohesion check

inventory repository-check:
	$(GOLIB) repository check

specification-check:
	$(GOLIB) specification check
