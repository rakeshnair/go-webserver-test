go run main.go &
sleep 5
/usr/local/heka/bin/hekad -config="/usr/local/etc/heka/hekad.toml" >> /var/log/heka.log 2>&1
