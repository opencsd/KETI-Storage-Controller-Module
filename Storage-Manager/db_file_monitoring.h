#pragma once
#include <vector>
#include <unordered_map>
#include <mutex>
#include <fcntl.h>
#include <unistd.h>
#include <iomanip>
#include <algorithm>
#include <condition_variable>
#include <iostream>
#include <fstream>
#include <string>
#include <map>

#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include "rapidjson/prettywriter.h" 

#include "storage-manager.grpc.pb.h"

using StorageEngineInstance::SSTList;
using StorageEngineInstance::DataFileInfo;
using StorageEngineInstance::DataFileInfo_CSD;

using namespace std;
using namespace rapidjson;

class DBFileMonitoring {
    public:
        struct DataBlockHandle {
            string block_index_handle;
            off64_t offset;
            off64_t length;
        };

        struct SSTInfo {
            vector<string> csd_list;
            vector<struct DataBlockHandle> lba_block_list;
        };
        
		static DBFileMonitoring& GetInstance() {
			static DBFileMonitoring DBFileMonitoring;
			return DBFileMonitoring;
		}

        static void InitDBFileMonitoring(){
			GetInstance().initDBFileMonitoring();
		}

        static DataFileInfo GetDataFileInfo(SSTList sstList){
            return GetInstance().getDataFileInfo(sstList);
        }

        static unordered_map<string,struct SSTInfo> GetSSTCSDMap(){
            return GetInstance().getSSTCSDMap();
        }

	private:
		DBFileMonitoring(){};
		DBFileMonitoring(const DBFileMonitoring&);
		~DBFileMonitoring() {};
		DBFileMonitoring& operator=(const DBFileMonitoring&){
			return *this;
		};

        static unordered_map<string,struct SSTInfo> getSSTCSDMap(){
            return GetInstance().sst_info_map_;
        }

        void initDBFileMonitoring();
        DataFileInfo getDataFileInfo(SSTList sstList);

		unordered_map<string,struct SSTInfo> sst_info_map_;
};