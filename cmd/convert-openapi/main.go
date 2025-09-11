package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"

    openapi2 "github.com/getkin/kin-openapi/openapi2"
    openapi2conv "github.com/getkin/kin-openapi/openapi2conv"
    openapi3 "github.com/getkin/kin-openapi/openapi3"
    "gopkg.in/yaml.v3"
    k8syaml "sigs.k8s.io/yaml"
)

func main() {
    in := flag.String("in", "docs/swagger.yaml", "Path to Swagger 2.0 file (yaml or json)")
    out := flag.String("out", "docs/openapi.yaml", "Path to write OpenAPI 3.0 yaml")
    flag.Parse()

    data, err := ioutil.ReadFile(*in)
    if err != nil {
        log.Fatalf("failed to read input: %v", err)
    }

    // Load swagger v2 from YAML/JSON
    var docV2 openapi2.T
    ext := filepath.Ext(*in)
    if ext == ".yaml" || ext == ".yml" {
        j, err := k8syaml.YAMLToJSON(data)
        if err != nil {
            log.Fatalf("failed to convert yaml to json: %v", err)
        }
        if err := docV2.UnmarshalJSON(j); err != nil {
            log.Fatalf("failed to parse swagger v2 json: %v", err)
        }
    } else {
        if err := docV2.UnmarshalJSON(data); err != nil {
            log.Fatalf("failed to parse swagger v2 json: %v", err)
        }
    }

    // Convert to v3
    docV3, err := openapi2conv.ToV3(&docV2)
    if err != nil {
        log.Fatalf("failed to convert to openapi v3: %v", err)
    }

    // Ensure servers are set using basePath & schemes
    if len(docV3.Servers) == 0 {
        // Build servers from schemes + basePath if available; default to http
        base := docV2.BasePath
        schemes := docV2.Schemes
        if len(schemes) == 0 {
            schemes = []string{"http"}
        }
        for _, s := range schemes {
            url := fmt.Sprintf("%s://localhost%s", s, base)
            docV3.Servers = append(docV3.Servers, &openapi3.Server{URL: url})
        }
    }

    // Marshal to YAML
    outData, err := yaml.Marshal(docV3)
    if err != nil {
        log.Fatalf("failed to marshal openapi v3: %v", err)
    }

    if err := os.WriteFile(*out, outData, 0644); err != nil {
        log.Fatalf("failed to write output: %v", err)
    }

    log.Printf("OpenAPI v3 written to %s", *out)
}
