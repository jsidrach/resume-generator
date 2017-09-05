BROWSER_REMOTE=(chromium --headless --disable-gpu --remote-debugging-address=127.0.0.1 --remote-debugging-port=9222)

.PHONY: start stop example resume

start:
	{ ${BROWSER_REMOTE} & echo $$! > remote.PID; }
	sleep 5

stop:
	kill `cat remote.PID` && rm remote.PID

example: start
	go run resume.go
	convert -density 500 output/example.pdf output/example.png
	$(MAKE) stop

resume: start
	go run resume.go -resume=resume.yaml -output=output/resume
	convert -density 500 output/resume.pdf output/resume.png
	$(MAKE) stop
