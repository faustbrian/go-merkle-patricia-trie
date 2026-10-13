SHELL := /usr/bin/env bash

GOLIB ?= golib

.PHONY: source-docs-check check check-full ci cohesion config inventory repository-check specification-check workflows

check: source-docs-check
	$(GOLIB) check --local --all

check-full: source-docs-check
	$(GOLIB) specification check --online
	$(GOLIB) check --all

ci: source-docs-check config inventory cohesion repository-check specification-check workflows check

cohesion:
	$(GOLIB) cohesion check

config:
	$(GOLIB) config validate

inventory:
	$(GOLIB) inventory

repository-check:
	$(GOLIB) repository check

specification-check:
	$(GOLIB) specification check

workflows:
	$(GOLIB) workflows check

source-docs-check:
	$(GOLIB) source-docs check
