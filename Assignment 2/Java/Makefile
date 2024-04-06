JC = javac

TARGET_TCP_CLI = EasyTCPClient
TARGET_TCP_SRV = EasyTCPServer
TARGET_UDP_CLI = EasyUDPClient
TARGET_UDP_SRV = EasyUDPServer
TARGET = $(TARGET_TCP_CLI) \
        $(TARGET_TCP_SRV) \
        $(TARGET_UDP_CLI) \
        $(TARGET_UDP_SRV)

SRC_TCP_CLI = $(TARGET_TCP_CLI).java
SRC_TCP_SRV = $(TARGET_TCP_SRV).java
SRC_UDP_CLI = $(TARGET_UDP_CLI).java
SRC_UDP_SRV = $(TARGET_UDP_SRV).java
SRCS = $(SRC_TCP_CLI) \
        $(SRC_TCP_SRV) \
        $(SRC_UDP_CLI) \
        $(SRC_UDP_SRV)

$(TARGET):
	$(JC) $(SRCS)

all: $(TARGET)

clean:
	rm -rf $(SRCS:.java=.class)