del main.o sti.dll
g++  -O2 -c -fpic main.cpp -o main.o
g++  -static-libgcc -static-libstdc++ -shared  -s -o  sti.dll main.o