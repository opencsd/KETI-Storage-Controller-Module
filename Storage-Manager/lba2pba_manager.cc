#include "lba2pba_manager.h"

void LBA2PBAManager::initLBA2PBAManager(){
	const unordered_map<string,DBFileMonitoring::SSTInfo>& sst_csd_map = DBFileMonitoring::GetSSTCSDMap();

	for(const auto entry: sst_csd_map){
		string sst_name = entry.first;
		const DBFileMonitoring::SSTInfo &sstInfo = entry.second;
		
		CSDPBAMap csd_pba_map;
		
		for(int i=0; i<sstInfo.csd_list.size(); i++){
			string csd_name = sstInfo.csd_list[i];

			off64_t offset_buffer[128][3];

			char cmd[MAXLINE];
			string fd_name = "newport_" + csd_name;

			sprintf(cmd,"filefrag -e /mnt/%s/sst/%s 2> /dev/null",fd_name.c_str(),sst_name.c_str());
			cout << cmd << endl;//file frag 실행

			char buf[MAXLINE];
			int flag = 0;
			FILE *fp=popen(cmd,"r");
			int index = 0;

			while(fgets(buf, MAXLINE, fp) != NULL){
				std::string line(buf);

				auto token = split(line,' ');// tokenize

				for(int i=0;i<token.size();i++){// trim
					token[i] = trim(token[i],".:\n\r");
				}
			
				if(flag){
					offset_buffer[index][0] = (off64_t)4096 * atoi(token[1].c_str()); //lba offset start
					offset_buffer[index][1] = (off64_t)4096 * atoi(token[3].c_str()); //pba offset start
					offset_buffer[index][2] = (off64_t)4096 * atoi(token[5].c_str()); //block length
					index++;
				}

				// start check
				if(token[0] == "ext"){
					flag = 1;
				}

				// last check
				for(auto e:token){
					if(e=="last,eof"){
						flag = 0;
					}
				}
			}

			pclose(fp);

			int tbl_size = index;

			off64_t req_offset;
			off64_t req_length;
			PBAList pba_list;

			for(int j=0; sstInfo.lba_block_list.size(); j++){
				flag = 0;
				req_offset = sstInfo.lba_block_list[j].offset;
				req_length = sstInfo.lba_block_list[j].length;

				for(int l = 0; l < tbl_size; l++){
					DataBlockHandle pba_chunk;

					if(flag || req_offset >= offset_buffer[l][0] && req_offset < offset_buffer[l][0] + offset_buffer[l][2]){
						flag = 1;
						if(req_length > offset_buffer[l][2]){ // X
							// printf("{\n\t\"Offset\" : %ld,\n\t\"Length\" : %ld\n},\n",offset_buffer[k][1] + req_offset - offset_buffer[k][0],offset_buffer[k][2]);
							pba_chunk.offset = offset_buffer[l][1] + req_offset - offset_buffer[l][0];
							pba_chunk.length = offset_buffer[l][2];

							pba_list.pba_block_list.push_back(pba_chunk);

							req_length -= offset_buffer[l][2];
							req_offset += offset_buffer[l][2];
						} else { // here
							// printf("{\n\t\"Offset\" : %ld,\n\t\"Length\" : %ld\n},\n",offset_buffer[k][1] + req_offset - offset_buffer[k][0],req_length);
							pba_chunk.offset = offset_buffer[l][1] + req_offset - offset_buffer[l][0];
							pba_chunk.length = req_length;

							pba_list.pba_block_list.push_back(pba_chunk);

							break;
						}
					}
				}
			}

			csd_pba_map.csd_pba_map_[csd_name] = pba_list;
		}
		sst_pba_map_[sst_name] = csd_pba_map;
	}
}

LBA2PBAResponse LBA2PBAManager::getPBAData(LBA2PBARequest request){
	LBA2PBAResponse lba2pbaResponse;

	int table_index_number = request.table_index_number();
	for(const auto entry: request.sst_csd_map()){
		string sst_name = entry.first;
		string csd_name = entry.second;

		LBA2PBAResponse_PBA pba;

		const LBA2PBAManager::PBAList& pba_list = LBA2PBAManager::GetPBAList(sst_name, csd_name);
		//table index number 확인하는 과정 필요
		for(int i=0; i<pba_list.pba_block_list.size(); i++){
			LBA2PBAResponse_Chunk chunk;

			chunk.set_offset(pba_list.pba_block_list[i].offset);
			chunk.set_length(pba_list.pba_block_list[i].length);

			pba.add_chunks()->CopyFrom(chunk);
		}
		pba.set_csd_id(csd_name);
		lba2pbaResponse.mutable_sst_pba_map()->insert({sst_name, pba});
	}

	return lba2pbaResponse;
}