PROGRAM_DIR="example-program"
if [ "$1" -ge 1 ]; then
  echo "$1"
fi
if [ "$1" = 0 ]; then
  systemctl daemon-reload
  rm -rf "/opt/$PROGRAM_DIR"
  rm -rf "/var/run/$PROGRAM_DIR"
fi
exit 0