import os
import struct

from celery.utils.log import get_task_logger
from time import time
from librtmp import RTMP, RTMPTimeoutError, PACKET_TYPE_INVOKE, RTMPPacket, PACKET_TYPE_VIDEO, PACKET_TYPE_AUDIO
from librtmp import RTMPError
from librtmp.amf import decode_amf
from librtmp.exceptions import AMFError

logger = get_task_logger(__name__)


class RTMPClient(RTMP):
    def __init__(self, uri, **kwargs):
        super().__init__(uri, **kwargs)

        self.fp = None
        self.recording = False
        self.target = None
        self.bytes_written = 0

    def handle_media(self, packet: RTMPPacket) -> None:

        # Check we have a file resource open
        if not self.fp:
            return

        length = len(packet.body)
        timestamp = packet.timestamp

        self.bytes_written += length

        data = struct.pack('>BBHBHB', packet.type, (length >> 16) & 0xff, length & 0x0ffff, (timestamp >> 16) & 0xff,
                           timestamp & 0x0ffff, (timestamp >> 24) & 0xff) + b'\x00\x00\x00' + packet.body
        data += struct.pack('>I', len(data))

        self.fp.write(data)

        # This will terminate the recording cleanly after 2G
        if self.bytes_written >= 2.0 * pow(1024, 3):
            self.close()

        statvfs = os.statvfs(os.path.dirname(self.target))

        # If we only have one gb of disk space left terminate this recording
        if statvfs.f_frsize * statvfs.f_bfree < pow(1024, 3):
            self.close()

    def record(self, file: str):

        self.target = file
        self.bytes_written = 0

        if not os.path.exists(os.path.dirname(self.target)):
            os.makedirs(os.path.dirname(self.target))

        self.fp = open(self.target, 'wb')
        self.fp.write(b'FLV\x01\x05\x00\x00\x00\x09\x00\x00\x00\x00')

        while self.connected:
            self.process_packets()

        # Close the file resource
        self.fp.close()
        self.fp = None

    def process_packets(self, transaction_id=None, invoked_method=None, timeout=None, debug=False):
        """
       Wait for packets and process them as needed.

       :param debug: bool, Log the packet
       :param transaction_id: int, Wait until the result of this
                              transaction ID is recieved.
       :param invoked_method: int, Wait until this method is invoked
                              by the server.
       :param timeout: int, The time to wait for a result from the server.
                            Note: This is the timeout used by this method only,
                            the connection timeout is still used when reading
                            packets.

       Raises :exc:`RTMPError` on error.
       Raises :exc:`RTMPTimeoutError` on timeout.

       Usage::

         >>> @conn.invoke_handler
         ... def add(x, y):
         ...   return x + y

         >>> @conn.process_packets()

       """

        start = time()

        while self.connected and transaction_id not in self._invoke_results:
            if timeout and (time() - start) >= timeout:
                raise RTMPTimeoutError("Timeout")

            try:
                packet = self.read_packet()
            except RTMPError:
                logger.warning("Failed to read packet")
                continue

            if packet.type == PACKET_TYPE_INVOKE:
                try:
                    decoded = decode_amf(packet.body)
                except AMFError:
                    continue

                try:
                    method, transaction_id_, obj = decoded[:3]
                    args = decoded[3:]
                except ValueError:
                    continue

                if debug:
                    logger.info("Decoded rtmp packet. method: {}, transaction_id: {}, obj: {}, args: {}".format(
                        method,
                        transaction_id_,
                        obj,
                        args
                    ))

                if method == "_result":
                    if len(args) > 0:
                        result = args[0]
                    else:
                        result = None

                    self._invoke_results[transaction_id_] = result
                else:
                    handler = self._invoke_handlers.get(method)
                    if handler:
                        res = handler(*args)
                        if res is not None:
                            self.call("_result", res,
                                      transaction_id=transaction_id_)

                    if method == invoked_method:
                        self._invoke_args[invoked_method] = args
                        break

                if transaction_id_ == 1.0:
                    self._connect_result = packet
                else:
                    self.handle_packet(packet)

            elif packet.type == PACKET_TYPE_VIDEO or packet.type == PACKET_TYPE_AUDIO:
                self.handle_media(packet)
            else:
                self.handle_packet(packet)

        if transaction_id:
            result = self._invoke_results.pop(transaction_id, None)

            return result

        if invoked_method:
            args = self._invoke_args.pop(invoked_method, None)

            return args
