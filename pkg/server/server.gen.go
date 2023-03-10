// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Defines values for DesiredVMStatus.
const (
	ForceStop DesiredVMStatus = "forceStop"
	Run       DesiredVMStatus = "run"
	Stop      DesiredVMStatus = "stop"
)

// Defines values for VMStatus.
const (
	Running VMStatus = "running"
	Stoped  VMStatus = "stoped"
)

// CloudConfig defines model for CloudConfig.
type CloudConfig struct {
	Users []CloudConfigUser `json:"users"`
}

// CloudConfigUser defines model for CloudConfigUser.
type CloudConfigUser struct {
	Name              string   `json:"name"`
	SshAuthorizedKeys []string `json:"ssh_authorized_keys"`
}

// DesiredVMStatus defines model for DesiredVMStatus.
type DesiredVMStatus string

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// VMSpec defines model for VMSpec.
type VMSpec struct {
	Ip                *string  `json:"ip,omitempty"`
	MemorySizeMb      int      `json:"memory_size_mb"`
	Name              string   `json:"name"`
	Rootfs            string   `json:"rootfs"`
	SshAuthorizedKeys []string `json:"ssh_authorized_keys"`
	StorageGb         int      `json:"storage_gb"`
	VcpuCount         int      `json:"vcpu_count"`
	Vmlinux           string   `json:"vmlinux"`
}

// VMStatus defines model for VMStatus.
type VMStatus string

// VirtualMachine defines model for VirtualMachine.
type VirtualMachine struct {
	CloudConfig  CloudConfig `json:"cloud_config"`
	Id           string      `json:"id"`
	Ip           string      `json:"ip"`
	MemorySizeMb int         `json:"memory_size_mb"`
	Name         string      `json:"name"`
	Vmlinux      string      `json:"vmlinux"`
}

// CreateVMJSONRequestBody defines body for CreateVM for application/json ContentType.
type CreateVMJSONRequestBody = VMSpec

// SetVMStatusJSONRequestBody defines body for SetVMStatus for application/json ContentType.
type SetVMStatusJSONRequestBody = DesiredVMStatus

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List all VMs
	// (GET /v0/vm)
	ListVMs(ctx echo.Context) error
	// Create a VM
	// (POST /v0/vm)
	CreateVM(ctx echo.Context) error
	// Delete VM
	// (DELETE /v0/vm/{vmId})
	DeleteVM(ctx echo.Context, vmId int) error
	// Show VM details
	// (GET /v0/vm/{vmId})
	GetVM(ctx echo.Context, vmId int) error
	// Get VM Log
	// (GET /v0/vm/{vmId}/log)
	GetVMLog(ctx echo.Context, vmId int) error
	// Get VM Status
	// (GET /v0/vm/{vmId}/status)
	GetVMStatus(ctx echo.Context, vmId int) error
	// Update VM Status
	// (PUT /v0/vm/{vmId}/status)
	SetVMStatus(ctx echo.Context, vmId int) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListVMs converts echo context to params.
func (w *ServerInterfaceWrapper) ListVMs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListVMs(ctx)
	return err
}

// CreateVM converts echo context to params.
func (w *ServerInterfaceWrapper) CreateVM(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateVM(ctx)
	return err
}

// DeleteVM converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteVM(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "vmId" -------------
	var vmId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "vmId", runtime.ParamLocationPath, ctx.Param("vmId"), &vmId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter vmId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteVM(ctx, vmId)
	return err
}

// GetVM converts echo context to params.
func (w *ServerInterfaceWrapper) GetVM(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "vmId" -------------
	var vmId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "vmId", runtime.ParamLocationPath, ctx.Param("vmId"), &vmId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter vmId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetVM(ctx, vmId)
	return err
}

// GetVMLog converts echo context to params.
func (w *ServerInterfaceWrapper) GetVMLog(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "vmId" -------------
	var vmId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "vmId", runtime.ParamLocationPath, ctx.Param("vmId"), &vmId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter vmId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetVMLog(ctx, vmId)
	return err
}

// GetVMStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetVMStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "vmId" -------------
	var vmId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "vmId", runtime.ParamLocationPath, ctx.Param("vmId"), &vmId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter vmId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetVMStatus(ctx, vmId)
	return err
}

// SetVMStatus converts echo context to params.
func (w *ServerInterfaceWrapper) SetVMStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "vmId" -------------
	var vmId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "vmId", runtime.ParamLocationPath, ctx.Param("vmId"), &vmId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter vmId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SetVMStatus(ctx, vmId)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v0/vm", wrapper.ListVMs)
	router.POST(baseURL+"/v0/vm", wrapper.CreateVM)
	router.DELETE(baseURL+"/v0/vm/:vmId", wrapper.DeleteVM)
	router.GET(baseURL+"/v0/vm/:vmId", wrapper.GetVM)
	router.GET(baseURL+"/v0/vm/:vmId/log", wrapper.GetVMLog)
	router.GET(baseURL+"/v0/vm/:vmId/status", wrapper.GetVMStatus)
	router.PUT(baseURL+"/v0/vm/:vmId/status", wrapper.SetVMStatus)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"UUcIqbuX3b8CAAD///PwkLHQFgAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
