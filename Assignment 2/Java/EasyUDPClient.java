import java.io.*;
import java.net.*;
import java.nio.charset.StandardCharsets;

public class EasyUDPClient {
    static String SERVER_NAME = "localhost";
    static int SERVER_PORT = 14094;

    public static void main(String[] args) {
        DatagramSocket serverConnection = makeConnection();
        if (serverConnection == null) {
            printError("Failed to init client.");
            return;
        }

        try {
            serverConnection.setSoTimeout(5000);
            while (true) {
                printMenu();

                int cmd = readCommand();
                String text = "";
                if (cmd == 0) {
                    continue;
                } else if (cmd == 5) {
                    break;
                } else if (cmd == 1) {
                    System.out.print("Input lowercase sentence: ");
                    BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));
                    text = reader.readLine();

                    if (text.length() >= 1024) {
                        printError("Text too long.");
                        continue;
                    }
                }

                byte[] msgData = (cmd + text).getBytes(StandardCharsets.UTF_8);
                DatagramPacket requestData = new DatagramPacket(
                        msgData, msgData.length,
                        InetAddress.getByName(SERVER_NAME), SERVER_PORT);
                serverConnection.send(requestData);

                byte[] responseBuffer = new byte[1024];
                DatagramPacket responseData = new DatagramPacket(responseBuffer, responseBuffer.length);
                serverConnection.receive(responseData);

                System.out.printf("\nReply from server: %s\n", new String(responseData.getData()).trim());
            }
        } catch (SocketTimeoutException e) {
            printError("Server Timeout.");
            closeConnection(serverConnection);
            return;
        } catch (IOException e) {
            printError("IO Error.");
            closeConnection(serverConnection);
            return;
        }

        closeConnection(serverConnection);
    }

    private static DatagramSocket makeConnection() {
        try{
            return new DatagramSocket();
        }catch (IOException e){
            return null;
        }
    }

    private static void closeConnection(DatagramSocket conn) {
        System.out.println("\rClosing Client Program...\nBye bye~");
        conn.close();
    }

    private static void printMenu() {
        System.out.println();
        System.out.println("< Select Menu. >");
        System.out.println("1) Convert Text to UPPER-case Letters");
        System.out.println("2) Get Server Uptime");
        System.out.println("3) Get Client IP / Port");
        System.out.println("4) Get Count of Requests Server Got");
        System.out.println("5) Exit Client");
    }

    private static int readCommand() {
        int cmd;
        String input;

        try {
            System.out.print("Input Command: ");
            BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));
            input = reader.readLine();
            cmd = Integer.parseInt(input);
        } catch (NumberFormatException e) {
            printError("Invalid Command.");
            return 0;
        } catch (IOException e) {
            printError("System Error.");
            return 0;
        }


        if (cmd < 1 || cmd > 5) {
            printError("Invalid Command.");
            return 0;
        }

        return cmd;
    }

    private static void printError(String msg) {
        System.err.printf("Error: %s\n", msg);
    }
}