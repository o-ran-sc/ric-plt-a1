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
import os
import queue
import time
import json
from threading import Thread
from rmr import rmr, helpers
from a1 import get_module_logger
from a1 import data
from a1.exceptions import PolicyTypeNotFound, PolicyInstanceNotFound

logger = get_module_logger(__name__)


RETRY_TIMES = int(os.environ.get("RMR_RETRY_TIMES", 4))

_RMR_LOOP = None
_RMR_THREAD = None


def _init_rmr():
    """
    init an rmr context
    This gets monkeypatched out for unit testing
    """
    # rmr.RMRFL_MTCALL puts RMR into a multithreaded mode, where a receiving thread populates an
    # internal ring of messages, and receive calls read from that
    # currently the size is 2048 messages, so this is fine for the foreseeable future
    logger.debug("Waiting for rmr to initialize..")
    mrc = rmr.rmr_init(b"4562", rmr.RMR_MAX_RCV_BYTES, rmr.RMRFL_MTCALL)
    while rmr.rmr_ready(mrc) == 0:
        time.sleep(0.5)

    return mrc


def _send(mrc, payload, message_type=0):
    """
    Sends a message up to RETRY_TIMES
    If the message is sent successfully, it returns the transactionid
    Does nothing otherwise
    """
    # TODO: investigate moving this below and allocating the space based on the payload size
    sbuf = rmr.rmr_alloc_msg(mrc, 4096)
    payload = payload if isinstance(payload, bytes) else payload.encode("utf-8")

    # retry RETRY_TIMES to send the message
    for _ in range(0, RETRY_TIMES):
        # setup the send message
        rmr.set_payload_and_length(payload, sbuf)
        rmr.generate_and_set_transaction_id(sbuf)
        sbuf.contents.state = 0
        sbuf.contents.mtype = message_type
        pre_send_summary = rmr.message_summary(sbuf)
        logger.debug("Pre message send summary: %s", pre_send_summary)
        transaction_id = pre_send_summary["transaction id"]  # save the transactionid because we need it later

        # send
        sbuf = rmr.rmr_send_msg(mrc, sbuf)
        post_send_summary = rmr.message_summary(sbuf)
        logger.debug("Post message send summary: %s", rmr.message_summary(sbuf))

        # check success or failure
        if post_send_summary["message state"] == 0 and post_send_summary["message status"] == "RMR_OK":
            # we are good
            logger.debug("Message sent successfully!")
            rmr.rmr_free_msg(sbuf)
            return transaction_id

    # we failed all RETRY_TIMES
    logger.debug("Send failed all %s times, stopping", RETRY_TIMES)
    rmr.rmr_free_msg(sbuf)
    return None


class _RmrLoop:
    """
    class represents an rmr loop meant to be called as a longstanding separate thread
    """

    def __init__(self, init_func_override=None, rcv_func_override=None):
        self._keep_going = True
        self._rcv_func = None
        self._last_ran = time.time()
        self._work_queue = queue.Queue()  # thread safe queue https://docs.python.org/3/library/queue.html

        # get a context
        self._mrc = init_func_override() if init_func_override else _init_rmr()

        # set the receive function
        self._rcv_func = rcv_func_override if rcv_func_override else lambda: helpers.rmr_rcvall_msgs(self._mrc, [21024])

    def last_ran(self):
        """return the unix time of the last completed work loop"""
        return self._last_ran

    def stop(self):
        """sets a flag for the loop to end"""
        self._keep_going = False

    def queue_work(self, work_item):
        """adds work for the loop"""
        self._work_queue.put(work_item)

    def loop(self):
        """
        This loop runs in an a1 thread forever, and has 3 jobs:
        - send out any messages that have to go out (create instance, delete instance)
        - read a1s mailbox and update the status of all instances based on acks from downstream policy handlers
        - clean up the database (eg delete the instance) under certain conditions based on those statuses (NOT DONE YET)
        """
        # loop forever
        logger.debug("Work loop starting")
        while self._keep_going:
            # send out all messages waiting for us
            while not self._work_queue.empty():
                work_item = self._work_queue.get(block=False, timeout=None)
                _send(self._mrc, payload=work_item["payload"], message_type=work_item["msg type"])

            # read our mailbox and update statuses
            updated_instances = set()
            for msg in self._rcv_func():
                try:
                    pay = json.loads(msg["payload"])
                    pti = pay["policy_type_id"]
                    pii = pay["policy_instance_id"]
                    data.set_status(pti, pii, pay["handler_id"], pay["status"])
                    updated_instances.add((pti, pii))
                except (PolicyTypeNotFound, PolicyInstanceNotFound, KeyError, json.decoder.JSONDecodeError):
                    # TODO: in the future we may also have to catch SDL errors
                    logger.debug(("Dropping malformed or non applicable message", msg))

            # for all updated instances, see if we can trigger a delete
            # should be no catch needed here, since the status update would have failed if it was a bad pair
            for ut in updated_instances:
                data.clean_up_instance(ut[0], ut[1])

            # TODO: what's a reasonable sleep time? we don't want to hammer redis too much, and a1 isn't a real time component
            self._last_ran = time.time()
            time.sleep(1)


# Public


def queue_work(item):
    """
    push an item into the work queue
    currently the only type of work is to send out messages
    """
    _RMR_LOOP.queue_work(item)


def start_rmr_thread(init_func_override=None, rcv_func_override=None):
    """
    Start a1s rmr thread
    """
    global _RMR_LOOP
    global _RMR_THREAD
    _RMR_LOOP = _RmrLoop(init_func_override, rcv_func_override)
    _RMR_THREAD = Thread(target=_RMR_LOOP.loop)
    _RMR_THREAD.start()


def stop_rmr_thread():
    """
    stops the rmr thread
    """
    _RMR_LOOP.stop()


def healthcheck_rmr_thread(seconds=30):
    """
    returns a boolean representing whether the rmr loop is healthy, by checking two attributes:
    1. is it running?,
    2. is it stuck in a long (> seconds) loop?
    """
    return _RMR_THREAD.is_alive() and ((time.time() - _RMR_LOOP.last_ran()) < seconds)
