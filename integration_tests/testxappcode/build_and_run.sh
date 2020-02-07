#!/bin/sh
git clone --branch 1.13.1 https://gerrit.oran-osc.org/r/ric-plt/lib/rmr \
    && cd rmr \
    && mkdir .build; cd .build \
    && echo "<<<installing rmr devel headers>>>" \
    && cmake .. -DDEV_PKG=1; make install \
    && echo "<<< installing rmr .so>>>" \
    && cmake .. -DPACK_EXTERNALS=1; sudo make install \
    && echo "cleanup" \
    && cd ../.. \
    && rm -rf rmr

GO111MODULE=on GO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o example-xapp example-xapp.go

LD_LIBRARY_PATH=/usr/local/lib/:/usr/local/lib64 RMR_SEED_RT=uta_rtg.rt ./example-xapp -f config-file.yaml
