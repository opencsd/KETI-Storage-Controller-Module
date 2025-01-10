#include "wal_manager.grpc.pb.h"
#include <grpc++/grpc++.h>
#include <memory>
#include <iostream>
#include <string>

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using wal_manager::query_req;
using wal_manager::query_res;
using wal_manager::WalManager;

class WalManagerImpl final : public WalManager::Service {
    Status process_query(ServerContext* context, const query_req* request, query_res* response) override {
        response-> wal_list();
        return Status::OK;
    }
};

void RunServer() {
    std::string server_address{"localhost:50051"};
    WalManagerImpl service;

    // Build server
    ServerBuilder builder;
    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);
    std::unique_ptr<Server> server{builder.BuildAndStart()};

    // Run server
    std::cout << "Server listening on " << server_address << std::endl;
    server->Wait();
}

int main(int argc, char** argv) {
    RunServer();
    return 0;
}