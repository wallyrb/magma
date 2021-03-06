// Code generated by go-swagger; DO NOT EDIT.

package federated_lte_networks

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

// NewGetFegLTENetworkIDFederationParams creates a new GetFegLTENetworkIDFederationParams object
// with the default values initialized.
func NewGetFegLTENetworkIDFederationParams() *GetFegLTENetworkIDFederationParams {
	var ()
	return &GetFegLTENetworkIDFederationParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetFegLTENetworkIDFederationParamsWithTimeout creates a new GetFegLTENetworkIDFederationParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetFegLTENetworkIDFederationParamsWithTimeout(timeout time.Duration) *GetFegLTENetworkIDFederationParams {
	var ()
	return &GetFegLTENetworkIDFederationParams{

		timeout: timeout,
	}
}

// NewGetFegLTENetworkIDFederationParamsWithContext creates a new GetFegLTENetworkIDFederationParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetFegLTENetworkIDFederationParamsWithContext(ctx context.Context) *GetFegLTENetworkIDFederationParams {
	var ()
	return &GetFegLTENetworkIDFederationParams{

		Context: ctx,
	}
}

// NewGetFegLTENetworkIDFederationParamsWithHTTPClient creates a new GetFegLTENetworkIDFederationParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetFegLTENetworkIDFederationParamsWithHTTPClient(client *http.Client) *GetFegLTENetworkIDFederationParams {
	var ()
	return &GetFegLTENetworkIDFederationParams{
		HTTPClient: client,
	}
}

/*GetFegLTENetworkIDFederationParams contains all the parameters to send to the API endpoint
for the get feg LTE network ID federation operation typically these are written to a http.Request
*/
type GetFegLTENetworkIDFederationParams struct {

	/*NetworkID
	  Network ID

	*/
	NetworkID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get feg LTE network ID federation params
func (o *GetFegLTENetworkIDFederationParams) WithTimeout(timeout time.Duration) *GetFegLTENetworkIDFederationParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get feg LTE network ID federation params
func (o *GetFegLTENetworkIDFederationParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get feg LTE network ID federation params
func (o *GetFegLTENetworkIDFederationParams) WithContext(ctx context.Context) *GetFegLTENetworkIDFederationParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get feg LTE network ID federation params
func (o *GetFegLTENetworkIDFederationParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get feg LTE network ID federation params
func (o *GetFegLTENetworkIDFederationParams) WithHTTPClient(client *http.Client) *GetFegLTENetworkIDFederationParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get feg LTE network ID federation params
func (o *GetFegLTENetworkIDFederationParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithNetworkID adds the networkID to the get feg LTE network ID federation params
func (o *GetFegLTENetworkIDFederationParams) WithNetworkID(networkID string) *GetFegLTENetworkIDFederationParams {
	o.SetNetworkID(networkID)
	return o
}

// SetNetworkID adds the networkId to the get feg LTE network ID federation params
func (o *GetFegLTENetworkIDFederationParams) SetNetworkID(networkID string) {
	o.NetworkID = networkID
}

// WriteToRequest writes these params to a swagger request
func (o *GetFegLTENetworkIDFederationParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
