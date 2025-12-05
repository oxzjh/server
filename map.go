package server

type Map map[string]any

func (m Map) GetBool(key string) bool {
	val, _ := m[key].(bool)
	return val
}

func (m Map) GetInt8(key string) int8 {
	return int8(m.GetFloat64(key))
}

func (m Map) GetUint8(key string) uint8 {
	return uint8(m.GetFloat64(key))
}

func (m Map) GetInt16(key string) int16 {
	return int16(m.GetFloat64(key))
}

func (m Map) GetUint16(key string) uint16 {
	return uint16(m.GetFloat64(key))
}

func (m Map) GetInt32(key string) int32 {
	return int32(m.GetFloat64(key))
}

func (m Map) GetUint32(key string) uint32 {
	return uint32(m.GetFloat64(key))
}

func (m Map) GetInt64(key string) int64 {
	return int64(m.GetFloat64(key))
}

func (m Map) GetUint64(key string) uint64 {
	return uint64(m.GetFloat64(key))
}

func (m Map) GetInt(key string) int {
	return int(m.GetFloat64(key))
}

func (m Map) GetUint(key string) uint {
	return uint(m.GetFloat64(key))
}

func (m Map) GetFloat32(key string) float32 {
	return float32(m.GetFloat64(key))
}

func (m Map) GetFloat64(key string) float64 {
	val, _ := m[key].(float64)
	return val
}

func (m Map) GetString(key string) string {
	val, _ := m[key].(string)
	return val
}

func (m Map) GetSlice(key string) []any {
	val, _ := m[key].([]any)
	return val
}

func (m Map) GetMap(key string) map[string]any {
	val, _ := m[key].(map[string]any)
	return val
}

func (m Map) GetBoolWithDefault(key string, defaultVal bool) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return defaultVal
}

func (m Map) GetInt8WithDefault(key string, defaultVal int8) int8 {
	if val, ok := m[key].(float64); ok {
		return int8(val)
	}
	return defaultVal
}

func (m Map) GetUint8WithDefault(key string, defaultVal uint8) uint8 {
	if val, ok := m[key].(float64); ok {
		return uint8(val)
	}
	return defaultVal
}

func (m Map) GetInt16WithDefault(key string, defaultVal int16) int16 {
	if val, ok := m[key].(float64); ok {
		return int16(val)
	}
	return defaultVal
}

func (m Map) GetUint16WithDefault(key string, defaultVal uint16) uint16 {
	if val, ok := m[key].(float64); ok {
		return uint16(val)
	}
	return defaultVal
}

func (m Map) GetInt32WithDefault(key string, defaultVal int32) int32 {
	if val, ok := m[key].(float64); ok {
		return int32(val)
	}
	return defaultVal
}

func (m Map) GetUint32WithDefault(key string, defaultVal uint32) uint32 {
	if val, ok := m[key].(float64); ok {
		return uint32(val)
	}
	return defaultVal
}

func (m Map) GetInt64WithDefault(key string, defaultVal int64) int64 {
	if val, ok := m[key].(float64); ok {
		return int64(val)
	}
	return defaultVal
}

func (m Map) GetUint64WithDefault(key string, defaultVal uint64) uint64 {
	if val, ok := m[key].(float64); ok {
		return uint64(val)
	}
	return defaultVal
}

func (m Map) GetIntWithDefault(key string, defaultVal int) int {
	if val, ok := m[key].(float64); ok {
		return int(val)
	}
	return defaultVal
}

func (m Map) GetUintWithDefault(key string, defaultVal uint) uint {
	if val, ok := m[key].(float64); ok {
		return uint(val)
	}
	return defaultVal
}

func (m Map) GetFloat32WithDefault(key string, defaultVal float32) float32 {
	if val, ok := m[key].(float64); ok {
		return float32(val)
	}
	return defaultVal
}

func (m Map) GetFloat64WithDefault(key string, defaultVal float64) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return defaultVal
}

func (m Map) GetStringWithDefault(key, defaultVal string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultVal
}

func (m Map) GetSliceWithDefault(key string, defaultVal []any) []any {
	if val, ok := m[key].([]any); ok {
		return val
	}
	return defaultVal
}

func (m Map) GetMapWithDefault(key string, defaultVal map[string]any) map[string]any {
	if val, ok := m[key].(map[string]any); ok {
		return val
	}
	return defaultVal
}
