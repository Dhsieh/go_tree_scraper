package scraper

const (
	ForestryImagePhotoUrl string = "https://www.forestryimages.org/browse/subthumb.cfm?sub=%s&cat=58&systemid=2"
	coniferTreeListUrl    string = "https://api.bugwood.org/rest/api/subject/.json?fmt=datatable&include=count&cat=58&systemid=2&draw=1&columns%5B0%5D%5Bdata%5D=0&columns%5B0%5D%5Bsearchable%5D=false&columns%5B0%5D%5Borderable%5D=false&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bdata%5D=1&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=true&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bdata%5D=2&columns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=true&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bdata%5D=3&columns%5B3%5D%5Bsearchable%5D=false&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&order%5B0%5D%5Bcolumn%5D=1&order%5B0%5D%5Bdir%5D=asc&start=0&length=126&search%5Bvalue%5D=&_=1616476429774"
	hardwoodtreeListUrl   string = "https://api.bugwood.org/rest/api/subject/.json?fmt=datatable&include=count&cat=57&systemid=2&draw=1&columns%5B0%5D%5Bdata%5D=0&columns%5B0%5D%5Bsearchable%5D=false&columns%5B0%5D%5Borderable%5D=false&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bdata%5D=1&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=true&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bdata%5D=2&columns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=true&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bdata%5D=3&columns%5B3%5D%5Bsearchable%5D=false&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&order%5B0%5D%5Bcolumn%5D=1&order%5B0%5D%5Bdir%5D=asc&start=0&length=126&search%5Bvalue%5D=&_=1616477521940"

	plantImageListUrl string = "https://api.bugwood.org/rest/api/image.json?includeonly=imgnum&length=100&sub=%s"
	imageUrl          string = "https://bugwoodcloud.org/images/384x256/%s.jpg"

	imagenum string = "imgnum"

	BingPhotoUrl string = "https://www.bing.com/images/async?q={}&first=1&mmasync=1"

	DownloadFolder string = "downloads"
)
