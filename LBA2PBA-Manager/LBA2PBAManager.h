#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <sys/types.h>
#include <unordered_map>
#include <iostream>
#include <string>
#include <thread>
#include <vector>
#include <sstream>

#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include "rapidjson/prettywriter.h" 

#include <grpcpp/grpcpp.h>
#include "lba2pba.grpc.pb.h"

using namespace rapidjson;
using namespace std;

#define STORAGE_CLUSTER_MASTER_IP "10.0.4.83"
#define LBA2PBA_Manager_Port 40302
#define MAXLINE 256
#define LOGTAG "LBA2PBA Manager"

// trim left 
inline std::string& ltrim(std::string& s, const char* t = " \t\n\r\f\v"){
	s.erase(0, s.find_first_not_of(t));
	return s;
}
// trim right 
inline std::string& rtrim(std::string& s, const char* t = " \t\n\r\f\v"){
	s.erase(s.find_last_not_of(t) + 1);
	return s;
}
// trim left & right 
inline std::string& trim(std::string& s, const char* t = " \t\n\r\f\v"){
	return ltrim(rtrim(s, t), t);
}

inline std::vector<std::string> split(std::string str, char delimiter){
    std::vector<std::string> answer;
    std::stringstream ss(str);
    std::string temp;
 
    while (getline(ss, temp, delimiter)) {
        if(temp == "" || temp == "\n" || temp == "\r\n"){
            continue;
        }
        answer.push_back(temp);
    }

    return answer;
}


inline void KETILOG(std::string id, std::string msg){
    std::cout << "[" << id << "] " << msg << std::endl;
}