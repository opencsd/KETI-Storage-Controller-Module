# Copyright 2018 gRPC authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# cmake build file for C++ wal_manager example.
# Assumes protobuf and gRPC have been installed using cmake.
# See cmake_externalproject/CMakeLists.txt for all-in-one cmake build
# that automatically builds all the dependencies before building wal_manager.

cmake_minimum_required(VERSION 3.5.1)

project(wal_manager C CXX)

include(./cmake/common.cmake)

# Proto file
get_filename_component(wal_proto "../protos/wal_manager.proto" ABSOLUTE)
get_filename_component(wal_proto_path "${wal_proto}" PATH)

# Generated sources
set(wal_proto_srcs "${CMAKE_CURRENT_BINARY_DIR}/wal_manager.pb.cc")
set(wal_proto_hdrs "${CMAKE_CURRENT_BINARY_DIR}/wal_manager.pb.h")
set(wal_grpc_srcs "${CMAKE_CURRENT_BINARY_DIR}/wal_manager.grpc.pb.cc")
set(wal_grpc_hdrs "${CMAKE_CURRENT_BINARY_DIR}/wal_manager.grpc.pb.h")
add_custom_command(
      OUTPUT "${wal_proto_srcs}" "${wal_proto_hdrs}" "${wal_grpc_srcs}" "${wal_grpc_hdrs}"
      COMMAND ${_PROTOBUF_PROTOC}
      ARGS --grpc_out "${CMAKE_CURRENT_BINARY_DIR}"
        --cpp_out "${CMAKE_CURRENT_BINARY_DIR}"
        -I "${wal_proto_path}"
        --plugin=protoc-gen-grpc="${_GRPC_CPP_PLUGIN_EXECUTABLE}"
        "${wal_proto}"
      DEPENDS "${wal_proto}")

# Include generated *.pb.h files
include_directories("${CMAKE_CURRENT_BINARY_DIR}")

# wal_grpc_proto
add_library(wal_grpc_proto
  ${wal_grpc_srcs}
  ${wal_grpc_hdrs}
  ${wal_proto_srcs}
  ${wal_proto_hdrs})
target_link_libraries(wal_grpc_proto
  ${_REFLECTION}
  ${_GRPC_GRPCPP}
  ${_PROTOBUF_LIBPROTOBUF})

# Targets wal_manager_[async_](client|server)
foreach(_target
  wal_client wal_server)
  add_executable(${_target} "${_target}.cc")
  target_link_libraries(${_target}
    wal_grpc_proto
    ${_REFLECTION}
    ${_GRPC_GRPCPP}
    ${_PROTOBUF_LIBPROTOBUF})
endforeach()
