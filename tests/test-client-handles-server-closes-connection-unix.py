#!/usr/bin/env python3

"""
Ensures when server disconnects that the client connection also disconnects, with UNIX sockets.
"""

from common import LOCALHOST, RootCert, STATUS_PORT, SocketPair, UnixClient, TlsServer, print_ok, run_ghostunnel, terminate

if __name__ == "__main__":
    ghostunnel = None
    try:
        # create certs
        root = RootCert('root')
        root.create_signed_cert('server')
        root.create_signed_cert('client')

        # start ghostunnel
        client = UnixClient()
        ghostunnel = run_ghostunnel(['client',
                                     '--listen=unix:{0}'.format(client.get_socket_path()),
                                     '--target={0}:13002'.format(LOCALHOST),
                                     '--keystore=client.p12',
                                     '--status={0}:{1}'.format(LOCALHOST,
                                                               STATUS_PORT),
                                     '--cacert=root.crt'])

        # connect with client, confirm that the tunnel is up
        pair = SocketPair(client, TlsServer('server', 'root', 13002))
        pair.validate_can_send_from_server(
            "hello world", "1: server -> client")
        pair.validate_can_send_from_client(
            "hello world", "1: client -> server")
        pair.validate_closing_server_closes_client(
            "1: server closed -> client closed")

        print_ok("OK")
    finally:
        terminate(ghostunnel)
