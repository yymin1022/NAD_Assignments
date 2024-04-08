import java.io.IOException;
import java.net.*;
import java.time.Duration;
import java.time.LocalDateTime;

public class EasyUDPServer {
    static int SERVER_PORT = 14094;

    static int serverResponseCnt;
    static LocalDateTime serverStartTime = null;

    public static void main(String[] args) {
        DatagramSocket serverConnection = initServer();
        if (serverConnection == null) {
            printError("Failed to init server.");
            return;
        }

        Runtime.getRuntime().addShutdownHook(new Thread(() -> closeConnection(serverConnection)));

        serverResponseCnt = 0;
        serverStartTime = LocalDateTime.now();

        while(true){
            try {
                byte[] requestBuffer = new byte[1024];
                DatagramPacket requestPacket = new DatagramPacket(requestBuffer, requestBuffer.length);
                serverConnection.receive(requestPacket);

                String requestData = new String(requestPacket.getData()).trim();
                String requestIP = requestPacket.getAddress().toString().substring(1);
                int requestPort = requestPacket.getPort();

                System.out.printf("UDP Connection Request from %s:%d\n", requestIP, requestPort);
                System.out.printf("Command %c\n", requestData.charAt(0));
                String responseData = getResponse(requestData.charAt(0),
                        requestData.substring(1),
                        requestIP,
                        requestPort);

                byte[] responseBuffer = responseData.getBytes();
                DatagramPacket responsePacket = new DatagramPacket(
                        responseBuffer, responseBuffer.length,
                        InetAddress.getByName(requestIP),
                        requestPort);
                serverConnection.send(responsePacket);
                serverResponseCnt++;
            } catch (SocketException e) {
                printError("Socket Error.");
                break;
            } catch (UnknownHostException e) {
                printError("Client Error.");
                break;
            } catch (IOException e) {
                printError("IO Error.");
                break;
            }
        }
    }

    private static DatagramSocket initServer() {
        try {
            System.out.printf("Server is ready to receive on port %s\n", SERVER_PORT);
            return new DatagramSocket(SERVER_PORT);
        } catch (SocketException e) {
            return null;
        }
    }

    private static void closeConnection(DatagramSocket conn) {
        System.out.println("\rClosing Server Program...\nBye bye~");
        conn.close();
    }

    private static String getResponse(char cmd, String data, String addr, int port) {
        switch (cmd) {
            case '1':
                return data.toUpperCase();
            case '2':
                LocalDateTime curTime = LocalDateTime.now();
                Duration upTime = Duration.between(serverStartTime, curTime);

                long upTimeSeconds = upTime.getSeconds();
                return String.format("run time = %02d:%02d:%02d",
                                    upTimeSeconds / 3600,
                                    (upTimeSeconds % 3600) / 60,
                                    (upTimeSeconds % 3600) % 60);

            case '3':
                return String.format("client IP = %s, port = %d", addr, port);
            case '4':
                return String.format("requests served = %d", serverResponseCnt);
        }
        return "";
    }

    private static void printError(String msg) {
        System.err.printf("Error: %s\n", msg);
    }
}