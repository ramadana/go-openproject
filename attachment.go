package openproject

import (
	"context"
	"fmt"
	"io/ioutil"
)

// AttachmentService handles attachments for the OpenProject instance / API.
type AttachmentService struct {
	client *Client
}

// Attachment is the object representing OpenProject attachments.
// TODO: Complete fields and complex fields (user, links, downloadlocation, container...)
type Attachment struct {
	Type        string               `json:"_type,omitempty" structs:"_type,omitempty"`
	ID          int                  `json:"id,omitempty" structs:"id,omitempty"`
	FileName    string               `json:"filename,omitempty" structs:"filename,omitempty"`
	FileSize    int                  `json:"filesize,omitempty" structs:"filesize,omitempty"`
	Description OPGenericDescription `json:"description,omitempty" structs:"description,omitempty"`
	ContentType string               `json:"contentType,omitempty" structs:"contentType,omitempty"`
	Digest      AttachmentDigest     `json:"digest,omitempty" structs:"digest,omitempty"`
	Links       *AttachmentLinks     `json:"_links,omitempty" structs:"_links,omitempty"`
}

// AttachmentList is only a small wrapper containing AttachmentElements
type AttachmentList struct {
	Embedded AttachmentElements `json:"_embedded" structs:"_embedded"`
	Total    int                `json:"total" structs:"total"`
	Count    int                `json:"count" structs:"count"`
}

// AttachmentElements represent elements within AttachmentList list
type AttachmentElements struct {
	Elements []Attachment `json:"elements" structs:"elements"`
}

// AttachmentDigest wraps algorithm and hash
type AttachmentDigest struct {
	Algorithm string `json:"algorithm,omitempty" structs:"algorithm,omitempty"`
	Hash      string `json:"hash,omitempty" structs:"hash,omitempty"`
}

// AttachmentLinks contains several links to object that related to this attachment.
// TODO: For now it is only download location. For details on the other object, please refer to Attachment API response.
type AttachmentLinks struct {
	DownloadLocation OPGenericLink `json:"downloadLocation,omitempty" structs:"downloadLocation,omitempty"`
}

// GetWithContext gets a wiki page from OpenProject using its ID
func (s *AttachmentService) GetWithContext(ctx context.Context, attachmentID string) (*Attachment, *Response, error) {
	apiEndPoint := fmt.Sprintf("api/v3/attachments/%s", attachmentID)
	Obj, Resp, err := GetWithContext(ctx, s, apiEndPoint)
	return Obj.(*Attachment), Resp, err
}

// GetList wraps GetWithContext using the background context.
func (s *AttachmentService) GetList(workpackageID string) (*AttachmentList, *Response, error) {
	apiEndPoint := fmt.Sprintf("api/v3/work_packages/%s/attachments", workpackageID)
	objList, resp, err := GetListWithContext(context.Background(), s, apiEndPoint, &FilterOptions{})
	if err != nil {
		return &AttachmentList{}, resp, err
	}
	return objList.(*AttachmentList), resp, err
}

// Get wraps GetWithContext using the background context.
func (s *AttachmentService) Get(attachmentID string) (*Attachment, *Response, error) {
	return s.GetWithContext(context.Background(), attachmentID)
}

// DownloadWithContext downloads a file from attachment using attachment ID
func (s *AttachmentService) DownloadWithContext(ctx context.Context, attachmentID string) (*[]byte, error) {
	apiEndpoint := fmt.Sprintf("api/v3/attachments/%s/content", attachmentID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Download(req)

	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	return &respBytes, err
}

// Download wraps DownloadWithContext using the background context.
func (s *AttachmentService) Download(attachmentID string) (*[]byte, error) {
	return s.DownloadWithContext(context.Background(), attachmentID)
}
