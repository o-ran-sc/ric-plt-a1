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

# This container uses a 2 stage build!
# Tips and tricks were learned from: https://pythonspeed.com/articles/multi-stage-docker-python/
FROM python:3.7-alpine AS compile-image
# Gevent needs gcc
RUN apk update && apk add gcc musl-dev
# I'd prefer to only do this in the second stage, but then we'd need chown, so easier just to add this here; this is thrown away anyway
ENV A1USER a1user
RUN addgroup -S $A1USER && adduser -S -G $A1USER $A1USER
USER $A1USER
# do the install of a1
# Speed hack; we install gevent FIRST because when building repeatedly (eg during dev) and only changing a1 code, we do not need to keep compiling gevent which takes forever
RUN pip install --upgrade pip && pip install --user gevent
COPY setup.py tox.ini .
COPY a1/ .
RUN pip install --user .

###########
# 2nd stage
FROM python:3.7-alpine
# dir that rmr routing file temp goes into
RUN mkdir -p /opt/route/
# python copy; this basically makes the 2 stage python build work
COPY --from=compile-image /root/.local /root/.local
# copy rmr .sos from the builder image
COPY --from=nexus3.o-ran-sc.org:10004/bldr-alpine3-go:1-rmr1.13.1 /usr/local/lib64/libnng.so /usr/local/lib64/libnng.so
COPY --from=nexus3.o-ran-sc.org:10004/bldr-alpine3-go:1-rmr1.13.1 /usr/local/lib64/librmr_nng.so /usr/local/lib64/librmr_nng.so
# Switch to a non-root user for security reasons. a1 does not currently write into any dirs so no chowns are needed at this time.
ENV A1USER a1user
RUN addgroup -S $A1USER && adduser -S -G $A1USER $A1USER
USER $A1USER
# misc setups
EXPOSE 10000
ENV LD_LIBRARY_PATH /usr/local/lib/:/usr/local/lib64
ENV RMR_SEED_RT /opt/route/local.rt
ENV PYTHONUNBUFFERED 1
# This step is critical; fixes..: WARNING: The script jsonschema is installed in '/root/.local/bin' which is not on PATH. Consider adding this directory to PATH or, if you prefer to suppress this warning, use --no-warn-script-location.
ENV PATH=/root/.local/bin:$PATH

# Run!
CMD run.py
