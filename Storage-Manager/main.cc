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
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include "rapidjson/prettywriter.h" 

#include "lba2pba_manager.h"

#define STORAGE_CLUSTER_MASTER_IP "10.0.4.83"
#define STORAGE_MANAGER_PORT 40301

using namespace rapidjson;
using namespace std;

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;

class StorageManagerServiceImpl final : public StorageManager::Service {
  Status RequestPBA(ServerContext* context, const LBA2PBARequest* request, LBA2PBAResponse* response) override {
	cout << "[Storage Manager] # request pba called" << endl;
	
	{
	// Check LBA Request - Debug Code   
	std::string test_json;
	google::protobuf::util::JsonPrintOptions options;
	options.always_print_primitive_fields = true;
	options.always_print_enums_as_ints = true;
	google::protobuf::util::MessageToJsonString(*request,&test_json,options);
	std::cout << "LBA Request" << std::endl;
	// std::cout << test_json << std::endl << std::endl;
	}

	LBA2PBAResponse lba2pbaResponse;
	lba2pbaResponse = LBA2PBAManager::GetPBAData(*request);
	response->CopyFrom(lba2pbaResponse);

	{
	// Check PBA Response - Debug Code   
	std::string test_json;
	google::protobuf::util::JsonPrintOptions options;
	options.always_print_primitive_fields = true;
	options.always_print_enums_as_ints = true;
	google::protobuf::util::MessageToJsonString(*response,&test_json,options);
	std::cout << "PBA Response" << std::endl;
	// std::cout << test_json << std::endl << std::endl;
	}

    return Status::OK;
  }

  Status GetDataFileInfo(ServerContext* context, const SSTList* request, DataFileInfo* response) override {
	cout << "[Storage Manager] # get data file info called" << endl;
	
	DataFileInfo response_;
	response_ = DBFileMonitoring::GetDataFileInfo(*request);
	response->CopyFrom(response_);

    return Status::OK;
  }
};

void RunServer() {
  std::string server_address((std::string)STORAGE_CLUSTER_MASTER_IP+":"+std::to_string(STORAGE_MANAGER_PORT));
  StorageManagerServiceImpl service;

  ServerBuilder builder;
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);
  std::unique_ptr<Server> server(builder.BuildAndStart());

  cout << "[Storage Manager] Server Listening on " << server_address << endl;

  server->Wait();
}

int main(int argc, char* argv[]){	
	DBFileMonitoring::InitDBFileMonitoring();
	LBA2PBAManager::InitLBA2PBAManager();

	RunServer();
	return 0;
}