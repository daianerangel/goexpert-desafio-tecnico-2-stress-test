run-app:
	go run main.go test --url http://google.com --requests 100 --concurrency 10

build-image:
	docker build -t load-tester .

run-with-docker-google-test:
	docker run -it load-tester test --url http://google.com --requests 300 --concurrency 20

run-with-docker-fullcycle-test:
	docker run -it load-tester test --url http://fullcycle.com.br --requests 30 --concurrency 20