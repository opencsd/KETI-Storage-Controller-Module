#include "WalServer.h"

#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include "rapidjson/prettywriter.h"

#include <iostream>
#include <vector>
#include <sstream>
#include <unordered_map>

#include <sys/socket.h>
#include <netinet/in.h>
#include <sys/types.h>

#define MAXLINE 256

using namespace rapidjson;

extern std::unordered_map<std::string,std::string> indexnum2tbl;

typedef struct unflushedrow_{
    std::string type;
    std::string key;
    std::string value;
} unflushedrow;
typedef std::unordered_map<std::string,unflushedrow> unflushedrows;
typedef std::unordered_map<std::string,unflushedrows> unflushedmap;

unflushedmap my_map;

WalServer::WalServer()
{
    //ctor
}
WalServer::WalServer(utility::string_t url) : m_listener(url)
{
    m_listener.support(methods::GET, std::bind(&WalServer::handle_get, this, std::placeholders::_1));
    m_listener.support(methods::PUT, std::bind(&WalServer::handle_put, this, std::placeholders::_1));
    m_listener.support(methods::POST, std::bind(&WalServer::handle_post, this, std::placeholders::_1));
    m_listener.support(methods::DEL, std::bind(&WalServer::handle_delete, this, std::placeholders::_1));
}
WalServer::~WalServer()
{
    //dtor
}

void WalServer::handle_error(pplx::task<void>& t)
{
    try
    {
        t.get();
    }
    catch(...)
    {
        // Ignore the error, Log it if a logger is available
    }
}


std::vector<std::string> split(std::string str, char delimiter){
    std::vector<std::string> answer;
    std::stringstream ss(str);
    std::string temp;
 
    while (getline(ss, temp, delimiter)) {
        answer.push_back(temp);
    }

    return answer;
}

//
// Get Request 
//
void WalServer::handle_get(http_request message)
{
    auto body_json = message.extract_string();
    std::string json = utility::conversions::to_utf8string(body_json.get());
    
    Document document;
    document.Parse(json.c_str());

    std::string tbl_name = document["tbl_name"].GetString();
    { // read Wal data
        // char cmd[256] = "ldb dump_wal --print_value --walfile=/usr/local/mysql/data/.rocksdb/001348.log";
        char cmd[256] = "ldb dump_wal --print_value --walfile=/root/workspace/keti/WAL-Manager/WAL-Manager_keti/log/050094.log";

        char buf[MAXLINE];
        std::string s_buf;

        FILE *fp=popen(cmd,"r");
        while(fgets(buf, MAXLINE, fp) != NULL){
            s_buf += buf;
            if(!strchr(buf,'\n')){
                continue;
            }
            auto token = split(s_buf,' ');
            for(int i=0; i<token.size();i++){
                unflushedrow row;
                if(token[i] == "PUT(0)"){
                    row.type = token[i];
                    row.key = token[i+2].substr(2);
                    row.value = token[i+4].substr(2);
                    
                    my_map[indexnum2tbl[row.key.substr(0,8)]][row.key] = row;  
                } else if(token[i] == "DELETE(0)") {
                    row.type = token[i];
                    row.key = token[i+2].substr(2);
                    row.value = "";

                    my_map[indexnum2tbl[row.key.substr(0,8)]][row.key] = row;
                }
            }
            s_buf = "";
        }
    }

    std::string rep;

    { // gen resule json
	    Document res_document;
	    res_document.SetObject();
        rapidjson::Document::AllocatorType& allocator = res_document.GetAllocator();

        Value deleteKey(kArrayType);
        Value unflushedRows(kObjectType);
        Value unflushedKeys(kArrayType);
        Value unflushedValues(kArrayType);

        Value str_val(kObjectType);

        auto rows = my_map[tbl_name];
        for(std::pair<std::string,unflushedrow> e : rows){
            auto row = e.second;
            if(row.type == "PUT(0)") {
                str_val.SetString(row.key.c_str(),static_cast<SizeType>(strlen(row.key.c_str())),allocator);
                unflushedKeys.PushBack(str_val,allocator);
                str_val.SetString(row.key.c_str(),static_cast<SizeType>(strlen(row.key.c_str())),allocator);
                deleteKey.PushBack(str_val,allocator);
                str_val.SetString(row.value.c_str(),static_cast<SizeType>(strlen(row.value.c_str())),allocator);
                unflushedValues.PushBack(str_val,allocator);
            } else if(row.type == "DELETE(0)") {
                str_val.SetString(row.key.c_str(),static_cast<SizeType>(strlen(row.key.c_str())),allocator);
                deleteKey.PushBack(str_val,allocator);
            }
        }
        unflushedRows.AddMember("key",unflushedKeys,allocator);
        unflushedRows.AddMember("value",unflushedValues,allocator);
        res_document.AddMember("deleteKey",deleteKey,allocator);
        res_document.AddMember("unflushedRows",unflushedRows,allocator);
        
        rapidjson::StringBuffer strbuf;
        rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(strbuf);
        res_document.Accept(writer);
        rep = std::string(strbuf.GetString());
    }    
    
    message.reply(status_codes::OK,rep);
    return;

};

//
// A POST request
//
void WalServer::handle_post(http_request message)
{
    ucout <<  message.to_string() << endl;
    message.reply(status_codes::NotFound,U("SUPPORT ONLY GET API"));
    return ;
};

//
// A DELETE request
//
void WalServer::handle_delete(http_request message)
{
    ucout <<  message.to_string() << endl;
    message.reply(status_codes::NotFound,U("SUPPORT ONLY GET API"));
    return;
};


//
// A PUT request 
//
void WalServer::handle_put(http_request message)
{
    ucout <<  message.to_string() << endl;
    message.reply(status_codes::NotFound,U("SUPPORT ONLY GET API"));
    return;
};
