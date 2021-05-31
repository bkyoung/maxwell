[Unit]
Description={{ .Executable }} secrets agent systemd service

[Service]
Type=simple
ExecStart={{ .Path }}

[Install]
WantedBy=multi-user.target
