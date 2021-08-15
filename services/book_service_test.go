package services

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	goisbn "github.com/abx123/go-isbn"
	"github.com/stretchr/testify/assert"

	"github.com/abx123/library/constant"
	"github.com/abx123/library/entities"
)

var book *entities.Book = &entities.Book{
	BookID:          0,
	ISBN:            "ISBN13",
	Title:           "DUMMY",
	Authors:         "kitefishBB",
	ImageURL:        "imageURL",
	SmallImageURL:   "smallImageURL",
	PublicationYear: 2021,
	Publisher:       "BB publishing house",
	Status:          1,
	Description:     "dummy description",
	PageCount:       999,
	Categories:      "dummy, category",
	Language:        "en",
	Source:          "google",
}

var giBook *goisbn.Book = &goisbn.Book{
	Title:         "DUMMY",
	PublishedYear: "2021",
	Authors:       []string{"kitefishBB"},
	Description:   "dummy description",
	IndustryIdentifiers: &goisbn.Identifier{
		ISBN13: "ISBN13",
	},
	PageCount:  999,
	Categories: []string{"dummy", "category"},
	ImageLinks: &goisbn.ImageLinks{
		SmallImageURL: "smallImageURL",
		ImageURL:      "imageURL",
		LargeImageURL: "largeImageURL",
	},
	Publisher: "BB publishing house",
	Language:  "en",
	Source:    "google",
}

// Custom type that allows setting the func that our Mock Get func will run instead
type MockGetType func(string) (*goisbn.Book, error)

// Custom type that allows setting the func that our Mock ValidateISBN func will run instead
type MockValidateISBNType func(string) bool

// Custom type that allows setting the func that our Mock Do func will run instead
type MockDoType func(req *http.Request) (*http.Response, error)

type MockGOISBN struct {
	MockGet          MockGetType
	MockValidateISBN MockValidateISBNType
}

// MockClient is the mock client
type MockClient struct {
	MockDo MockDoType
}

// Overriding what the Do function should "do" in our MockClient
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

// Overriding what the Get function should do in our MockClient
func (m *MockGOISBN) Get(isbn string) (*goisbn.Book, error) {
	return m.MockGet(isbn)
}

// Overriding what the ValidateISBN function should do in our MockClient
func (m *MockGOISBN) ValidateISBN(isbn string) bool {
	return m.MockValidateISBN(isbn)
}

