.PHONY: example resume

example:
	go run resume-generator.go

resume:
	go run resume-generator.go resume.yaml output/resume
