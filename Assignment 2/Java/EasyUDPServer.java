import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;

public class EasyUDPServer {
    static int SERVER_PORT = 14094;

    public static void main(String[] args) throws IOException {
        DatagramSocket ds = new DatagramSocket(SERVER_PORT);

        while(true){
            byte[] data = new byte[1024];
            DatagramPacket dp = new DatagramPacket(data, data.length);

            ds.receive(dp);
            System.out.println(new String(dp.getData()).trim());
        }
    }
}