package db

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"strings"
)

const batchMaxSize = 1000

func bulkExec(db *gorm.DB, sql string, values ...any) error {
	return bulkRaw(db, sql, nil, values...)
}

func bulkRaw(db *gorm.DB, sql string, dest any, values ...any) error {
	var totalSize int
	var slices []reflect.Value
	currIndex := 0
	for _, value := range values {
		val := reflect.ValueOf(value)
		currIndex += strings.Index(sql[currIndex:], "?") + 1

		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			strType := typeToString(value)
			sql = sql[:currIndex-1] + strType + sql[currIndex:]
			currIndex = currIndex - len(strType) - 1
			continue
		}

		if totalSize != 0 && totalSize != val.Len() {
			return gorm.ErrInvalidData
		}

		totalSize = val.Len()
		slices = append(slices, val)
	}
	sql = strings.ReplaceAll(sql, "?", "%s")

	var sliceDest reflect.Value
	var isSliceDest bool
	if dest != nil {
		destVal := reflect.ValueOf(dest)
		if destVal.Kind() != reflect.Ptr {
			return fmt.Errorf("dest must be a pointer, got %T", dest)
		}
		sliceDest = destVal.Elem()
		if sliceDest.Kind() == reflect.Slice {
			isSliceDest = true
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < totalSize; i += batchMaxSize {
			end := i + batchMaxSize
			if end > totalSize {
				end = totalSize
			}
			var params []any
			for _, slice := range slices {
				sub := slice.Slice(i, end)
				params = append(params, joinSlice(sub.Interface()))
			}
			sqlTmp := fmt.Sprintf(sql, params...)
			if dest != nil {
				var batchDest any
				if isSliceDest {
					batchDest = reflect.New(sliceDest.Type()).Interface()
				} else {
					batchDest = dest
				}
				if err := tx.Raw(sqlTmp).Scan(batchDest).Error; err != nil {
					return err
				}
				if isSliceDest {
					sliceDest.Set(reflect.AppendSlice(sliceDest, reflect.ValueOf(batchDest).Elem()))
				}
			} else {
				if err := tx.Exec(sqlTmp).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func joinSlice(slice any) string {
	v := reflect.ValueOf(slice)

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return ""
	}

	var elems []string
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		elems = append(elems, typeToString(elem))
	}
	return strings.Join(elems, ",")
}

func typeToString(val any) string {
	if isNil(val) {
		return "NULL"
	}
	switch v := val.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case fmt.Stringer:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v.String(), "'", "''"))
	case bool:
		return strconv.FormatBool(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	default:
		return fmt.Sprintf("'%v'", val)
	}
}

func isNil(val any) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}
