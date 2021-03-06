// Code generated by go-swagger; DO NOT EDIT.

package federation_gateways

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetFegNetworkIDGatewaysParams creates a new GetFegNetworkIDGatewaysParams object
// with the default values initialized.
func NewGetFegNetworkIDGatewaysParams() *GetFegNetworkIDGatewaysParams {
	var ()
	return &GetFegNetworkIDGatewaysParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetFegNetworkIDGatewaysParamsWithTimeout creates a new GetFegNetworkIDGatewaysParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetFegNetworkIDGatewaysParamsWithTimeout(timeout time.Duration) *GetFegNetworkIDGatewaysParams {
	var ()
	return &GetFegNetworkIDGatewaysParams{

		timeout: timeout,
	}
}

// NewGetFegNetworkIDGatewaysParamsWithContext creates a new GetFegNetworkIDGatewaysParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetFegNetworkIDGatewaysParamsWithContext(ctx context.Context) *GetFegNetworkIDGatewaysParams {
	var ()
	return &GetFegNetworkIDGatewaysParams{

		Context: ctx,
	}
}

// NewGetFegNetworkIDGatewaysParamsWithHTTPClient creates a new GetFegNetworkIDGatewaysParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetFegNetworkIDGatewaysParamsWithHTTPClient(client *http.Client) *GetFegNetworkIDGatewaysParams {
	var ()
	return &GetFegNetworkIDGatewaysParams{
		HTTPClient: client,
	}
}

/*GetFegNetworkIDGatewaysParams contains all the parameters to send to the API endpoint
for the get feg network ID gateways operation typically these are written to a http.Request
*/
type GetFegNetworkIDGatewaysParams struct {

	/*NetworkID
	  Network ID

	*/
	NetworkID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get feg network ID gateways params
func (o *GetFegNetworkIDGatewaysParams) WithTimeout(timeout time.Duration) *GetFegNetworkIDGatewaysParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get feg network ID gateways params
func (o *GetFegNetworkIDGatewaysParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get feg network ID gateways params
func (o *GetFegNetworkIDGatewaysParams) WithContext(ctx context.Context) *GetFegNetworkIDGatewaysParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get feg network ID gateways params
func (o *GetFegNetworkIDGatewaysParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get feg network ID gateways params
func (o *GetFegNetworkIDGatewaysParams) WithHTTPClient(client *http.Client) *GetFegNetworkIDGatewaysParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get feg network ID gateways params
func (o *GetFegNetworkIDGatewaysParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithNetworkID adds the networkID to the get feg network ID gateways params
func (o *GetFegNetworkIDGatewaysParams) WithNetworkID(networkID string) *GetFegNetworkIDGatewaysParams {
	o.SetNetworkID(networkID)
	return o
}

// SetNetworkID adds the networkId to the get feg network ID gateways params
func (o *GetFegNetworkIDGatewaysParams) SetNetworkID(networkID string) {
	o.NetworkID = networkID
}

// WriteToRequest writes these params to a swagger request
func (o *GetFegNetworkIDGatewaysParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param network_id
	if err := r.SetPathParam("network_id", o.NetworkID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
