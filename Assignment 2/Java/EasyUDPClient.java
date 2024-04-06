import java.io.IOException;
import java.net.*;

public class EasyUDPClient {
    static String SERVER_NAME = "localhost";
    static int SERVER_PORT = 14094;

    public static void main(String[] args) throws IOException {
        DatagramSocket ds = new DatagramSocket();

        String msg = "Hello, World!";
        DatagramPacket dp = new DatagramPacket(
                msg.getBytes(), msg.getBytes().length, InetAddress.getByName(SERVER_NAME), SERVER_PORT);
        ds.send(dp);
    }
}