go run main.go &
/usr/local/heka/bin/hekad -config="/usr/local/etc/heka/hekad.toml" > /var/log/heka.log 2>&1 &


