// Generated by the gRPC C++ plugin.
// If you make any local change, they will be lost.
// source: lba2pba.proto

#include "lba2pba.pb.h"
#include "lba2pba.grpc.pb.h"

#include <functional>
#include <grpcpp/support/async_stream.h>
#include <grpcpp/support/async_unary_call.h>
#include <grpcpp/impl/channel_interface.h>
#include <grpcpp/impl/client_unary_call.h>
#include <grpcpp/support/client_callback.h>
#include <grpcpp/support/message_allocator.h>
#include <grpcpp/support/method_handler.h>
#include <grpcpp/impl/rpc_service_method.h>
#include <grpcpp/support/server_callback.h>
#include <grpcpp/impl/server_callback_handlers.h>
#include <grpcpp/server_context.h>
#include <grpcpp/impl/service_type.h>
#include <grpcpp/support/sync_stream.h>
namespace StorageEngineInstance {

static const char* LBA2PBAManager_method_names[] = {
  "/StorageEngineInstance.LBA2PBAManager/RequestPBA",
};

std::unique_ptr< LBA2PBAManager::Stub> LBA2PBAManager::NewStub(const std::shared_ptr< ::grpc::ChannelInterface>& channel, const ::grpc::StubOptions& options) {
  (void)options;
  std::unique_ptr< LBA2PBAManager::Stub> stub(new LBA2PBAManager::Stub(channel, options));
  return stub;
}

LBA2PBAManager::Stub::Stub(const std::shared_ptr< ::grpc::ChannelInterface>& channel, const ::grpc::StubOptions& options)
  : channel_(channel), rpcmethod_RequestPBA_(LBA2PBAManager_method_names[0], options.suffix_for_stats(),::grpc::internal::RpcMethod::NORMAL_RPC, channel)
  {}

::grpc::Status LBA2PBAManager::Stub::RequestPBA(::grpc::ClientContext* context, const ::StorageEngineInstance::LBARequest& request, ::StorageEngineInstance::PBAResponse* response) {
  return ::grpc::internal::BlockingUnaryCall< ::StorageEngineInstance::LBARequest, ::StorageEngineInstance::PBAResponse, ::grpc::protobuf::MessageLite, ::grpc::protobuf::MessageLite>(channel_.get(), rpcmethod_RequestPBA_, context, request, response);
}

void LBA2PBAManager::Stub::async::RequestPBA(::grpc::ClientContext* context, const ::StorageEngineInstance::LBARequest* request, ::StorageEngineInstance::PBAResponse* response, std::function<void(::grpc::Status)> f) {
  ::grpc::internal::CallbackUnaryCall< ::StorageEngineInstance::LBARequest, ::StorageEngineInstance::PBAResponse, ::grpc::protobuf::MessageLite, ::grpc::protobuf::MessageLite>(stub_->channel_.get(), stub_->rpcmethod_RequestPBA_, context, request, response, std::move(f));
}

void LBA2PBAManager::Stub::async::RequestPBA(::grpc::ClientContext* context, const ::StorageEngineInstance::LBARequest* request, ::StorageEngineInstance::PBAResponse* response, ::grpc::ClientUnaryReactor* reactor) {
  ::grpc::internal::ClientCallbackUnaryFactory::Create< ::grpc::protobuf::MessageLite, ::grpc::protobuf::MessageLite>(stub_->channel_.get(), stub_->rpcmethod_RequestPBA_, context, request, response, reactor);
}

::grpc::ClientAsyncResponseReader< ::StorageEngineInstance::PBAResponse>* LBA2PBAManager::Stub::PrepareAsyncRequestPBARaw(::grpc::ClientContext* context, const ::StorageEngineInstance::LBARequest& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderHelper::Create< ::StorageEngineInstance::PBAResponse, ::StorageEngineInstance::LBARequest, ::grpc::protobuf::MessageLite, ::grpc::protobuf::MessageLite>(channel_.get(), cq, rpcmethod_RequestPBA_, context, request);
}

::grpc::ClientAsyncResponseReader< ::StorageEngineInstance::PBAResponse>* LBA2PBAManager::Stub::AsyncRequestPBARaw(::grpc::ClientContext* context, const ::StorageEngineInstance::LBARequest& request, ::grpc::CompletionQueue* cq) {
  auto* result =
    this->PrepareAsyncRequestPBARaw(context, request, cq);
  result->StartCall();
  return result;
}

LBA2PBAManager::Service::Service() {
  AddMethod(new ::grpc::internal::RpcServiceMethod(
      LBA2PBAManager_method_names[0],
      ::grpc::internal::RpcMethod::NORMAL_RPC,
      new ::grpc::internal::RpcMethodHandler< LBA2PBAManager::Service, ::StorageEngineInstance::LBARequest, ::StorageEngineInstance::PBAResponse, ::grpc::protobuf::MessageLite, ::grpc::protobuf::MessageLite>(
          [](LBA2PBAManager::Service* service,
             ::grpc::ServerContext* ctx,
             const ::StorageEngineInstance::LBARequest* req,
             ::StorageEngineInstance::PBAResponse* resp) {
               return service->RequestPBA(ctx, req, resp);
             }, this)));
}

LBA2PBAManager::Service::~Service() {
}

::grpc::Status LBA2PBAManager::Service::RequestPBA(::grpc::ServerContext* context, const ::StorageEngineInstance::LBARequest* request, ::StorageEngineInstance::PBAResponse* response) {
  (void) context;
  (void) request;
  (void) response;
  return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
}


}  // namespace StorageEngineInstance
