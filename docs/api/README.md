# API collection

`survey-spark.postman_collection.json` — Postman v2.1. Imports into Postman,
Insomnia, Bruno and Thunder Client.

## Import and run

1. Postman → **Import** → select the JSON file.
2. Open the collection's **Variables** tab and set `admin_email` and
   `admin_password` to the `ADMIN_EMAIL` / `ADMIN_PASSWORD` you configured on
   the server. `base_url` defaults to the Render deployment.
3. Run the folders **in numbered order**. They are ordered by dependency:

   | Folder | Why it is where it is |
   | --- | --- |
   | 1 Setup | A country must be active before anyone can register |
   | 2 Auth | Creates the interviewer and respondent, stores their tokens |
   | 3 Admin | Grants points; publishing escrows them and a new account has none |
   | 4 Shared | `/me`, regions |
   | 5 Interviewer | Authoring, publishing, analytics |
   | 6 Respondent | Answering, saved list, completed |
   | 7 Cleanup | Deletes the survey. Kept separate so a full run does not destroy what folder 6 needs |

Each login stores its own token (`admin_token`, `interviewer_token`,
`respondent_token`) and every folder is bound to the right one, so there is no
token copying. `survey_id`, `question_id` and `category_id` are captured from
responses as you go.

## Running the whole thing headless

```bash
npx newman run docs/api/survey-spark.postman_collection.json \
  --env-var base_url=http://127.0.0.1:8080 \
  --env-var admin_email=admin@example.com \
  --env-var admin_password=<password>
```

A clean run is 42 requests, 0 failures. It is safe to re-run: registration
tolerates 409, and the survey is recreated each pass.

## Notes

- Numeric fields (`country_id`, `survey_id`, `question_id`, ...) are unquoted
  in request bodies on purpose. Go binds them as `uint`, so a quoted
  `"{{country_id}}"` is sent as a string and rejected with 400.
- `region_id` defaults to 1 (Yangon). A respondent only sees surveys matching
  their own `country_id` **and** `region_id`.
