# shrink
# See LICENSE for copyright and license details.
.POSIX:

PREFIX ?= /usr/local
GO ?= go
GOFLAGS ?=
RM ?= rm -f

all: shrink

shrink:
	$(GO) build $(GOFLAGS)

clean:
	$(RM) shrink

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f shrink $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/shrink

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/bin/shrink

.DEFAULT_GOAL := all

.PHONY: all shrink clean install uninstall
