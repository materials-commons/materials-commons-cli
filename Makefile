test:
	rm -rf test_data/.materials
	rm -rf test_data/corrupted
	mkdir test_data/.materials
	cp test_data/projects test_data/.materials
	mkdir -p test_data/corrupted/.materials
	cp test_data/projects_corrupted test_data/corrupted/.materials/projects
	go test

fmt:
	go fmt
	(cd wsmaterials; go fmt)
	(cd materials; go fmt)
	(cd website; go fmt)
	(cd desktop; go fmt)