package valourgo

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/auroradevllc/apiclient/multipart"
	log "github.com/sirupsen/logrus"
)

type MessageAttachmentType int

const (
	AttachmentTypeImage MessageAttachmentType = iota
	AttachmentTypeVideo
	AttachmentTypeAudio
	AttachmentTypeFile

	// Future embedded attachments for Valour functions
	AttachmentTypeValourMessage
	AttachmentTypeValourInvite
	AttachmentTypeValourPlanet
	AttachmentTypeValourChannel
	AttachmentTypeValourItem
	AttachmentTypeValourEcoAccount
	AttachmentTypeValourEcoTrade
	AttachmentTypeValourReceipt
	AttachmentTypeValourBot

	// Generic link preview using Open Graph
	AttachmentTypeSitePreview

	// Video platforms
	AttachmentTypeYouTube
	AttachmentTypeVimeo
	AttachmentTypeTwitch
	AttachmentTypeTikTok

	// Social platforms
	AttachmentTypeTwitter
	AttachmentTypeReddit
	AttachmentTypeInstagram
	AttachmentTypeBluesky

	// Music platforms
	AttachmentTypeSpotify
	AttachmentTypeSoundCloud

	// Developer platforms
	AttachmentTypeGitHub
)

type MessageAttachment struct {
	Location string                `json:"Location"`
	MimeType string                `json:"MimeType"`
	FileName string                `json:"FileName"`
	Width    int                   `json:"Width"`
	Height   int                   `json:"Height"`
	Inline   bool                  `json:"Inline"`
	Type     MessageAttachmentType `json:"Type"`
}

// UploadImage uploads an image to Valour
func (n *Node) UploadImage(fileName string, r io.Reader, size int64) (*MessageAttachment, error) {
	// Upload to app.valour.gg/image/upload
	// Base64 decode response for url
	s := multipart.New()

	var header bytes.Buffer

	tee := io.TeeReader(r, &header)

	cfg, format, err := image.DecodeConfig(tee)

	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"fileName": fileName,
		"size":     size,
		"format":   format,
		"width":    cfg.Width,
		"height":   cfg.Height,
	}).Debug("Image ready for upload")

	attachment := &MessageAttachment{
		FileName: fileName,
		Width:    cfg.Width,
		Height:   cfg.Height,
		Type:     AttachmentTypeImage,
	}

	// Re-combine bytes we read and the rest of the data
	mr := io.MultiReader(&header, r)

	h := make(textproto.MIMEHeader)

	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="image"; filename="%s"`, escapeQuotes(fileName)))
	h.Set("Content-Type", "image/"+format) // format will be jpeg/png/gif

	if err := s.CreatePart(h, mr, size); err != nil {
		return nil, err
	}

	res, err := n.request(http.MethodPost, "upload/image", s)

	if err != nil {
		return nil, err
	}

	b, err := res.Bytes()

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed with status code %d: %s", res.StatusCode, string(b))
	}

	attachment.Location = string(b)

	return attachment, nil
}

// UploadFile uploads a file to Valour
func (n *Node) UploadFile(fileName string, r io.Reader, size int64) (*MessageAttachment, error) {
	s := multipart.New()

	var header bytes.Buffer

	tee := io.TeeReader(r, &header)

	headerReadSize := 512

	if size < int64(headerReadSize) {
		headerReadSize = int(size)
	}

	fileHeader := make([]byte, headerReadSize)

	if _, err := io.ReadFull(tee, fileHeader); err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(fileHeader)

	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = contentType[:idx]
	}

	log.WithFields(log.Fields{
		"name":        fileName,
		"contentType": contentType,
		"size":        size,
	}).Debug("File is ready for uploading")

	attachment := &MessageAttachment{
		FileName: fileName,
		MimeType: contentType,
		Type:     AttachmentTypeFile,
	}

	// Re-combine bytes we read and the rest of the data
	mr := io.MultiReader(&header, r)

	h := make(textproto.MIMEHeader)

	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="image"; filename="%s"`, escapeQuotes(fileName)))
	h.Set("Content-Type", contentType)

	if err := s.CreatePart(h, mr, size); err != nil {
		return nil, err
	}

	res, err := n.request(http.MethodPost, "upload/file", s)

	if err != nil {
		return nil, err
	}

	b, err := res.Bytes()

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed with status code %d: %s", res.StatusCode, string(b))
	}

	attachment.Location = string(b)

	return attachment, err
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
