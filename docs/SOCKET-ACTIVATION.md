Socket Activation
=================

Ghostunnel supports socket activation via both systemd (on Linux) and launchd
(on macOS). Socket activation is support for the `--listen` and `--status`
flags, and can be used by passing an address of the form `systemd:<name>` or
`launchd:<name>`, where `<name>` should be the name of the socket as defined in
your systemd/launchd configuration.

For example, a launchd plist to launch ghostunnel in server mode on :8081,
listening for status connections on :8082, and forwarding connections to :8083
could look like this:

```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>com.square.ghostunnel</string>
    <key>ProgramArguments</key>
    <array>
      <string>/usr/bin/ghostunnel</string>
      <string>server</string>
      <string>--keystore</string>
      <string>/etc/ghostunnel/server-keystore.p12</string>
      <string>--cacert</string>
      <string>/etc/ghostunnel/cacert.pem</string>
      <string>--target</string>
      <string>localhost:8083</string>
      <string>--listen</string>
      <string>launchd:Listener</string>
      <string>--status</string>
      <string>launchd:Status</string>
      <string>--allow-cn</string>
      <string>client</string>
    </array>
    <key>StandardOutPath</key>
    <string>/var/log/ghostunnel.out.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/ghostunnel.err.log</string>
    <key>Sockets</key>
    <dict>
      <key>Listener</key>
      <dict>
        <key>SockServiceName</key>
        <string>8081</string>
        <key>SockType</key>
        <string>stream</string>
        <key>SockFamily</key>
        <string>IPv4</string>
      </dict>
      <key>Status</key>
      <dict>
        <key>SockServiceName</key>
        <string>8082</string>
        <key>SockType</key>
        <string>stream</string>
        <key>SockFamily</key>
        <string>IPv4</string>
      </dict>
    </dict>
  </dict>
</plist>
```

Note that in the launchd case *both* `SockType` and `SockFamily` need to be
defined for each socket. If for example the family were to be left out, launchd
would open two sockets (IPv4 and IPv6) for the given key (like `Listener`) and
pass them to ghostunnel which is not currently supported.
