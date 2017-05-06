#!/bin/bash

echo "mode: atomic" > coverage.out

PACKAGES=`govendor list -no-status +local`
EXIT_CODE=0

for PKG in $PACKAGES; do
  echo =-= $PKG

  govendor test -v -coverprofile=profile.out -covermode=atomic $PKG; __EXIT_CODE__=$?

  if [ "$__EXIT_CODE__" -ne "0" ]; then
    EXIT_CODE=$__EXIT_CODE__
  fi

  if [ -f profile.out ]; then
    tail -n +2 profile.out >> coverage.out; rm profile.out
  fi
done

exit $EXIT_CODE
