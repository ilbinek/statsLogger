CC=g++

all: main.o kill.o player.o mission.o
#$(CC) -c -o main.o main.cpp
	$(CC) -shared -o Stats_x64.so main.o kill.o player.o mission.o
	rm *.o
#cp Stats_x64.so /home/steam/arma3/

%.o: %.cpp
	$(CC) $(CFLAGS) -fPIC -c $<

