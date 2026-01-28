package main

import (
	"io"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-smtp"
)

// Backend는 SMTP 서버의 백엔드를 구현합니다
type Backend struct{}

// NewSession은 새로운 SMTP 세션을 생성합니다
func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

// Session은 SMTP 세션을 나타냅니다
type Session struct {
	From string
	To   []string
}

// Mail은 메일 발신자를 설정합니다
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Printf("메일 발신자: %s\n", from)
	s.From = from
	return nil
}

// Rcpt는 메일 수신자를 추가합니다
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	log.Printf("메일 수신자: %s\n", to)
	s.To = append(s.To, to)
	return nil
}

// Data는 메일 본문을 읽어들입니다
func (s *Session) Data(r io.Reader) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	log.Printf("=== 메일 수신 ===\n")
	log.Printf("발신자: %s\n", s.From)
	log.Printf("수신자: %v\n", s.To)

	// 수신자 도메인 추출
	for _, recipient := range s.To {
		if idx := strings.Index(recipient, "@"); idx != -1 {
			domain := recipient[idx+1:]
			log.Printf("수신 도메인: %s (전체: %s)\n", domain, recipient)
		}
	}

	log.Printf("본문:\n%s\n", string(body))
	log.Printf("================\n")

	return nil
}

// Reset은 세션을 초기화합니다
func (s *Session) Reset() {
	s.From = ""
	s.To = nil
}

// Logout은 세션을 종료합니다
func (s *Session) Logout() error {
	return nil
}

func main() {
	backend := &Backend{}

	server := smtp.NewServer(backend)
	server.Addr = ":2525" // SMTP 포트 (25번 대신 2525 사용)
	server.Domain = "localhost"
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024 // 1MB
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true // 테스트용이므로 비보안 인증 허용

	log.Printf("SMTP 서버 시작: %s\n", server.Addr)
	log.Printf("도메인: %s\n", server.Domain)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
