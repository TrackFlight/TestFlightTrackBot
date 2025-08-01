// Code generated by sqlgen. DO NOT EDIT.

package db

import (
{{- range .Imports }}
    "{{.}}"
{{- end }}
    "github.com/valkey-io/valkey-go"
    "context"
)

{{$rawName := .StoreName | ToPascalCase -}}
{{$name := printf "%sStore" $rawName -}}
type {{$name}} struct {
    db DBTX
    cx context.Context
    redis valkey.Client
}
{{ range .Queries -}}
{{- $isBulk := IsBulkQuery .}}
{{- $queryOptions := GetQueryOptions .}}
{{- $filteredColumns := FilterColumnsByKeys .Columns $queryOptions.Exclude -}}
{{- $hasResults := gt (len $filteredColumns) 0 -}}
{{- $cacheOptions := $queryOptions.Cache}}
{{- $tableSingularName := $cacheOptions.Table | Singular}}
{{- $allowedGetCache := and $cacheOptions.Allow (not (eq .Cmd ":exec")) (eq $cacheOptions.Kind "get") (not $isBulk)}}
{{- $allowedVersioning := and $allowedGetCache $cacheOptions.VersionBy}}
{{- $allowedSplitCacheSave := and $allowedGetCache $cacheOptions.KeyColumn.IsArray}}
{{- $filteredSplitCacheName := ""}}
{{- if $allowedSplitCacheSave}}
{{- $filteredSplitCacheName = printf "%sFiltered" ($cacheOptions.Key | ToCamelCase)}}
{{- end}}
{{- $allowedUpdateVersions := and (eq $cacheOptions.Kind "update_version") $isBulk}}
{{- $isCacheKeyDummy := and (contains $queryOptions.Exclude ($cacheOptions.Key | Singular)) $allowedUpdateVersions}}
{{- $allowedDummy := and (gt (sub (len .Columns) (len $filteredColumns) (bool_to_int $isCacheKeyDummy)) 0)}}
{{- $cacheDummyKeyType := ""}}
{{- $versioningField := ""}}
{{- $fieldName := ""}}
{{- $fieldUpdateName := ""}}
{{- $versionType := ""}}
{{- if $cacheOptions.KeyColumn}}
{{- $cacheDummyKeyType = ToGoType $cacheOptions.KeyColumn true}}
{{- end}}
{{- if $allowedUpdateVersions}}
{{- $fieldUpdateName = $cacheOptions.Key | ToPascalCase | Singular | ToGoCase}}
{{- end}}
{{- if $allowedVersioning}}
{{- $versioningField = printf "fmt.Sprintf(\"%s_version:%s\"" ($cacheOptions.VersionBy.Table | Singular) (GetSprintfFormat $cacheOptions.VersionBy.Column)}}
{{- $fieldName = $cacheOptions.VersionBy.Column.Name | ToPascalCase | ToGoCase}}
{{- $versionType = ToGoType $cacheOptions.VersionBy.Column true}}
{{- end}}
{{$queryName := .Name | ToCamelCase}}
{{- $isMany := eq .Cmd ":many"}}
{{- $arraysList := GetArrays .}}
{{- $isSingleBulk := eq (len $arraysList) 1}}
const {{$queryName}} = `{{.Text}}`
{{- $returnName := printf "%s%sRow" .Name $rawName}}
{{- $allowPointer := not $isMany}}
{{- if eq (len $filteredColumns) 1 -}}
{{- $allowPointer = false}}
{{- $returnName = ToGoType (index $filteredColumns 0) true -}}
{{- end }}
{{- $checkCompatible := ""}}
{{- $bulkParamsName := printf "%s%sParams" (.Name | ToPascalCase) $rawName}}
{{- if and $isBulk (not $isSingleBulk) }}
type {{$bulkParamsName}} struct {
{{- range .Params}}
{{- if .Column.IsArray}}
    {{.Column.Name | ToPascalCase | Singular | ToGoCase}} {{ToGoType .Column true}}
{{- end }}
{{- end}}
}
{{- end}}
{{- if gt (len $filteredColumns) 1 }}
{{- $checkCompatible = FindCompatible (index $filteredColumns 0).Table.Name .Columns}}
{{- if gt (len $checkCompatible) 0 }}
{{- $returnName = printf "models.%s" ($checkCompatible | ToPascalCase | Singular) -}}
{{- else}}

