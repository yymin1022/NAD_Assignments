# TCP Makefile
# ID : 20194094
# Name : Yongmin Yoo

CC = gcc
CFLAGS = -Wall -Werror -Wextra

TARGET = MultiTCPServer
SRC = MultiTCPServer.c
OBJ = $(SRC:.c=.o)

all: $(TARGET)

$(TARGET): $(OBJ)
	$(CC) $(CFLAGS) $(OBJ) -o $(TARGET)

%.o: %.c
	$(CC) $(CFLAGS) -c $^

clean:
	rm -rf *.o
	rm -rf $(TARGET)

re: clean all

test: all
	./MultiTCPServer

.PHONY: $(TARGET) clean re test