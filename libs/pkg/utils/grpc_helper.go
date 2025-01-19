package utils

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func marshalValue[T protoreflect.ProtoMessage](src *anypb.Any) (*T, error) {
	var res T
	err := src.UnmarshalTo(res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func ExtractGrpcAnyValue(src *anypb.Any) (interface{}, error) {
	switch src.TypeUrl {
	case "type.googleapis.com/google.protobuf.StringValue":
		res, err := marshalValue[*wrapperspb.StringValue](src)
		if err != nil {
			return nil, err
		}
		return (*res).Value, nil
	case "type.googleapis.com/google.protobuf.Int32Value":
		res, err := marshalValue[*wrapperspb.Int32Value](src)
		if err != nil {
			return nil, err
		}
		return (*res).Value, nil
	case "type.googleapis.com/google.protobuf.Int64Value":
		res, err := marshalValue[*wrapperspb.Int64Value](src)
		if err != nil {
			return nil, err
		}
		return (*res).Value, nil
	case "type.googleapis.com/google.protobuf.BoolValue":
		res, err := marshalValue[*wrapperspb.BoolValue](src)
		if err != nil {
			return nil, err
		}
		return (*res).Value, nil
	default:
		return nil, fmt.Errorf("Unkown type")
	}
}
