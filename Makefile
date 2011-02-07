build: clean libs coord testbot master

libs :
	6g lib/easynet.go
	6g lib/ttypes.go

	gopack crg easynet.a easynet.6
	gopack crg ttypes.a ttypes.6

	mv easynet.a lib
	mv ttypes.a lib

coord : libs
	6g -I "lib/" src/coord/coord.go src/coord/protocol.go src/coord/botmotion.go src/coord/types.go
	6l -L "lib/" -o build/coord coord.6

testbot : libs
	6g -I "lib/" src/bots/test/test.go
	6l -L "lib/" -o build/test test.6

master : libs
	6g -I "lib/" src/server/tecellate.go src/server/grid.go
	6l -L "lib/" -o build/tecellate tecellate.6

run: coord testbot master
	./build/coord 127.0.0.1:8002 &
	./build/coord 127.0.0.1:8102 &
	(sleep 0.5; ./build/tecellate testgrid.txt)

fancyrun:
	# For when you want to have the coordinators log in separate windows (use tail)
	./build/coord 127.0.0.1:8002 >> out/coord1.txt &
	./build/coord 127.0.0.1:8102 >> out/coord2.txt &
	(sleep 0.5; ./build/tecellate testgrid.txt)

kill:
	killall coord & killall tecellate & killall test

paper_concept:
	cd papers/eecs423_concept; $(MAKE) build

paper_final:
	cd papers/eecs423_final; $(MAKE) build

.PHONY : clean
clean :
	-find . -name "*.6" | xargs -I"%s" rm %s
