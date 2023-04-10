package sitemap

import (
	"encoding/xml"
	"time"
)

type Image struct {
	XMLName xml.Name `xml:"image:image"`

	Location string `xml:"image:loc"`
	Title    string `xml:"image:title"`
}

type Video struct {
	XMLName xml.Name `xml:"image:image"`

	Thumbnail       string    `xml:"video:thumbnail_loc"`
	Title           string    `xml:"video:title"`
	Duration        uint64    `xml:"video:duration"`
	ViewCount       uint64    `xml:"video:view_count"`
	PublicationDate time.Time `xml:"video:publication_date"`
	FamilyFriendly  string    `xml:"video:family_friendly"`
	GalleryLocation string    `xml:"video:gallery_loc"`
	Live            string    `xml:"video:live"`
}

type Url struct {
	Location        string    `xml:"location"`
	LastMod         time.Time `xml:"lastmod"`
	ChangeFrequency string    `xml:"changefrequency"`
	Priority        float32   `xml:"priority"`

	Images []*Image `xml:",omitempty"`
	Video  *Video   `xml:",omitempty"`
}

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`

	Namespace      string `xml:"xmlns,attr"`
	ImageNamespace string `xml:"xmlns:image,attr"`
	VideoNamespace string `xml:"xmlns:video,attr"`

	Urls []*Url `xml:"url"`
}

type SiteMapIndex struct {
	XMLName   xml.Name `xml:"sitemapindex"`
	Namespace string   `xml:"xmlns,attr"`

	SiteMaps []*UrlSet `xml:"sitemap"`
}

type SiteMapBuilder struct {
}
