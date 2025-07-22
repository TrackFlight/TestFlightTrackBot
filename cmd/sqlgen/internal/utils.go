package internal

import (
	"fmt"
	"github.com/Laky-64/gologging"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func ToPascalCase(name string) string {
	nameSplit := strings.Split(name, "_")
	for i, x := range nameSplit {
		nameSplit[i] = strings.ToUpper(x[:1]) + x[1:]
	}
	return strings.Join(nameSplit, "")
}

func ToCamelCase(name string) string {
	pascal := ToPascalCase(name)
	if len(pascal) == 0 {
		return ""
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

var commonCases = []string{
	"id",
	"url",
	"http",
	"json",
	"api",
	"xml",
	"ip",
	"uuid",
	"uid",
	"html",
	"sql",
}

func ToGoCase(name string) string {
	for _, v := range commonCases {
		re := regexp.MustCompile(`(?m)(?i)\b` + v + `\b`)
		name = re.ReplaceAllString(name, strings.ToUpper(v))
	}

	for _, v := range commonCases {
		re := regexp.MustCompile(`(?m)(?i)(` + v + `)(?-i)([A-Z]|$)`)
		name = re.ReplaceAllString(name, strings.ToUpper(v))
	}
	return name
}

func ToGoType(col Column, isExternal bool) string {
	return internalGoType(col, false, isExternal)
}

func internalGoType(col Column, isSprint bool, isExternal bool) string {
	name := col.Type.Name
	name = strings.TrimPrefix(name, "pg_catalog.")
	if name == "text" || name == "varchar" {
		if isSprint {
			return "%s"
		}
		return "string"
	}
	if name == "int4" || name == "int8" || name == "integer" || name == "bigint" || name == "serial" || name == "bigserial" {
		if isSprint {
			return "%d"
		}
		if col.NotNull {
			return "int64"
		} else {
			return "pgtype.Int8"
		}
	}
	if name == "bool" {
		if isSprint {
			return "%t"
		}
		return "bool"
	}
	if name == "float4" || name == "float8" {
		if isSprint {
			return "%f"
		}
		return "float64"
	}
	if name == "timestamptz" {
		if isSprint {
			return "%s"
		}
		return "pgtype.Timestamptz"
	}
	if isSprint {
		return "%v"
	}
	typeName := ToPascalCase(name)
	if isExternal {
		return fmt.Sprintf("models.%s", typeName)
	}
	return typeName
}

func DetectQueryImports(tables []Table, queries []Query) []string {
	var imports = make(map[string]struct{})

	for _, query := range queries {
		if len(query.Columns) > 0 {
			if len(internalFindCompatible(tables, query.Columns[0].Table.Name, query.Columns)) > 0 {
				imports[fmt.Sprintf("%s/models", currentPackage)] = struct{}{}
			} else {
				tmpImports := DetectImports(query.Columns, true)
				for _, imp := range tmpImports {
					imports[imp] = struct{}{}
				}
			}
			if GetQueryOptions(&query).Cache.Allow {
				imports["fmt"] = struct{}{}
				imports["encoding/json"] = struct{}{}
			}
		}
		var columnParams []Column
		for _, param := range query.Params {
			columnParams = append(columnParams, param.Column)
		}
		tmpParamsImports := DetectImports(columnParams, true)
		for _, imp := range tmpParamsImports {
			imports[imp] = struct{}{}
		}
		if query.Cmd == ":one" || query.Cmd == ":many" {
			imports["github.com/jackc/pgx/v5"] = struct{}{}
			imports["errors"] = struct{}{}
		}
	}

	var list []string
	for imp := range imports {
		list = append(list, imp)
	}
	sort.Strings(list)
	return list
}

func DetectImports(columns []Column, isExternal bool) []string {
	var imports = make(map[string]struct{})
	for _, col := range columns {
		goType := ToGoType(col, isExternal)
		if strings.HasPrefix(goType, "pgtype.") {
			imports["github.com/jackc/pgx/v5/pgtype"] = struct{}{}
		} else if strings.HasPrefix(goType, "models.") {
			imports[fmt.Sprintf("%s/models", currentPackage)] = struct{}{}
		}
	}

	var list []string
	for imp := range imports {
		list = append(list, imp)
	}
	sort.Strings(list)
	return list
}

func Singular(s string) string {
	switch {
	case strings.HasSuffix(s, "ves") && len(s) > 3:
		return s[:len(s)-3] + "f"
	case strings.HasSuffix(s, "ies") && len(s) > 3:
		return s[:len(s)-3] + "y"
	case strings.HasSuffix(s, "ches") || strings.HasSuffix(s, "shes") ||
		strings.HasSuffix(s, "ses") || strings.HasSuffix(s, "xes") ||
		strings.HasSuffix(s, "zes") && len(s) > 4:
		return s[:len(s)-2]
	case strings.HasSuffix(s, "s") && len(s) > 1:
		return s[:len(s)-1]
	}
	return s
}

func FindCompatible(tables []Table) func(tableName string, cols []Column) string {
	return func(tableName string, cols []Column) string {
		return internalFindCompatible(tables, tableName, cols)
	}
}

func GetQueryOptions(query *Query) *QueryOptions {
	options := &QueryOptions{
		Cache: CacheOptions{
			TTL: 60 * 15,
		},
	}
	for _, comment := range query.Comments {
		commentData := strings.SplitN(strings.TrimSpace(comment), ":", 2)
		switch commentData[0] {
		case "cache":
			options.Cache.Allow = true

			paramStr := strings.TrimSpace(commentData[1])
			i := 0
			length := len(paramStr)
			for i < length {
				for i < length && unicode.IsSpace(rune(paramStr[i])) {
					i++
				}
				keyStart := i
				for i < length && paramStr[i] != ':' {
					i++
				}
				if i >= length {
					break
				}
				key := strings.TrimSpace(paramStr[keyStart:i])
				i++
				var value []string
				if i < length && paramStr[i] == '"' {
					i++
					valStart := i
					for i < length && paramStr[i] != '"' {
						if unicode.IsSpace(rune(paramStr[i])) {
							i++
							continue
						}
						if paramStr[i] == ',' {
							tmpValue := paramStr[valStart:i]
							if len(tmpValue) > 0 {
								value = append(value, tmpValue)
							}
							valStart = i + 1
						}
						i++
					}
					tmpValue := strings.TrimSpace(paramStr[valStart:i])
					if len(tmpValue) > 0 {
						value = append(value, tmpValue)
					}
					i++
				} else {
					valStart := i
					for i < length && !unicode.IsSpace(rune(paramStr[i])) {
						i++
					}
					value = append(value, paramStr[valStart:i])
				}
				if len(value) == 0 {
					continue
				}
				switch key {
				case "type":
					options.Cache.Kind = value[0]
				case "key":
					options.Cache.Key = value[0]
				case "table":
					options.Cache.Table = value[0]
				case "fields":
					options.Cache.Fields = value
				case "ttl":
					seconds, err := parseDurationToSeconds(value[0])
					if err == nil && seconds > 0 {
						options.Cache.TTL = seconds
					}
				}
			}

			if options.Cache.Kind == "" || options.Cache.Table == "" {
				gologging.FatalF("cache options must have type and table defined in query %s", query.Name)
			} else if options.Cache.Kind == "remove" && len(options.Cache.Fields) == 0 {
				gologging.FatalF("cache options with type 'remove' must have fields defined in query %s", query.Name)
			}
		case "order":
			order := strings.TrimSpace(commentData[1])
			i := 0
			length := len(order)
			valStart := 0
			for i < length {
				for i < length && unicode.IsSpace(rune(order[i])) {
					i++
				}
				if order[i] == ',' {
					tmpValue := strings.TrimSpace(order[valStart:i])
					if len(tmpValue) > 0 {
						options.Order = append(options.Order, tmpValue)
					}
					valStart = i + 1
				}
				i++
			}
			tmpValue := strings.TrimSpace(order[valStart:i])
			if len(tmpValue) > 0 {
				options.Order = append(options.Order, tmpValue)
			}
			if len(options.Order) > 0 && len(options.Order) != len(query.Params) {
				gologging.FatalF("order options must have the same number of fields as parameters in query %s", query.Name)
			}
		}
	}
	return options
}

func GetSprintfFormat(query *Query, from string) string {
	for _, param := range query.Params {
		if from == param.Column.Name {
			return internalGoType(param.Column, true, false)
		}
	}
	gologging.Fatal("GetSprintfFormat: no parameter found for column %s in query %s", from, query.Name)
	return ""
}

func IsBulkQuery(query *Query) bool {
	foundBulk := false
	invalidBulkQuery := false
	for i, params := range query.Params {
		if params.Column.IsArray {
			foundBulk = true
		} else if query.Params[min(len(query.Params)-1, i+1)].Column.IsArray && query.Params[max(0, i-1)].Column.IsArray {
			invalidBulkQuery = true
		}
	}
	resulBulk := foundBulk && !invalidBulkQuery
	if query.Cmd == ":one" && resulBulk {
		gologging.FatalF("query %s is marked as :one but contains array parameters, which is not allowed", query.Name)
	}
	return resulBulk
}

func GetParamsOrdered(query *Query, order []string) []Param {
	if len(query.Params) <= 1 || len(order) == 0 {
		return query.Params
	}

	paramMap := make(map[string]Param, len(query.Params))
	for _, p := range query.Params {
		paramMap[p.Column.Name] = p
	}

	orderedParams := make([]Param, 0, len(order))
	for _, name := range order {
		if p, ok := paramMap[name]; ok {
			orderedParams = append(orderedParams, p)
		}
	}

	return orderedParams
}

func internalFindCompatible(tables []Table, tableName string, cols []Column) string {
	for _, table := range tables {
		if table.Rel.Name == tableName {
			if len(table.Columns) != len(cols) {
				return ""
			}
			for i, col := range cols {
				if col.Name != table.Columns[i].Name || col.Type.Name != table.Columns[i].Type.Name {
					return ""
				}
			}
			return table.Rel.Name
		}
	}
	return ""
}

func parseDurationToSeconds(input string) (int64, error) {
	re := regexp.MustCompile(`(?i)(\d+)\s*([wdhms])`)
	matches := re.FindAllStringSubmatch(input, -1)

	var totalSec int64
	for _, m := range matches {
		val, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return 0, err
		}
		switch strings.ToLower(m[2]) {
		case "w":
			totalSec += val * 7 * 24 * 3600
		case "d":
			totalSec += val * 24 * 3600
		case "h":
			totalSec += val * 3600
		case "m":
			totalSec += val * 60
		case "s":
			totalSec += val
		}
	}
	return totalSec, nil
}
