package main

import (
    "fmt"
    "strings"
    "unicode"
    "unicode/utf8"

    "google.golang.org/genproto/googleapis/api/annotations"
    "google.golang.org/protobuf/compiler/protogen"
    "google.golang.org/protobuf/proto"
)

func protocVersion(gen *protogen.Plugin) string {
    v := gen.Request.GetCompilerVersion()
    if v == nil {
        return "(unknown)"
    }
    var suffix string
    if s := v.GetSuffix(); s != "" {
        suffix = "-" + s
    }
    return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }

func buildAccessor(field string, rawFinal bool) string {
    // Corner case if passed the result of strings.Join on an empty slice.
    if field == "" {
        return ""
    }

    var ax strings.Builder
    split := strings.Split(field, ".")
    idx := len(split)
    if rawFinal {
        idx--
    }
    for _, s := range split[:idx] {
        fmt.Fprintf(&ax, ".Get%s()", snakeToCamel(s))
    }
    if rawFinal {
        fmt.Fprintf(&ax, ".%s", snakeToCamel(split[len(split)-1]))
    }
    return ax.String()
}

// Given a chained description for a field in a proto message,
// e.g. squid.mantle.mass_kg
// return the string description of the go expression
// describing idiomatic access to the terminal field
// i.e. .GetSquid().GetMantle().GetMassKg()
//
// This is the normal way to retrieve values.
func fieldGetter(field string) string {
    return buildAccessor(field, false)
}

// snakeToCamel converts snake_case and SNAKE_CASE to CamelCase.
func snakeToCamel(s string) string {
    var sb strings.Builder
    up := true
    for _, r := range s {
        if r == '_' {
            up = true
        } else if up && unicode.IsDigit(r) {
            sb.WriteRune('_')
            sb.WriteRune(r)
            up = false
        } else if up {
            sb.WriteRune(unicode.ToUpper(r))
            up = false
        } else {
            sb.WriteRune(unicode.ToLower(r))
        }
    }
    return sb.String()
}

func strContains(a []string, s string) bool {
    for _, as := range a {
        if as == s {
            return true
        }
    }
    return false
}

// isRequired returns if a field is annotated as REQUIRED or not.
func isRequired(field *protogen.Field) bool {
    if field.Desc.Options() == nil {
        return false
    }

    eBehav := proto.GetExtension(field.Desc.Options(), annotations.E_FieldBehavior)

    behaviors := eBehav.([]annotations.FieldBehavior)
    for _, b := range behaviors {
        if b == annotations.FieldBehavior_REQUIRED {
            return true
        }
    }

    return false
}

func lowerFirst(s string) string {
    if s == "" {
        return ""
    }
    r, w := utf8.DecodeRuneInString(s)
    return string(unicode.ToLower(r)) + s[w:]
}
