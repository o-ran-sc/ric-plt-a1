# RIC A1 Mediator

The xApp A1 mediator exposes a generic REST API by which xApps can
receive and send northbound messages.  The A1 mediator will take
the payload from such generic REST messages, validate the payload,
and then communicate the payload to the xApp via RMR messaging.

Please see documentation in the docs/ subdirectory.


### Building RIC A1 Mediator arm64 docker image

docker build -f Dockerfile_alpine_arm64 -t ric-plt-rtmgr:0.6.3  .

NOTE: Requires an alpine builder image from dev repo https://gerrit.o-ran-sc.org/r/admin/repos/it/dev

docker build -f bldr-imgs/bldr-alpine3-rmr/Dockerfile -t bldr-alpine3-rmr:4.5.2 .
