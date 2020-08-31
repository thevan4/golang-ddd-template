if [ "$1" -ge 1 ]; then
  systemctl stop example-program.service
fi
if [ "$1" = 0 ]; then
  systemctl stop example-program.service
  systemctl disable example-program.service
fi
exit 0