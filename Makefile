TARGETS = unzippa unzippall
PKGNAME = unzippa
ARCH = $$(dpkg --print-architecture)

SHELL := /bin/bash

all: $(TARGETS)

$(TARGETS): %: cmd/%/main.go
	go get -v ./...
	go build -ldflags="-s -w" -v -o $@ $<

flow.png: flow.dot
	dot -Tpng < $< > $@

fixtures/fake.zip:
	python fixtures/fake.py && mv fake.zip fixtures/fake.zip && mv fake.txt fixtures/fake.txt

clean:
	rm -f $(TARGETS)
	rm -f fixtures/fake.zip
	rm -f fixtures/fake.txt
	rm -f $(PKGNAME)_*deb
	rm -f $(PKGNAME)-*rpm

deb: $(TARGETS)
	mkdir -p packaging/deb/$(PKGNAME)/usr/sbin
	cp $(TARGETS) packaging/deb/$(PKGNAME)/usr/sbin
	find packaging/deb/$(PKGNAME)/usr -type d -exec chmod 0755 {} \;
	find packaging/deb/$(PKGNAME)/usr -type f -exec chmod 0644 {} \;
	mkdir -p packaging/deb/$(PKGNAME)/DEBIAN/
	cp packaging/deb/control.$(ARCH) packaging/deb/$(PKGNAME)/DEBIAN/control
	cd packaging/deb && fakeroot dpkg-deb --build $(PKGNAME) .
	mv packaging/deb/$(PKGNAME)_*.deb .

rpm: $(TARGETS)
	mkdir -p $(HOME)/rpmbuild/{BUILD,SOURCES,SPECS,RPMS}
	cp ./packaging/rpm/$(PKGNAME).spec $(HOME)/rpmbuild/SPECS
	cp $(TARGETS) $(HOME)/rpmbuild/BUILD
	./packaging/rpm/buildrpm.sh $(PKGNAME)
	cp $(HOME)/rpmbuild/RPMS/x86_64/$(PKGNAME)*.rpm .
