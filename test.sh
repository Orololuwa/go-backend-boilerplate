# !/bin/bash
# run /bin/chmod +x test.sh #✅ to make the file executable on mac

# go test -v ./... #✅ to get the verbose

# Other test commands includes 👇🏾
# go test -cover ./... #✅ to get the coverage
go test -coverprofile=coverage.out  ./... && go tool cover -html=coverage.out #✅ to get the output as an html file
