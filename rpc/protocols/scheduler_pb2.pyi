import service_pb2 as _service_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

BUILD_DOCKER: TaskType
DESCRIPTOR: _descriptor.FileDescriptor
IDLE: TaskType

class Task(_message.Message):
    __slots__ = ["message", "taskType"]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    TASKTYPE_FIELD_NUMBER: _ClassVar[int]
    message: str
    taskType: TaskType
    def __init__(self, taskType: _Optional[_Union[TaskType, str]] = ..., message: _Optional[str] = ...) -> None: ...

class TaskType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
