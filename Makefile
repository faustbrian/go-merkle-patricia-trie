SHELL := /usr/bin/env bash

GOLIB ?= golib

.PHONY: check ci inventory repository-check specification-check

check:
	$(GOLIB) check --all

ci:
	$(GOLIB) repository check
	$(GOLIB) specification check
	$(GOLIB) check --all

inventory repository-check:
	$(GOLIB) repository check

specification-check:
	$(GOLIB) specification check
