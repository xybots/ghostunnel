#!/usr/bin/env python3

"""
Test a client which has disabled authentication. It will not send a TLS Client
Certificate, but should still perform verification.
"""

from common import LOCALHOST, RootCert, STATUS_PORT, SocketPair, TcpClient, TlsServer, print_ok, run_ghostunnel, terminate
import ssl

if __name__ == "__main__":
    ghostunnel = None
    try:
        # create certs
        root = RootCert('root')
        root.create_signed_cert('server1')
        root.create_signed_cert(
            'server2', san="IP:127.0.0.1,IP:::1,DNS:foobar")

        other_root = RootCert('other_root')
        other_root.create_signed_cert('other_server')

        # start ghostunnel
        ghostunnel = run_ghostunnel(['client',
                                     '--listen={0}:13001'.format(LOCALHOST),
                                     '--target=localhost:13002',
                                     '--cacert=root.crt',
                                     '--disable-authentication',
                                     '--status={0}:{1}'.format(LOCALHOST,
                                                               STATUS_PORT)])

        # connect to server1, confirm that the tunnel is up
        pair = SocketPair(TcpClient(13001), TlsServer(
            'server1', 'root', 13002, cert_reqs=ssl.CERT_NONE))
        pair.validate_can_send_from_client(
            "hello world", "1: client -> server")
        pair.validate_can_send_from_server(
            "hello world", "1: server -> client")
        pair.validate_closing_client_closes_server(
            "1: client closed -> server closed")

        # connect to other_server, confirm that the tunnel isn't up
        try:
            pair = SocketPair(TcpClient(13001), TlsServer(
                'other_server', 'other_root', 13002, cert_reqs=ssl.CERT_NONE))
            raise Exception('failed to reject other_server')
        except ssl.SSLError:
            print_ok("other_server with unknown CA correctly rejected")

        # connect to server2, confirm that the tunnel isn't up
        try:
            pair = SocketPair(TcpClient(13001), TlsServer(
                'server2', 'root', 13002, cert_reqs=ssl.CERT_NONE))
            raise Exception('failed to reject serve2')
        except ssl.SSLError:
            print_ok("server2 with incorrect CN correctly rejected")

        print_ok("OK")
    finally:
        terminate(ghostunnel)
