build: clean
	gobuild -a -ignore="logflow/logflow_testmain.go"

test: build 
	-gobuild -t -run -ignore="logflow/logflow_testmain.go"
	rm _testmain _testmain.6 _testmain.go

test_log: 
	6g -I . -o logflow.6 logflow/sink.go logflow/source.go logflow/all_test.go
	6g -I . -o logflow_testmain.6 logflow/logflow_testmain.go
	6l -o logflow_testmain -L . logflow_testmain.6
	./logflow_testmain
	rm logflow_testmain logflow.6 logflow_testmain.6

.PHONY : clean
clean :
	-rm -r master/main coord/main agents/agent1
	-rm _testmain _testmain.6 _testmain.go
	-find . -name "*.6" | xargs -I"%s" rm %s
	-find . -name "*.a" | xargs -I"%s" rm %s

