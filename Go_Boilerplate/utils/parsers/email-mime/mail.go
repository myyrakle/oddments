package parsers

import (
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

type ParsedEmail struct {
	From        string
	To          []string
	Subject     string
	TextBody    string
	HTMLBody    string
	Attachments []Attachment
}

type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

func ParseEmail(rawEmail string) (ParsedEmail, error) {
	msg, err := mail.ReadMessage(strings.NewReader(rawEmail))
	if err != nil {
		return ParsedEmail{}, err
	}

	// MIME encoded-word 디코딩
	dec := &mime.WordDecoder{}
	from, _ := dec.DecodeHeader(msg.Header.Get("From"))
	subject, _ := dec.DecodeHeader(msg.Header.Get("Subject"))

	// To 헤더 파싱 (RFC 5322 호환)
	toList := []string{}
	toHeader := msg.Header.Get("To")
	if toHeader != "" {
		if addrs, err := mail.ParseAddressList(toHeader); err == nil {
			for _, addr := range addrs {
				toList = append(toList, addr.String())
			}
		} else {
			// ParseAddressList 실패시 fallback (전체 헤더를 하나의 수신자로)
			toList = append(toList, strings.TrimSpace(toHeader))
		}
	}

	email := ParsedEmail{
		From:    from,
		To:      toList,
		Subject: subject,
	}

	// Content-Type 파싱
	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		// Content-Type이 없으면 단순 텍스트로 처리
		body, _ := io.ReadAll(msg.Body)
		email.TextBody = string(body)
		return email, nil
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary := params["boundary"]
		if boundary == "" {
			// boundary가 없으면 전체를 TextBody로
			body, _ := io.ReadAll(msg.Body)
			email.TextBody = string(body)
			return email, nil
		}
		mr := multipart.NewReader(msg.Body, boundary)

		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return ParsedEmail{}, err
			}

			contentType := part.Header.Get("Content-Type")
			contentTransferEncoding := part.Header.Get("Content-Transfer-Encoding")

			body, err := io.ReadAll(part)
			if err != nil {
				return ParsedEmail{}, err
			}

			// Base64 디코딩
			if strings.EqualFold(contentTransferEncoding, "base64") {
				decoded, err := base64.StdEncoding.DecodeString(string(body))
				if err == nil {
					body = decoded
				}
			}

			if strings.Contains(contentType, "text/plain") {
				email.TextBody = string(body)
			} else if strings.Contains(contentType, "text/html") {
				email.HTMLBody = string(body)
			} else {
				// 첨부파일
				filename := part.FileName()
				if filename != "" {
					email.Attachments = append(email.Attachments, Attachment{
						Filename:    filename,
						ContentType: contentType,
						Data:        body,
					})
				}
			}
		}
	} else {
		// single part
		body, err := io.ReadAll(msg.Body)
		if err != nil {
			return ParsedEmail{}, err
		}

		if strings.Contains(mediaType, "text/html") {
			email.HTMLBody = string(body)
		} else {
			email.TextBody = string(body)
		}
	}

	return email, nil
}
