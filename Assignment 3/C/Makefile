CC = gcc
CFLAGS = -Wall -Werror -Wextra

TARGET = MultiTCPServer
SRC = MultiTCPServer.c
OBJ = $(SRC:.c=.o)

$(TARGET): $(OBJ)
	$(CC) $(CFLAGS) $(OBJ) -o $(TARGET)

$(OBJ):
	$(CC) $(CFLAGS) -c $(SRC)

clean:
	rm -rf *.o
	rm -rf $(TARGET)

re: clean $(TARGET)

.PHONY: $(TARGET) clean re