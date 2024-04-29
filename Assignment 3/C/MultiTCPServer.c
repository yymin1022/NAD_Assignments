/**
 * TCPServer.c
 * ID : 20194094
 * Name : Yongmin Yoo
 **/

#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <unistd.h>

#define BUF_SIZE 1024
#define SERVER_PORT 24094

char    *get_response(int cmd, char *data, int fd);
int     exit_error(char *err_msg);
void    print_time();
void    str_toupper(char *str);

int         server_response_cnt = 0;
time_t      server_start_time_data;

int main()
{
    int                 client_cnt;
    int                 client_id;
    int                 server_binder;
    int                 server_option = 1;
    int                 server_socket_fd;
    fd_set              client_fds;
    struct sockaddr_in  server_addr;
    struct tm           server_start_time;

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

    server_start_time_data = time(NULL);
    gmtime_r(&server_start_time_data, &server_start_time);
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
        if ((cur_time.tm_sec - server_start_time.tm_sec) % 10 == 0)
        {
            print_time();
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
                    print_time();
                    printf("Client %d connected. Number of clients connected = %d\n", client_socket_fd, client_cnt);
                }
                else
                {
                    char    client_req_val[BUF_SIZE];
                    char    *server_res_val;
                    int     client_req_cmd;
                    ssize_t client_req_len;
                    ssize_t server_res_len;

                    client_req_len = read(fd, client_req_val, BUF_SIZE);
                    if (client_req_len == 0)
                    {
                        FD_CLR(fd, &client_fds);
                        close(fd);
                        client_cnt--;
                        print_time();
                        printf("Client %d disconnected. Number of clients connected = %d\n", fd, client_cnt);
                        continue;
                    }

                    client_req_cmd = client_req_val[0] - '0';
                    server_res_val = get_response(client_req_cmd, client_req_val + 1, fd);
                    server_res_len = strlen(server_res_val);
                    print_time();
                    printf("TCP Connection Request from %s\n", "ADDR");
                    printf("Command %d\n", client_req_cmd);
                    write(fd, server_res_val, server_res_len);
                    server_response_cnt++;
                }
            }
        }
    }
}

char    *get_response(int cmd, char *data, int fd)
{
    char                *res;
    socklen_t           addr_size;
    time_t              cur_time_data;
    struct sockaddr_in  addr_info;
    struct tm           cur_time;

    switch (cmd)
    {
        case 1:
            str_toupper(data);
            return (data);
        case 2:
            cur_time_data = time(NULL) - server_start_time_data;
            gmtime_r(&cur_time_data, &cur_time);
            res = malloc(20 * sizeof(char));
            sprintf(res, "run time = %02d:%02d:%02d", cur_time.tm_hour, cur_time.tm_min, cur_time.tm_sec);
            return (res);
        case 3:
            addr_size = sizeof(struct sockaddr_in);
            getsockname(fd, (struct sockaddr *)&addr_info, &addr_size);
            char client_ip[30];
            char client_port[6];
            strcpy(client_ip, inet_ntoa(addr_info.sin_addr));
            sprintf(client_port, ":%d", addr_info.sin_port);
            return (strdup(strcat(client_ip, client_port)));
        case 4:
            res = malloc(20 * sizeof(char));
            sprintf(res, "requests served = %d", server_response_cnt);
            return (res);
        default:
            return ("");
    }
}

int exit_error(char *err_msg)
{
    write(2, "Error: ", 7);
    write(2, err_msg, strlen(err_msg));
    write(2, "\n", 1);
    exit(-1);
}

void    str_toupper(char *str)
{
    size_t  i;

    i = 0;
    while(str[i])
    {
        if (str[i] >= 'a' && str[i] <= 'z')
            str[i] = str[i] - 'a' + 'A';
        i++;
    }
}

void    print_time()
{
    time_t      time_data;
    struct tm   up_time;

    time_data = time(NULL);
    localtime_r(&time_data, &up_time);
    printf("[Time: %02d:%02d:%02d] ", up_time.tm_hour, up_time.tm_min, up_time.tm_sec);
}