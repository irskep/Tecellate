build: clean
	6g lib/easynet.go
	6g lib/ttypes.go

	gopack crg easynet.a easynet.6
	gopack crg ttypes.a ttypes.6

	mv easynet.a lib
	mv ttypes.a lib

	6g -I "lib/" src/bots/test/test.go
	6l -L "lib/" -o build/test test.6

	6g -I "lib/" src/coord/coord.go
	6l -L "lib/" -o build/coord coord.6

	6g -I "lib/" src/server/tecellate.go src/server/grid.go
	6l -L "lib/" -o build/tecellate tecellate.6

	rm *.6

run:
	./build/coord 127.0.0.1:8002 &
	./build/coord 127.0.0.1:8102 &
	(sleep 1; ./build/tecellate)

kill:
	killall coord

.PHONY : clean
clean :
	-find . -name "*.6" | xargs -I"%s" rm %s
	# -rm -f time _testmain 2> /dev/null
	ls
