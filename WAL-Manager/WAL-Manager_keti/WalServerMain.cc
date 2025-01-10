// C++ header include
#include <iostream>
#include <memory>
#include <string>
#include <unordered_map>

// db connect instance include
#include "stdafx.h"
#include "WalServer.h"

std::unique_ptr<WalServer> g_httpHandler;

std::unordered_map<std::string,std::string> indexnum2tbl;

void init_map(){
    indexnum2tbl.insert(std::make_pair("0000011D","lineitem"));
    indexnum2tbl.insert(std::make_pair("00000124","customer"));
    indexnum2tbl.insert(std::make_pair("00000123","nation"));
    indexnum2tbl.insert(std::make_pair("0000011E","orders"));
    indexnum2tbl.insert(std::make_pair("00000120","part"));
    indexnum2tbl.insert(std::make_pair("00000121","partsupp"));
    //indexnum2tbl.insert(std::make_pair("00000122","region"));
    indexnum2tbl.insert(std::make_pair("0000010D","region"));
    indexnum2tbl.insert(std::make_pair("0000011F","supplier"));
}

void on_initialize(const string_t& address)
{
    web::uri_builder uri(address);  

    auto addr = uri.to_uri().to_string();
     g_httpHandler = std::unique_ptr<WalServer>(new WalServer(addr));
     g_httpHandler->open().wait();

    ucout << utility::string_t(U("Listening for requests at: ")) << addr << std::endl;

    init_map();

    return;
}

void on_shutdown()
{
	g_httpHandler->close().wait();
    return;
}

int main(int argc, char *argv[])
{
    // wal request 수신
    // * wal request json 형태
    //  - volumn_id : string - [volume id]
    //  - tbl_name : string - [table name]
    // [table name]과 매핑되는 index num 획득
    // [volume id]-파일 이름(csd로 부터 상대 경로)-[csd id] 매핑 정보를 통해 모든 wal 파일의 절대 경로 획득
    // $ ldb dump_wal --print_value --walfile=[wal 파일의 절대 경로]
    // index num과 매칭되는 PUT(0), DELETE(0) 동작 획득
    // 획득한 PUT(0), DELETE(0) 동작 정보를 사용해 result json 생성
    // 생성한 result json 반환
    utility::string_t port = U("12345");
    if(argc == 2)
    {
        port = argv[1];
    }

    // utility::string_t address = U("http://10.0.5.121:");
        utility::string_t address = U("http://10.0.4.83:");

    address.append(port);

    on_initialize(address);
    std::cout << "Press ENTER to exit." << std::endl;

    std::string line;
    std::getline(std::cin, line);

    on_shutdown();
    return 0;
}