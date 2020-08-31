PROGRAM_DIR="example-program"

mkdir -p "/var/run/$PROGRAM_DIR"

systemctl daemon-reload
systemctl enable example-program.service
# systemctl start example-program.service

exit 0