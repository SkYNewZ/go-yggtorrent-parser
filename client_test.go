package yggtorrent

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var benchURL = "https://www5.yggtorrent.fi/torrent/filmvid%C3%A9o/film/808831-to+wong+foo+thanks+for+everything+julie+newmar+1995+multi+1080p+web+h264-none"

func Test_getIDFromLink(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "expected",
			args: args{benchURL},
			want: "808831",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIDFromLink(tt.args.u); got != tt.want {
				t.Errorf("getIDFromLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkLastIndex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		idx := strings.LastIndex(benchURL, "/")
		if idx == -1 {
			b.Fatal("oops")
		}

		end := benchURL[idx+1:]
		_ = strings.Split(end, "-")[0]
	}
}

func BenchmarkSplit(b *testing.B) {
	for n := 0; n < b.N; n++ {
		splits := strings.Split(benchURL, "/")
		end := splits[len(splits)-1]
		_ = strings.Split(end, "-")[0]
	}
}

func Test_strToUInt(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			name: "expected",
			args: args{"10"},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strToUInt(tt.args.str); got != tt.want {
				t.Errorf("strToUInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_ParseResults(t *testing.T) {
	type fields struct {
		baseURL string
	}
	type args struct {
		data io.Reader
	}

	file, err := os.Open("testdata/ygg.html")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Result
		wantErr bool
	}{
		{
			name:   "expected",
			fields: fields{"foo"},
			args:   args{file},
			want: []*Result{
				{
					ID:          "808831",
					Name:        "To.Wong.Foo.Thanks.for.Everything.Julie.Newmar.1995.MULTi.1080p.WEB.H264-NoNE",
					PublishedAt: time.Date(2021, 10, 0o6, 16, 49, 51, 0o0, time.Local),
					Size:        "11.59Go",
					Seeders:     2,
					Leechers:    0,
					InfoURL:     "https://www5.yggtorrent.fi/torrent/filmvidéo/film/808831-to+wong+foo+thanks+for+everything+julie+newmar+1995+multi+1080p+web+h264-none",
					DownloadURL: "foo/engine/download_torrent?id=808831",
				},
				{
					ID:          "688934",
					Name:        "To.Wong.Foo.Thanks.for.Everything.Julie.Newmar.(1995).Multi.VFI.DVDRip.480p.MPEG2.AC3",
					PublishedAt: time.Date(2020, 11, 20, 0o4, 58, 0o4, 0o0, time.Local),
					Size:        "6.63Go",
					Seeders:     1,
					Leechers:    0,
					InfoURL:     "https://www5.yggtorrent.fi/torrent/filmvidéo/film/688934-to+wong+foo+thanks+for+everything+julie+newmar+1995+multi+vfi+dvdrip+480p+mpeg2+ac3",
					DownloadURL: "foo/engine/download_torrent?id=688934",
				},
				{
					ID:          "808074",
					Name:        "To.Wong.Foo.Thanks.For.Everything.1995.MULTi.1080p.WEB-DL.DD5.1.AC3.AVC",
					PublishedAt: time.Date(2021, 10, 0o3, 23, 20, 30, 0o0, time.Local),
					Size:        "4.66Go",
					Seeders:     7,
					Leechers:    0,
					InfoURL:     "https://www5.yggtorrent.fi/torrent/filmvidéo/film/808074-to+wong+foo+thanks+for+everything+1995+multi+1080p+web-dl+dd5+1+ac3+avc",
					DownloadURL: "foo/engine/download_torrent?id=808074",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{
				baseURL: tt.fields.baseURL,
			}
			got, err := c.ParseResults(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ParseResults() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
