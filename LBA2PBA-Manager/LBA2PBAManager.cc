#include "LBA2PBAManager.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::ServerReaderWriter;
using grpc::Status;

using StorageEngineInstance::StorageManager;
using StorageEngineInstance::LBARequest;
using StorageEngineInstance::LBARequest_SST;
using StorageEngineInstance::PBAResponse;
using StorageEngineInstance::PBAResponse_SST;
using StorageEngineInstance::TableBlock;
using StorageEngineInstance::Chunks;
using StorageEngineInstance::Chunk;

PBAResponse RunLBA2PBA(LBARequest request){
	PBAResponse response;

	off64_t offset_buffer[128][3];

	for(const auto sst: request.sst_list()){	
		string sst_name = sst.first;
		KETILOG("LBA2PBA Manager","File Name: "+sst_name);

		PBAResponse_SST response_sst;

		for(int i=0; i<sst.second.csd_list_size(); i++){
			//do hdparm
			char cmd[256];
			string csd_id = sst.second.csd_list(i);
			string fdName = "newport_" + csd_id; 

			sprintf(cmd,"filefrag -e /mnt/%s/sst/%s 2> /dev/null",fdName.c_str(),sst_name.c_str());
			cout << cmd << endl;//file frag 실행
			KETILOG("LBA2PBA Manager","CSD ID: "+csd_id);
			
			char buf[MAXLINE];
			int flag = 0;
			FILE *fp=popen(cmd,"r");
			int index = 0;
			while(fgets(buf, MAXLINE, fp) != NULL){
				std::string line(buf);
				// tokenize
				auto token = split(line,' ');
				// trim
				for(int i=0;i<token.size();i++){
					token[i] = trim(token[i],".:\n\r");
				}
			
				if(flag){
					// for(auto e:token){
					// 	printf("%s\n",e.c_str());
					// }
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
						// end
						// printf("end\n");
						flag = 0;
					}
				}
			}
			int tbl_size = index;
			pclose(fp);
			// for(int k = 0; k < index; k++){
			// 	cout << offset_buffer[k][0] << " " << offset_buffer[k][1] << " " << offset_buffer[k][2] << endl;
			// }

			off64_t lba_offset;
			off64_t lba_length;

			TableBlock table_block;

			for(const auto table_block_chunk: sst.second.table_lba_block().table_block_chunks()){
				int table_index_number = table_block_chunk.first;

				Chunks chunks;	

				for(const auto lba: table_block_chunk.second.chunks()){		
					// std::cout << "Offset : " << table_block_chunk.second.chunks(j).offset() << std::endl;
					// std::cout << "Length : " << table_block_chunk.second.chunks(j).length() << std::endl;

					flag = 0;
					lba_offset = lba.offset();
					lba_length = lba.length();

					for(int k = 0; k < tbl_size; k++){
						Chunk pba_chunk;
						if(flag || lba_offset >= offset_buffer[k][0] && lba_offset < offset_buffer[k][0] + offset_buffer[k][2]){
							flag = 1;
							if(lba_length > offset_buffer[k][2]){ // X
								// printf("{\n\t\"Offset\" : %ld,\n\t\"Length\" : %ld\n},\n",offset_buffer[k][1] + lba_length - offset_buffer[k][0],offset_buffer[k][2]);
								pba_chunk.set_offset(offset_buffer[k][1] + lba_offset - offset_buffer[k][0]);
								pba_chunk.set_length(offset_buffer[k][2]);
								chunks.add_chunks()->CopyFrom(pba_chunk); //push back res offset to offset list

								lba_length -= offset_buffer[k][2];
								lba_offset += offset_buffer[k][2];
							} else { // here
								// printf("{\n\t\"Offset\" : %ld,\n\t\"Length\" : %ld\n},\n",offset_buffer[k][1] + lba_length - offset_buffer[k][0],lba_length);
								pba_chunk.set_offset(offset_buffer[k][1] + lba_offset - offset_buffer[k][0]);
								pba_chunk.set_length(lba_length);
								chunks.add_chunks()->CopyFrom(pba_chunk); //push back res offset to offset list
								break;
							}
						}
					}
				}

				table_block.mutable_table_block_chunks()->insert({table_index_number,chunks});
			}	

			response_sst.mutable_table_pba_block()->insert({csd_id,table_block});
		}

		response.mutable_sst_list()->insert({sst_name,response_sst});
	}
	
	return response;
}

class StorageManagerServiceImpl final : public StorageManager::Service {
  Status RequestPBA(ServerContext* context, const LBARequest* request, PBAResponse* response) override {
    KETILOG("LBA2PBA Manager", "# called pba request");
	
	// {
	// std::string test_json;
	// google::protobuf::util::JsonPrintOptions options;
	// options.always_print_primitive_fields = true;
	// options.always_print_enums_as_ints = true;
	// google::protobuf::util::MessageToJsonString(*request,&test_json,options);
	// std::cout << test_json << std::endl << std::endl;
	// }

	PBAResponse res;
	res = RunLBA2PBA(*request);
	response->CopyFrom(res);

	// {
	// std::string test_json;
	// google::protobuf::util::JsonPrintOptions options;
	// options.always_print_primitive_fields = true;
	// options.always_print_enums_as_ints = true;
	// google::protobuf::util::MessageToJsonString(*response,&test_json,options);
	// std::cout << test_json << std::endl << std::endl;
	// }

    return Status::OK;
  }
};

void RunServer() {
  std::string server_address((std::string)STORAGE_CLUSTER_MASTER_IP+":"+std::to_string(LBA2PBA_Manager_Port));
  StorageManagerServiceImpl service;

  ServerBuilder builder;
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);
  std::unique_ptr<Server> server(builder.BuildAndStart());

  KETILOG("LBA2PBA Manager", "LBA2PBA Manager Server listening on "+server_address);

  server->Wait();
}

int main(int argc, char* argv[]){	
	RunServer();
	return 0;
}