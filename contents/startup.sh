echo "Starting Main GO application..."
go run main.go &
sleep 5
echo "Starting Heka..."
/usr/local/heka/bin/hekad -config="/usr/local/etc/heka/hekad.toml" >> /var/log/heka.log 2>&1
