// Code generated by go-swagger; DO NOT EDIT.

package carrier_wifi_networks

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

// NewGetCwfNetworkIDSubscriberConfigBaseNamesParams creates a new GetCwfNetworkIDSubscriberConfigBaseNamesParams object
// with the default values initialized.
func NewGetCwfNetworkIDSubscriberConfigBaseNamesParams() *GetCwfNetworkIDSubscriberConfigBaseNamesParams {
	var ()
	return &GetCwfNetworkIDSubscriberConfigBaseNamesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetCwfNetworkIDSubscriberConfigBaseNamesParamsWithTimeout creates a new GetCwfNetworkIDSubscriberConfigBaseNamesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetCwfNetworkIDSubscriberConfigBaseNamesParamsWithTimeout(timeout time.Duration) *GetCwfNetworkIDSubscriberConfigBaseNamesParams {
	var ()
	return &GetCwfNetworkIDSubscriberConfigBaseNamesParams{

		timeout: timeout,
	}
}

// NewGetCwfNetworkIDSubscriberConfigBaseNamesParamsWithContext creates a new GetCwfNetworkIDSubscriberConfigBaseNamesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetCwfNetworkIDSubscriberConfigBaseNamesParamsWithContext(ctx context.Context) *GetCwfNetworkIDSubscriberConfigBaseNamesParams {
	var ()
	return &GetCwfNetworkIDSubscriberConfigBaseNamesParams{

		Context: ctx,
	}
}

// NewGetCwfNetworkIDSubscriberConfigBaseNamesParamsWithHTTPClient creates a new GetCwfNetworkIDSubscriberConfigBaseNamesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetCwfNetworkIDSubscriberConfigBaseNamesParamsWithHTTPClient(client *http.Client) *GetCwfNetworkIDSubscriberConfigBaseNamesParams {
	var ()
	return &GetCwfNetworkIDSubscriberConfigBaseNamesParams{
		HTTPClient: client,
	}
}

/*GetCwfNetworkIDSubscriberConfigBaseNamesParams contains all the parameters to send to the API endpoint
for the get cwf network ID subscriber config base names operation typically these are written to a http.Request
*/
type GetCwfNetworkIDSubscriberConfigBaseNamesParams struct {

	/*NetworkID
	  Network ID

	*/
	NetworkID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get cwf network ID subscriber config base names params
func (o *GetCwfNetworkIDSubscriberConfigBaseNamesParams) WithTimeout(timeout time.Duration) *GetCwfNetworkIDSubscriberConfigBaseNamesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get cwf network ID subscriber config base names params
func (o *GetCwfNetworkIDSubscriberConfigBaseNamesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get cwf network ID subscriber config base names params
func (o *GetCwfNetworkIDSubscriberConfigBaseNamesParams) WithContext(ctx context.Context) *GetCwfNetworkIDSubscriberConfigBaseNamesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get cwf network ID subscriber config base names params
func (o *GetCwfNetworkIDSubscriberConfigBaseNamesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get cwf network ID subscriber config base names params
func (o *GetCwfNetworkIDSubscriberConfigBaseNamesParams) WithHTTPClient(client *http.Client) *GetCwfNetworkIDSubscriberConfigBaseNamesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get cwf network ID subscriber config base names params
func (o *GetCwfNetworkIDSubscriberConfigBaseNamesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithNetworkID adds the networkID to the get cwf network ID subscriber config base names params
func (o *GetCwfNetworkIDSubscriberConfigBaseNamesParams) WithNetworkID(networkID string) *GetCwfNetworkIDSubscriberConfigBaseNamesParams {
	o.SetNetworkID(networkID)
	return o
}

// SetNetworkID adds the networkId to the get cwf network ID subscriber config base names params
func (o *GetCwfNetworkIDSubscriberConfigBaseNamesParams) SetNetworkID(networkID string) {
	o.NetworkID = networkID
}

// WriteToRequest writes these params to a swagger request
func (o *GetCwfNetworkIDSubscriberConfigBaseNamesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
