# TCP Makefile
# ID : 20194094
# Name : Yongmin Yoo

CC = gcc
CFLAGS = -Wall -Werror -Wextra

TARGET_CLIENT = MyTCPClient
TARGET_SERVER = MultiTCPServer
SRC_CLIENT = MyTCPClient.c
SRC_SERVER = MultiTCPServer.c
OBJ_CLIENT = $(SRC_CLIENT:.c=.o)
OBJ_SERVER = $(SRC_SERVER:.c=.o)

all: $(TARGET_CLIENT) $(TARGET_SERVER)

$(TARGET_CLIENT): $(OBJ_CLIENT)
	$(CC) $(CFLAGS) $(OBJ_CLIENT) -o $(TARGET_CLIENT)

$(TARGET_SERVER): $(OBJ_SERVER)
	$(CC) $(CFLAGS) $(OBJ_SERVER) -o $(TARGET_SERVER)

%.o: %.c
	$(CC) $(CFLAGS) -c $^

clean:
	rm -rf *.o
	rm -rf $(TARGET_CLIENT) $(TARGET_SERVER)

re: clean all

.PHONY: $(TARGET_CLIENT) $(TARGET_SERVER) clean re