func TestGet(t *testing.T) {

	type testCase struct {
		name       string
		desc       string
		isbnValid  bool
		goIsbnBook *goisbn.Book
		expRes     *entities.Book
		expErr     error
		goIsbnErr  error
		crawlerErr error
		jsonResp   string
		respCode   int
	}
	testCases := []testCase{
		{
			name:       "Happy Case",
			desc:       "all ok",
			goIsbnBook: giBook,
			expRes:     book,
			isbnValid:  true,
		},
		{
			name: "Happy Case",
			desc: "return by crawler",
			expRes: &entities.Book{
				ISBN:      "9781784756055",
				Title:     "Unlucky 13",
				Authors:   "James Patterson",
				Publisher: "BB Books",
				ImageURL:  "https://images.isbndb.com/covers/60/55/9781784756055.jpg",
				Source:    "isbndb_crawl",
				Status:    1,
			},
			isbnValid: true,
			goIsbnErr: constant.ErrBookNotFound,
			jsonResp: `<!DOCTYPE html>
				<html lang="en" dir="ltr" prefix="content: http://purl.org/rss/1.0/modules/content/  dc: http://purl.org/dc/terms/  foaf: http://xmlns.com/foaf/0.1/  og: http://ogp.me/ns#  rdfs: http://www.w3.org/2000/01/rdf-schema#  schema: http://schema.org/  sioc: http://rdfs.org/sioc/ns#  sioct: http://rdfs.org/sioc/types#  skos: http://www.w3.org/2004/02/skos/core#  xsd: http://www.w3.org/2001/XMLSchema# ">
				  <head>
					<meta charset="utf-8" />
				<script>(function(i,s,o,g,r,a,m){i["GoogleAnalyticsObject"]=r;i[r]=i[r]||function(){(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)})(window,document,"script","https://www.google-analytics.com/analytics.js","ga");ga("create", "UA-20601258-1", {"cookieDomain":"auto"});ga("set", "anonymizeIp", true);ga("send", "pageview");</script>
				<meta name="Generator" content="Drupal 8 (https://www.drupal.org)" />
				<meta name="MobileOptimized" content="width" />
				<meta name="HandheldFriendly" content="true" />
				<meta name="viewport" content="width=device-width, initial-scale=1.0" />
				<link rel="shortcut icon" href="/sites/default/files/favicon_0.ico" type="image/vnd.microsoft.icon" />
				
					<title>Unlucky 13 | ISBNdb</title>
					<link rel="stylesheet" media="all" href="/sites/default/files/css/css_kMUjRg96yM727DXSMBZEtPNDnuK8QrcmP-J72nlVHH4.css" />
				<link rel="stylesheet" media="all" href="/sites/default/files/css/css_c79fBV0hFcvqEa8YBSWssZgTGTHqeoIRBNKlMJ6nwYs.css" />
				
					
				<!--[if lte IE 8]>
				<script src="/sites/default/files/js/js_VtafjXmRvoUgAzqzYTA3Wrjkx9wcWhjP0G4ZnnqRamA.js"></script>
				<![endif]-->
				
						  <link rel="canonical" href="https://isbndb.com/book/9781784756055" />
					  
				  </head>
				  <body class="path-book">
					<a href="#main-content" class="visually-hidden focusable skip-link">
					  Skip to main content
					</a>
					
					  <div class="dialog-off-canvas-main-canvas" data-off-canvas-main-canvas>
					
				<!-- Header and Navbar -->
				<header class="main-header">
				   
				  <div class="topnav-wrap">
					<div class="container">
					  <div class="row">
										  <div class="col-sm-8 col-md-8 top-user-menu col-sm-offset-4 col-md-offset-4">
							<ul class="list-inline">
									  
								<li>  <div class="region region-user-menu">
					<nav role="navigation" aria-labelledby="block-multipurpose-business-theme-account-menu-menu" id="block-multipurpose-business-theme-account-menu">
							
				  <h2 class="visually-hidden" id="block-multipurpose-business-theme-account-menu-menu">User account menu</h2>
				  
				
						
				
							  <ul class="menu">
										  <li class="menu-item"
									  >
									<a href="/user/login" data-drupal-link-system-path="user/login">Log in</a>
									  </li>
									  <li class="menu-item"
									  >
									<a href="/isbn-database" data-drupal-link-system-path="node/23">Register</a>
									  </li>
						</ul>
				  
				
				  </nav>
				
				  </div>
				</li>
													
								<li>  <div class="region region-search">
					<div id="block-isbndbsearchblock-2" class="block block-isbndb block-isbndb-menu-block">
				  
					
					  <div id="main-search block--isbndb_search">
					<form action="/search/books/" method="GET" class="isbndb_search_block">
						<div class="input-group">
							<input type="hidden" name="search_param" value="books" id="search_param">         
							<input id="search_query" type="text" class="search_query form-control" name="x" placeholder="Search ISBN or Title...">
							<span class="input-group-btn">
								<button id="search-button" class="search-button btn btn-default" type="submit"><span class="glyphicon glyphicon-search"></span></button>
							</span>
						</div>
					</form>
				</div>
				
				  </div>
				
				  </div>
				</li>
										  </ul>            
						  </div>
							  </div>
					</div>
				  </div>
				  
				  <div class="nav-hr hidden-sm"></div>
				
				  <nav class="navbar navbar-default" role="navigation">
					<div class="container">
					  <div class="row">
					  <div class="navbar-header col-md-3">
						<button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#main-navigation">
						  <span></span>
						  <!-- <span class="sr-only">Toggle navigation</span>
						  <span class="icon-bar"></span>
						  <span class="icon-bar"></span>
						  <span class="icon-bar"></span> -->
						</button>
									<div class="region region-header">
					<div id="block-multipurpose-business-theme-branding" class="site-branding block block-system block-system-branding-block">
				  
					
						<div class="brand logo">
					  <a href="/" title="Home" rel="home" class="site-branding__logo">
						<img src="/sites/default/files/ISBN-295x62_0.jpg" alt="Home" />
					  </a>
					</div>
					</div>
				<div id="block-isbndbsearchblock-3" class="block block-isbndb block-isbndb-menu-block">
				  
					
					  <div id="main-search block--isbndb_search">
					<form action="/search/books/" method="GET" class="isbndb_search_block">
						<div class="input-group">
							<input type="hidden" name="search_param" value="books" id="search_param">         
							<input id="search_query" type="text" class="search_query form-control" name="x" placeholder="Search ISBN or Title...">
							<span class="input-group-btn">
								<button id="search-button" class="search-button btn btn-default" type="submit"><span class="glyphicon glyphicon-search"></span></button>
							</span>
						</div>
					</form>
				</div>
				
				  </div>
				
				  </div>
				
							  </div>
				
					  <!-- Navigation -->
					  <div class="col-md-9">
									<div class="region region-primary-menu">
					<nav role="navigation" aria-labelledby="block-multipurpose-business-theme-main-menu-menu" id="block-multipurpose-business-theme-main-menu">
							
				  <h2 class="visually-hidden" id="block-multipurpose-business-theme-main-menu-menu">Main navigation</h2>
				  
				
						
							  <ul class="sm menu-base-theme" id="main-menu"  class="menu nav navbar-nav">
									  <li>
						<a href="/" data-drupal-link-system-path="&lt;front&gt;">Home</a>
								  </li>
								  <li>
						<a href="/isbn-database" data-drupal-link-system-path="node/23">ISBN Database</a>
								  </li>
								  <li>
						<a href="/apidocs/v2" data-target="#" data-toggle="dropdown">Documentation</a>
												  <ul>
									  <li>
						<a href="/apidocs/v2" data-drupal-link-system-path="apidocs/v2">API v2</a>
								  </li>
						</ul>
				  
							</li>
								  <li>
						<a href="/articles" data-drupal-link-system-path="articles">Articles</a>
								  </li>
								  <li>
						<a href="/news" title="News" data-drupal-link-system-path="news">News</a>
								  </li>
								  <li>
						<a href="/contact" data-drupal-link-system-path="contact">Contact</a>
								  </li>
						</ul>
				  
				
				
				  </nav>
				
				  </div>
				
							  
					  </div>
					  <!--End Navigation -->
				
					  </div>
					</div>
				  </nav>
				
				</header>
				<!--End Header & Navbar -->
				
				<div id="feedback-link" class="hidden">
					<a class="use-ajax" data-dialog-options="{&quot;title&quot;:&quot;Feedback &amp; Support&quot;,&quot;width&quot;:500}" data-dialog-type="modal" href="/contact">Feedback &amp; Support</a>
				</div>
				
				<!-- Start Slider -->
				<!-- End Slider -->
				
				
				<!-- Banner -->
				  <!-- End Banner -->
				
				
				<!--Highlighted-->
					  <div class="container">
					  <div class="row">
						<div class="col-md-12">
							<div class="region region-highlighted">
					<div data-drupal-messages-fallback class="hidden"></div>
				
				  </div>
				
						</div>
					  </div>
					</div>
				  <!--End Highlighted-->
				
				
				<!-- Page Title -->
				  <div id="page-title">
					<div id="page-title-inner">
					  <!-- start: Container -->
					  <div class="container">
						  <div class="region region-page-title">
					<div id="block-multipurpose-business-theme-page-title" class="block block-core block-page-title-block">
				  
					
					  
				  <h1>Unlucky 13</h1>
				
				
				  </div>
				
				  </div>
				
					  </div>
					</div>
				  </div>
				<!-- End Page Title ---- >
				
				
				<!-- layout -->
				<div id="wrapper">
				  <!-- start: Container -->
				  <div class="container">
					
					<!--Content top-->
						  <!--End Content top-->
					
					<!--start:content -->
				
					<div class="row layout">
					  <!--- Start Left SideBar -->
							<!---End Right SideBar -->
				
					  <!--- Start content -->
							  <div class="content_layout">
						  <div class=col-md-12>
							  <div class="region region-content">
					<div id="block-multipurpose-business-theme-content" class="block block-system block-system-main-block">
				  
					
						  <div class="container">
						<div class="row">
							<div class="artwork col-xs-12 col-md-3">
				
							<!-- <img src="/sites/default/files/default-book-cover.jpg" style="height:250px; width:190px; background-color:#dddddd"/> -->
								<object height="250px" width="190px" data="https://images.isbndb.com/covers/60/55/9781784756055.jpg" type="image/png">
								 <img height="250px" width="190px" src="/modules/isbndb/img/default-book-cover.jpg" />
								</object>
							</div>
							<div class="book-table col-xs-12 col-md-6">
							  <table class="table table-hover table-responsive ">
												<tr> <th>Full Title</th> <td>Unlucky 13</td> </tr>
																<tr> <th>ISBN</td> <th>1784756059</td> </tr>
																<tr> <th>ISBN13</th> <td>9781784756055</td> </tr>
																<tr> <th>Publisher</th> <td>BB Books</td> </tr>
																<tr> <th>Authors</th> <td>James Patterson</td> </tr>
																								
								
																																								
								
								
											  
							  
							  </table>
							</div>
						</div>
				
					
						<div class="row">
							<div class="col-md-offset-3 col-md-6 col-xs-12">
								<br />
								<br />
								<br />
								<p class="text-center">
								  Need more data?   Get a FREE 7 day trial and get access to the full database of 24 million
								  books and all data points including title, author, publisher, publish date, binding, pages,
								  list price, and more.
									<br /><br /><br />
								 <a href="/isbn-database"><button class="btn btn-green btn-lg">Get Started Free</button></a>                 
								</p>
							</div>
						</div>
				
					
				
				  </div>
				
				  </div>
				
						  </div>
						</div>
							<!---End content -->
				
					  <!--- Start Right SideBar -->
							<!---End Right SideBar -->
					  
					</div>
					<!--End Content -->
				
					<!--Start Content Bottom-->
						<!--End Content Bottom-->
				  </div>
				</div>
				<!-- End layout -->
				
				
				<!-- topwidget Widget -->
				<!--End topwidget Widget -->
				
				
				<!-- Portfolio Widget -->
				<!--End Portfolio Widget -->
				
				
				<!-- Start Skills -->
				<!--End skills -->
				
				
				<!--showcase-->
				<!--End showcase-->
				
				
				<!-- Start pricetable -->
				<!--End pricetable -->
				
				
				<!--- Start testimonials -->
				<!---End testimonials -->
				
				
				<!-- Start features -->
				<!--End features -->
				
				
				<!-- Start bottom -->
				<!--End Bottom -->
				
				
				<!-- start: Footer -->
				  <div class="footerwidget widget-group">
					<div class="container">
					  
					  <div class="row">
				
						<!-- Start Footer First Region -->
						<div class = col-md-6>
										<div class="region region-footer-first">
					<div id="block-footerfirst" class="block block-block-content block-block-content72b2c59c-8ba3-4237-b402-023073b98801">
				  
					  <h2>Quick links</h2>
					
					  
							<div><ul class="footer_links">
					<li><a href="/isbn-database">ISBN Database</a></li>
					<li><a href="/faq">FAQs</a></li>
					<!-- <li><a href="/api/v3/docs">How it Works</a></li> -->
					<li><a href="/user/register">Register</a></li>
					<li><a href="/contact">Contact</a></li>
					<li><a href="/isbndb/books/request">Report ISBN</a></li>
				</ul>
				</div>
					  
				  </div>
				
				  </div>
				
								  </div>
						<!-- End Footer First Region -->
				
						<!-- Start Footer Second Region -->
						<div class = col-md-6>
								  </div>
						<!-- End Footer Second Region -->
				
						<!-- Start Footer third Region -->
						<div class = col-md-6>
										<div class="region region-footer-third">
					<div id="block-footercta" class="block block-block-content block-block-content76716c99-42dd-488d-8589-0967229a7da2">
				  
					
					  
							<div><a href="/user/register" class="btn btn-info" role="button">Subscribe Now</a></div>
					  
				  </div>
				
				  </div>
				
								  </div>
						<!-- End Footer Third Region -->
					  </div>
					</div>
				  </div>
				<!--End Footer -->
				
				<div class="nav-hr hidden-sm"></div>
				
				<div class="copyright">
				  <div class="container">
					<div class="row">
				
					  <!-- Copyright -->
					  <div class="col-sm-6 col-md-6">
						<p>Copyright © 2021. All rights reserved</p>
					  </div>
					  <!-- End Copyright -->
				
					  <!-- Credit link -->
							<!-- End Credit link -->
					  
							  <div class="col-sm-6 col-md-6">
							<div class="region region-footer-menu">
					<nav role="navigation" aria-labelledby="block-multipurpose-business-theme-footer-menu" id="block-multipurpose-business-theme-footer">
							
				  <h2 class="visually-hidden" id="block-multipurpose-business-theme-footer-menu">Footer menu</h2>
				  
				
						
				
							  <ul class="menu">
										  <li class="menu-item"
									  >
									<a href="/terms-and-conditions" title="Terms and Conditions" rel="nofollow" data-drupal-link-system-path="node/36">Terms and Conditions</a>
									  </li>
									  <li class="menu-item"
									  >
									<a href="/privacy-policy" title="Privacy Policy" rel="nofollow" data-drupal-link-system-path="node/32">Privacy Policy</a>
									  </li>
						</ul>
				  
				
				  </nav>
				
				  </div>
				
						</div>
							
					</div>
				  </div>
				</div>
				
				<!-- Google map -->
				<!-- End Google map -->
				
				<!-- QR Scan Modal -->
				<div id="qrscan-modal" class="modal bs-example-modal-sm" data-backdrop="false" tabindex="-1" role="dialog" aria-labelledby="mySmallModalLabel">
				  <div class="modal-dialog modal-sm" role="document">
					<div class="modal-content">
						<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
						<div id="qrscan-target"></div>
					</div>
				  </div>
				</div>
				<!-- End QR Scan Modal -->
				
				  </div>
				
					
					<script type="application/json" data-drupal-selector="drupal-settings-json">{"path":{"baseUrl":"\/","scriptPath":null,"pathPrefix":"","currentPath":"book\/9781784756055","currentPathIsAdmin":false,"isFront":false,"currentLanguage":"en"},"pluralDelimiter":"\u0003","suppressDeprecationErrors":true,"ajaxPageState":{"libraries":"core\/drupal.dialog.ajax,core\/html5shiv,geshifilter\/geshifilter,google_analytics\/google_analytics,isbndb\/isbndb,multipurpose_business_theme\/bootstrap,multipurpose_business_theme\/colorbox,multipurpose_business_theme\/flexslider,multipurpose_business_theme\/global-components,multipurpose_business_theme\/owl,multipurpose_business_theme\/quicksand,system\/base","theme":"multipurpose_business_theme","theme_token":null},"ajaxTrustedUrl":[],"google_analytics":{"trackOutbound":true,"trackMailto":true,"trackDownload":true,"trackDownloadExtensions":"7z|aac|arc|arj|asf|asx|avi|bin|csv|doc(x|m)?|dot(x|m)?|exe|flv|gif|gz|gzip|hqx|jar|jpe?g|js|mp(2|3|4|e?g)|mov(ie)?|msi|msp|pdf|phps|png|ppt(x|m)?|pot(x|m)?|pps(x|m)?|ppam|sld(x|m)?|thmx|qtm?|ra(m|r)?|sea|sit|tar|tgz|torrent|txt|wav|wma|wmv|wpd|xls(x|m|b)?|xlt(x|m)|xlam|xml|z|zip"},"user":{"uid":0,"permissionsHash":"d96912554d22b28e774a5b40ea8b7ba6d5153d050a45da9ef109d68fb4e2c73c"}}</script>
				<script src="/sites/default/files/js/js_F5tv6lobw_RPDEzeSvhnuuoxlgyGY7-qsFZdGHYP6Hc.js"></script>
				
				  </body>
				</html>
				`,
			respCode: 200,
		},
		{
			name:      "Sad Case",
			desc:      "goisbn returns error, crawler return non 2xx response code",
			isbnValid: true,
			goIsbnErr: constant.ErrBookNotFound,
			respCode:  999,
			expErr:    constant.ErrRetrievingBookDetails,
		},
		{
			name:      "Sad Case",
			desc:      "crawler returns empty page",
			isbnValid: true,
			goIsbnErr: constant.ErrBookNotFound,
			jsonResp:  `<html></html>`,
			respCode:  200,
			expErr:    constant.ErrBookNotFound,
		},
		{
			name:       "Sad Case",
			desc:       "goisbn and crawler returns error",
			isbnValid:  true,
			goIsbnErr:  constant.ErrBookNotFound,
			crawlerErr: fmt.Errorf("mock error"),
			respCode:   200,
			expErr:     constant.ErrRetrievingBookDetails,
		},
	}

	for _, v := range testCases {
		gi := &MockGOISBN{
			MockGet: func(string) (*goisbn.Book, error) {
				return v.goIsbnBook, v.goIsbnErr
			},
			MockValidateISBN: func(string) bool {
				return v.isbnValid
			},
		}
		svc := NewBookService(gi)
		svc.client = &MockClient{
			MockDo: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: v.respCode,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(v.jsonResp))),
				}, v.crawlerErr
			},
		}
		actRes, actErr := svc.Get(context.Background(), "dummy isbn")
		assert.Equal(t, v.expRes, actRes)
		assert.Equal(t, v.expErr, actErr)
	}
}

