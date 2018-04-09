SHELL := /bin/bash

unzippa: cmd/unzippa/main.go
	go build -o $@ $<

flow.png: flow.dot
	dot -Tpng < $< > $@

fixtures/fake.zip:
	python fixtures/fake.py && mv fake.zip fixtures/fake.zip && mv fake.txt fixtures/fake.txt

clean:
	rm -f flow.png
	rm -f unzippa
	rm -f fixtures/fake.zip
