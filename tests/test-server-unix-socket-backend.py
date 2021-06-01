#!/usr/bin/env python3

"""
Ensures ghostunnel can connect to a unix socket.
"""

from common import LOCALHOST, RootCert, STATUS_PORT, SocketPair, TlsClient, UnixServer, print_ok, run_ghostunnel, terminate

if __name__ == "__main__":
    ghostunnel = None
    try:
        # create certs
        root = RootCert('root')
        root.create_signed_cert('server')
        root.create_signed_cert('client')

        # start ghostunnel
        socket = UnixServer()
        ghostunnel = run_ghostunnel(['server',
                                     '--listen={0}:13001'.format(LOCALHOST),
                                     '--target=unix:{0}'.format(socket.get_socket_path()),
                                     '--keystore=server.p12',
                                     '--cacert=root.crt',
                                     '--allow-ou=client',
                                     '--status={0}:{1}'.format(LOCALHOST,
                                                               STATUS_PORT)])

        # connect with client, confirm that the tunnel is up
        pair = SocketPair(TlsClient('client', 'root', 13001), socket)
        pair.validate_can_send_from_server(
            "hello world", "1: server -> client")
        pair.validate_can_send_from_client(
            "hello world", "1: client -> server")
        pair.validate_closing_server_closes_client(
            "1: server closed -> client closed")

        print_ok("OK")
    finally:
        terminate(ghostunnel)
