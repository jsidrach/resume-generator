BROWSER_BIN=/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome
BROWSER_ADDRESS=127.0.0.1
BROWSER_PORT=9222
BROWSER_PAGE=http://${BROWSER_ADDRESS}:${BROWSER_PORT}
BROWSER_REMOTE=(${BROWSER_BIN} --headless --disable-gpu \
                 --remote-debugging-address="${BROWSER_ADDRESS}" \
								 --remote-debugging-port="${BROWSER_PORT}")
BUILD_COMMAND=go run "resume.go"
CONVERT_COMMAND=convert -density 500

.PHONY: start stop example resume

start:
	{ ${BROWSER_REMOTE} & echo $$! > "remote.PID"; }
	sleep 5

stop:
	kill `cat remote.PID` && rm "remote.PID"

example: start
	${BUILD_COMMAND} -browser="${BROWSER_PAGE}"
	${CONVERT_COMMAND} "output/example.pdf" "output/example.png"
	$(MAKE) stop

resume: start
	${BUILD_COMMAND} -browser="${BROWSER_PAGE}" -resume="resume.yaml" -output="output/resume"
	${CONVERT_COMMAND} "output/resume.pdf" "output/resume.png"
	$(MAKE) stop
