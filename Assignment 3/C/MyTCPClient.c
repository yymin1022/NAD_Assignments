/**
 * TCPClient.c
 * ID : 20194094
 * Name : Yongmin Yoo
 **/

#include <arpa/inet.h>
#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#define BUF_SIZE 1024
#define SERVER_NAME "127.0.0.1"
#define SERVER_PORT 24094

int exit_error(char *err_msg);

int main()
{
    ssize_t             socket_fd;
    struct sockaddr_in  server_addr;

    socket_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (socket_fd < 0)
        return (exit_error("Socket Error"));

    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(SERVER_PORT);
    if (inet_pton(AF_INET, SERVER_NAME, &server_addr.sin_addr) <= 0)
        return (exit_error("inet_pton Error"));
    if (connect(socket_fd, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0)
        return (exit_error("Server Connection Failed"));

    while (1)
    {
        char    client_req_val[BUF_SIZE];
        char    server_res_val[BUF_SIZE];
        ssize_t  client_req_len;
        ssize_t  server_res_len;

        client_req_len = read(0, client_req_val, BUF_SIZE);
        if (client_req_len < 0)
            break;

        send(socket_fd, client_req_val, client_req_len, 0);
        server_res_len = read(socket_fd, server_res_val, BUF_SIZE);
        write(1, server_res_val, server_res_len);
        write(1, "\n", 1);
    }
    return 0;
}

int exit_error(char *err_msg)
{
    write(2, "Error: ", 7);
    write(2, err_msg, strlen(err_msg));
    write(2, "\n", 1);
    return -1;
}