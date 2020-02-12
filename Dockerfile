# ==================================================================================
#       Copyright (c) 2019 Nokia
#       Copyright (c) 2018-2019 AT&T Intellectual Property.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#          http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
# ==================================================================================
FROM python:3.7-alpine

# dir that rmr routing file temp goes into
RUN mkdir -p /opt/route/

# Gevent needs gcc
RUN apk update && apk add bash gcc musl-dev

# Speed hack; we install gevent here because when building repeatedly (eg during dev) and only changing a1 code,
# we do not need to keep compiling gevent which takes forever
RUN pip install --upgrade pip && pip install gevent

# copies
COPY --from=nexus3.o-ran-sc.org:10004/bldr-alpine3-go:1-rmr1.13.1 /usr/local/lib64/libnng.so /usr/local/lib64/libnng.so
COPY --from=nexus3.o-ran-sc.org:10004/bldr-alpine3-go:1-rmr1.13.1 /usr/local/lib64/librmr_nng.so /usr/local/lib64/librmr_nng.so
COPY a1/ /tmp/a1
COPY setup.py tox.ini /tmp/
WORKDIR /tmp

# do the actual install; this writes into /usr/local, need root
RUN pip install .

# If we got REALLY fancy, we would use a two stage here and find where gevent goes and copy it over.

# Switch to a non-root user for security reasons.
# a1 does not currently write into any dirs so no chowns are needed at this time.
ENV A1USER a1user
RUN addgroup -S $A1USER && adduser -S -G $A1USER $A1USER
USER $A1USER

# misc setups
EXPOSE 10000
ENV LD_LIBRARY_PATH /usr/local/lib/:/usr/local/lib64
ENV RMR_SEED_RT /opt/route/local.rt
# dont buffer logging
ENV PYTHONUNBUFFERED 1

CMD run.py
