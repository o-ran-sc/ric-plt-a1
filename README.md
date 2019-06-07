# RIC A1 Mediator

The xApp A1 mediator exposes a generic REST API by which xApps can
receive and sent northbound messages, e.g., for receiving policy
intents (e.g., "isolate UEs from cell 5"), for receiving machine
learning models, or for sending current application state like a
neighbor cell relations table learned from the RAN. The A1 mediator
will take the payload from such generic REST messages, validate the
payload, and then communicate the payload to the xApp via RMR
messaging.

Please see documentation in the docs/ subdirectory.
