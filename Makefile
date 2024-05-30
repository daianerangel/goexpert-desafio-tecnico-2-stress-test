run-app:
	go run main.go test --url http://google.com --requests 100 --concurrency 10

run-with-docker:
	docker build -t load-tester .
	docker run -it load-tester test --url http://google.com --requests 100 --concurrency 10