type {{$returnName}} struct {
{{- range $filteredColumns }}
    {{.Name | ToPascalCase | ToGoCase}} {{ToGoType . true}} `json:"{{.Name}}{{if not .NotNull}},omitempty{{end}}"`
{{- end }}
}
{{- end }}
{{- end }}

func (ctx *{{$name}}) {{.Name}}(
{{- $first := true -}}
{{- $addedBulk := true -}}
{{- range (GetParamsOrdered . $queryOptions.Order) -}}
{{- if not (and $isBulk .Column.IsArray (not $isSingleBulk)) -}}
{{- if not $first}}, {{end -}}
{{- .Column.Name | ToCamelCase}} {{ if .Column.IsArray -}}[]{{- end -}}{{ToGoType .Column true}}
{{- else if $addedBulk -}}
{{- if not $first}}, {{end -}}
bulkParams []{{$bulkParamsName}}
{{- $addedBulk = false -}}
{{- end -}}
{{- $first = false -}}
{{ end -}}
) {{if $hasResults -}}({{ if $isMany -}}[]{{else if $allowPointer -}}*{{- end -}}
{{$returnName -}}, error)
{{- else -}} error
{{- end }} {
    {{- if and (not (eq .Cmd ":exec")) (gt (len $filteredColumns) 0) }}
    var i {{ if $isMany -}}[]{{- end -}}{{$returnName}}
    {{- end}}
    {{- if $isCacheKeyDummy}}
    var iVersions []{{$cacheDummyKeyType}}
    {{- end}}
    {{- if $isBulk }}
    {{- if not $isSingleBulk }}
    {{- range .Params}}
    {{- if .Column.IsArray}}
    var {{.Column.Name | ToCamelCase}} []{{ToGoType .Column true}}
    {{- end}}
    {{- end}}
    for _, params := range bulkParams {
        {{- range .Params}}
        {{- if .Column.IsArray}}
        {{- $listName := .Column.Name | ToCamelCase}}
        {{$listName}} = append({{$listName}}, params.{{.Column.Name | ToPascalCase | Singular | ToGoCase}})
        {{- end}}
        {{- end}}
    }
    {{- end}}
    tx, err := ctx.db.Begin(ctx.cx)
    if err != nil {
        return {{ if and (not (eq .Cmd ":exec")) $hasResults -}}nil, {{end}}err
    }
    defer tx.Rollback(ctx.cx)
    totalSize := len({{- if not $isSingleBulk }}bulkParams{{else}}{{(index $arraysList 0).Name | ToCamelCase}}{{end}})
    for start := 0; start < totalSize; start += batchMaxSize {
        end := start + batchMaxSize
        if end > totalSize {
        	end = totalSize
        }
    {{- end}}
    {{- if or $allowedGetCache (eq $cacheOptions.Kind "remove")}}
    {{- if not (or (and (contains $cacheOptions.Fields "all") (eq (len $cacheOptions.Fields) 1)) $allowedSplitCacheSave)}}
    {{- if $cacheOptions.Key}}
    valkeyKey := fmt.Sprintf("{{$tableSingularName}}:{{GetSprintfFormatFromKey . $cacheOptions.Key}}", {{$cacheOptions.Key | ToCamelCase}})
    {{- else}}
    valkeyKey := "{{$tableSingularName}}:list"
    {{- end }}
    {{- end }}
    {{- end }}
    {{- if $allowedGetCache}}
    valkeyField := "{{if eq (len .Columns) 1}}{{(index .Columns 0).Name}}{{else}}{{- if and $cacheOptions.Key (or (not $isMany) $allowedSplitCacheSave)}}__row{{else}}__all{{- end -}}{{- end -}}"
    {{- if $allowedVersioning}}
    valkeyVersionField := "__versions"
    {{- end}}
    {{- if $allowedSplitCacheSave}}
    var cmdList []valkey.Completed
    for _, x := range {{$cacheOptions.Key | ToCamelCase}} {
        cmdList = append(
            cmdList,
            ctx.redis.B().Hget().Key(
                fmt.Sprintf("{{$tableSingularName}}:{{GetSprintfFormatFromKey . $cacheOptions.Key}}", x),
            ).Field(valkeyField).Build(),
        )
    }
    {{- end}}
    res := ctx.redis.DoMulti(
        ctx.cx,
        {{- if $allowedSplitCacheSave}}
        cmdList...
        {{- else}}
        ctx.redis.B().Hget().Key(valkeyKey).Field(valkeyField).Build(),
        {{- end}}
        {{- if $allowedVersioning}}
        ctx.redis.B().Hget().Key(valkeyKey).Field(valkeyVersionField).Build(),
        {{- end}}
    )
    {{- if $allowedSplitCacheSave}}
    var {{$filteredSplitCacheName}} []{{ToGoType $cacheOptions.KeyColumn true}}
    for idx, r := range res {
    	if r.Error() == nil {
    		var item {{$returnName}}
    		if err := r.DecodeJSON(&item); err == nil {
    			i = append(i, item)
    			continue
    		}
    	}
    	{{$filteredSplitCacheName}} = append({{$filteredSplitCacheName}}, {{$cacheOptions.Key | ToCamelCase}}[idx])
    }
    {{- end}}
    {{- if not $allowedSplitCacheSave}}
    if res[0].Error() == nil {{- if $allowedVersioning}} && res[1].Error() == nil{{- end}} {
        {{- if $allowedVersioning}}
        var vMap map[{{$versionType}}]int64
        if err := res[1].DecodeJSON(&vMap); err == nil {
            vMapIDs := slices.Collect(maps.Keys(vMap))
            versionKeys := make([]string, len(vMapIDs))
            for idx, vMapID := range vMapIDs {
                versionKeys[idx] = {{$versioningField}}, vMapID)
            }
            verRes, errVer := ctx.redis.Do(
                ctx.cx,
                ctx.redis.B().Mget().Key(versionKeys...).Build(),
            ).AsIntSlice()
            if errVer == nil {
                valid := true
                for idx, vMapID := range vMapIDs {
					if vMap[vMapID] != verRes[idx] && verRes[idx] != 0 {
                    	valid = false
                    	break
                    }
				}
				if valid {
				    if errAll := res[0].DecodeJSON(&i); errAll == nil {
				        return {{if $allowPointer -}}&{{- end -}}i, nil
				    }
                }
            }
        }
        {{- else}}
        if err := res[0].DecodeJSON(&i); err == nil {
            return {{if $allowPointer -}}&{{- end -}}i, nil
        }
        {{- end}}
    }
    {{- end }}
    {{- end }}
    {{- if $allowedVersioning}}
    var versionKeys []string
    {{- end}}
    {{- if $allowedSplitCacheSave}}
    if len(entityIdsFiltered) > 0 {
    {{- end}}
    {{if eq .Cmd ":one" -}}
    row
    {{- else if $isMany -}}
    rows, errQuery
    {{- else -}}
    tag, errScan
    {{- end}} := {{ if $isBulk }}tx{{else}}ctx.db{{end}}.{{- if eq .Cmd ":one" -}}
    QueryRow
    {{- else if $isMany -}}
    Query
    {{- else -}}
    Exec
    {{- end -}}(
        ctx.cx,
        {{$queryName}},
        {{- range .Params}}
        {{.Column.Name | ToCamelCase}}{{ if and $isBulk .Column.IsArray }}[start:end]{{end}}{{if and $allowedSplitCacheSave (eq .Column.Name $cacheOptions.Key)}}Filtered{{- end}},
        {{- end }}
    )
    {{- if not (eq .Cmd ":exec")}}
    {{- $tmpVarName := "i"}}
    {{- if $isMany}}
    {{- $tmpVarName = "item"}}
    defer rows.Close()
    if errQuery != nil {
        return {{if $hasResults -}}nil,{{end}} errQuery
    }
    for rows.Next() {
        {{- if gt (len $filteredColumns) 0}}
        var item {{$returnName}}
        {{- end}}
        {{- if $allowedDummy}}
        var dummy any
        {{- end}}
        {{- if $isCacheKeyDummy }}
        var itemVersion {{$cacheDummyKeyType}}
        {{- end}}
        errScan := rows.Scan(
    {{- else if eq .Cmd ":one"}}
    errScan := row.Scan(
    {{- end}}
    {{- if gt (len .Columns) 1}}
    {{- range .Columns}}
    {{- if and $isCacheKeyDummy (eq ($cacheOptions.Key | Singular) .Name)}}
    &itemVersion,
    {{- else if and $allowedDummy (contains $queryOptions.Exclude .Name)}}
    &dummy,
    {{- else if eq (len $filteredColumns) 1}}
    &{{$tmpVarName}},
    {{- else}}
    &{{$tmpVarName}}.{{.Name | ToPascalCase | ToGoCase}},
    {{- end }}
    {{- end }}
    {{- else if and $isCacheKeyDummy (eq (len .Columns) 1)}}
    &itemVersion,
    {{- else }}
    &{{$tmpVarName}},
    {{- end }}
    )
    if errScan != nil {
        if errors.Is(errScan, pgx.ErrNoRows) {
        	errScan = nil
        }
        return {{if $hasResults -}}{{ if or $isMany $allowPointer }}nil{{else}}i{{- end}},{{end}} errScan
    }
    {{- if $isMany }}
        {{- if gt (len $filteredColumns) 0}}
        i = append(i, item)
        {{- end}}
        {{- if $allowedVersioning}}
        versionKeys = append(versionKeys, {{$versioningField}}, item.{{$fieldName}}))
        {{- end}}
        {{- if $isCacheKeyDummy}}
        iVersions = append(iVersions, itemVersion)
        {{- end}}
    }
    if rows.Err() != nil {
        return {{if $hasResults -}}nil,{{end}} rows.Err()
    }
    {{- end}}
    {{- if $allowedVersioning }}
    {{- if not $isMany}}
    versionKeys = append(versionKeys, {{$versioningField}}, i.{{$fieldName}}))
    {{- end}}
    verRes, _ := ctx.redis.Do(
        ctx.cx,
        ctx.redis.B().Mget().Key(versionKeys...).Build(),
    ).AsIntSlice()
    vMap := make(map[{{$versionType}}]int64)
    {{- if $isMany}}
    for idx, item := range i {
        vMap[item.{{$fieldName}}] = verRes[idx]
    }
    {{- else}}
    vMap[i.{{$fieldName}}] = verRes[0]
    {{- end}}
    jsonVer, err := json.Marshal(vMap)
    if err != nil {
        jsonVer = []byte("{}")
    }
    {{- end}}
    {{- end}}
    {{- if eq .Cmd ":exec" }}
        if errScan != nil {
            return errScan
        }
        if tag.RowsAffected() == 0 {
            return pgx.ErrNoRows
        }
    {{- end}}
    {{- if $allowedGetCache}}
    {{- if not $allowedSplitCacheSave}}
    jsonData, err := json.Marshal(i)
    if err != nil {
        return {{if or $allowPointer $isMany}}nil{{else}}i{{end}}, err
    }
    {{- else}}
    var orderedI []{{$returnName}}
    var cmdSaveList []valkey.Completed
    for _, r := range {{$cacheOptions.Key | ToCamelCase}} {
        for _, item := range i {
		    {{- $keyElement := printf "item.%s" ($cacheOptions.Key | ToPascalCase | Singular | ToGoCase) }}
            if {{$keyElement}} == r {
		        if slices.Contains({{$filteredSplitCacheName}}, {{$keyElement}}) {
		            jsonData, err := json.Marshal(item)
                    if err != nil {
                        return nil, err
                    }
                    valkeyKey := fmt.Sprintf("{{$tableSingularName}}:{{GetSprintfFormatFromKey . $cacheOptions.Key}}", {{$keyElement}})
		            cmdSaveList = append(
		                cmdSaveList,
		                ctx.redis.B().Hset().Key(valkeyKey).FieldValue().FieldValue(valkeyField, string(jsonData)).Build(),
		                ctx.redis.B().Expire().Key(valkeyKey).Seconds({{$cacheOptions.TTL}}).Build(),
		            )
		        }
		        orderedI = append(orderedI, item)
		        break
		    }
		}
    }
    {{- end}}
    ctx.redis.DoMulti(
        ctx.cx,
        {{- if $allowedSplitCacheSave}}
        cmdSaveList...,
        {{- else}}
        ctx.redis.B().Hset().Key(valkeyKey).FieldValue().FieldValue(valkeyField, string(jsonData)).Build(),
        {{- if $allowedVersioning }}
        ctx.redis.B().Hset().Key(valkeyKey).FieldValue().FieldValue(valkeyVersionField, string(jsonVer)).Build(),
        {{- end}}
        ctx.redis.B().Expire().Key(valkeyKey).Seconds({{$cacheOptions.TTL}}).Build(),
        {{- end}}
    )
    {{- else if eq $cacheOptions.Kind "remove"}}
    ctx.redis.Do{{- if gt (len $cacheOptions.Fields) 1 }}Multi{{end}}(
        ctx.cx,
        {{- range $cacheOptions.Fields }}
        {{- if eq . "all"}}
        ctx.redis.B().Hdel().Key("{{$tableSingularName}}:list").Field("__all").Build(),
        {{- else if eq . "all_by_key"}}
        ctx.redis.B().Hdel().Key(valkeyKey).Field("__all").Build(),
        {{- else if eq . "row"}}
        ctx.redis.B().Hdel().Key(valkeyKey).Field("__row").Build(),
        {{- else}}
        ctx.redis.B().Hdel().Key(valkeyKey).Field("{{.}}").Build(),
        {{- end }}
        {{- end }}
    )
    {{- end }}
    {{- if $isBulk }}
    }
    if errComm := tx.Commit(ctx.cx); errComm != nil {
        return {{ if and (not (eq .Cmd ":exec")) $hasResults -}}nil, {{end}}errComm
    }
    {{- if $allowedUpdateVersions}}
    var cmdList []valkey.Completed
    for _, params := range bulkParams {
        {{- if le (len $filteredColumns) 1}}
        if !slices.Contains({{if $isCacheKeyDummy}}iVersions{{else}}i{{end}}, params.{{$fieldUpdateName}})
        {{- else}}
        if !slices.ContainsFunc(i, func(x {{$returnName}}) bool {
        	return x.{{$fieldUpdateName}} == params.{{$fieldUpdateName}}
        })
        {{- end}} {
            continue
        }
        updateKey := {{ printf "fmt.Sprintf(\"%s_version:%s\", params.%s)" $tableSingularName (GetSprintfFormat $cacheOptions.KeyColumn) $fieldUpdateName}}
        cmdList = append(
            cmdList,
            ctx.redis.B().Incr().Key(updateKey).Build(),
        )
        cmdList = append(
            cmdList,
            ctx.redis.B().Expire().Key(updateKey).Seconds({{$cacheOptions.TTL}}).Build(),
        )
    }
    ctx.redis.DoMulti(ctx.cx, cmdList...)
    {{- end}}
    {{- end}}
    {{- if $allowedSplitCacheSave}}
    }
    {{- end}}
    return {{ if and (not (eq .Cmd ":exec")) $hasResults -}}{{if $allowPointer -}}&{{- end -}}i, {{end}}nil
}
{{ end -}}