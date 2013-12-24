GODIRS = . wsmaterials materials site autoupdate send mcfs
P = github.com/materials-commons/materials

all: fmt test
	(cd materials; go build materials.go)
	(cd mcfs; go build mcfs.go)

test:
	rm -rf test_data/.materials
	rm -rf test_data/corrupted
	rm -rf test_data/conversion
	mkdir -p test_data/.materials/projectdb
	mkdir -p test_data/conversion/.materials
	cp test_data/*.project test_data/.materials/projectdb
	cp test_data/projects test_data/conversion/.materials/projects
	cp test_data/.user test_data/.materials
	mkdir -p /tmp/tproj/a
	touch /tmp/tproj/a/a.txt
	-./runtests.sh
	rm -rf /tmp/tproj

fmt:
	-for d in $(GODIRS); do (cd $$d; go fmt); done
