# tgMessagesRemover

In Telegram you should use this option: 
Group Privacy = Turn off (bot receiving all messages, not only started with slash.)

Ps. Table in your database schema will be created automatically. 

### Daemonize
* sudo nano /lib/systemd/system/messagesremover.service
* insert
```
[Unit]
Description=Messages Remover

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/home/ec2-user/messages_remover_bot
Restart=always
RestartSec=3
ExecStart=/home/ec2-user/messages_remover_bot/messages_remover

[Install]
WantedBy=multi-user.target
```

* sudo systemctl enable messagesremover.service
* sudo service messagesremover start

### Cron
```
* * * * *	cd /home/ec2-user/messages_remover_bot && ./messages_remover cron
```

### Compilation on MacOS

`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o messages_remover`