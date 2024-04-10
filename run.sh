# !/bin/bash
# run /bin/chmod +x run.sh

go build -o build cmd/main/*.go 
./build -goenv=development -dbname=bookings -dbuser=orololuwa