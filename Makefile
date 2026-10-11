SHELL := /usr/bin/env bash

GOLIB ?= golib

.PHONY: check check-full ci cohesion config inventory repository-check specification-check workflows

check:
	$(GOLIB) check --local --all

check-full:
	$(GOLIB) specification check --online
	$(GOLIB) check --all

ci: config inventory cohesion repository-check specification-check workflows check

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
