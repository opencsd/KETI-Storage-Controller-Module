#pragma once
#include <iostream>
#include <unordered_map>
#include <sstream>

#include <grpcpp/grpcpp.h>
#include "storage-manager.grpc.pb.h"

#include "db_file_monitoring.h"

#define MAXLINE 256

using StorageEngineInstance::StorageManager;
using StorageEngineInstance::LBA2PBARequest;
using StorageEngineInstance::LBA2PBAResponse;
using StorageEngineInstance::LBA2PBAResponse_Chunk;
using StorageEngineInstance::LBA2PBAResponse_PBA;

using namespace std;

class LBA2PBAManager {
	public:
		struct DataBlockHandle {
            string block_index_handle;
            off64_t offset;
            off64_t length;
        };

        struct PBAList {
            vector<struct DataBlockHandle> pba_block_list;
        };

		struct CSDPBAMap {
			unordered_map<string, struct PBAList> csd_pba_map_;
		};
		
		static LBA2PBAManager& GetInstance() {
			static LBA2PBAManager LBA2PBAManager;
			return LBA2PBAManager;
		}

		static void InitLBA2PBAManager(){
			GetInstance().initLBA2PBAManager();
		}

		static LBA2PBAResponse GetPBAData(LBA2PBARequest request){
            return GetInstance().getPBAData(request);
        }

		static PBAList GetPBAList(string sst_name, string csd_name){
			return GetInstance().getPBAList(sst_name, csd_name);
		}

	private:
		LBA2PBAManager(){};
		LBA2PBAManager(const LBA2PBAManager&);
		~LBA2PBAManager() {};
		LBA2PBAManager& operator=(const LBA2PBAManager&){
			return *this;
		};

		PBAList getPBAList(string sst_name, string csd_name){
			return GetInstance().sst_pba_map_[sst_name].csd_pba_map_[csd_name];
		}

		void initLBA2PBAManager();
		LBA2PBAResponse getPBAData(LBA2PBARequest request);

		unordered_map<string,struct CSDPBAMap> sst_pba_map_;
};

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