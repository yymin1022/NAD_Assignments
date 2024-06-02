# SplitFileServer / Client Makefile
# ID : 20194094
# Name : Yongmin Yoo

JC = javac

TARGET_C = SplitFileClient
TARGET_S = SplitFileServer
OBJ_C = $(TARGET_C).class
OBJ_S = $(TARGET_S).class
SRC_C = $(TARGET_C).java
SRC_S = $(TARGET_S).java

all: $(TARGET_C) $(TARGET_S)

$(TARGET_C):
	$(JC) $(SRC_C)

$(TARGET_S):
	$(JC) $(SRC_S)

clean:
	rm -rf $(OBJ_C) $(OBJ_S)
	rm -rf $(TARGET_C) $(TARGET_S)

re: clean all

.PHONY: $(TARGET_C) $(TARGET_S) clean re test
