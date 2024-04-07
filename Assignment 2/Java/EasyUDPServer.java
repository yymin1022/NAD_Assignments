import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.SocketException;

public class EasyUDPServer {
    static int SERVER_PORT = 14094;

    static int serverResponseCnt = 0;

    public static void main(String[] args) throws IOException {
        DatagramSocket serverConnection = initServer();
        if (serverConnection == null) {
            printError("Failed to init server.");
            return;
        }

        Runtime.getRuntime().addShutdownHook(new Thread(() -> closeConnection(serverConnection)));

        while(true){
            byte[] requestBuffer = new byte[1024];
            DatagramPacket requestPacket = new DatagramPacket(requestBuffer, requestBuffer.length);
            serverConnection.receive(requestPacket);

            String requestData = new String(requestPacket.getData()).trim();
            String requestIP = requestPacket.getAddress().toString().substring(1);
            int requestPort = requestPacket.getPort();

            String responseData = getResponse(requestData.charAt(0),
                                                requestData.substring(1),
                                                requestIP,
                                                String.valueOf(requestPort));

            DatagramPacket responsePacket = new DatagramPacket(
                    responseData.getBytes(),
                    responseData.getBytes().length,
                    InetAddress.getByName(requestIP),
                    requestPort);
            serverConnection.send(responsePacket);
            serverResponseCnt++;
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

    private static String getResponse(char cmd, String data, String addr, String port) {
        switch (cmd) {
            case '1':
                return data.toUpperCase();
            case '2':
                return "UpTime";
            case '3':
                return "Client IP";
            case '4':
                return String.format("requests served = %d", serverResponseCnt);
        }
        return "";
    }

    private static void printError(String msg) {
        System.err.printf("Error: %s\n", msg);
    }
}