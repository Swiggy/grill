PROJECTNAME := grill

.PHONY: setup
setup:
	git config --global url.git@bitbucket.org:.insteadOf https://bitbucket.org

.PHONY: coverage
coverage: setup
	go test -count=1 -covermode=atomic -coverprofile=coverage.out ./...

clean:
	rm -rf coverage*.out



