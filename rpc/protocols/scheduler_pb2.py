# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: scheduler.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from . import service_pb2 as service__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0fscheduler.proto\x1a\rservice.proto\"4\n\x04Task\x12\x1b\n\x08taskType\x18\x02 \x01(\x0e\x32\t.TaskType\x12\x0f\n\x07message\x18\x01 \x01(\t*&\n\x08TaskType\x12\x08\n\x04IDLE\x10\x00\x12\x10\n\x0c\x42UILD_DOCKER\x10\x01\x32)\n\tScheduler\x12\x1c\n\x07GetTask\x12\x08.Service\x1a\x05.Task\"\x00\x42\x0fZ\rrpc/protocolsb\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'scheduler_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\rrpc/protocols'
  _TASKTYPE._serialized_start=88
  _TASKTYPE._serialized_end=126
  _TASK._serialized_start=34
  _TASK._serialized_end=86
  _SCHEDULER._serialized_start=128
  _SCHEDULER._serialized_end=169
# @@protoc_insertion_point(module_scope)
