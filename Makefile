SHELL := /bin/bash

flow.png: flow.dot
	dot -Tpng < $< > $@

clean:
	rm -f flow.png

