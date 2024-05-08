# Chat Makefile
# ID : 20194094
# Name : Yongmin Yoo

RC = rustc

TARGET_C = ChatClient
TARGET_S = ChatServer
SRC_C = ChatClient.rs
SRC_S = ChatServer.rs

all: $(TARGET_C) $(TARGET_S)

$(TARGET_C):
	$(RC) $(SRC_C)

$(TARGET_S):
	$(RC) $(SRC_S)

clean:
	rm -rf $(TARGET_C) $(TARGET_S)

re: clean all

.PHONY: $(TARGET_C) $(TARGET_S) clean re test
