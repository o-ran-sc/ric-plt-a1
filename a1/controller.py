"""
Main a1 controller
"""
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
import json
from flask import Response
from jsonschema import validate
from jsonschema.exceptions import ValidationError
import connexion
from a1 import get_module_logger
from a1 import a1rmr, exceptions, data


logger = get_module_logger(__name__)


def _try_func_return(func):
    """
    generic caller that returns the apporp http response if exceptions are raised
    """
    try:
        return func()
    except (ValidationError, exceptions.PolicyTypeAlreadyExists) as exc:
        logger.exception(exc)
        return "", 400
    except (exceptions.PolicyTypeNotFound, exceptions.PolicyInstanceNotFound) as exc:
        logger.exception(exc)
        return "", 404
    except BaseException as exc:
        # catch all, should never happen...
        logger.exception(exc)
        return Response(status=500)


def _gen_body_to_handler(operation, policy_type_id, policy_instance_id, payload=None):
    """
    used to create the payloads that get sent to downstream policy handlers
    """
    return {
        "operation": operation,
        "policy_type_id": policy_type_id,
        "policy_instance_id": policy_instance_id,
        "payload": payload,
    }


# Healthcheck


def get_healthcheck():
    """
    Handles healthcheck GET
    Currently, this basically checks the server is alive.a1rmr
    """
    return "", 200


# Policy types


def get_all_policy_types():
    """
    Handles GET /a1-p/policytypes
    """
    return _try_func_return(data.get_type_list)


def create_policy_type(policy_type_id):
    """
    Handles PUT /a1-p/policytypes/policy_type_id
    """

    def _put_type_handler(policy_type_id, body):
        data.store_policy_type(policy_type_id, body)
        return "", 201

    body = connexion.request.json
    return _try_func_return(lambda: _put_type_handler(policy_type_id, body))


def get_policy_type(policy_type_id):
    """
    Handles GET /a1-p/policytypes/policy_type_id
    """
    return _try_func_return(lambda: data.get_policy_type(policy_type_id))


def delete_policy_type(policy_type_id):
    """
    Handles DELETE /a1-p/policytypes/policy_type_id
    """
    logger.error(policy_type_id)
    return "", 501


# Policy instances


def get_all_instances_for_type(policy_type_id):
    """
    Handles GET /a1-p/policytypes/policy_type_id/policies
    """
    return _try_func_return(lambda: data.get_instance_list(policy_type_id))


def get_policy_instance(policy_type_id, policy_instance_id):
    """
    Handles GET /a1-p/policytypes/polidyid/policies/policy_instance_id
    """
    # 200 is automatic here
    return _try_func_return(lambda: data.get_policy_instance(policy_type_id, policy_instance_id))


def get_policy_instance_status(policy_type_id, policy_instance_id):
    """
    Handles GET /a1-p/policytypes/polidyid/policies/policy_instance_id/status
    """

    def get_status_handler(policy_type_id, policy_instance_id):
        """
        Pop trough A1s mailbox, insert the latest status updates into the database, and then return the status vector

        NOTE: this is done lazily. Meaning, when someone performs a GET on this API, we pop through a1s mailbox.
        THis may not work in the future if there are "thousands" of policy acknowledgements that hit a1 before this is called,
        because the rmr mailbox may fill. However, in the near term, we do not expect this to happen.
        """
        # check validity to 404 first:
        data.instance_is_valid(policy_type_id, policy_instance_id)

        # pop a1s mailbox, looking for policy notifications
        new_messages = a1rmr.dequeue_all_waiting_messages(21024)

        # try to parse the messages as responses. Drop those that are malformed
        for msg in new_messages:
            # note, we don't use the parameters "policy_type_id, policy_instance" from above here,
            # because we are popping the whole mailbox, which might include other statuses
            pay = json.loads(msg["payload"])
            if "policy_type_id" in pay and "policy_instance_id" in pay and "handler_id" in pay and "status" in pay:
                data.set_policy_instance_status(
                    pay["policy_type_id"], pay["policy_instance_id"], pay["handler_id"], pay["status"]
                )
            else:
                logger.debug("Dropping message")
                logger.debug(pay)

        # return the status vector
        return data.get_policy_instance_statuses(policy_type_id, policy_instance_id)

    return _try_func_return(lambda: get_status_handler(policy_type_id, policy_instance_id))


def create_or_replace_policy_instance(policy_type_id, policy_instance_id):
    """
    Handles PUT /a1-p/policytypes/polidyid/policies/policy_instance_id
    """

    def put_instance_handler(policy_type_id, policy_instance_id, instance):
        """
        Handles policy instance put

        For now, policy_type_id is used as the message type
        """
        #  validate the PUT against the schema
        schema = data.get_policy_type(policy_type_id)["create_schema"]
        validate(instance=instance, schema=schema)

        # store the instance
        data.store_policy_instance(policy_type_id, policy_instance_id, instance)

        # send rmr (best effort)
        body = _gen_body_to_handler("CREATE", policy_type_id, policy_instance_id, payload=instance)
        a1rmr.send(json.dumps(body), message_type=policy_type_id)

        return "", 201

    instance = connexion.request.json
    return _try_func_return(lambda: put_instance_handler(policy_type_id, policy_instance_id, instance))


def delete_policy_instance(policy_type_id, policy_instance_id):
    """
    Handles DELETE /a1-p/policytypes/polidyid/policies/policy_instance_id
    """

    def delete_instance_handler(policy_type_id, policy_instance_id):
        # delete the instance
        # HOLDUP: do we really want to blow this away? this means you can never get statusus again
        data.delete_policy_instance(policy_type_id, policy_instance_id)

        # send rmr (best effort)
        body = _gen_body_to_handler("DELETE", policy_type_id, policy_instance_id)
        a1rmr.send(json.dumps(body), message_type=policy_type_id)

        return "", 204

    # return _try_func_return(lambda: delete_instance_handler(policy_type_id, policy_instance_id))
    return "", 501
