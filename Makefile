package_name := sonargo
target_dir := sonar
endpoint := http://127.0.0.1:9000/api
username := admin
password := admin
init-clean:
	rm -f ${target_dir}/*.go
	rm -rf integration_testing
	echo "package $(package_name)" > doc.go

update: init-clean
	go mod tidy
	go run ./cmd/main/main.go -f assets/api.json -n ${package_name}  -o ${target_dir} -e ${endpoint} -logtostderr=true -u ${username} -p ${password}
gen: init-clean
	go mod tidy
	go run ./cmd/main/main.go -f assets/api.json -n ${package_name} -o ${target_dir} -e ${endpoint} -logtostderr=true -u ${username} -p ${password}
