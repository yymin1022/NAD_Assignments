import java.io.*;
import java.net.Socket;
import java.nio.file.*;

public class SplitFileClient {
    private static final String SERVER_ADDRESS_1 = "localhost";
    private static final int SERVER_PORT_1 = 44094;
    private static final String SERVER_ADDRESS_2 = "localhost";
    private static final int SERVER_PORT_2 = 54094;

    public static void main(String[] args) {
        if (args.length != 2) {
            System.out.println("Usage: java SplitFileClient <put|get> <filename>");
            return;
        }

        String cmd = args[0];
        String filename = args[1];

        try {
            if ("put".equalsIgnoreCase(cmd)) {
                String filePart1 = splitFile(filename, 1);
                String filePart2 = splitFile(filename, 2);

                sendPart(filename, filePart1, SERVER_ADDRESS_1, SERVER_PORT_1);
                Files.delete(Paths.get(filePart1));

                sendPart(filename, filePart2, SERVER_ADDRESS_2, SERVER_PORT_2);
                Files.delete(Paths.get(filePart2));

                System.out.println("File successfully split and sent to servers.");
            } else if ("get".equalsIgnoreCase(cmd)) {
                String filePart1 = getPart(filename, SERVER_ADDRESS_1, SERVER_PORT_1, 1);
                String filePart2 = getPart(filename, SERVER_ADDRESS_2, SERVER_PORT_2, 2);

                String outputFilename = getMergedFilename(filename);
                mergeFiles(filePart1, filePart2, outputFilename);

                Files.delete(Paths.get(filePart1));
                Files.delete(Paths.get(filePart2));

                System.out.println("File successfully retrieved and merged: " + outputFilename);
            } else {
                System.out.println("Usage: java SplitFileClient <put|get> <filename>");
            }
        } catch (IOException e) {
            System.err.println("Error: " + e.getMessage());
            e.printStackTrace();
        }
    }

    private static void sendPart(String filename, String partFilename, String serverAddress, int serverPort) throws IOException {
        try (Socket socket = new Socket(serverAddress, serverPort);
             OutputStream socketOutput = socket.getOutputStream();
             InputStream socketInput = socket.getInputStream();
             FileInputStream partFile = new FileInputStream(partFilename)) {

            socketOutput.write(("PUT:" + filename).getBytes());
            socketOutput.flush();

            BufferedReader reader = new BufferedReader(new InputStreamReader(socketInput));
            String response = reader.readLine();
            if (!"READY".equals(response.trim())) {
                throw new IOException("Server not ready for file content");
            }

            byte[] buffer = new byte[1024];
            int bytesRead;
            while ((bytesRead = partFile.read(buffer)) != -1) {
                socketOutput.write(buffer, 0, bytesRead);
            }
            socketOutput.flush();
        }
    }

    private static String getPart(String filename, String serverAddress, int serverPort, int partNum) throws IOException {
        String partFilename = filename + "-part" + partNum + ".tmp";
        try (Socket socket = new Socket(serverAddress, serverPort);
             OutputStream socketOutput = socket.getOutputStream();
             InputStream socketInput = socket.getInputStream();
             FileOutputStream partFile = new FileOutputStream(partFilename)) {

            socketOutput.write(("GET:" + filename).getBytes());
            socketOutput.flush();

            byte[] buffer = new byte[1025];
            int bytesRead;
            while ((bytesRead = socketInput.read(buffer)) != -1) {
                if (new String(buffer, 0, bytesRead).contains("EOF")) {
                    int eofIndex = new String(buffer, 0, bytesRead).indexOf("EOF");
                    if (eofIndex > 1) {
                        partFile.write(buffer, 1, eofIndex - 1);
                    }
                    break;
                } else {
                    partFile.write(buffer, 1, bytesRead - 1);
                }
            }
        }
        return partFilename;
    }

    private static String splitFile(String filename, int partNum) throws IOException {
        File inputFile = new File(filename);
        String partFilename = filename + "-part" + partNum + ".tmp";
        try (FileInputStream inputFileStream = new FileInputStream(inputFile);
             FileOutputStream partFileStream = new FileOutputStream(partFilename)) {

            byte[] buffer = new byte[1];
            boolean writeToPart1 = (partNum == 1);

            while (inputFileStream.read(buffer) != -1) {
                if (writeToPart1) {
                    partFileStream.write(buffer);
                }
                writeToPart1 = !writeToPart1;
            }
        }
        return partFilename;
    }

    private static void mergeFiles(String part1, String part2, String outputFile) throws IOException {
        try (FileOutputStream outFile = new FileOutputStream(outputFile);
             FileInputStream partFile1 = new FileInputStream(part1);
             FileInputStream partFile2 = new FileInputStream(part2)) {

            byte[] buffer1 = new byte[1];
            byte[] buffer2 = new byte[1];
            int bytesRead1, bytesRead2;

            while (true) {
                bytesRead1 = partFile1.read(buffer1);
                bytesRead2 = partFile2.read(buffer2);

                if (bytesRead1 == -1 && bytesRead2 == -1) {
                    break;
                }

                if (bytesRead1 != -1) {
                    outFile.write(buffer1, 0, bytesRead1);
                }
                if (bytesRead2 != -1) {
                    outFile.write(buffer2, 0, bytesRead2);
                }
            }
        }
    }

    private static String getMergedFilename(String filename) {
        int extIndex = filename.lastIndexOf('.');
        if (extIndex == -1) {
            return filename + "-merged";
        }
        return filename.substring(0, extIndex) + "-merged" + filename.substring(extIndex);
    }
}
