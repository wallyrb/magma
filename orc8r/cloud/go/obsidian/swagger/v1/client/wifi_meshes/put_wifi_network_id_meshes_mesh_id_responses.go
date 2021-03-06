// Code generated by go-swagger; DO NOT EDIT.

package wifi_meshes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "magma/orc8r/cloud/go/obsidian/swagger/v1/models"
)

// PutWifiNetworkIDMeshesMeshIDReader is a Reader for the PutWifiNetworkIDMeshesMeshID structure.
type PutWifiNetworkIDMeshesMeshIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PutWifiNetworkIDMeshesMeshIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewPutWifiNetworkIDMeshesMeshIDCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewPutWifiNetworkIDMeshesMeshIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPutWifiNetworkIDMeshesMeshIDCreated creates a PutWifiNetworkIDMeshesMeshIDCreated with default headers values
func NewPutWifiNetworkIDMeshesMeshIDCreated() *PutWifiNetworkIDMeshesMeshIDCreated {
	return &PutWifiNetworkIDMeshesMeshIDCreated{}
}

/*PutWifiNetworkIDMeshesMeshIDCreated handles this case with default header values.

New mesh ID
*/
type PutWifiNetworkIDMeshesMeshIDCreated struct {
	Payload models.MeshID
}

func (o *PutWifiNetworkIDMeshesMeshIDCreated) Error() string {
	return fmt.Sprintf("[PUT /wifi/{network_id}/meshes/{mesh_id}][%d] putWifiNetworkIdMeshesMeshIdCreated  %+v", 201, o.Payload)
}

func (o *PutWifiNetworkIDMeshesMeshIDCreated) GetPayload() models.MeshID {
	return o.Payload
}

func (o *PutWifiNetworkIDMeshesMeshIDCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutWifiNetworkIDMeshesMeshIDDefault creates a PutWifiNetworkIDMeshesMeshIDDefault with default headers values
func NewPutWifiNetworkIDMeshesMeshIDDefault(code int) *PutWifiNetworkIDMeshesMeshIDDefault {
	return &PutWifiNetworkIDMeshesMeshIDDefault{
		_statusCode: code,
	}
}

/*PutWifiNetworkIDMeshesMeshIDDefault handles this case with default header values.

Unexpected Error
*/
type PutWifiNetworkIDMeshesMeshIDDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the put wifi network ID meshes mesh ID default response
func (o *PutWifiNetworkIDMeshesMeshIDDefault) Code() int {
	return o._statusCode
}

func (o *PutWifiNetworkIDMeshesMeshIDDefault) Error() string {
	return fmt.Sprintf("[PUT /wifi/{network_id}/meshes/{mesh_id}][%d] PutWifiNetworkIDMeshesMeshID default  %+v", o._statusCode, o.Payload)
}

func (o *PutWifiNetworkIDMeshesMeshIDDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *PutWifiNetworkIDMeshesMeshIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
