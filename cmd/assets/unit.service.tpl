[Unit]
Description={{ .Executable }} secrets agent systemd service

[Service]
Type=simple
ExecStart={{ .Path }} --config {{ .Config }}

[Install]
WantedBy=multi-user.target
