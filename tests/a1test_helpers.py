import msgpack
from ricsdl.exceptions import RejectedByBackend, NotConnected, BackendError


class MockSDLWrapper:
    """
    Mock wrapper for SDL that uses a dict so we do not rely on Redis for unit tests.
    Would be really nice if SDL itself came with a "standalone: dictionary" mode for this purpose...
    """

    def __init__(self):
        self.POLICY_DATA = {}

    def set(self, key, value):
        """set a key"""

        # these are for unit testing that the handler works on various SDL errors
        if key == "a1.policy_type.111":
            raise RejectedByBackend()
        if key == "a1.policy_type.112":
            raise NotConnected()
        if key == "a1.policy_type.113":
            raise BackendError()

        self.POLICY_DATA[key] = msgpack.packb(value, use_bin_type=True)

    def get(self, key):
        """get a key"""
        if key in self.POLICY_DATA:
            return msgpack.unpackb(self.POLICY_DATA[key], raw=False)
        return None

    def find_and_get(self, prefix):
        """get all k v pairs that start with prefix"""
        return {k: msgpack.unpackb(v, raw=False) for k, v in self.POLICY_DATA.items() if k.startswith(prefix)}

    def delete(self, key):
        """ delete a key"""
        del self.POLICY_DATA[key]
