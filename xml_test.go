package utils

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

type serversXMLStructure struct {
	XMLName xml.Name             `xml:"servers"`
	Version string               `xml:"version,attr"`
	Servers []serverXMLStructure `xml:"server"`
}

type serverXMLStructure struct {
	XMLName xml.Name `xml:"server"`
	Name    string   `xml:"name"`
}

var xmlString = `<servers version="1"><server><name>fiber one</name></server><server><name>fiber two</name></server></servers>`

func Test_GolangXMLEncoder(t *testing.T) {
	t.Parallel()

	var (
		ss = &serversXMLStructure{
			Version: "1",
			Servers: []serverXMLStructure{
				{Name: "fiber one"},
				{Name: "fiber two"},
			},
		}
		xmlEncoder XMLMarshal = xml.Marshal
	)

	raw, err := xmlEncoder(ss)
	require.NoError(t, err)

	require.Equal(t, string(raw), xmlString)
}

func Test_DefaultXMLEncoder(t *testing.T) {
	t.Parallel()

	var (
		ss = &serversXMLStructure{
			Version: "1",
			Servers: []serverXMLStructure{
				{Name: "fiber one"},
				{Name: "fiber two"},
			},
		}
		xmlEncoder XMLMarshal = xml.Marshal
	)

	raw, err := xmlEncoder(ss)
	require.NoError(t, err)
	require.Equal(t, string(raw), xmlString)
}

func Test_DefaultXMLDecoder(t *testing.T) {
	t.Parallel()

	var (
		ss         serversXMLStructure
		xmlBytes                = []byte(xmlString)
		xmlDecoder XMLUnmarshal = xml.Unmarshal
	)

	err := xmlDecoder(xmlBytes, &ss)
	require.NoError(t, err)
	require.Len(t, ss.Servers, 2)
	require.Equal(t, "1", ss.Version)
	require.Equal(t, "fiber one", ss.Servers[0].Name)
	require.Equal(t, "fiber two", ss.Servers[1].Name)
}

func Benchmark_GolangXMLEncoder(b *testing.B) {
	var (
		ss = &serversXMLStructure{
			Version: "1",
			Servers: []serverXMLStructure{
				{Name: "fiber one"},
				{Name: "fiber two"},
			},
		}
		xmlEncoder XMLMarshal = xml.Marshal
	)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := xmlEncoder(ss)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_DefaultXMLEncoder(b *testing.B) {
	var (
		ss = &serversXMLStructure{
			Version: "1",
			Servers: []serverXMLStructure{
				{Name: "fiber one"},
				{Name: "fiber two"},
			},
		}
		xmlEncoder XMLMarshal = xml.Marshal
	)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := xmlEncoder(ss)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_DefaultXMLDecoder(b *testing.B) {
	var (
		ss         serversXMLStructure
		xmlBytes                = []byte(xmlString)
		xmlDecoder XMLUnmarshal = xml.Unmarshal
	)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := xmlDecoder(xmlBytes, &ss)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Test_XMLDecodeInvalid(t *testing.T) {
	t.Parallel()

	var ss serversXMLStructure
	err := xml.Unmarshal([]byte("<invalid>"), &ss)
	require.Error(t, err)
}
