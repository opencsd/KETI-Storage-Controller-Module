# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: wal_manager.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x11wal_manager.proto\x12\x0bwal_manager\"5\n\tquery_req\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\x0c\n\x04type\x18\x02 \x01(\t\x12\r\n\x05value\x18\x03 \x01(\t\"i\n\x06\x43olumn\x12\x0e\n\x06\x63olumn\x18\x01 \x01(\t\x12\x10\n\x08\x64\x61tatype\x18\x02 \x01(\t\x12\r\n\x05\x63type\x18\x03 \x01(\t\x12\x0c\n\x04\x63len\x18\x04 \x01(\x05\x12\x11\n\tprecision\x18\x05 \x01(\x05\x12\r\n\x05value\x18\x06 \x01(\t\"V\n\x03Wal\x12\x12\n\ntable_name\x18\x01 \x01(\t\x12\x11\n\tindex_val\x18\x02 \x01(\t\x12(\n\x0b\x63olumn_list\x18\x03 \x03(\x0b\x32\x13.wal_manager.Column\"<\n\tquery_res\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\"\n\x08wal_list\x18\x02 \x03(\x0b\x32\x10.wal_manager.Wal2O\n\nWalManager\x12\x41\n\rprocess_query\x12\x16.wal_manager.query_req\x1a\x16.wal_manager.query_res\"\x00\x62\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'wal_manager_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  _QUERY_REQ._serialized_start=34
  _QUERY_REQ._serialized_end=87
  _COLUMN._serialized_start=89
  _COLUMN._serialized_end=194
  _WAL._serialized_start=196
  _WAL._serialized_end=282
  _QUERY_RES._serialized_start=284
  _QUERY_RES._serialized_end=344
  _WALMANAGER._serialized_start=346
  _WALMANAGER._serialized_end=425
# @@protoc_insertion_point(module_scope)
