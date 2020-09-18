.PHONY: help install dependencies clean

help:
	@cat $(firstword $(MAKEFILE_LIST))

install: \
	dependencies \
	../docker-compose.override.yml

dependencies:
	type docker > /dev/null

../docker-compose.override.yml: docker-compose.override.yml
	cp $< $@

clean:
	rm ../docker-compose.override.yml
