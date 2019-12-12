.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. SPDX-License-Identifier: CC-BY-4.0

A1 Mediator
===========

Code
----
https://gerrit.o-ran-sc.org/r/admin/repos/ric-plt/a1

API
---

You can see the API (OpenAPI3 spec) at ``a1/openapi.yml``. You can also
see the “pretty” version if you run the container at
``http://localhost:10000/ui/``.

Policy Overview
----------------
There are two "object types" associated with policy: policy types and policy instances.

Policy types define the name, description, and most importantly the schema of all instances of that type. Think of policy types as defining a JSON schema for the messages sent from A1 to xapps.

Xapps do not receive policy types from A1; types are used only by A1 to validate instance creation requests. However, xapps must register to receive instances of type ids in their xapp descriptor.

Xapp developers can also create new policy types, though the exact process of where these are stored is still TBD. For practical purposes, when the RIC is running, A1s API needs to be invoked to load the policy types before instances can be created.

Policy instances are concrete instantiations of a policy type. They give concrete values of a policy. There may be many instances of a single type. Whenever a policy instance is created in A1, messages are sent over RMR to all xapps registered for that policy type; see below.

Xapps can "sign up" for multiple policy types using their xapp descriptor.

Xapps are expected to handle multiple simultaneous instances of each type that they are registered for.

Xapps supporting A1


Integrating Xapps with A1
-------------------------

A1 to Xapps
~~~~~~~~~~~
When A1 sends a message to xapps, the schema for messages from A1 to the xapp is defined by ``downstream_message_schema`` in ``docs/a1_xapp_contract_openapi.yaml``

All policy instance requests get sent from A1 using message type 20010

Xapps to A1
~~~~~~~~~~~
There are three scenarios in which Xapps are to send a message to A1:

1. When an xapp receives a CREATE or UPDATE message for a policy instance. Xapps must respond to these requests by sending a message of type 20011 to A1. The schema for that message is defined by ``downstream_notification_schema`` in ``docs/a1_xapp_contract_openapi.yaml``
2. Since policy instances can "deprecate" other instances, there are times when xapps need to asyncronously tell A1 that a policy is no longer active. Same message type and schema. The only difference between case 1 and 2 is that case 1 is a "reply" and case 2 is "unsolicited".
3. Xapps can request A1 to re-send all instances of a type using a query, message 20012. When A1 receives this (TBD HERE, STILL BE WORKED OUT)


Known differences from A1 1.0.0 spec
------------------------------------
This is a list of some of the known differences between the API here and the a1 spec dated 2019.09.30.
In some cases, the spec is deficient and we are "ahead", in other cases this does not yet conform to recent spec changes.

1. [RIC is ahead] There is no notion of policy types in the spec, however this aspect is quite critical for the intended use of the RIC, where many Xapps may implement the same policy, and new Xapps may be created often that define new types.
   Moreover, policy types define the schema for policy instances, and without types, A1 cannot validate whether instances are valid, which the RIC A1m does. The RIC A1m view of things is that, there are a set of schemas, called policy types, and one or more instances of each schema can be created. Instances are validated against types.
The spec currently provides no mechanism for the implementation of A1 to know whether policy [instances] are correct since there is no schema for them.
This difference has the rather large consequence that none of the RIC A1m URLs match the spec.

2. [RIC is ahead] There is a rich status URL in the RIC A1m for policy instances, but this is not in the spec.

3. [RIC is ahead] There is a state machine for when instances are actually deleted from the RIC (at which point all URLs referencing it are a 404); this is a configurable option when deploying the RIC A1m.

4. [CR coming to spec] The spec contains a PATCH for partially updating a policy instance, and creating/deleting multiple instances, however the team agreed to remove this from a later version of the Spec. The RIC A1m does not have this operation.

5. [Spec is ahead] The RIC A1 PUT bodies for policy instances do not exactly conform to the "scope" and "statements" block that the spec defines. They are very close otherwise, however.
   (I would argue some of the spec is redundant; for example "policy [instance] id" is a key inside the PUT body to create an instance, but it is already in the URL.)

6. [Spec is ahead] The RIC A1m does not yet notify external clients when instance statuses change.

7. [Spec is ahead] The spec defines that a query of all policy instances should return the full bodies, however right now the RIC A1m returns a list of IDs (assuming subsequent queries can fetch the bodies).

8. [?] The spec document details some very specific "types", but the RIC A1m allows these to be loaded in (see #1). For example, spec section 4.2.6.2. We believe this should be removed from the spec and rather defined as a type. Xapps can be created that define new types, so the spec will quickly become "stale" if "types" are defined in the spec.
