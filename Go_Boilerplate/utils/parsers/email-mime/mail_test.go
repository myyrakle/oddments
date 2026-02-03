package parsers_test

import (
	"fmt"
	"testing"

	"github.com/looko-corp/acloset-api/pkg/parsers"
	"github.com/stretchr/testify/assert"
)

func TestParseEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rawEmail string
		want     parsers.ParsedEmail
		wantErr  bool
	}{
		{
			name:     "Simple plain text email",
			rawEmail: "Return-Path: <sssang97@naver.com>\r\nReceived: from cvsmtppost13.nm.naver.com (cvsmtppost13.nm.naver.com [114.111.35.30])\r\n by inbound-smtp.us-west-1.amazonaws.com with SMTP id 93mnchu4njsnvipqfa74ta3nrovpu3el3bk7qgo1\r\n for test@add-dev.acloset.net;\r\n Mon, 02 Feb 2026 06:39:27 +0000 (UTC)\r\nX-SES-Spam-Verdict: PASS\r\nX-SES-Virus-Verdict: PASS\r\nReceived-SPF: pass (spfCheck: domain of naver.com designates 114.111.35.30 as permitted sender) client-ip=114.111.35.30; envelope-from=sssang97@naver.com; helo=cvsmtppost13.nm.naver.com;\r\nAuthentication-Results: amazonses.com;\r\n spf=pass (spfCheck: domain of naver.com designates 114.111.35.30 as permitted sender) client-ip=114.111.35.30; envelope-from=sssang97@naver.com; helo=cvsmtppost13.nm.naver.com;\r\n dkim=pass header.i=@naver.com;\r\n dmarc=pass header.from=naver.com;\r\nX-SES-RECEIPT: AEFBQUFBQUFBQUFHSzgwMVVlN3pYTTdKbmkzdGE0RHIxNFhuSW52SE5FdWFTNys2OHRqUXhRbWJEa21rWGVFbFlETWFNWTNLSlQ4N3A0U3lPUVVzYjNUUUVhYllQUmxwR0hRNDlQY011dCs3WWE2anByMmwrV3dCNG1OUndtTVdKRGZBRTEvcTBPUFc2ZWU0MWQwdUJqQ2hwM1RwL0JNZWdJZVc4aHZiM0pSVVFXRkdKUlFsa0RZNlJDRnJ0SUdpWVlZbDVtU0t2SWdUQmxQdnYzUnczc2hMM0hNZjBreGJENFFvMjhmeFNuNHpIdHVKTUFrQXdtR1paSTViZVZqTGZBZEhkMDZrZ3NDN0xUQ2JBRWpoV3JMSmljUXYyUDIxRFRMSUNzdm5GYms4Vys3cTlrekowOVE9PQ==\r\nX-SES-DKIM-SIGNATURE: a=rsa-sha256; q=dns/txt; b=AWZITDdONgvkNWfM+QwMVTatvaYmZk5DLgzQrQ9xvKaYeseDdzeFYZAJ9svm66EKFwLFElhZPYn5NxrQqyHoh5ynAD/jcsgUL6hku/e/I8WfibSDvFM8ntDMJQGDvkyuCotLiWHyWINiF1srDB2lSVwu2JnU4vUo261SnXTSLro=; c=relaxed/simple; s=voqdhibj3ww47nmt5hkqcrgg7xiynmza; d=amazonses.com; t=1770014368; v=1; bh=JnzimfB0MSsEiUKhTXDY/e0w5Zb0Z03ROeYoGiBlv+8=; h=From:To:Cc:Bcc:Subject:Date:Message-ID:MIME-Version:Content-Type:X-SES-RECEIPT;\r\nX-Originating-IP: 121.134.168.27\r\nX-Works-Smtp-Source: UXb/axvrFqJCpNmrFovwWHF0\r\nReceived: from cvsendbo021.nm ([10.112.18.58])\r\n  by cvsmtppost13.nm.naver.com with ESMTP id eftCyxJARuWDPBwOv5LcSA\r\n  for <test@add-dev.acloset.net>;\r\n  Mon, 02 Feb 2026 06:39:26 -0000\r\nDKIM-Signature: v=1; a=rsa-sha256; c=relaxed/simple; d=naver.com; s=s20171208;\r\n\tt=1770014366; bh=JnzimfB0MSsEiUKhTXDY/e0w5Zb0Z03ROeYoGiBlv+8=;\r\n\th=Message-ID:Date:From:To:Subject:From:Subject:Feedback-ID:\r\n\t X-Works-Security;\r\n\tb=b0iD6JPxdH1Vp2MAeXnvv+NZmE4gvtvXAS8zcdlHvgIQIgIFc7gwU2DsFSrQClMva\r\n\t GJN0qRp8G+2Lbbx3B6ybe7jbGln1fiwaS2wU12xV6KQfA3gfWqkoEvsEAi8D481/ID\r\n\t repf70IVuTuDWxbPMzpdHwYoFJ1nuMz5hBDnY0cUdhuZVfQ/WrD/Qj8Rrohr4klUUZ\r\n\t XuQ9Py6zHorptSzUzdK1mu5KhT1L8EygnQrqx+gBXIuJhrMHjHGMmK+llGLh2ziP79\r\n\t ai1m7PZQpfSpnFSuslxKxDsoFWEOptACvQS5Auw9f3DDB6HvDRmEbxXbKg07ehN8A6\r\n\t DHJaDYkzcCs0A==\r\nX-Session-ID: 6042c7ef7a21ca7a59593a1210a31fbcfbb3fb745305da7afdf70bc47f086826\r\nMIME-Version: 1.0\r\nMessage-ID: <65997ca68c3ff04a7479e4cd56378fc8@cweb004.nm>\r\nDate: Mon, 02 Feb 2026 15:39:25 +0900\r\nFrom: =?utf-8?B?7KCV7IOB6rG0?= <sssang97@naver.com>\r\nImportance: normal\r\nTo: <test@add-dev.acloset.net>\r\nSubject: =?utf-8?B?7YWM7Iqk7Yq47KCc66qp?=\r\nX-Originating-IP: 121.134.168.27\r\nX-Works-Send-Opt: 3rnwjAIYjHmqFqgrKxJYFx2XKXwYKBmm\r\nContent-Type: multipart/alternative;\r\n\tboundary=\"-----Boundary-WM=_7f7eb30db700.1770014366150\"\r\n\r\n-------Boundary-WM=_7f7eb30db700.1770014366150\r\nContent-Type: text/plain;\r\n\tcharset=\"utf-8\"\r\nContent-Transfer-Encoding: base64\r\n\r\n7YWM7Iqk7Yq464K07JqpCg==\r\n\r\n-------Boundary-WM=_7f7eb30db700.1770014366150\r\nContent-Type: text/html;\r\n\tcharset=\"utf-8\"\r\nContent-Transfer-Encoding: base64\r\n\r\nPGh0bWw+PGhlYWQ+PHN0eWxlPnB7bWFyZ2luLXRvcDowcHg7bWFyZ2luLWJvdHRvbTowcHg7fTwv\r\nc3R5bGU+PC9oZWFkPjxib2R5PjxkaXYgc3R5bGU9ImZvbnQtc2l6ZToxNHB4OyBmb250LWZhbWls\r\neTpHdWxpbSzqtbTrprwsc2Fucy1zZXJpZjsiPu2FjOyKpO2KuOuCtOyaqTwvZGl2PjwvYm9keT48\r\nL2h0bWw+PHRhYmxlIHN0eWxlPSdkaXNwbGF5Om5vbmUnPjx0cj48dGQ+PGltZyBzcmM9Imh0dHBz\r\nOi8vbWFpbC5uYXZlci5jb20vcmVhZFJlY2VpcHQvbm90aWZ5Lz9pbWc9aGVlTlc0SnFiWEtsRnhN\r\nWWFxYmxheCUyQm9NeE0lMkZNcUY0cG9nZE14YmRGcTAwRnpGdkZ4TXFGcUM0TXFDZ01YJTJCME1v\r\nZ21GVmw1V3glMkZzJTJCemtxJTJCdUlDcHp0UnB6a3I3NEpvV3plcXBCdDVXNGtkLmdpZiIgYm9y\r\nZGVyPSIwIi8+PC90ZD48L3RyPjwvdGFibGU+\r\n\r\n-------Boundary-WM=_7f7eb30db700.1770014366150--\r\n",
			want: parsers.ParsedEmail{
				From:     "정상건 <sssang97@naver.com>",
				To:       []string{"<test@add-dev.acloset.net>"},
				Subject:  "테스트제목",
				TextBody: "테스트내용\n",
				HTMLBody: `<html><head><style>p{margin-top:0px;margin-bottom:0px;}</style></head><body><div style="font-size:14px; font-family:Gulim,굴림,sans-serif;">테스트내용</div></body></html><table style='display:none'><tr><td><img src="https://mail.naver.com/readReceipt/notify/?img=heeNW4JqbXKlFxMYaqblax%2BoMxM%2FMqF4pogdMxbdFq00FzFvFxMqFqC4MqCgMX%2B0MogmFVl5Wx%2Fs%2Bzkq%2BuICpztRpzkr74JoWzeqpBt5W4kd.gif" border="0"/></td></tr></table>`,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsers.ParseEmail(tt.rawEmail)

			if err != nil {
				t.Logf("Error: %v", err)
			}

			assert.Equalf(t, tt.wantErr, err != nil, fmt.Sprintf("%v", err))
			assert.Equal(t, tt.want, got)
		})
	}
}
