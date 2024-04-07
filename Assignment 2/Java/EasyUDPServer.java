import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.SocketException;

public class EasyUDPServer {
    static int SERVER_PORT = 14094;

    public static void main(String[] args) throws IOException {
        DatagramSocket serverConnection = initServer();
        if (serverConnection == null) {
            printError("Failed to init server.");
            return;
        }

        Runtime.getRuntime().addShutdownHook(new Thread(() -> closeConnection(serverConnection)));

        while(true){
            byte[] data = new byte[1024];
            DatagramPacket dp = new DatagramPacket(data, data.length);

            serverConnection.receive(dp);
            System.out.println(new String(dp.getData()).trim());
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

    private static void printError(String msg) {
        System.err.printf("Error: %s\n", msg);
    }
}