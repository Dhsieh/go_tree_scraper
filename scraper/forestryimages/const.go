package forestryscraper

const (
	// Photo url template to scrape a image
	PhotoUrl string = "https://www.forestryimages.org/browse/subthumb.cfm?sub=%s&cat=58&systemid=2"
	// url to get all the confiers
	coniferTreeListUrl string = "https://api.bugwood.org/rest/api/subject/.json?fmt=datatable&include=count&cat=58&columns%5B0%%205D%5Bdata%5D=0&columns%5B0%5D%5Bsearchable%5D=false&columns%5B0%5D%5Borderable%5D=false&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bda%20ta%5D=1&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=true&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bdata%5D=2&c%20olumns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=true&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bdata%5D=3&columns%5B%203%5D%5Bsearchable%5D=false&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&order%5B0%5D%5Bcolumn%5D=1&page=1"
	// url to get part 1 of all the deciduous (hardwood)
	hardwoodtreeListUrl1 string = "https://api.bugwood.org/rest/api/subject/.json?fmt=datatable&include=count&cat=57&columns%5B0%%205D%5Bdata%5D=0&columns%5B0%5D%5Bsearchable%5D=false&columns%5B0%5D%5Borderable%5D=false&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bda%20ta%5D=1&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=true&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bdata%5D=2&c%20olumns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=true&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bdata%5D=3&columns%5B%203%5D%5Bsearchable%5D=false&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&order%5B0%5D%5Bcolumn%5D=1&page=1"
	// url to get part 2 of all the deciduous (hardwood)
	hardwoodtreeListUrl2 string = "https://api.bugwood.org/rest/api/subject/.json?fmt=datatable&include=count&cat=57&columns%5B0%%205D%5Bdata%5D=0&columns%5B0%5D%5Bsearchable%5D=false&columns%5B0%5D%5Borderable%5D=false&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bda%20ta%5D=1&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=true&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bdata%5D=2&c%20olumns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=true&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bdata%5D=3&columns%5B%203%5D%5Bsearchable%5D=false&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&order%5B0%5D%5Bcolumn%5D=1&page=2"

	// Url template that will list all the photos
	plantImageListUrl string = "https://api.bugwood.org/rest/api/image.json?includeonly=imgnum&length=100&sub=%d"
	imageUrl          string = "https://bugwoodcloud.org/images/384x256/%s.jpg"

	imagenum string = "imgnum"

	// old configerTreeListUrl
	// https://api.bugwood.org/rest/api/subject/.json?fmt=datatable&include=count&cat=58&systemid=2&draw=1&columns%5B0%    5D%5Bdata%5D=0&columns%5B0%5D%5Bsearchable%5D=false&columns%5B0%5D%5Borderable%5D=false&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bda    ta%5D=1&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=true&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bdata%5D=2&c    olumns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=true&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bdata%5D=3&columns%5B    3%5D%5Bsearchable%5D=false&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&order%5B0%5D%5Bcolumn%5D=1&order%5B0%5D%5Bdir    %5D=asc&start=0&length=126&search%5Bvalue%5D=&_=161647642977

)
