# coding: utf-8

"""
Zipfile with 100000 entry, extracting 10000, random order.

    $ time unzip -p fixtures/fake.zip $(cat fixtures/fake.txt | tr '\n' ' ')
    real    0m20.564s
    user    0m19.978s
    sys     0m0.146s

    $ time unzippa -m fixtures/fake.txt fixtures/fake.zip > xxxx

    real    0m0.138s
    user    0m0.136s
    sys     0m0.038s

With unzip, we also get some strange errors like:

    caution: filename not matched:  file-71204.txt

Although the files are in the archive:

    $ unzip -l fixtures/fake.zip | grep file-71204.txt
        23  04-09-2018 15:04   file-71204.txt
"""

import random
import zipfile

N = 100000

with zipfile.ZipFile('fake.zip', 'w') as zf:
    for i in range(N):
        name = 'file-%s.txt' % i
        content = 'hello from file #%s\n' % i
        zf.writestr(name, content, zipfile.ZIP_DEFLATED)

SAMPLE_SIZE = int(N / 10)

with open('fake.txt', 'w') as output:
    for i in range(SAMPLE_SIZE):
        filename = 'file-%s.txt' % random.randint(0, N - 1)
        output.write("%s\n" % filename)
