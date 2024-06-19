// SPDX-FileCopyrightText: 2022-2023 The go-mail Authors
//
// SPDX-License-Identifier: MIT

package mail

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	exampleMailPlainNoEnc = `Date: Wed, 01 Nov 2023 00:00:00 +0000
MIME-Version: 1.0
Message-ID: <1305604950.683004066175.AAAAAAAAaaaaaaaaB@go-mail.dev>
Subject: Example mail // plain text without encoding
User-Agent: go-mail v0.4.0 // https://github.com/wneessen/go-mail
X-Mailer: go-mail v0.4.0 // https://github.com/wneessen/go-mail
From: "Toni Tester" <go-mail@go-mail.dev>
To: <go-mail+test@go-mail.dev>
Cc: <go-mail+cc@go-mail.dev>
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 8bit

Dear Customer,

This is a test mail. Please do not reply to this. Also this line is very long so it
should be wrapped.


Thank your for your business!
The go-mail team

--
This is a signature`
	exampleMailPlainNoEncInvalidDate = `Date: Inv, 99 Nov 9999 99:99:00 +0000
MIME-Version: 1.0
Message-ID: <1305604950.683004066175.AAAAAAAAaaaaaaaaB@go-mail.dev>
Subject: Example mail // plain text without encoding
User-Agent: go-mail v0.4.0 // https://github.com/wneessen/go-mail
X-Mailer: go-mail v0.4.0 // https://github.com/wneessen/go-mail
From: "Toni Tester" <go-mail@go-mail.dev>
To: <go-mail+test@go-mail.dev>
Cc: <go-mail+cc@go-mail.dev>
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 8bit

Dear Customer,

This is a test mail. Please do not reply to this. Also this line is very long so it
should be wrapped.


Thank your for your business!
The go-mail team

--
This is a signature`
	exampleMailPlainNoEncNoDate = `MIME-Version: 1.0
Message-ID: <1305604950.683004066175.AAAAAAAAaaaaaaaaB@go-mail.dev>
Subject: Example mail // plain text without encoding
User-Agent: go-mail v0.4.0 // https://github.com/wneessen/go-mail
X-Mailer: go-mail v0.4.0 // https://github.com/wneessen/go-mail
From: "Toni Tester" <go-mail@go-mail.dev>
To: <go-mail+test@go-mail.dev>
Cc: <go-mail+cc@go-mail.dev>
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 8bit

Dear Customer,

This is a test mail. Please do not reply to this. Also this line is very long so it
should be wrapped.


Thank your for your business!
The go-mail team

--
This is a signature`
	exampleMailPlainQP = `Date: Wed, 01 Nov 2023 00:00:00 +0000
MIME-Version: 1.0
Message-ID: <1305604950.683004066175.AAAAAAAAaaaaaaaaB@go-mail.dev>
Subject: Example mail // plain text quoted-printable
User-Agent: go-mail v0.4.0 // https://github.com/wneessen/go-mail
X-Mailer: go-mail v0.4.0 // https://github.com/wneessen/go-mail
From: "Toni Tester" <go-mail@go-mail.dev>
To: <go-mail+test@go-mail.dev>
Cc: <go-mail+cc@go-mail.dev>
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: quoted-printable

Dear Customer,

This is a test mail. Please do not reply to this. Also this line is very lo=
ng so it
should be wrapped.


Thank your for your business!
The go-mail team

--
This is a signature`
	exampleMailPlainB64 = `Date: Wed, 01 Nov 2023 00:00:00 +0000
MIME-Version: 1.0
Message-ID: <1305604950.683004066175.AAAAAAAAaaaaaaaaB@go-mail.dev>
Subject: Example mail // plain text base64
User-Agent: go-mail v0.4.0 // https://github.com/wneessen/go-mail
X-Mailer: go-mail v0.4.0 // https://github.com/wneessen/go-mail
From: "Toni Tester" <go-mail@go-mail.dev>
To: <go-mail+test@go-mail.dev>
Cc: <go-mail+cc@go-mail.dev>
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: base64

RGVhciBDdXN0b21lciwKClRoaXMgaXMgYSB0ZXN0IG1haWwuIFBsZWFzZSBkbyBub3QgcmVwbHkg
dG8gdGhpcy4gQWxzbyB0aGlzIGxpbmUgaXMgdmVyeSBsb25nIHNvIGl0CnNob3VsZCBiZSB3cmFw
cGVkLgoKClRoYW5rIHlvdXIgZm9yIHlvdXIgYnVzaW5lc3MhClRoZSBnby1tYWlsIHRlYW0KCi0t
ClRoaXMgaXMgYSBzaWduYXR1cmU=`
)

