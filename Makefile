GODIRS = . wsmaterials materials site desktop autoupdate

all: fmt test
	(cd materials; go build materials.go)

test:
	rm -rf test_data/.materials
	rm -rf test_data/corrupted
	mkdir test_data/.materials
	cp test_data/projects test_data/.materials
	cp test_data/.user test_data/.materials
	mkdir -p test_data/corrupted/.materials
	cp test_data/projects_corrupted test_data/corrupted/.materials/projects
	-for d in $(GODIRS); do (cd $$d; go test -v); done

fmt:
	-for d in $(GODIRS); do (cd $$d; go fmt); done
