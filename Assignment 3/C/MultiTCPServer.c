/**
 * TCPServer.c
 * ID : 20194094
 * Name : Yongmin Yoo
 **/

#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <unistd.h>

#define BUF_SIZE 1024
#define SERVER_PORT 24094

int     exit_error(char *err_msg);
void    print_time(time_t *start_time);

int main() {
    int                 client_cnt;
    int                 client_id;
    int                 server_binder;
    int                 server_option = 1;
    int                 server_socket_fd;
    fd_set              client_fds;
    time_t              start_time_data;
    struct tm           start_time;
    struct sockaddr_in  server_addr;

    server_socket_fd = socket(AF_INET, SOCK_STREAM, 0);
    setsockopt(server_socket_fd, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, &server_option, sizeof(server_option));

    bzero(&server_addr, sizeof(server_addr));
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(SERVER_PORT);

    server_binder = bind(server_socket_fd, (struct sockaddr *)&server_addr, sizeof(server_addr));
    if (server_binder != 0)
        return (exit_error("Socket Binding Error"));
    server_binder = listen(server_socket_fd, 128);
    if (server_binder != 0)
        return (exit_error("Socket Listening Error"));
    printf("Server is ready to receive on port %d\n", SERVER_PORT);

    FD_ZERO(&client_fds);
    FD_SET(server_socket_fd, &client_fds);
    client_id = server_socket_fd;

    start_time_data = time(NULL);
    localtime_r(&start_time_data, &start_time);
    client_cnt = 0;
    while (1)
    {
        fd_set          tmp_fds;
        time_t          cur_time_data;
        struct tm       cur_time;
        struct timeval  timeout_val;

        tmp_fds = client_fds;
        timeout_val.tv_sec = 0;
        timeout_val.tv_usec = 500000;

        cur_time_data = time(NULL);
        gmtime_r(&cur_time_data, &cur_time);
        if ((cur_time.tm_sec - start_time.tm_sec) % 10 == 0)
        {
            print_time(&start_time_data);
            printf("Number of clients connected = %d\n", client_cnt);
            usleep(500000);
        }

        if (select(client_id + 1, &tmp_fds, 0, 0, &timeout_val) < 0)
            exit_error("Select Error");
        for (int fd = 0; fd < client_id + 1; fd++)
        {
            if (FD_ISSET(fd, &tmp_fds))
            {
                if (fd == server_socket_fd)
                {
                    int                 client_socket_fd;
                    socklen_t           client_len;
                    struct sockaddr_in  client_addr;

                    client_len = sizeof(client_addr);
                    client_socket_fd = accept(server_socket_fd, (struct sockaddr *)&client_addr, &client_len);

                    FD_SET(client_socket_fd, &client_fds);
                    if (client_id < client_socket_fd)
                        client_id = client_socket_fd;
                    client_cnt++;
                }
                else
                {
                    char    client_req_val[BUF_SIZE];
//                    char    server_res_val[BUF_SIZE];
                    ssize_t client_req_len;
//                    ssize_t server_res_len;

                    client_req_len = read(fd, client_req_val, BUF_SIZE);
                    if (client_req_len == 0)
                    {
                        FD_CLR(fd, &client_fds);
                        close(fd);
                        client_cnt--;
                    }
                    else
                    {
                        write(1, "Client Message : ", 17);
                        write(1, client_req_val, client_req_len);
                        write(1, "\n", 1);
                        write(fd, "Hello", 5);
                    }
                }
            }
        }
    }
    return 0;
}

void    print_time(time_t *start_time)
{
    time_t      time_data;
    struct tm   up_time;

    time_data = time(NULL) - *start_time;
    gmtime_r(&time_data, &up_time);
    printf("[Time: %02d:%02d:%02d] ", up_time.tm_hour, up_time.tm_min, up_time.tm_sec);
}

int exit_error(char *err_msg)
{
    write(2, "Error: ", 7);
    write(2, err_msg, strlen(err_msg));
    write(2, "\n", 1);
    exit(-1);
}