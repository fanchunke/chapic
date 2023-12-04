package runtime

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Bind(msg proto.Message, req *http.Request, pathParams map[string]string) error {
	if err := bindPathParams(msg, pathParams); err != nil {
		return err
	}

	if err := req.ParseForm(); err != nil {
		return err
	}
	var seqs [][]string
	if err := runtime.PopulateQueryParameters(msg, req.Form, utilities.NewDoubleArray(seqs)); err != nil {
		return err
	}
	return nil
}

func bindPathParams(msg proto.Message, pathParams map[string]string) error {
	message := msg.ProtoReflect()

	for jsonField, value := range pathParams {
		// 获取 Protobuf 字段描述符
		fieldDescriptor := message.Descriptor().Fields().ByJSONName(jsonField)
		if fieldDescriptor != nil {
			// 将字符串值转换为目标字段类型并设置到消息中
			fieldValue, err := convertStringToFieldValue(fieldDescriptor.Kind(), value)
			if err != nil {
				return err
			}
			message.Set(fieldDescriptor, fieldValue)
		}
	}

	return nil
}

func convertStringToFieldValue(kind protoreflect.Kind, value string) (protoreflect.Value, error) {
	switch kind {
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(value), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		// 使用 strconv 进行转换
		val, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt32(int32(val)), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		// 使用 strconv 进行转换
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt64(val), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		// 使用 strconv 进行转换
		val, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfUint32(uint32(val)), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		// 使用 strconv 进行转换
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfUint64(val), nil
	case protoreflect.FloatKind:
		// 使用 strconv 进行转换
		val, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfFloat32(float32(val)), nil
	case protoreflect.DoubleKind:
		// 使用 strconv 进行转换
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfFloat64(val), nil
	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported field type: %v", kind)
	}
}
