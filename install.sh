#!/bin/bash -eu
# tfmodblock install script

UNAME="$(uname -s | tr '[A-Z]' '[a-z]')"

if echo $UNAME | grep -q -vE 'linux|darwin'; then
  echo "$UNAME is not supported" 1>&2
  exit 1
fi

_ARCH="$(uname -m | tr '[A-Z]' '[a-z]')"

case $_ARCH in
  arm64|aarch64)
    ARCH='arm64' ;;
  x64|x86_64|amd64)
    ARCH='x86_64' ;;
  x86)
    ARCH='i386' ;;
  *)
    echo "$_ARCH is not supported" 1>&2
    exit 1 ;;
esac

DEST_PATH='/usr/local/bin'
if [[ $(whoami) != 'root' ]]; then
  DEST_PATH="/home/$(whoami)/.local/bin"
  [[ ! -e $DEST_PATH ]] && mkdir -p $DEST_PATH
fi

TARBALL_URL="https://github.com/tsubasaogawa/tfmodblock/releases/latest/download/tfmodblock_${UNAME}_${ARCH}.tar.gz"
curl -L --silent $TARBALL_URL | tar zx -O tfmodblock > $DEST_PATH/tfmodblock
chmod 755 $DEST_PATH/tfmodblock
echo "Installed tfmodblock to $DEST_PATH"
