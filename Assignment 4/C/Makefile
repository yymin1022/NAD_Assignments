# Chat Makefile
# ID : 20194094
# Name : Yongmin Yoo

CC = gcc
CFLAGS = -Wall -Werror -Wextra

TARGET_C = ChatClient
TARGET_S = ChatServer
SRC_C = ChatClient.c
SRC_S = ChatServer.c
OBJ_C = $(SRC_C:.c=.o)
OBJ_S = $(SRC_S:.c=.o)

all: $(TARGET_C) $(TARGET_S)

$(TARGET_C): $(OBJ_C)
	$(CC) $(CFLAGS) $(OBJ_C) -o $(TARGET_C)

$(TARGET_S): $(OBJ_S)
	$(CC) $(CFLAGS) $(OBJ_S) -o $(TARGET_S)

%.o: %.c
	$(CC) $(CFLAGS) -c $^

clean:
	rm -rf $(OBJ_C) $(OBJ_S)
	rm -rf $(TARGET_C) $(TARGET_S)

re: clean all

.PHONY: $(TARGET_C) $(TARGET_S) clean re test
