# Instructions for cross-compiling

Ghostunnel has support for loading private keys from PKCS#11 modules, which
should work with any hardware security module that exposes a PKCS#11 interface.
A working CGO toolchain is required in order to compile with PKCS#11 support
enabled.

One way to cross-compile ghostunnel is with [karalabe/xgo][xgo]. Note that
libtool is a required build dependency, and libltdl needs to be available at
runtime. You can build a static binary to avoid the libltdl runtime dependency
by passing appropriate ldflags to the compiler. 

For example, to build a static 64-bit Windows binary:

    xgo \
      -deps https://ftp.gnu.org/pub/gnu/libtool/libtool-2.4.6.tar.gz \
      -branch master \
      -targets 'windows/amd64' \
      -ldflags '-w -extldflags "-static" -extld x86_64-w64-mingw32-gcc' \
      github.com/ghostunnel/ghostunnel

Ghostunnel ships with a `Makefile.dist` that will cross-compile for
darwin/amd64, linux/amd64, windows/386 and windows/amd64 when asked to build
the `dist` target. Note that [xgo][xgo] (and Docker) must already be installed
to run this. For more info, see also [xgo][xgo]'s README on GitHub.

[xgo]: https://github.com/karalabe/xgo
