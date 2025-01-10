#include "wal_manager.grpc.pb.h"
#include <grpc++/grpc++.h>
#include <memory>
#include <iostream>

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using wal_manager::query_req;
using wal_manager::query_res;
using wal_manager::WalManager;
using wal_manager::Wal;

class wal_managerClient {
public:
    wal_managerClient(std::shared_ptr<Channel> channel) : _stub{WalManager::NewStub(channel)} {}

    std::string process_query(const std::string& key, const std::string& type, const std::string& value) {
        // Prepare request
        query_req request;
        request.set_key(key);
        request.set_type(type);
        request.set_value(value);

        // Send request
        query_res response;
        ClientContext context;
        Status status;
        status = _stub->process_query(&context, request, &response);

        // Handle response
        if (status.ok()) {
            return response.key();
        } else {
            std::cerr << status.error_code() << ": " << status.error_message() << std::endl;
            return "RPC failed";
        }
    }

private:
    std::unique_ptr<WalManager::Stub> _stub;
};

int main(int argc, char** argv) {
    std::string server_address{"localhost:12345"};
    wal_managerClient client{grpc::CreateChannel(server_address, grpc::InsecureChannelCredentials())};
    std::string key{"0|0|nation"}, type{"TABLE"}, value{"nation"};
    std::string response = client.process_query(key, type, value);
    std::cout << "Client received: " << response<< std::endl;
    return 0;
}