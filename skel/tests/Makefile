.PHONY: test golden golden-root runner

test: runner
	mkdir -p test
	gotestsum --format testname --junitfile test/junit-report.xml \
		-- . --tags=integration -v --skel.no-parallel

golden-root: runner
	go test . --tags=integration -v -test.run '^TestSkelTargets$$/^archetype-app$$' --skel.generate-golden \
		--skel.golden-archive=before/plain-root.txtar --skel.no-parallel

golden: golden-root
	go test . --tags=integration -v --skel.generate-golden --skel.no-parallel

runner:
	go install gotest.tools/gotestsum@latest
