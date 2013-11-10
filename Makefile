test:
	rm -rf test_data/.materials
	mkdir test_data/.materials
	cp test_data/projects test_data/.materials
	go test
