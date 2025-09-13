# ⚙️ setting for bot

rename `.env.example` on `.env` and update settings for you bot

```
TOKEN=7969376429:CBEWj2Bm5-if5G194gT7r2kAQjn3MORyhVE
URL_BOT=https://ysqwp-83-219-149-215.a.free.pinggy.link
DB_DSN=host=localhost user=www dbname=go_tg_bot password=314159 port=5432 sslmode=disable
DEBUG=false
```

## 🌵 webhook in telegram bot

```
go mod init
go mod tidy
```

## 🍐 webhook start pinggy.io

🍎 add url webhook in .env

```
ssh -p 443 -R0:127.0.0.1:8080 -L4300:localhost:4300 free.pinggy.io
```

## 🍏 Run

```
killall -9 go
go run .
```

## 🌶️ nginx

```
go mod init go_tg
go build

systemctl restart nginx
```

## 🍎 systemd

```
sudo nano /etc/systemd/system/go_tg.service
```

add this command
```
[Unit]
Description=go_tg

[Service]
User=www
Group=www
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/home/www/goproject/go_tg/
ExecStart=/home/www/goproject/go_tg/go_tg

[Install]
WantedBy=multi-user.target
```

command for start
```
sudo systemctl start go_tg
sudo systemctl enable go_tg
sudo systemctl status go_tg
sudo systemctl restart go_tg

sudo systemctl stop go_tg
sudo systemctl disable go_tg
```

## 🍎 PostgreSQL

add table
```
sudo -u postgres psql
CREATE DATABASE go_tg;
ALTER DATABASE go_tg OWNER TO www;
```

Important: if no tables are found, you need to make sure that the user is in the correct database — to switch, you can use the command

\c database_name

command in psql for add database_name

```
sudo -u postgres psql
\c go_tg
\dt table_name_users
SELECT * FROM table_name_users;
```