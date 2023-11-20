//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

var data = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

func TestTableGetDomainStat(t *testing.T) {
	type args struct {
		r      io.Reader
		domain string
	}
	tests := []struct {
		name           string
		args           args
		wantDomainStat DomainStat
		wantErr        error
	}{
		{
			name: "find 'com'",
			args: args{
				r:      bytes.NewBufferString(data),
				domain: "com",
			},
			wantDomainStat: DomainStat{
				"browsecat.com": 2,
				"linktype.com":  1,
			},
			wantErr: nil,
		},
		{
			name: "find 'gov'",
			args: args{
				r:      bytes.NewBufferString(data),
				domain: "gov",
			},
			wantDomainStat: DomainStat{
				"browsedrive.gov": 1,
			},
			wantErr: nil,
		},
		{
			name: "find 'unknown'",
			args: args{
				r:      bytes.NewBufferString(data),
				domain: "unknown",
			},
			wantDomainStat: DomainStat{},
			wantErr:        nil,
		},
		{
			name: "get users error",
			args: args{
				r:      bytes.NewBufferString("text that cannot be marshalled in Users"),
				domain: "unknown",
			},
			wantDomainStat: nil,
			wantErr:        ErrGetUsers,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDomainStat, gotErr := GetDomainStat(tt.args.r, tt.args.domain)
			require.ErrorIs(t, gotErr, tt.wantErr)

			require.True(t, reflect.DeepEqual(gotDomainStat, tt.wantDomainStat),
				fmt.Sprintf("GetDomainStat() got = %v, want %v", gotDomainStat, tt.wantDomainStat))
		})
	}
}

func BenchmarkGetDomainStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(b, err)
	}
}
