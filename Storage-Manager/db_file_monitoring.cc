#include "db_file_monitoring.h"

void DBFileMonitoring::initDBFileMonitoring(){
    string json = "";
	std::ifstream openFile("../initdata/db_file_monitoring_init.json");
	if(openFile.is_open() ){
		std::string line;
		while(getline(openFile, line)){
			json += line;
		}
		openFile.close();
	}
	
	//parse json	
	Document document;
	document.Parse(json.c_str());

	Value &SSTList = document["sstList"];

    for(int i=0; i<SSTList.Size(); i++){
        Value &SST = SSTList[i];
        SSTInfo sstInfo;
        string sst_name = SST["sstName"].GetString();

        Value &CSDList = SST["csdList"];
        for(int j=0; j<CSDList.Size(); j++){
            string csd_name = CSDList[j].GetString();
            sstInfo.csd_list.push_back(csd_name);
        }

        Value &LBAList = SST["lbaList"];
        for(int j=0; j<LBAList.Size(); j++){
            Value &BlockHandleObject = LBAList[j];
		    DataBlockHandle data_block_handle;
            data_block_handle.offset = BlockHandleObject["offset"].GetInt();
			data_block_handle.length = BlockHandleObject["length"].GetInt();
            sstInfo.lba_block_list.push_back(data_block_handle);
        }

        sst_info_map_[sst_name] = sstInfo;
    }

    cout << "[DB File Monitoring] Init Done" << endl;
}

DataFileInfo DBFileMonitoring::getDataFileInfo(SSTList sstList){
    DataFileInfo dataFileInfo;

    for(int i=0; i<sstList.sst_list_size(); i++){
        string sst_name = sstList.sst_list(i);
        
        DataFileInfo_CSD csd;
        for(const auto &csd_ : sst_info_map_[sst_name].csd_list){
            csd.add_csd_id(csd_);
        }

        dataFileInfo.mutable_sst_csd_map()->insert({sst_name,csd}); 
    }

    return dataFileInfo;
}