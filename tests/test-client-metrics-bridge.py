#!/usr/bin/env python3

"""
Ensures that metrics bridge submission works.
"""

from common import LOCALHOST, RootCert, STATUS_PORT, print_ok, run_ghostunnel, terminate
import time
import json
import http.server
import threading

received_metrics = None


class FakeMetricsBridgeHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        # pylint: disable=global-statement
        global received_metrics
        print_ok("handling POST to fake bridge")
        length = int(self.headers['Content-Length'])
        received_metrics = json.loads(self.rfile.read(length).decode('utf-8'))


if __name__ == "__main__":
    ghostunnel = None
    try:
        # create certs
        root = RootCert('root')
        root.create_signed_cert('client')

        httpd = http.server.HTTPServer(
            ('localhost', 13080), FakeMetricsBridgeHandler)
        server = threading.Thread(target=httpd.handle_request)
        server.start()

        # start ghostunnel
        ghostunnel = run_ghostunnel(['client',
                                     '--listen={0}:13001'.format(LOCALHOST),
                                     '--target={0}:13002'.format(LOCALHOST),
                                     '--keystore=client.p12',
                                     '--cacert=root.crt',
                                     '--metrics-interval=1s',
                                     '--status={0}:{1}'.format(LOCALHOST,
                                                               STATUS_PORT),
                                     '--metrics-url=http://localhost:13080/post'])

        # wait for metrics to post
        for i in range(0, 10):
            if received_metrics:
                break
            else:
                # wait a little longer...
                time.sleep(1)

        if not received_metrics:
            raise Exception("did not receive metrics from instance")

        if not isinstance(received_metrics, list):
            raise Exception("ghostunnel metrics expected to be JSON list")

        print_ok("OK")
    finally:
        terminate(ghostunnel)
