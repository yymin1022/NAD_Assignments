import java.io.IOException;
import java.net.*;

public class EasyUDPServer {
    static int SERVER_PORT = 14094;

    static int serverResponseCnt = 0;

    public static void main(String[] args) {
        DatagramSocket serverConnection = initServer();
        if (serverConnection == null) {
            printError("Failed to init server.");
            return;
        }

        Runtime.getRuntime().addShutdownHook(new Thread(() -> closeConnection(serverConnection)));

        while(true){
            try {
                byte[] requestBuffer = new byte[1024];
                DatagramPacket requestPacket = new DatagramPacket(requestBuffer, requestBuffer.length);
                serverConnection.receive(requestPacket);

                String requestData = new String(requestPacket.getData()).trim();
                String requestIP = requestPacket.getAddress().toString().substring(1);
                int requestPort = requestPacket.getPort();

                String responseData = getResponse(requestData.charAt(0),
                        requestData.substring(1),
                        requestIP,
                        requestPort);

                DatagramPacket responsePacket = new DatagramPacket(
                        responseData.getBytes(),
                        responseData.getBytes().length,
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
        return switch (cmd) {
            case '1' -> data.toUpperCase();
            case '2' -> "UpTime";
            case '3' -> String.format("client IP = %s, port = %d", addr, port);
            case '4' -> String.format("requests served = %d", serverResponseCnt);
            default -> "";
        };
    }

    private static void printError(String msg) {
        System.err.printf("Error: %s\n", msg);
    }
}