## Aboout the source 
- the code for scraping result data is in `/scrape`
- the code for the front-end is in `/web-app`

## How it works
- the captcha is a joke
- for not getting blocked, the requests for fetching html are spaced by a timeframe
- the code for branches and roll numbers is located in the file `/scrape/roll_numbers.go`, So if roll numbers are not apeearing, this file will be required to be modified
- Most of the scraping logic is in the `/scrape/cmd/fetch_result.go` that compiles to an executable
- that executable, when run, scrapes data and puts records in `/scrape/result.db` sqlite file. 
- This operation takes around 3-4 hours, So I've made a github action for it
- We can't directly use an sqlite file in a static webapp(with crazy workarounds probably we can), So the executable from `/scrape/cmd/gen_json.go` generates three json files from the sqlite db recors
	- `/scrape/detailed_result.json`
	- `/scrape/ranks_result.json`
	- `/scrape/result_config.json`
- The db to json step can be put in the github action: TODO
- These three files are required to be put in `/web-app/static/`
- The code is automatically deploy on vercel if setup correctly
- The code can be cleaned up: TODO
- The last change somehow made the desktop layout shitty TODO
- Don't know why web-app suddenly requires `export NODE_OPTIONS=--max_old_space_size=4096`
- There's a bug, branches are not getting set properly with new roll number scheme(like `20BCE`)