func TestMapBookToEntity(t *testing.T) {
	type testCase struct {
		name   string
		desc   string
		book   *goisbn.Book
		expRes *entities.Book
	}

	testCases := []testCase{
		{
			name:   "Happy case",
			desc:   "all ok",
			book:   giBook,
			expRes: book,
		},
		{
			name: "Happy case",
			desc: "alt mapping",
			book: &goisbn.Book{
				Title:         "DUMMY",
				PublishedYear: "2021",
				Authors:       []string{"kitefishBB"},
				Description:   "dummy description",
				IndustryIdentifiers: &goisbn.Identifier{
					ISBN: "ISBN10",
				},
				PageCount:  999,
				Categories: []string{"dummy", "category"},
				ImageLinks: &goisbn.ImageLinks{
					ImageURL:      "imageURL",
					LargeImageURL: "largeImageURL",
				},
				Publisher: "BB publishing house",
				Language:  "en",
				Source:    "google",
			},
			expRes: &entities.Book{
				ISBN:            "ISBN10",
				Title:           "DUMMY",
				Authors:         "kitefishBB",
				ImageURL:        "imageURL",
				SmallImageURL:   "imageURL",
				PublicationYear: 2021,
				Publisher:       "BB publishing house",
				Status:          1,
				Description:     "dummy description",
				PageCount:       999,
				Categories:      "dummy, category",
				Language:        "en",
				Source:          "google",
			},
		},
		{
			name:   "Happy case",
			desc:   "all ok",
			book:   giBook,
			expRes: book,
		},
		{
			name: "Happy case",
			desc: "alt mapping 2",
			book: &goisbn.Book{
				Title:         "DUMMY",
				PublishedYear: "2021",
				Authors:       []string{"kitefishBB"},
				Description:   "dummy description",
				IndustryIdentifiers: &goisbn.Identifier{
					ISBN: "ISBN10",
				},
				PageCount:  999,
				Categories: []string{"dummy", "category"},
				ImageLinks: &goisbn.ImageLinks{
					SmallImageURL: "smallImageURL",
					LargeImageURL: "largeImageURL",
				},
				Publisher: "BB publishing house",
				Language:  "en",
				Source:    "google",
			},
			expRes: &entities.Book{
				ISBN:            "ISBN10",
				Title:           "DUMMY",
				Authors:         "kitefishBB",
				ImageURL:        "smallImageURL",
				SmallImageURL:   "smallImageURL",
				PublicationYear: 2021,
				Publisher:       "BB publishing house",
				Status:          1,
				Description:     "dummy description",
				PageCount:       999,
				Categories:      "dummy, category",
				Language:        "en",
				Source:          "google",
			},
		},
	}

	for _, v := range testCases {
		actRes := mapBookToEnitiy(v.book)

		assert.Equal(t, v.expRes, actRes)
	}

}
