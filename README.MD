## About the source 
- the code for scraping result data is in `/scrape`
- the code for the front-end is in `/web-app`
- There's are two Github action jobs: one to refresh data from official website, another to update branch for a roll number. Since GitHub actions will be cancelled after a certain time, the former has a year parameter to load new data for that year only
## How it works
- the captcha is a joke
- for not getting blocked, the requests for fetching html are dispatched after a delay and retry logic handles if we fail because of too many requests.
- the code for branches and roll numbers is located in the file `/scrape/roll_numbers.go`, So if roll numbers are not appearing, this file will be required to be modified
- Most of the scraping logic is in the `/scrape/cmd/fetch_result.go` , which when run, scrapes data and puts records in `/scrape/result.db` sqlite file. It updates the data in db file, it doesn't rewrite the whole file.
- We can't directly use one sqlite file in a static webapp(with crazy workarounds probably we can), So the executable from `/scrape/cmd/gen_json.go` generates three json files from the sqlite db records
	- `/scrape/detailed_result.json`
	- `/scrape/ranks_result.json`
	- `/scrape/result_config.json`
- These three files are required to be put in `/web-app/static/`
- The code is automatically deployed through vercel if setup correctly
- TODO: Don't know why web-app suddenly requires bigger `--max_old_space_size`, any workaround for this?

## Branch change or missing roll numbers
- report on this [issue](https://github.com/rahul-gill/NITH-Result/issues/7)
### How to update branch change info
- update the db table first with correct branch with
```sql
UPDATE student
set branch = :new_branch
WHERE roll_number = :roll_number;
```
- re-run the `gen_json.go` executable to re-generate the json files
- copy json files from scrape to `web-app/static`
- modify the `scrape/roll_numbers.go` and  add this roll number and correct branch in `BranchExceptionRollNumbers` map


![Image](https://github.com/user-attachments/assets/4d1b97d4-7beb-49af-aaff-868579641a7d)
