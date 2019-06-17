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
FROM python:3.7

ADD receiver.py /

# Install RMr
# see https://packagecloud.io/o-ran-sc/master/packages/debian/stretch/rmr_1.0.33_amd64.deb
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/master/packages/debian/stretch/rmr_1.0.33_amd64.deb/download.deb
RUN dpkg -i rmr_1.0.33_amd64.deb

# Install python-rmr
RUN pip install --upgrade pip 
RUN pip install rmr==0.10.1

# rmr setups
RUN mkdir -p /opt/route/
ENV LD_LIBRARY_PATH /usr/local/lib
ENV RMR_SEED_RT /opt/route/local.rt

WORKDIR /
CMD ["python","-u","receiver.py"]
