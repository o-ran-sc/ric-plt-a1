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
from flask import Response
import connexion
import json
from jsonschema.exceptions import ValidationError
from a1 import get_module_logger
from a1 import a1rmr, exceptions, utils, data


logger = get_module_logger(__name__)


def _get_policy_definition(policy_type_id):
    # Currently we read the manifest on each call, which would seem to allow updating A1 in place. Revisit this?
    manifest = utils.get_ric_manifest()
    for m in manifest["controls"]:
        if m["name"] == policy_type_id:
            return m

    raise exceptions.PolicyTypeNotFound()


def _get_needed_policy_info(policy_type_id):
    """
    Get the needed info for a policy
    """
    m = _get_policy_definition(policy_type_id)
    return (
        utils.rmr_string_to_int(m["message_receives_rmr_type"]),
        m["message_receives_payload_schema"] if "message_receives_payload_schema" in m else None,
        utils.rmr_string_to_int(m["message_sends_rmr_type"]),
    )


def _get_needed_policy_fetch_info(policy_type_id):
    """
    Get the needed info for fetching a policy state
    """
    m = _get_policy_definition(policy_type_id)
    req_k = "control_state_request_rmr_type"
    ack_k = "control_state_request_reply_rmr_type"
    return (
        utils.rmr_string_to_int(m[req_k]) if req_k in m else None,
        utils.rmr_string_to_int(m[ack_k]) if ack_k in m else None,
    )


def _try_func_return(func):
    """
    generic caller that returns the apporp http response if exceptions are raised
    """
    try:
        return func()
    except ValidationError as exc:
        logger.exception(exc)
        return "", 400
    except exceptions.PolicyTypeNotFound as exc:
        logger.exception(exc)
        return "", 404
    except exceptions.PolicyInstanceNotFound as exc:
        logger.exception(exc)
        return "", 404
    except exceptions.MissingManifest as exc:
        logger.exception(exc)
        return "A1 was unable to find the required RIC manifest. report this!", 500
    except exceptions.MissingRmrString as exc:
        logger.exception(exc)
        return "A1 does not have a mapping for the desired rmr string. report this!", 500
    except BaseException as exc:
        # catch all, should never happen...
        logger.exception(exc)
        return Response(status=500)


def _put_handler(policy_type_id, policy_instance_id, instance):
    """
    Handles policy put

    For now, policy_type_id is used as the message type
    """

    _, schema, _ = _get_needed_policy_info(policy_type_id)

    # TODO: check for 404

    # validate the PUT against the schema, or if there is no shema, make sure the pUT is empty
    if schema:
        utils.validate_json(instance, schema)
    elif instance != {}:
        return "BODY SUPPLIED BUT POLICY HAS NO EXPECTED BODY", 400

    # store the instance
    data.store_policy_instance(policy_type_id, policy_instance_id, instance)

    body = {
        "operation": "CREATE",
        "policy_type_id": policy_type_id,
        "policy_instance_id": policy_instance_id,
        "payload": instance,
    }

    # send rmr (best effort)
    a1rmr.send(json.dumps(body), message_type=policy_type_id)

    return "", 201


# def _get_handler(policy_type_id):
#     """
#     Handles policy GET
#     """
#     mtype_send, mtype_return = _get_needed_policy_fetch_info(policy_type_id)
#
#     if not (mtype_send and mtype_return):
#         return "POLICY DOES NOT SUPPORT FETCHING", 501
#
#     # send rmr, wait for ACK
#     return_payload = a1rmr.send_ack_retry("", message_type=mtype_send, expected_ack_message_type=mtype_return)
#
#     # right now it is assumed that xapps respond with JSON payloads
#     try:
#         return (json.loads(return_payload), 200)
#     except json.decoder.JSONDecodeError:
#         return {"reason": "NOT JSON", "return_payload": return_payload}, 502
#

# Public


def get_handler(policy_type_id):
    """
    Handles policy GET
    """
    return _try_func_return(lambda: _get_handler(policy_type_id))


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
    return "", 501


def create_policy_type(policy_type_id):
    """
    Handles PUT /a1-p/policytypes/policy_type_id
    """
    return "", 501


def get_policy_type(policy_type_id):
    """
    Handles GET /a1-p/policytypes/policy_type_id
    """
    return "", 501


def delete_policy_type(policy_type_id):
    """
    Handles DELETE /a1-p/policytypes/policy_type_id
    """
    return "", 501


# Policy instances


def get_all_instances_for_type(policy_type_id):
    """
    Handles GET /a1-p/policytypes/policy_type_id/policies
    """
    return "", 501


def get_policy_instance(policy_type_id, policy_instance_id):
    """
    Handles GET /a1-p/policytypes/polidyid/policies/policy_instance_id
    """
    return _try_func_return(lambda: data.get_policy_instance(policy_type_id, policy_instance_id))


def get_policy_instance_status(policy_type_id, policy_instance_id):
    """
    Handles GET /a1-p/policytypes/polidyid/policies/policy_instance_id/status
    """

    # TODO: Make a new function that calls rmr and the below, then pass it to trY_func

    return _try_func_return(lambda: data.get_policy_instance_statuses(policy_type_id, policy_instance_id))


def create_or_replace_policy_instance(policy_type_id, policy_instance_id):
    """
    Handles PUT /a1-p/policytypes/polidyid/policies/policy_instance_id
    """
    instance = connexion.request.json
    return _try_func_return(lambda: _put_handler(policy_type_id, policy_instance_id, instance))


def delete_policy_instance(policy_type_id, policy_instance_id):
    """
    Handles DELETE /a1-p/policytypes/polidyid/policies/policy_instance_id
    """
    return "", 501
