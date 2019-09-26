#!/bin/sh
git clone --branch 1.8.1 https://gerrit.oran-osc.org/r/ric-plt/lib/rmr \
    && cd rmr \
    && mkdir .build; cd .build; cmake .. -DPACK_EXTERNALS=1; sudo make install \
    && cd ../.. \
    && rm -rf rmr