func TestEMLToMsgFromString(t *testing.T) {
	tests := []struct {
		name string
		eml  string
		enc  string
		sub  string
	}{
		{
			"Plain text no encoding", exampleMailPlainNoEnc, "8bit",
			"Example mail // plain text without encoding",
		},
		{
			"Plain text quoted-printable", exampleMailPlainQP, "quoted-printable",
			"Example mail // plain text quoted-printable",
		},
		{
			"Plain text base64", exampleMailPlainB64, "base64",
			"Example mail // plain text base64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := EMLToMsgFromString(tt.eml)
			if err != nil {
				t.Errorf("failed to parse EML: %subject", err)
			}
			if msg.Encoding() != tt.enc {
				t.Errorf("EMLToMsgFromString failed: expected encoding: %subject, but got: %subject", tt.enc, msg.Encoding())
			}
			if subject := msg.GetGenHeader(HeaderSubject); len(subject) > 0 && !strings.EqualFold(subject[0], tt.sub) {
				t.Errorf("EMLToMsgFromString failed: expected subject: %subject, but got: %subject",
					tt.sub, subject[0])
			}
		})
	}
}

func TestEMLToMsgFromFile(t *testing.T) {
	tests := []struct {
		name string
		eml  string
		enc  string
		sub  string
	}{
		{
			"Plain text no encoding", exampleMailPlainNoEnc, "8bit",
			"Example mail // plain text without encoding",
		},
		{
			"Plain text quoted-printable", exampleMailPlainQP, "quoted-printable",
			"Example mail // plain text quoted-printable",
		},
		{
			"Plain text base64", exampleMailPlainB64, "base64",
			"Example mail // plain text base64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", fmt.Sprintf("*-%s", tt.name))
			if err != nil {
				t.Errorf("failed to create temp dir: %s", err)
				return
			}
			defer func() {
				if err = os.RemoveAll(tempDir); err != nil {
					t.Error("failed to remove temp dir:", err)
				}
			}()
			err = os.WriteFile(fmt.Sprintf("%s/%s.eml", tempDir, tt.name), []byte(tt.eml), 0666)
			if err != nil {
				t.Error("failed to write mail to temp file:", err)
				return
			}
			msg, err := EMLToMsgFromFile(fmt.Sprintf("%s/%s.eml", tempDir, tt.name))
			if err != nil {
				t.Errorf("failed to parse EML: %subject", err)
			}
			if msg.Encoding() != tt.enc {
				t.Errorf("EMLToMsgFromString failed: expected encoding: %subject, but got: %subject", tt.enc, msg.Encoding())
			}
			if subject := msg.GetGenHeader(HeaderSubject); len(subject) > 0 && !strings.EqualFold(subject[0], tt.sub) {
				t.Errorf("EMLToMsgFromString failed: expected subject: %subject, but got: %subject",
					tt.sub, subject[0])
			}
		})
	}
}

func TestEMLToMsgFromStringBrokenDate(t *testing.T) {
	_, err := EMLToMsgFromString(exampleMailPlainNoEncInvalidDate)
	if err == nil {
		t.Error("EML with invalid date was supposed to fail, but didn't")
	}
	now := time.Now()
	m, err := EMLToMsgFromString(exampleMailPlainNoEncNoDate)
	if err != nil {
		t.Errorf("EML with no date parsing failed: %s", err)
	}
	da := m.GetGenHeader(HeaderDate)
	if len(da) < 1 {
		t.Error("EML with no date expected current date, but got nothing")
		return
	}
	d := da[0]
	if d != now.Format(time.RFC1123Z) {
		t.Errorf("EML with no date expected: %s, got: %s", now.Format(time.RFC1123Z), d)
	}
}
